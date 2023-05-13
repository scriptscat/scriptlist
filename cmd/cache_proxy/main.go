package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type cacheItem struct {
	content      []byte
	lastModified string
}

var (
	proxyURL = "http://43.155.80.229"
	cacheMap = make(map[string]*cacheItem)
	mu       sync.RWMutex
)

// http://127.0.0.1:8080/scripts/code/367/OCS%20%E7%BD%91%E8%AF%BE%E5%8A%A9%E6%89%8B.user.js
func main() {
	http.HandleFunc("/scripts/", handleRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Create a new request based on the original to modify headers for proxying
	newReq, err := http.NewRequestWithContext(r.Context(),
		r.Method, proxyURL+r.RequestURI, r.Body,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range r.Header {
		newReq.Header[k] = v
	}
	newReq.Host = "scriptcat.org"
	newReq.Header.Set("X-Real-IP", r.RemoteAddr)
	newReq.Header.Set("X-Forwarded-For", r.Header.Get("X-Forwarded-For"))
	newReq.Header.Set("X-Forwarded-Proto", r.Header.Get("X-Forwarded-Proto"))
	newReq.Header.Set("REMOTE-HOST", r.RemoteAddr)
	newReq.Header.Set("Upgrade", r.Header.Get("Upgrade"))
	newReq.Header.Set("Connection", r.Header.Get("Connection"))
	newReq.Header.Set("Cookie", r.Header.Get("Cookie"))

	mu.RLock()
	item, exists := cacheMap[r.RequestURI]
	mu.RUnlock()

	if exists {
		newReq.Header.Set("If-Modified-Since", item.lastModified)
	}
	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified && exists {
		w.Write(item.content)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastModified := resp.Header.Get("Last-Modified")

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		w.Write(body)
		return
	}

	if lastModified == "" {
		w.Write(body)
		return
	}

	mu.Lock()
	cacheMap[r.RequestURI] = &cacheItem{
		content:      body,
		lastModified: lastModified,
	}
	mu.Unlock()
	w.Write(body)
}
