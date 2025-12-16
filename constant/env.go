package constant

// 这个env 都是我们后续需要放入到配置文件中的，是我自己整理的
const (
	RetryTimes       = 3   //重试次数
	StreamingTimeout = 300 //stream响应超时时间
	RelayTimeout     = 2   //http请求下游大模型的超时时间
)
