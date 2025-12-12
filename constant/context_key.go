package constant

type ContextKey string

const (
	ContextKeyOriginalModel  ContextKey = "original_model"
	ContextKeyUserId         ContextKey = "id"
	ContextKeyChannelId      ContextKey = "channel_id"
	ContextKeyChannelName    ContextKey = "channel_name"
	ContextKeyChannelKey     ContextKey = "channel_key"
	ContextKeyChannelType    ContextKey = "channel_type"
	ContextKeyChannelBaseUrl ContextKey = "base_url"
)
