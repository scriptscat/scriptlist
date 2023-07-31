package open_ctr

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
)

type Open struct {
}

func NewOpen() *Open {
	return &Open{}
}

// CrxDownload 谷歌crx下载服务
func (o *Open) CrxDownload() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		version := ctx.Query("version")
		// 从google chrome商店查询版本信息
		info, err := o.crxDetail(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "查询版本信息失败",
			})
			return
		}
		filepath := "./resource/open/crx/" + info.Name + "_" + id + "_" + info.Version + ".crx"
		filename := info.Name + "_" + id + "_" + info.Version + ".crx"
		if version != "" && version != info.Version {
			// 从仓库中查询历史版本
			filepath = "./resource/open/crx/" + info.Name + "_" + id + "_" + version + ".crx"
			filename = info.Name + "_" + id + "_" + version + ".crx"
			f, err := os.Open(filepath)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"message": "未找到该版本",
				})
				return
			}
			defer f.Close()
			ctx.Writer.Header().
				Set("Content-Disposition", `attachment; filename="`+
					url.PathEscape(filename)+`"`)
			_, err = io.Copy(ctx.Writer, f)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
			return
		}
		var r io.Reader
		// 判断是否有文件缓存,有则直接返回,没有则下载
		_, err = os.Stat(filepath)
		if err != nil {
			if !os.IsNotExist(err) {
				httputils.HandleResp(ctx, err)
				return
			}
			// 下载crx文件
			if err := os.MkdirAll("./resource/open/crx/", 0755); err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
			resp, err := o.crxDownload(ctx, id)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()
			f, err := os.Create(filepath)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
			defer f.Close()
			r = io.TeeReader(resp.Body, f)
		} else {
			f, err := os.Open(filepath)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
			defer f.Close()
			r = f
		}
		ctx.Writer.Header().
			Set("Content-Disposition", `attachment; filename="`+
				url.PathEscape(filename)+`"`)

		_, err = io.Copy(ctx.Writer, r)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}

	}
}

type detail struct {
	Name    string
	Version string
}

func (o *Open) crxDownload(ctx context.Context, id string) (*http.Response, error) {
	// https://clients2.google.com/service/update2/crx?response=redirect&prodversion=114.0.1823.43&acceptformat=crx2,crx3&x=id%3D${result[1]}%26uc&nacl_arch=${nacl_arch}
	resp, err := http.Get("https://clients2.google.com/service/update2/crx?response=redirect&prodversion=114.0.1823.43&acceptformat=crx2,crx3&x=id%3D" + id + "%26uc&nacl_arch=")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *Open) crxDetail(ctx context.Context, id string) (*detail, error) {
	ret := &detail{}
	if err := cache.Ctx(ctx).GetOrSet("open:crx:detail:"+id, func() (interface{}, error) {
		resp, err := http.Get("https://chrome.google.com/webstore/detail/" + id)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		version := regexp.MustCompile(`<meta\s*itemprop="version"\s*content="(.*?)"\s*/>`).FindStringSubmatch(string(body))
		if len(version) > 0 {
			ret.Version = version[1]
		}
		name := regexp.MustCompile(`<meta\s*itemprop="name"\s*content="(.*?)"\s*/>`).FindStringSubmatch(string(body))
		if len(name) > 0 {
			ret.Name = name[1]
		}
		return ret, nil
	}, cache.Expiration(time.Hour)).Scan(ret); err != nil {
		return nil, err
	}
	return ret, nil
}
