package relay

import (
	"new-api-demo/constant"
	"new-api-demo/relay/channel"
	"new-api-demo/relay/channel/ali"
	"new-api-demo/relay/channel/aws"
	"new-api-demo/relay/channel/baidu"
	"new-api-demo/relay/channel/baidu_v2"
	"new-api-demo/relay/channel/claude"
	"new-api-demo/relay/channel/cloudflare"
	"new-api-demo/relay/channel/cohere"
	"new-api-demo/relay/channel/coze"
	"new-api-demo/relay/channel/deepseek"
	"new-api-demo/relay/channel/dify"
	"new-api-demo/relay/channel/gemini"
	"new-api-demo/relay/channel/jimeng"
	"new-api-demo/relay/channel/jina"
	"new-api-demo/relay/channel/minimax"
	"new-api-demo/relay/channel/mistral"
	"new-api-demo/relay/channel/mokaai"
	"new-api-demo/relay/channel/moonshot"
	"new-api-demo/relay/channel/ollama"
	"new-api-demo/relay/channel/openai"
	"new-api-demo/relay/channel/palm"
	"new-api-demo/relay/channel/perplexity"
	"new-api-demo/relay/channel/replicate"
	"new-api-demo/relay/channel/siliconflow"
	"new-api-demo/relay/channel/submodel"
	"new-api-demo/relay/channel/tencent"
	"new-api-demo/relay/channel/vertex"
	"new-api-demo/relay/channel/volcengine"
	"new-api-demo/relay/channel/xai"
	"new-api-demo/relay/channel/xunfei"
	"new-api-demo/relay/channel/zhipu"
	"new-api-demo/relay/channel/zhipu_4v"
)

func GetAdaptor(apiType int) channel.Adaptor {
	switch apiType {
	case constant.APITypeAli:
		return &ali.Adaptor{}
	case constant.APITypeAnthropic:
		return &claude.Adaptor{}
	case constant.APITypeBaidu:
		return &baidu.Adaptor{}
	case constant.APITypeGemini:
		return &gemini.Adaptor{}
	case constant.APITypeOpenAI:
		return &openai.Adaptor{}
	case constant.APITypePaLM:
		return &palm.Adaptor{}
	case constant.APITypeTencent:
		return &tencent.Adaptor{}
	case constant.APITypeXunfei:
		return &xunfei.Adaptor{}
	case constant.APITypeZhipu:
		return &zhipu.Adaptor{}
	case constant.APITypeZhipuV4:
		return &zhipu_4v.Adaptor{}
	case constant.APITypeOllama:
		return &ollama.Adaptor{}
	case constant.APITypePerplexity:
		return &perplexity.Adaptor{}
	case constant.APITypeAws:
		return &aws.Adaptor{}
	case constant.APITypeCohere:
		return &cohere.Adaptor{}
	case constant.APITypeDify:
		return &dify.Adaptor{}
	case constant.APITypeJina:
		return &jina.Adaptor{}
	case constant.APITypeCloudflare:
		return &cloudflare.Adaptor{}
	case constant.APITypeSiliconFlow:
		return &siliconflow.Adaptor{}
	case constant.APITypeVertexAi:
		return &vertex.Adaptor{}
	case constant.APITypeMistral:
		return &mistral.Adaptor{}
	case constant.APITypeDeepSeek:
		return &deepseek.Adaptor{}
	case constant.APITypeMokaAI:
		return &mokaai.Adaptor{}
	case constant.APITypeVolcEngine:
		return &volcengine.Adaptor{}
	case constant.APITypeBaiduV2:
		return &baidu_v2.Adaptor{}
	case constant.APITypeOpenRouter:
		return &openai.Adaptor{}
	case constant.APITypeXinference:
		return &openai.Adaptor{}
	case constant.APITypeXai:
		return &xai.Adaptor{}
	case constant.APITypeCoze:
		return &coze.Adaptor{}
	case constant.APITypeJimeng:
		return &jimeng.Adaptor{}
	case constant.APITypeMoonshot:
		return &moonshot.Adaptor{} // Moonshot uses Claude API
	case constant.APITypeSubmodel:
		return &submodel.Adaptor{}
	case constant.APITypeMiniMax:
		return &minimax.Adaptor{}
	case constant.APITypeReplicate:
		return &replicate.Adaptor{}
	}
	return nil
}
