package main

import (
	"fmt"
	"log"
	"net/http"
	"new-api-demo/model"
	"new-api-demo/route"
	"new-api-demo/service"
	"time"

	"github.com/gin-gonic/gin"
)

// 简化的请求结构
type SimpleRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages,omitempty"`
	Prompt   string    `json:"prompt,omitempty"`
	Stream   bool      `json:"stream,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 简化的响应结构
type SimpleResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// 处理 /v1/completions 接口
func handleCompletions(c *gin.Context) {
	var req SimpleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid request format",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 验证必需字段
	if req.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "model is required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 生成简单的响应
	response := SimpleResponse{
		ID:      fmt.Sprintf("cmpl-%d", time.Now().Unix()),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "This is a demo response for completions endpoint.",
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}

	c.JSON(http.StatusOK, response)
}

// 处理 /v1/chat/completions 接口
func handleChatCompletions(c *gin.Context) {
	var req SimpleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid request format",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 验证必需字段
	if req.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "model is required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "messages is required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 生成简单的响应
	response := SimpleResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "This is a demo response for chat completions endpoint.",
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}

	c.JSON(http.StatusOK, response)
}

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
