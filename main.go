package main

import (
	"log"
	"new-api-demo/model"
	"new-api-demo/route"
	"new-api-demo/service"

	"github.com/gin-gonic/gin"
)

func main() {
	model.InitChannelCache()
	service.InitHttpClient()
	service.InitTokenEncoders()
	// 创建 Gin 路由
	router := gin.Default()
	// 设置路由组
	route.SetRouter(router)

	// 启动服务器
	port := ":13000"
	log.Printf("Demo server starting on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  POST http://localhost%s/v1/completions", port)
	log.Printf("  POST http://localhost%s/v1/chat/completions", port)

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
