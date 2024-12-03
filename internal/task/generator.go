package task

import (
	"fmt"
	"monkey-test-api/internal/types"
	"time"
)

// Generator 测试任务生成器
type Generator struct {
	paramConfig *types.ParamConfig
}

// NewGenerator 创建新的任务生成器
func NewGenerator(config *types.ParamConfig) *Generator {
	return &Generator{
		paramConfig: config,
	}
}

// GenerateTestObjects 生成测试对象
func (g *Generator) GenerateTestObjects(req *types.Request) ([]TestObject, error) {
	var testObjects []TestObject

	// 处理字符串参数
	for key, value := range req.Params {
		if _, ok := value.(string); ok {
			for _, test := range g.paramConfig.String.Tests {
				testObj := TestObject{
					TestID:          fmt.Sprintf("%s-%s", req.RequestID, generateID()),
					ParentRequestID: req.RequestID,
					Method:         req.Method,
					URL:            req.URL,
					Headers:        req.Headers,
					Body:           req.Body,
					Params:         copyParams(req.Params),
					Reason:         fmt.Sprintf("%s: %s", test.Description, test.Reason),
					Status:         "待处理",
					Risk:           test.Risk,
					Timestamp:      time.Now(),
				}
				testObj.Params[key] = test.Value
				testObjects = append(testObjects, testObj)
			}
		}
	}

	// 处理数字参数
	for key, value := range req.Params {
		if _, ok := value.(float64); ok {
			for _, test := range g.paramConfig.Number.Tests {
				testObj := TestObject{
					TestID:          fmt.Sprintf("%s-%s", req.RequestID, generateID()),
					ParentRequestID: req.RequestID,
					Method:         req.Method,
					URL:            req.URL,
					Headers:        req.Headers,
					Body:           req.Body,
					Params:         copyParams(req.Params),
					Reason:         fmt.Sprintf("%s: %s", test.Description, test.Reason),
					Status:         "待处理",
					Risk:           test.Risk,
					Timestamp:      time.Now(),
				}
				testObj.Params[key] = test.Value
				testObjects = append(testObjects, testObj)
			}
		}
	}

	// 处理布尔参数
	for key, value := range req.Params {
		if _, ok := value.(bool); ok {
			for _, test := range g.paramConfig.Bool.Tests {
				testObj := TestObject{
					TestID:          fmt.Sprintf("%s-%s", req.RequestID, generateID()),
					ParentRequestID: req.RequestID,
					Method:         req.Method,
					URL:            req.URL,
					Headers:        req.Headers,
					Body:           req.Body,
					Params:         copyParams(req.Params),
					Reason:         fmt.Sprintf("%s: %s", test.Description, test.Reason),
					Status:         "待处理",
					Risk:           test.Risk,
					Timestamp:      time.Now(),
				}
				testObj.Params[key] = test.Value
				testObjects = append(testObjects, testObj)
			}
		}
	}

	return testObjects, nil
}

// 辅助函数：生成唯一ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 辅助函数：复制参数映射
func copyParams(params map[string]interface{}) map[string]interface{} {
	newParams := make(map[string]interface{})
	for k, v := range params {
		newParams[k] = v
	}
	return newParams
} 