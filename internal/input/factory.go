package input

import (
	"errors"
	"monkey-test-api/internal/input/curl"
)

// ParserType 表示解析器类型
type ParserType string

const (
	ParserTypeCurl ParserType = "curl"
)

// GetParser 根据类型返回对应的解析器
func GetParser(parserType ParserType) (Parser, error) {
	switch parserType {
	case ParserTypeCurl:
		return curl.NewParser(), nil
	default:
		return nil, errors.New("不支持的解析器类型")
	}
} 