package task

import (
	"encoding/json"
	"time"
)

// ParamType 参数类型
type ParamType string

const (
	ParamTypeString  ParamType = "string"
	ParamTypeNumber  ParamType = "number"
	ParamTypeBoolean ParamType = "boolean"
	ParamTypeObject  ParamType = "object"
	ParamTypeArray   ParamType = "array"
	ParamTypeDate    ParamType = "date"
	ParamTypeTime    ParamType = "time"
)

// ParamTestValue 参数测试值配置
type ParamTestValue struct {
	Value  interface{} `toml:"value"`  // 测试值
	Reason string      `toml:"reason"` // 测试原因
}

// ParamTestConfig 参数测试配置
type ParamTestConfig struct {
	Values []ParamTestValue `toml:"values"`
}

// Config 测试配置
type Config struct {
	String  ParamTestConfig `toml:"string"`
	Number  ParamTestConfig `toml:"number"`
	Boolean ParamTestConfig `toml:"boolean"`
	Object  ParamTestConfig `toml:"object"`
	Array   ParamTestConfig `toml:"array"`
	Date    ParamTestConfig `toml:"date"`
	Time    ParamTestConfig `toml:"time"`
}

// TaskGenerator 任务生成器接口
type TaskGenerator interface {
	GenerateTestObjects(parentRequest *Request) ([]TestObject, error)
}

// Request 父请求
type Request struct {
	RequestID  string                 `json:"requestID"`
	Method     string                 `json:"method"`
	URL        string                 `json:"url"`
	Headers    map[string]string      `json:"headers"`
	Body       interface{}            `json:"body"`
	Params     map[string]interface{} `json:"params"`
	Timestamp  time.Time              `json:"timestamp"`
}

// TestObject 测试对象（子请求）
type TestObject struct {
	TestID         string                 `json:"testID"`
	ParentRequestID string                `json:"parentRequestID"`
	Method         string                 `json:"method"`
	URL            string                 `json:"url"`
	Headers        map[string]string      `json:"headers"`
	Body           interface{}            `json:"body"`
	Params         map[string]interface{} `json:"params"`
	Reason         string                 `json:"reason"`
	Status         string                 `json:"status"`
	Timestamp      time.Time              `json:"timestamp"`
} 