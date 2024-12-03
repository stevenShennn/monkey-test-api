package types

import (
	"time"
)

// Request 表示一个 API 请求
type Request struct {
	RequestID  string                 `json:"request_id"`
	Method     string                 `json:"method"`
	URL        string                 `json:"url"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
	Params     map[string]interface{} `json:"params"`
	Timestamp  time.Time             `json:"timestamp"`
} 