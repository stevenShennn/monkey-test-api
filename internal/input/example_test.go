package input_test

import (
	"encoding/json"
	"fmt"
	"monkey-test-api/internal/input"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputExample(t *testing.T) {
	// 1. 准备测试数据：一个包含多个 cURL 命令的字符串
	curlCommands := `
curl 'https://api.example.com/users' \
  -H 'accept: application/json' \
  -H 'authorization: Bearer token123' \
  -H 'content-type: application/json' \
  -X POST \
  -d '{"name":"test","age":25}'

curl 'https://api.example.com/products?category=electronics' \
  -H 'accept: application/json' \
  -H 'authorization: Bearer token123'
`

	// 2. 获取 cURL 解析器
	parser, err := input.GetParser(input.ParserTypeCurl)
	assert.NoError(t, err, "获取解析器应该成功")

	// 3. 解析 cURL 命令
	requests, err := parser.Parse(curlCommands)
	assert.NoError(t, err, "解析 cURL 命令应该成功")
	assert.Len(t, requests, 2, "应该解析出两个请求")

	// 4. 验证解析结果并展示如何使用
	for i, req := range requests {
		// 将请求转换为 JSON 以便查看
		jsonData, err := json.MarshalIndent(req, "", "  ")
		assert.NoError(t, err, "转换 JSON 应该成功")

		fmt.Printf("Request %d:\n%s\n\n", i+1, string(jsonData))

		// 验证具体字段
		switch i {
		case 0: // 第一个请求（POST /users）
			assert.Equal(t, "POST", req.Method)
			assert.Equal(t, "https://api.example.com/users", req.URL)
			assert.Equal(t, "application/json", req.Headers["accept"])
			assert.Equal(t, "Bearer token123", req.Headers["authorization"])
			assert.Equal(t, `{"name":"test","age":25}`, req.Body)

		case 1: // 第二个请求（GET /products）
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, "https://api.example.com/products?category=electronics", req.URL)
			assert.Equal(t, "application/json", req.Headers["accept"])
			assert.Equal(t, "Bearer token123", req.Headers["authorization"])
		}

		// 示例：如何使用解析后的请求
		ExampleRequestHandler(t, req)
	}
}

// ExampleRequestHandler 展示如何处理解析后的请求
func ExampleRequestHandler(t *testing.T, req input.Request) {
	// 1. 检查认证信息
	authHeader, exists := req.Headers["authorization"]
	assert.True(t, exists, "请求应该包含认证信息")
	assert.Contains(t, authHeader, "Bearer ", "应该是 Bearer 认证")

	// 2. 根据请求方法进行不同处理
	switch req.Method {
	case "GET":
		// 处理 GET 请求
		fmt.Printf("处理 GET 请求: %s\n", req.URL)
		// 这里可以添加更多处理逻辑...

	case "POST":
		// 处理 POST 请求
		fmt.Printf("处理 POST 请求: %s\n", req.URL)
		// 验证请求体
		assert.NotNil(t, req.Body, "POST 请求应该包含请求体")
		// 这里可以添加更多处理逻辑...

	default:
		t.Errorf("未支持的请求方法: %s", req.Method)
	}

	// 3. 处理代理设置（如果有）
	if req.Proxy != nil && req.Proxy.Enabled {
		fmt.Printf("使用代理: %s\n", req.Proxy.ProxyURL)
	}
} 