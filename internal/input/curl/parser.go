package curl

import (
	"errors"
	"github.com/google/uuid"
	"monkey-test-api/internal/input"
	"regexp"
	"strings"
)

type Parser struct{}

// NewParser 创建一个新的 cURL 解析器
func NewParser() *Parser {
	return &Parser{}
}

// Parse 实现 input.Parser 接口，解析 cURL 命令
func (p *Parser) Parse(curlCmd string) ([]input.Request, error) {
	// 分割多个 cURL 命令
	commands := splitCurlCommands(curlCmd)
	requests := make([]input.Request, 0, len(commands))

	for _, cmd := range commands {
		req, err := p.parseSingleCommand(cmd)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// parseSingleCommand 解析单个 cURL 命令
func (p *Parser) parseSingleCommand(cmd string) (input.Request, error) {
	req := input.Request{
		RequestID: uuid.New().String(),
		Headers:   make(map[string]string),
		Params:    make(map[string]string),
	}

	// 提取 URL
	urlRegex := regexp.MustCompile(`curl ['"]([^'"]+)['"]`)
	if matches := urlRegex.FindStringSubmatch(cmd); len(matches) > 1 {
		req.URL = matches[1]
		req.Method = "GET" // 默认方法为 GET
	} else {
		return req, errors.New("无法解析 URL")
	}

	// 提取请求头
	headerRegex := regexp.MustCompile(`-H ['"]([^:]+):([^'"]+)['"]`)
	headerMatches := headerRegex.FindAllStringSubmatch(cmd, -1)
	for _, match := range headerMatches {
		if len(match) == 3 {
			req.Headers[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
		}
	}

	// 提取请求方法
	methodRegex := regexp.MustCompile(`-X ([A-Z]+)`)
	if matches := methodRegex.FindStringSubmatch(cmd); len(matches) > 1 {
		req.Method = matches[1]
	}

	// 提取请求体
	bodyRegex := regexp.MustCompile(`-d ['"](.+)['"]`)
	if matches := bodyRegex.FindStringSubmatch(cmd); len(matches) > 1 {
		req.Body = matches[1]
	}

	return req, nil
}

// splitCurlCommands 分割多个 cURL 命令
func splitCurlCommands(input string) []string {
	// 按照 "curl" 关键字分割，确保每个命令都是完整的
	commands := regexp.MustCompile(`(?m)^curl`).Split(input, -1)
	result := make([]string, 0)
	
	for _, cmd := range commands {
		if cmd = strings.TrimSpace(cmd); cmd != "" {
			result = append(result, "curl"+cmd)
		}
	}
	
	return result
} 