# New API Demo

这是一个简化的 API demo，只包含两个接口用于学习：

- `POST /v1/completions` - 文本补全接口
- `POST /v1/chat/completions` - 聊天补全接口

## 运行方式

```bash
# 初始化依赖
go mod tidy

# 运行服务
go run main.go
```

服务将在 `http://localhost:8080` 启动

## 测试示例

### 测试 /v1/completions

```bash
curl -X POST http://localhost:8080/v1/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "prompt": "Hello, world"
  }'
```

### 测试 /v1/chat/completions

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "user",
        "content": "Hello!"
      }
    ]
  }'
```

## 说明

这是一个极简版本，只保留了基本的请求解析和响应返回逻辑，用于学习这两个接口的基本结构。
