# New API Demo - 精简版大模型中转分发系统

本项目是基于 [New API](https://github.com/QuantumNous/new-api) 的精简版本。由于原项目功能极其丰富，代码逻辑相对复杂，本项目通过“大刀阔斧”的删减，仅保留最核心的 OpenAI 格式转发逻辑，旨在为开发者提供一个**极简的、易于学习和二开的 LLM 转发网关原型**。

## 🌟 项目特点

为了方便理解大模型中转的核心流程，本项目删减了以下功能：

* **单一协议**：仅保留 Chat OpenAI 格式转发。
* **无数据库依赖**：去掉了 SQL 和 Redis。渠道与模型信息在初始化时直接配置，无需部署复杂的数据库环境。
* **轻量化**：去掉了多 Key 模式、计费系统、限流逻辑以及令牌（Token）校验逻辑。
* **极简配置**：删除了所有繁杂的设置项，仅保留必要的模型参数。

---

## 🚀 快速开始

### 1. 修改配置（假数据初始化）

由于本项目去除了数据库，所有的渠道（Channel）信息均通过代码硬编码实现。

请在 `main.go` 中定位到以下代码行并进入：
`model.InitChannelCache()`

在 `func InitChannelCache()` 中，你可以手动配置你的模型映射、API Key 和 Base URL：

```go
func InitChannelCache() {
	// 配置权重与优先级
	weight1, weight2 := uint(10), uint(20)
	priority1, priority2 := int64(1), int64(2)

	// 模型与渠道 ID 的映射
	model2channels = map[string][]int{
		"deepseek-ai/DeepSeek-V3.1": {1, 2},
	}

	baseUrl := "https://your-api-proxy.com/v1"
	key := "sk-xxxxxxxxxxxxxxxxxxxxxxxx"
    
	modelMapping := `{"deepseek-ai/DeepSeek-V3.1": "DeepSeek-V3.1"}`

	// 模拟数据库中的渠道数据
	channelsIDM = map[int]*Channel{
		1: {Id: 1, Name: "渠道A", Type: 8, Weight: &weight1, Priority: &priority1, BaseURL: &baseUrl, Key: key, ModelMapping: &modelMapping},
		2: {Id: 2, Name: "渠道B", Type: 8, Weight: &weight2, Priority: &priority2, BaseURL: &baseUrl, Key: key, ModelMapping: &modelMapping},
	}
}

```

### 2. 编译并运行

```bash
go mod tidy
go run main.go

```

---

## 🛠 接口示例

你可以生成 Swagger 文档并导入到 **Apifox** 或 **Postman** 中进行调试。

### A. 普通对话（非流式）

**请求参数：**

```json
{
  "model": "deepseek-ai/DeepSeek-V3.1",
  "messages": [{"role": "user", "content": "请解释什么是递归。"}],
  "stream": false
}

```

### B. 思考模型/工具调用（流式）

本项目支持转发 DeepSeek-V3.1 等模型的 `reasoning_content`（思考过程）和 `tools`（工具调用）参数。

**请求示例（含 Tool Call）：**

```json
{
  "model": "deepseek-ai/DeepSeek-V3.1",
  "messages": [
    {
      "role": "system",
      "content": "你是一个编程专家。如果用户询问性能，请使用 get_code_runtime 工具。"
    },
    {
      "role": "user",
      "content": "对比下 Python 中递归和循环的运行耗时？"
    }
  ],
  "tools": [{
    "type": "function",
    "function": {
      "name": "get_code_runtime",
      "description": "获取特定编程语言的基础逻辑执行耗时",
      "parameters": {
        "type": "object",
        "properties": {
          "language": {"type": "string", "enum": ["python", "java", "cpp"]},
          "logic_type": {"type": "string"}
        },
        "required": ["language", "logic_type"]
      }
    }
  }],
  "stream": true,
  "enable_thinking": true
}

```

---

## 📈 进阶说明

* **负载均衡**：通过 `Weight`（权重）参数，系统可以自动在多个同模型渠道间分配流量。
* **异常重试**：当高优先级（Priority）渠道不可用时，系统会自动尝试后续渠道。

## 🔗 相关项目

* 原版项目：[New API](https://github.com/QuantumNous/new-api)
* 演示项目：[New API Demo](https://github.com/asjfoajs/new-api-demo)
