package types

type ChannelError struct {
	ChannelId   int    `json:"channel_id"`
	ChannelType int    `json:"channel_type"`
	ChannelName string `json:"channel_name"`
	UsingKey    string `json:"using_key"`
}

func NewChannelError(channelId int, channelType int, channelName string, usingKey string) *ChannelError {
	return &ChannelError{
		ChannelId:   channelId,
		ChannelType: channelType,
		ChannelName: channelName,
		UsingKey:    usingKey,
	}
}
