package route

import (
	"new-api-demo/controller"
	"new-api-demo/middleware"
	"new-api-demo/types"

	"github.com/gin-gonic/gin"
)

func SetRouter(router *gin.Engine) {
	SetRelayRouter(router)
}

func SetRelayRouter(router *gin.Engine) {
	//1.cors
	router.Use(middleware.CORS())
	//2.请求体比较大，可以压缩过来，这里解压
	router.Use(middleware.DecompressRequestMiddleware())
	//3.统计
	router.Use(middleware.StatsMiddleware())

	// 设置路由组
	relayV1Router := router.Group("/v1")
	//relayV1Router.Use(middleware.TokenAuth())
	//relayV1Router.Use(middleware.ModelRequestRateLimit())
	{

		//http router
		httpRouter := relayV1Router.Group("")
		httpRouter.Use(middleware.Distribute()) // 分发=>核心业务逻辑

		// chat related routes
		httpRouter.POST("/completions", func(c *gin.Context) {
			controller.Relay(c, types.RelayFormatOpenAI)
		})
		httpRouter.POST("/chat/completions", func(c *gin.Context) {
			controller.Relay(c, types.RelayFormatOpenAI)
		})
	}
}
