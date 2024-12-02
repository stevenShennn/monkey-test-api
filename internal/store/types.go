package store

import (
	"time"
)

// Request 父请求数据模型
type Request struct {
	RequestID  string            `bson:"request_id" json:"requestID"`
	Method     string            `bson:"method" json:"method"`
	URL        string            `bson:"url" json:"url"`
	Headers    map[string]string `bson:"headers" json:"headers"`
	Body       interface{}       `bson:"body" json:"body"`
	Params     map[string]string `bson:"params" json:"params"`
	Timestamp  time.Time         `bson:"timestamp" json:"timestamp"`
}

// TestObject 子请求数据模型
type TestObject struct {
	TestID         string            `bson:"test_id" json:"testID"`
	ParentRequestID string           `bson:"parent_request_id" json:"parentRequestID"`
	Method         string            `bson:"method" json:"method"`
	URL            string            `bson:"url" json:"url"`
	Headers        map[string]string `bson:"headers" json:"headers"`
	Body           interface{}       `bson:"body" json:"body"`
	Params         map[string]string `bson:"params" json:"params"`
	Reason         string            `bson:"reason" json:"reason"`
	Response       *Response         `bson:"response,omitempty" json:"response,omitempty"`
	Timestamp      time.Time         `bson:"timestamp" json:"timestamp"`
}

// Response 响应数据模型
type Response struct {
	Status    int               `bson:"status" json:"status"`
	Body      interface{}       `bson:"body" json:"body"`
	Headers   map[string]string `bson:"headers" json:"headers"`
	Timestamp time.Time         `bson:"timestamp" json:"timestamp"`
} 