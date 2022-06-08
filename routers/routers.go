package routers

import (
	"net/http"
	"web_app/logger"
	"web_app/settings"

	_ "web_app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/**
 * @Author ouzhsh
 * @Description //TODO 设置路由
 * @Date 21:58 2022/6/2
 **/
func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置成发布模式
	}

	// 初始化gin的引擎 r
	r := gin.New()

	// 注册中间件
	r.Use(gin.Logger(),logger.GinLogger(), logger.GinRecovery(true) )

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, settings.Conf.Version)
	})

	// 配置swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
