package input

import (
	"monkey-test-api/internal/types"
)

// Parser 定义解析器接口
type Parser interface {
	Parse(input string) ([]types.Request, error)
} 