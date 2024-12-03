package task

import (
	"monkey-test-api/internal/types"
	"time"
)

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

// TaskGenerator 定义任务生成器接口
type TaskGenerator interface {
	GenerateTestObjects(req *types.Request) ([]TestObject, error)
}
