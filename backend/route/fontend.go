package route

import (
	"crypto/md5"
	"encoding/hex"
	"io/fs"
	"net/http"
	"strings"

	"github.com/apicat/apicat/v2/frontend"

	"github.com/gin-gonic/gin"
)

func frontendHandle(r *gin.Engine) {

	indexBytes, err := frontend.Dist.ReadFile("dist/index.html")
	if err != nil {
		panic("")
	}

	assets, err := fs.Sub(frontend.Dist, "dist/assets")
	if err != nil {
		panic("")
	}
	r.StaticFS("/assets", http.FS(assets))

	// 配合前端单页模式history路由
	// index.html 缓存页面请求直接使用本地缓存
	// 计算页面md5 当做etag 用来告诉浏览器资源版本
	indexmd5 := md5.Sum(indexBytes)
	etag := hex.EncodeToString(indexmd5[:])
	r.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodGet && !strings.HasPrefix(ctx.Request.URL.Path, "/api") {
			ctx.Header("ETag", etag)
			ctx.Header("Cache-Control", "no-cache")
			if match := ctx.GetHeader("If-None-Match"); match != "" {
				if strings.Contains(match, etag) {
					ctx.Status(http.StatusNotModified)
					return
				}
			}
			t := struct {
				Title string
			}{getPageTitle(ctx.Request.URL.Path)}
			ctx.HTML(http.StatusOK, "index.html", t)
		} else {
			http.NotFound(ctx.Writer, ctx.Request)
		}
	})
}

func getPageTitle(path string) string {
	return "请修改 title" + path
}
