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
			//controller.Relay(c, types.RelayFormatOpenAI)
			chatCompletions(c)
		})
		httpRouter.POST("/chat/completions", func(c *gin.Context) {
			//controller.Relay(c, types.RelayFormatOpenAI)
			chatCompletions(c)
		})
	}
}

// ChatCompletions godoc
// @Summary Chat Completions
// @Description 支持流式与非流式聊天补全。
// @Tags 聊天
// @Accept json
// @Produce json
// @Param request body ChatCompletionRequest true "OpenAI ChatCompletion 请求体"
// @Success 200 {object} ChatCompletionResponse "非流式 JSON 响应"
// @Success 200 {string} ChatCompletionStreamResponse "text/event-stream 流式响应"
// @Router /v1/chat/completions [post]
func chatCompletions(c *gin.Context) {
	controller.Relay(c, types.RelayFormatOpenAI)
}
