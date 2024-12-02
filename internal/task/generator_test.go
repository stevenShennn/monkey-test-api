package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_GenerateTestObjects(t *testing.T) {
	// 创建测试配置
	config := &Config{
		String: ParamTestConfig{
			Values: []ParamTestValue{
				{Value: "test", Reason: "测试参数%s普通字符串"},
				{Value: "' OR '1'='1", Reason: "测试参数%sSQL注入"},
			},
		},
		Number: ParamTestConfig{
			Values: []ParamTestValue{
				{Value: 1, Reason: "测试参数%s最小值"},
				{Value: 9999, Reason: "测试参数%s最大值"},
			},
		},
	}

	generator := NewGenerator(config)

	// 创建测试请求
	parentRequest := &Request{
		RequestID: "test-123",
		Method:    "POST",
		URL:       "http://example.com/api",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Params: map[string]interface{}{
			"name": "john",
			"age":  25,
		},
		Timestamp: time.Now(),
	}

	// 生成测试对象
	testObjects, err := generator.GenerateTestObjects(parentRequest)
	assert.NoError(t, err)

	// 验证生成的测试对象
	assert.Len(t, testObjects, 4) // 2个参数 * 2个测试值 = 4个测试对象

	// 验证每个测试对象
	for _, obj := range testObjects {
		// 验证基本字段
		assert.NotEmpty(t, obj.TestID)
		assert.Equal(t, parentRequest.RequestID, obj.ParentRequestID)
		assert.Equal(t, parentRequest.Method, obj.Method)
		assert.Equal(t, parentRequest.URL, obj.URL)
		assert.Equal(t, parentRequest.Headers, obj.Headers)
		assert.Equal(t, "待处理", obj.Status)

		// 验证参数
		assert.Len(t, obj.Params, 2)
		assert.Contains(t, obj.Params, "name")
		assert.Contains(t, obj.Params, "age")

		// 验证原因字段
		assert.NotEmpty(t, obj.Reason)
	}
} 