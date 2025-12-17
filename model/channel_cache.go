package model

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	//"new-api-demo/setting/ratio_setting"
)

// var group2model2channels map[string]map[string][]int // enabled channel
var model2channels map[string][]int // enabled channel
var channelsIDM map[int]*Channel    // all channels include disabled
var channelSyncLock sync.RWMutex

func InitChannelCache() {
	//写点假数据不调用数据库
	weight1, weight2, weight3 := uint(10), uint(20), uint(30)
	priority1, priority2, priority3 := int64(1), int64(2), int64(3)

	model2channels = map[string][]int{
		"deepseek-ai/DeepSeek-V3.1": {1, 2, 3},
	}
	baseUrl := "https://www.sophnet.com/api/open-apis/v1/chat/completions"
	key := "**************************************************"
	modelMapping := `{
	"deepseek-ai/DeepSeek-V3.1": "DeepSeek-V3.1"
	}`
	channelsIDM = map[int]*Channel{
		1: {Id: 1, Name: "channel1", Type: 8, Weight: &weight1, Priority: &priority1, BaseURL: &baseUrl, Key: key, ModelMapping: &modelMapping},
		2: {Id: 2, Name: "channel2", Type: 8, Weight: &weight2, Priority: &priority2, BaseURL: &baseUrl, Key: key, ModelMapping: &modelMapping},
		3: {Id: 3, Name: "channel3", Type: 8, Weight: &weight3, Priority: &priority3, BaseURL: &baseUrl, Key: key, ModelMapping: &modelMapping},
	}
}

// GetRandomSatisfiedChannel 随机选择算法
func GetRandomSatisfiedChannel(model string, retry int) (*Channel, error) {
	channelSyncLock.RLock()
	defer channelSyncLock.RUnlock()

	//TODO：后续改成从redis中获取
	channels := model2channels[model]
	if len(channels) == 0 { // 如果没有找到指定模型的渠道，则返回 nil
		return nil, nil
	}

	if len(channels) == 1 { // 如果只有一个渠道，则直接返回
		if channel, ok := channelsIDM[channels[0]]; ok {
			return channel, nil
		}
		return nil, fmt.Errorf("数据库一致性错误，渠道# %d 不存在，请联系管理员修复", channels[0])
	}

	uniquePriorities := make(map[int]bool) // 用于存储渠道的优先级
	for _, channelId := range channels {
		if channel, ok := channelsIDM[channelId]; ok {
			uniquePriorities[int(channel.GetPriority())] = true
		} else {
			return nil, fmt.Errorf("数据库一致性错误，渠道# %d 不存在，请联系管理员修复", channelId)
		}
	}
	var sortedUniquePriorities []int // 存储渠道的优先级排序 降序0最高
	for priority := range uniquePriorities {
		sortedUniquePriorities = append(sortedUniquePriorities, priority)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedUniquePriorities))) // 降序

	if retry >= len(uniquePriorities) { //重试次数不能超过渠道的优先级数量
		retry = len(uniquePriorities) - 1
	}
	targetPriority := int64(sortedUniquePriorities[retry]) //如果只有两个优先级0和2，第0次重试是第0优先级，第一次重试是第2优先级，以此类推

	// get the priority for the given retry number // 获取指定重试次数的优先级
	var sumWeight = 0             //总权重的大小
	var targetChannels []*Channel // 目标渠道
	for _, channelId := range channels {
		if channel, ok := channelsIDM[channelId]; ok {
			if channel.GetPriority() == targetPriority {
				sumWeight += channel.GetWeight()
				targetChannels = append(targetChannels, channel)
			}
		} else {
			return nil, fmt.Errorf("数据库一致性错误，渠道# %d 不存在，请联系管理员修复", channelId)
		}
	}

	if len(targetChannels) == 0 {
		return nil, errors.New(fmt.Sprintf("no channel found, model: %s, priority: %d", model, targetPriority))
	}

	// smoothing factor and adjustment 平滑系数和调整
	smoothingFactor := 1
	smoothingAdjustment := 0

	if sumWeight == 0 { //如果权重为0
		// when all channels have weight 0, set sumWeight to the number of channels and set smoothing adjustment to 100  当所有通道的权重都为0时，将sumWeight设置为通道数，并将平滑调整设置为100
		// each channel's effective weight = 100 // 每个通道的有效权重=100
		sumWeight = len(targetChannels) * 100
		smoothingAdjustment = 100 // 平滑调整
	} else if sumWeight/len(targetChannels) < 10 { // 当平均权重小于10时，将平滑系数设置为100
		// when the average weight is less than 10, set smoothing factor to 100
		smoothingFactor = 100
	}

	// Calculate the total weight of all channels up to endIdx
	totalWeight := sumWeight * smoothingFactor // 总权重= 权重之和*平滑系数

	// Generate a random value in the range [0, totalWeight)
	randomWeight := rand.Intn(totalWeight) // 生成一个随机值，范围在[0, totalWeight)]

	// Find a channel based on its weight
	for _, channel := range targetChannels { //想象一下 假设有5个渠道，权重分别为10，20，30，40，50 比如落在55那么就是10不行，20不行，最后落在30上
		randomWeight -= channel.GetWeight()*smoothingFactor + smoothingAdjustment
		if randomWeight < 0 {
			return channel, nil
		}
	}
	// return null if no channel is not found // 如果没有找到渠道，则返回 null（理论上不走）
	return nil, errors.New("channel not found")
}
