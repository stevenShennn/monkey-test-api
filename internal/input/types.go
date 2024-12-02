package input

// Request 表示解析后的请求数据
type Request struct {
	RequestID string            `json:"request_id"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      interface{}       `json:"body"`
	Params    map[string]string `json:"params"`
	Proxy     *ProxyConfig     `json:"proxy,omitempty"`
}

// ProxyConfig 代理服务器配置
type ProxyConfig struct {
	Enabled   bool   `json:"enabled"`
	ProxyURL  string `json:"proxy_url"`
}

// Parser 定义请求解析器接口
type Parser interface {
	Parse(input string) ([]Request, error)
} 