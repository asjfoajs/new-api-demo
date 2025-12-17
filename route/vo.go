// Package route 无实际意义 swagger生成使用
package route

// ChatCompletionRequest 聊天接口请求参数
type ChatCompletionRequest struct {
	// 聊天上下文信息
	// 纯文本示例: [{"role": "system", "content": "You are a helpful assistant."}]
	Messages []ChatMessage `json:"messages" swagg:"required" description:"聊天上下文信息"`

	// 模型名称
	Model string `json:"model" swagg:"required" example:"DeepSeek-V3.1" description:"模型"`

	// 工具列表，只支持 function
	Tools []Tool `json:"tools,omitempty" description:"工具列表"`

	// 是否以流式接口的形式返回数据，默认 false
	Stream bool `json:"stream,omitempty" default:"false" description:"是否以流式模式返回"`

	// 模型回复最大长度（单位 token）
	MaxTokens int `json:"max_tokens,omitempty" description:"回复最大长度"`

	// 采样温度，0.2-2.0 之间。默认 1.0
	Temperature float64 `json:"temperature,omitempty" default:"1.0" example:"0.2" description:"采样温度"`

	// 影响输出文本的多样性，默认 1.0
	TopP float64 `json:"top_p,omitempty" default:"1.0" description:"核采样"`

	// 停止生成更多 Token 的字符串
	Stop []string `json:"stop,omitempty" description:"停止位"`

	// 通过对已生成的 token 增加惩罚减少重复，范围 [-2.0, 2.0]
	PresencePenalty float64 `json:"presence_penalty,omitempty" default:"0" description:"存在惩罚"`

	// 降低模型逐字重复同一行的可能性，范围 [-2.0, 2.0]
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty" default:"0" description:"频率惩罚"`

	// 是否返回输出 tokens 的对数概率，默认 false
	Logprobs bool `json:"logprobs,omitempty" default:"false" description:"返回对数概率"`

	// 指定每个输出位置返回的 token 数量，范围 [0, 20]
	TopLogprobs int `json:"top_logprobs,omitempty" description:"返回 top tokens 数量"`

	// 指定模型输出的格式对象，如设置 {"type": "json_object"} 开启 JSON 模式
	ResponseFormat *ResponseFormat `json:"response_format,omitempty" description:"输出格式控制"`

	// 是否开启思考模式
	EnableThinking bool `json:"enable_thinking,omitempty" description:"开启思考模式"`
}

// ChatMessage 消息体定义
type ChatMessage struct {
	Role    string      `json:"role" swagg:"required" example:"user" description:"角色: system, user, assistant"`
	Content interface{} `json:"content" swagg:"required" description:"消息内容，支持字符串或多模态数组"`
}

// Tool 工具定义
type Tool struct {
	Type     string       `json:"type" example:"function" description:"工具类型"`
	Function FunctionDesc `json:"function" description:"函数描述"`
}

// FunctionDesc 函数描述详情
type FunctionDesc struct {
	Name        string      `json:"name" description:"函数名称"`
	Description string      `json:"description" description:"函数描述"`
	Parameters  interface{} `json:"parameters" description:"函数参数，符合 JSON Schema"`
}

// ResponseFormat 格式控制
type ResponseFormat struct {
	// 指定格式类型，如 "text" 或 "json_object"
	Type string `json:"type" example:"json_object" description:"格式类型"`
}

// ChatCompletionResponse 非流式响应体
type ChatCompletionResponse struct {
	Object  string `json:"object" example:"chat.completion" description:"回包类型"`
	Created int64  `json:"created" example:"1677652288" description:"时间戳"`
	Model   string `json:"model" example:"Qwen/Qwen2.5-72B-Instruct" description:"模型名称"`
	Choices []struct {
		Index        int    `json:"index" description:"索引"`
		FinishReason string `json:"finish_reason" example:"stop" description:"结束原因: stop(正常结束), length(token超长)"`
		Message      struct {
			Role      string     `json:"role" example:"assistant" description:"角色"`
			Content   string     `json:"content" description:"模型回答内容"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty" description:"工具调用列表"`
		} `json:"message" description:"模型回答"`
		// 引用列表：调用自定义模型且模型输出包含文档引用时存在
		Refs []Reference `json:"refs,omitempty" description:"引用列表，包含此次回答的所有引用来源信息"`
	} `json:"choices" description:"结果列表"`
}

// ToolCall 工具调用详情
type ToolCall struct {
	Function struct {
		Name      string `json:"name" description:"函数名称"`
		Arguments string `json:"arguments" description:"函数参数(JSON字符串)"`
	} `json:"function" description:"函数调用信息"`
}

// ChatCompletionStreamResponse 流式响应分片 (Chunk)
type ChatCompletionStreamResponse struct {
	Object  string `json:"object" example:"chat.completion.chunk" description:"回包类型"`
	Created int64  `json:"created" description:"时间戳"`
	Model   string `json:"model" description:"模型名称"`
	Choices []struct {
		Index        int    `json:"index" description:"索引"`
		FinishReason string `json:"finish_reason" description:"结束原因"`
		Delta        struct {
			Role    string `json:"role,omitempty" description:"角色"`
			Content string `json:"content,omitempty" description:"分片内容"`
		} `json:"delta" description:"模型回答增量"`
		// 在流式响应过程中，会实时或在结束前输出引用来源信息
		Refs []Reference `json:"refs,omitempty" description:"引用列表"`
	} `json:"choices"`
}

// Reference 引用来源信息
type Reference struct {
	Index   int    `json:"index" description:"引用来源出现顺序"`
	Title   string `json:"title" description:"引用数据标题"`
	Content string `json:"content" description:"引用数据内容"`
	// 引用数据类型: file(文件知识), qa(Q&A Table), web(web搜索)
	Type string `json:"type" swagg:"enum:file,qa,web" description:"引用数据类型"`
	// 引用数据URL。注意：file和qa需配置apikey访问，web无需配置
	URL string `json:"url" description:"引用数据url"`
}
