package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/cago-frame/cago"
	"github.com/cago-frame/cago/configs"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/cago-frame/cago/server/mux"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type cacheItem struct {
	content      []byte
	lastModified string
}

var (
	proxyURL = "http://127.0.0.1"
	cacheMap = make(map[string]*cacheItem)
	mu       sync.RWMutex
)

// http://127.0.0.1:8080/scripts/code/367/OCS%20%E7%BD%91%E8%AF%BE%E5%8A%A9%E6%89%8B.user.js
func main() {
	cfg, err := configs.NewConfig("cache_proxy")
	if err != nil {
		log.Fatalf("new config err: %v", err)
	}
	proxyURL = cfg.String(context.Background(), "proxy_url")
	if err := cago.New(context.Background(), cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		RegistryCancel(mux.HTTP(func(ctx context.Context, r *mux.Router) error {
			r.Any("/scripts/*path", handleRequest)
			return nil
		})).Start(); err != nil {
		log.Fatalf("start err: %v", err)
	}
	if err != nil {
		panic(err)
	}
}

func handleRequest(c *gin.Context) {
	// Create a new request based on the original to modify headers for proxying
	newReq, err := http.NewRequestWithContext(c.Request.Context(),
		c.Request.Method, proxyURL+c.Request.RequestURI, c.Request.Body,
	)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	//for k, v := range c.Request.Header {
	//	newReq.Header[k] = v
	//}
	newReq.Host = "scriptcat.org"
	if cookie := c.Request.Header.Get("Cookie"); cookie != "" {
		newReq.Header.Set("Cookie", cookie)
	}
	if rip := c.Request.Header.Get("X-Real-IP"); rip != "" {
		newReq.Header.Set("X-Real-IP", rip)
	}
	if xff := c.Request.Header.Get("X-Forwarded-For"); xff != "" {
		newReq.Header.Set("X-Forwarded-For", xff)
	}
	if xfp := c.Request.Header.Get("X-Forwarded-Proto"); xfp != "" {
		newReq.Header.Set("X-Forwarded-Proto", xfp)
	}
	if ua := c.Request.Header.Get("User-Agent"); ua != "" {
		newReq.Header.Set("User-Agent", ua)
	}
	mu.RLock()
	item, exists := cacheMap[c.Request.RequestURI]
	mu.RUnlock()

	if exists {
		newReq.Header.Set("If-Modified-Since", item.lastModified)
	}
	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode == http.StatusNotModified && exists {
		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write(item.content)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	lastModified := resp.Header.Get("Last-Modified")

	if resp.StatusCode != http.StatusOK {
		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}
		c.Writer.WriteHeader(resp.StatusCode)
		_, _ = c.Writer.Write(body)
		return
	}

	if lastModified == "" {
		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}
		c.Writer.WriteHeader(resp.StatusCode)
		_, _ = c.Writer.Write(body)
		return
	}

	mu.Lock()
	logger.Default().
		Info("response content-type",
			zap.String("content-type", resp.Header.Get("Content-Type")),
			zap.String("uri", c.Request.RequestURI))
	// 如果respond type是gzip类型,则不缓存
	if resp.Header.Get("Content-Type") == "application/x-gzip" {
		mu.Unlock()
		logger.Default().Error("gzip type, not cache",
			zap.String("header", resp.Header.Get("Content-Type")),
			zap.ByteString("body", body),
		)
		// gzip解压
		r, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			logger.Default().Error("gzip.NewReader error", zap.Error(err))
			return
		}
		// 转发为test/plain
		for k, v := range resp.Header {
			for _, vv := range v {
				if k == "Content-Type" {
					c.Writer.Header().Add(k, "text/plain")
					break
				}
				c.Writer.Header().Add(k, vv)
			}
		}
		c.Writer.WriteHeader(resp.StatusCode)
		_, err = io.Copy(c.Writer, r) //nolint:gosec
		if err != nil {
			logger.Default().Error("io.Copy error", zap.Error(err))
		}
		return
	}
	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)

	cacheMap[c.Request.RequestURI] = &cacheItem{
		content:      body,
		lastModified: lastModified,
	}
	mu.Unlock()
	_, _ = c.Writer.Write(body)
}
