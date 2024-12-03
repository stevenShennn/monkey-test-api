package store

import (
	"time"
	"monkey-test-api/internal/types"
	"monkey-test-api/internal/task"
)

// Request 类型别名
type Request = types.Request

// TestObject 表示一个测试对象
type TestObject struct {
	TestID          string                 `json:"test_id"`
	ParentRequestID string                 `json:"parent_request_id"`
	Method          string                 `json:"method"`
	URL             string                 `json:"url"`
	Headers         map[string]string      `json:"headers"`
	Body            map[string]interface{} `json:"body"`
	Params          map[string]interface{} `json:"params"`
	Description     string                 `json:"description"`
	Reason          string                 `json:"reason"`
	Status          string                 `json:"status"`
	Risk            string                 `json:"risk"`
	Timestamp       time.Time             `json:"timestamp"`
}

// Response 响应数据模型
type Response struct {
	Status    int               `bson:"status" json:"status"`
	Body      interface{}       `bson:"body" json:"body"`
	Headers   map[string]string `bson:"headers" json:"headers"`
	Timestamp time.Time         `bson:"timestamp" json:"timestamp"`
}

// ConvertTaskTestObjects 将 task.TestObject 转换为 store.TestObject
func ConvertTaskTestObjects(taskObjs []task.TestObject) []TestObject {
	storeObjs := make([]TestObject, len(taskObjs))
	for i, obj := range taskObjs {
		storeObjs[i] = TestObject{
			TestID:          obj.TestID,
			ParentRequestID: obj.ParentRequestID,
			Method:          obj.Method,
			URL:             obj.URL,
			Headers:         obj.Headers,
			Body:            obj.Body,
			Params:          obj.Params,
			Description:     obj.Description,
			Reason:          obj.Reason,
			Status:          obj.Status,
			Risk:            obj.Risk,
			Timestamp:       obj.Timestamp,
		}
	}
	return storeObjs
} 