package router

import (
	web "EvoBot"
	"EvoBot/backend/global"
	_router "EvoBot/backend/router"
	"EvoBot/cmd/docs"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type StaticFSWrapper struct {
	http.FileSystem
	FixedModTime time.Time
}

func (f *StaticFSWrapper) Open(name string) (http.File, error) {
	file, err := f.FileSystem.Open(name)

	return &StaticFileWrapper{File: file, fixedModTime: f.FixedModTime}, err
}

type StaticFileWrapper struct {
	http.File
	fixedModTime time.Time
}

func (f *StaticFileWrapper) Stat() (os.FileInfo, error) {

	fileInfo, err := f.File.Stat()

	return &StaticFileInfoWrapper{FileInfo: fileInfo, fixedModTime: f.fixedModTime}, err
}

type StaticFileInfoWrapper struct {
	os.FileInfo
	fixedModTime time.Time
}

func (f *StaticFileInfoWrapper) ModTime() time.Time {
	return f.fixedModTime
}

func getUIAssetFs() http.FileSystem {
	uiAssetFs, err := fs.Sub(web.UIFs, "frontend/dist/assets")
	if err != nil {

		return &StaticFSWrapper{
			FileSystem:   http.Dir("./frontend/dist/assets"),
			FixedModTime: time.Now(),
		}
	}
	return &StaticFSWrapper{
		FileSystem:   http.FS(uiAssetFs),
		FixedModTime: time.Now(),
	}
}

func Init() *gin.Engine {
	if global.CONF.LogConfig.Model != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	eng := gin.New()
	eng.Use(gin.Recovery())
	eng.Use(gin.Logger())

	eng.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/ui/") {
			ctx.FileFromFS("frontend/dist/", http.FS(web.UIFs))
			return
		}
		ctx.JSON(404, gin.H{"error": "not found"})
	})

	// swaggerRouter := r.Group("evobot").Use(middleware.AuthRequired())
	docs.SwaggerInfo.BasePath = "/api/v1"
	eng.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	eng.StaticFS("/ui/assets", getUIAssetFs())
	eng.StaticFileFS("/ui/favicon.ico", "frontend/dist/favicon.ico", http.FS(web.UIFs))

	eng.Static("/api/v1/images", "./uploads")
	eng.GET("/health", func(c *gin.Context) {
		c.JSON(200, "ok")
	})

	indexGroup := eng.Group("/")
	{
		indexGroup.Use(gzip.Gzip(gzip.DefaultCompression))
		indexGroup.GET("/", func(ctx *gin.Context) {
			ctx.FileFromFS("frontend/dist/", http.FS(web.UIFs))
		})
	}

	uiGroup := eng.Group("/ui")
	{
		uiGroup.Use(gzip.Gzip(gzip.DefaultCompression))
		uiGroup.GET("/", func(ctx *gin.Context) {
			ctx.FileFromFS("frontend/dist/", http.FS(web.UIFs))
		})
	}

	PrivateGroup := eng.Group("/api/v1")
	for _, router := range _router.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	global.ZAPLOG.Info("init router success")
	return eng
}
