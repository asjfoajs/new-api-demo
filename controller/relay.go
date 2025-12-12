package controller

import (
	"bytes"
	"fmt"
	"io"
	"new-api-demo/constant"
	"new-api-demo/relay"

	"net/http"
	"new-api-demo/common"
	"new-api-demo/logger"
	"new-api-demo/middleware"
	"new-api-demo/model"
	relaycommon "new-api-demo/relay/common"
	"new-api-demo/relay/helper"
	"new-api-demo/service"
	"new-api-demo/types"
	"strings"

	"github.com/gin-gonic/gin"
)

func relayHandler(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError {
	return relay.TextHelper(c, info)
}

func Relay(c *gin.Context, relayFormat types.RelayFormat) {
	requestId := common.GetContextKeyString(c, common.RequestIdKey)                  //获取请求的id X-Oneapi-Request-Id
	originalModel := common.GetContextKeyString(c, constant.ContextKeyOriginalModel) //原始模型

	var (
		newAPIError *types.NewAPIError
	)

	defer func() { //处理错误
		if newAPIError != nil {
			logger.LogError(c, fmt.Sprintf("relay error: %s", newAPIError.Error()))
			newAPIError.SetMessage(fmt.Sprintf("%s (request id: %s)", newAPIError.Error(), requestId))
			c.JSON(newAPIError.StatusCode, gin.H{
				"error": newAPIError.ToOpenAIError(),
			})
		}
	}()

	request, err := helper.GetAndValidateRequest(c, relayFormat) //验证请求格式
	if err != nil {
		newAPIError = types.NewError(err, types.ErrorCodeInvalidRequest)
		return
	}

	relayInfo, err := relaycommon.GenRelayInfo(c, relayFormat, request) //对请求信息做一次包装转换
	if err != nil {
		newAPIError = types.NewError(err, types.ErrorCodeGenRelayInfoFailed)
		return
	}

	for i := 0; i <= common.RetryTimes; i++ {
		channel, err := getChannel(c, originalModel, i)
		if err != nil {
			logger.LogError(c, err.Error())
			newAPIError = err
			break
		}

		addUsedChannel(c, channel.Id) //记录使用过的渠道
		requestBody, _ := common.GetRequestBody(c)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		switch relayFormat {
		default:
			newAPIError = relayHandler(c, relayInfo) //默认请求
		}

		if newAPIError == nil {
			return
		}

		logger.LogError(c, fmt.Sprintf("channel error (channel #%d, status code: %d): %s", channel.Id, newAPIError.StatusCode, newAPIError.Error()))

		if !shouldRetry(c, newAPIError, common.RetryTimes-i) {
			break
		}
	}

	useChannel := c.GetStringSlice("use_channel")
	if len(useChannel) > 1 {
		retryLogStr := fmt.Sprintf("重试：%s", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(useChannel)), "->"), "[]"))
		logger.LogInfo(c, retryLogStr)
	}
}

func addUsedChannel(c *gin.Context, channelId int) {
	useChannel := c.GetStringSlice("use_channel")
	useChannel = append(useChannel, fmt.Sprintf("%d", channelId))
	c.Set("use_channel", useChannel)
}

func getChannel(c *gin.Context, originalModel string, retryCount int) (*model.Channel, *types.NewAPIError) {
	//if retryCount == 0 {
	//	return &model.Channel{
	//		Id:   c.GetInt("channel_id"),
	//		Type: c.GetInt("channel_type"),
	//		Name: c.GetString("channel_name"),
	//	}, nil
	//}
	channel, err := service.CacheGetRandomSatisfiedChannel(c, originalModel, retryCount) //加权随机选择
	if err != nil {
		return nil, types.NewError(fmt.Errorf("获取模型 %s 的可用渠道失败（retry）: %s", originalModel, err.Error()), types.ErrorCodeGetChannelFailed, types.ErrOptionWithSkipRetry())
	}
	if channel == nil {
		return nil, types.NewError(fmt.Errorf("模型 %s 的可用渠道不存在（retry）", originalModel), types.ErrorCodeGetChannelFailed, types.ErrOptionWithSkipRetry())
	}
	newAPIError := middleware.SetupContextForSelectedChannel(c, channel, originalModel)
	if newAPIError != nil {
		return nil, newAPIError
	}
	return channel, nil
}

func shouldRetry(c *gin.Context, openaiErr *types.NewAPIError, retryTimes int) bool {
	if openaiErr == nil {
		return false
	}
	if types.IsChannelError(openaiErr) {
		return true
	}
	if types.IsSkipRetryError(openaiErr) {
		return false
	}
	if retryTimes <= 0 {
		return false
	}
	if _, ok := c.Get("specific_channel_id"); ok {
		return false
	}
	if openaiErr.StatusCode == http.StatusTooManyRequests {
		return true
	}
	if openaiErr.StatusCode == 307 {
		return true
	}
	if openaiErr.StatusCode/100 == 5 {
		// 超时不重试
		if openaiErr.StatusCode == 504 || openaiErr.StatusCode == 524 {
			return false
		}
		return true
	}
	if openaiErr.StatusCode == http.StatusBadRequest {
		return false
	}
	if openaiErr.StatusCode == 408 {
		// azure处理超时不重试
		return false
	}
	if openaiErr.StatusCode/100 == 2 {
		return false
	}
	return true
}
