package open_ctr

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	// 注册常见图片格式解码器
	_ "image/gif"
	_ "image/jpeg"

	_ "github.com/biessek/golang-ico"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/cago-frame/cago/pkg/sync"
	"go.uber.org/zap"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/cago-frame/cago/database/cache"
	"github.com/gin-gonic/gin"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

// faviconData 存储缓存的favicon数据
type faviconData struct {
	Data        string `json:"data"`         // base64编码的图标数据
	ContentType string `json:"content_type"` // 图标的Content-Type
}

// 支持的图标尺寸
var allowedSizes = map[int]bool{
	16: true, 32: true, 48: true, 64: true, 128: true,
}

// linkIconRegexp 匹配HTML中的<link rel="icon">标签
var linkIconRegexp = regexp.MustCompile(`(?i)<link[^>]+rel=["'](?:shortcut icon|icon)["'][^>]*>`)
var hrefRegexp = regexp.MustCompile(`(?i)href=["']([^"']+)["']`)

// Favicon 获取网站favicon图标
func (o *Open) Favicon() gin.HandlerFunc {
	lock := sync.NewLocker("lock:open:favicon")
	return func(ctx *gin.Context) {
		domain := ctx.Query("domain")
		szStr := ctx.DefaultQuery("sz", "")

		// 校验sz参数
		sz := 0
		var err error
		if szStr != "" {
			sz, err = strconv.Atoi(szStr)
			if err != nil || !allowedSizes[sz] {
				sz = 32
			}
		}

		// 校验domain参数
		if domain == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "domain is required"})
			return
		}
		// 使用publicsuffix校验域名合法性
		_, err = publicsuffix.Domain(domain)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid domain"})
			return
		}

		// 从缓存获取或抓取favicon(缓存已缩放的图片)
		cacheKey := fmt.Sprintf("open:favicon:%s:%d", domain, sz)
		result := &faviconData{}

		// 对同一域名循环尝试加锁,避免并发请求同一个域名
		lockKey := domain
		for i := 0; i < 50; i++ {
			if err := lock.TryLockKey(ctx, lockKey, sync.WithLockTimeout(time.Second*10)); err != nil {
				// 获取锁失败,说明有其他请求正在抓取,等待一段时间后重试
				select {
				case <-ctx.Done():
					ctx.JSON(http.StatusNotFound, gin.H{"message": "favicon not found"})
					return
				case <-time.After(200 * time.Millisecond):
					continue
				}
			}
			break
		}
		defer func() {
			_ = lock.UnlockKey(ctx, lockKey)
		}()

		if err := cache.Ctx(ctx).GetOrSet(cacheKey, func() (any, error) {
			data, contentType, err := o.fetchFavicon(ctx, domain)
			if data == nil {
				// 获取失败,设置为空值,避免频繁抓取
				logger.Ctx(ctx).Warn("获取favicon失败",
					zap.String("domain", domain), zap.Error(err))
				return &faviconData{}, nil
			}
			// 缩放到目标尺寸
			if sz != 0 {
				resizedData, err := resizeFavicon(data, sz)
				if resizedData == nil {
					// 缩放失败,使用原始数据
					logger.Ctx(ctx).Warn("缩放favicon失败,使用原始数据",
						zap.String("domain", domain), zap.Int("size", sz), zap.Error(err))
					return &faviconData{
						Data:        base64.StdEncoding.EncodeToString(data),
						ContentType: contentType,
					}, nil
				}
				data = resizedData
			}
			return &faviconData{
				Data:        base64.StdEncoding.EncodeToString(data),
				ContentType: "image/png",
			}, nil
		}, cache.Expiration(24*time.Hour)).Scan(result); err != nil {
			// 缓存出错,直接返回404
			ctx.JSON(http.StatusNotFound, gin.H{"message": "favicon not found"})
			return
		}

		if result.Data == "" {
			// 没有数据,直接返回404
			ctx.JSON(http.StatusNotFound, gin.H{"message": "favicon not found"})
			return
		}

		// 解码base64数据
		imgData, err := base64.StdEncoding.DecodeString(result.Data)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "favicon not found"})
			return
		}

		// 设置响应头并返回图标
		ctx.Writer.Header().Set("Content-Type", result.ContentType)
		ctx.Writer.Header().Set("Cache-Control", "max-age=86400")
		_, _ = ctx.Writer.Write(imgData)
	}
}

// resizeFavicon 将图片数据缩放到指定尺寸,输出为PNG格式
func resizeFavicon(data []byte, sz int) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	// 如果原图尺寸已经和目标一致,直接返回原始数据
	bounds := src.Bounds()
	if bounds.Dx() == sz && bounds.Dy() == sz {
		return data, nil
	}

	// 使用高质量缩放(CatmullRom)
	dst := image.NewRGBA(image.Rect(0, 0, sz, sz))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}
	return buf.Bytes(), nil
}

// fetchFavicon 抓取网站favicon
func (o *Open) fetchFavicon(ctx context.Context, domain string) ([]byte, string, error) {
	// 1. 尝试直接获取 /favicon.ico
	data, contentType, err := o.downloadFavicon(ctx, "https://"+domain+"/favicon.ico")
	if err == nil && data != nil {
		return data, contentType, nil
	}

	// 2. 请求首页,解析HTML中的<link>标签
	data, contentType, err = o.fetchFaviconFromHTML(ctx, domain)
	if err == nil && data != nil {
		return data, contentType, nil
	}

	// 3. 所有方式都失败
	return nil, "", fmt.Errorf("failed to fetch favicon for %s", domain)
}

// fetchFaviconFromHTML 从网站首页HTML中解析favicon链接并下载
func (o *Open) fetchFaviconFromHTML(ctx context.Context, domain string) ([]byte, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+domain, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ScriptCat/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 限制读取大小为512KB(只需要head部分)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return nil, "", err
	}

	html := string(body)

	// 解析<link>标签
	matches := linkIconRegexp.FindAllString(html, -1)
	for _, match := range matches {
		hrefMatch := hrefRegexp.FindStringSubmatch(match)
		if len(hrefMatch) < 2 {
			continue
		}
		iconURL := hrefMatch[1]

		// 处理相对URL
		iconURL = o.resolveURL(iconURL, domain)
		if iconURL == "" {
			continue
		}

		// 下载图标
		data, contentType, err := o.downloadFavicon(ctx, iconURL)
		if err == nil && data != nil {
			return data, contentType, nil
		}
	}

	return nil, "", fmt.Errorf("no favicon found in HTML")
}

// resolveURL 将相对URL解析为绝对URL
func (o *Open) resolveURL(href string, domain string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}
	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}
	if strings.HasPrefix(href, "/") {
		return "https://" + domain + href
	}
	return "https://" + domain + "/" + href
}

// downloadFavicon 下载favicon图标文件
func (o *Open) downloadFavicon(ctx context.Context, url string) ([]byte, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ScriptCat/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 检查Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !isImageContentType(contentType) {
		// 尝试通过内容检测
		data, err := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024)) // 1MB限制
		if err != nil {
			return nil, "", err
		}
		detected := http.DetectContentType(data)
		if !isImageContentType(detected) {
			return nil, "", fmt.Errorf("not an image: %s", detected)
		}
		return data, detected, nil
	}

	// 限制图标大小为1MB
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024))
	if err != nil {
		return nil, "", err
	}

	if len(data) == 0 {
		return nil, "", fmt.Errorf("empty response")
	}

	return data, contentType, nil
}

// isImageContentType 检查Content-Type是否为图片类型
func isImageContentType(ct string) bool {
	ct = strings.ToLower(ct)
	return strings.HasPrefix(ct, "image/") || strings.Contains(ct, "application/octet-stream")
}
