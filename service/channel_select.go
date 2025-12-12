package service

import (
	"new-api-demo/model"

	"github.com/gin-gonic/gin"
)

func CacheGetRandomSatisfiedChannel(c *gin.Context, modelName string, retry int) (*model.Channel, error) {
	return model.GetRandomSatisfiedChannel(modelName, retry)
}
