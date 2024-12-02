package task

import (
	"fmt"
	"github.com/google/uuid"
	"monkey-test-api/internal/logger"
	"reflect"
	"time"
)

type Generator struct {
	config *Config
}

// NewGenerator 创建新的任务生成器
func NewGenerator(config *Config) TaskGenerator {
	return &Generator{
		config: config,
	}
}

// GenerateTestObjects 生成测试对象
func (g *Generator) GenerateTestObjects(parentRequest *Request) ([]TestObject, error) {
	logger.Infof("开始为请求 %s 生成测试对象", parentRequest.RequestID)
	
	var testObjects []TestObject

	// 遍历请求参数
	for paramName, paramValue := range parentRequest.Params {
		paramType := g.detectParamType(paramValue)
		testValues := g.getTestValues(paramType)

		// 为每个测试值生成一个测试对象
		for _, testValue := range testValues {
			testObj := g.createTestObject(parentRequest, paramName, testValue)
			testObjects = append(testObjects, testObj)
		}
	}

	logger.Infof("为请求 %s 生成了 %d 个测试对象", parentRequest.RequestID, len(testObjects))
	return testObjects, nil
}

// detectParamType 检测参数类型
func (g *Generator) detectParamType(value interface{}) ParamType {
	if value == nil {
		return ParamTypeString
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return ParamTypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return ParamTypeNumber
	case reflect.Bool:
		return ParamTypeBoolean
	case reflect.Map:
		return ParamTypeObject
	case reflect.Slice, reflect.Array:
		return ParamTypeArray
	default:
		return ParamTypeString
	}
}

// getTestValues 获取测试值
func (g *Generator) getTestValues(paramType ParamType) []ParamTestValue {
	switch paramType {
	case ParamTypeString:
		return g.config.String.Values
	case ParamTypeNumber:
		return g.config.Number.Values
	case ParamTypeBoolean:
		return g.config.Boolean.Values
	case ParamTypeObject:
		return g.config.Object.Values
	case ParamTypeArray:
		return g.config.Array.Values
	case ParamTypeDate:
		return g.config.Date.Values
	case ParamTypeTime:
		return g.config.Time.Values
	default:
		return g.config.String.Values
	}
}

// createTestObject 创建测试对象
func (g *Generator) createTestObject(parent *Request, paramName string, testValue ParamTestValue) TestObject {
	// 复制父请求的参数
	params := make(map[string]interface{})
	for k, v := range parent.Params {
		params[k] = v
	}

	// 替换测试参数
	params[paramName] = testValue.Value

	return TestObject{
		TestID:         uuid.New().String(),
		ParentRequestID: parent.RequestID,
		Method:         parent.Method,
		URL:            parent.URL,
		Headers:        parent.Headers,
		Body:           parent.Body,
		Params:         params,
		Reason:         fmt.Sprintf(testValue.Reason, paramName),
		Status:         "待处理",
		Timestamp:      time.Now(),
	}
} 