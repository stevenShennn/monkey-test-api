package curl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"monkey-test-api/internal/types"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) ([]types.Request, error) {
	// 1. 分割多个 cURL 命令
	commands := splitCurlCommands(input)
	if len(commands) == 0 {
		return nil, fmt.Errorf("未找到有效的 cURL 命令")
	}

	// 2. 解析每个命令
	var requests []types.Request
	for _, cmd := range commands {
		req, err := p.parseSingleCommand(cmd)
		if err != nil {
			return nil, fmt.Errorf("解析命令失败: %v", err)
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// splitCurlCommands 分割多个 cURL 命令
func splitCurlCommands(input string) []string {
	// 按换行符分割
	lines := strings.Split(input, "\n")
	var commands []string
	var currentCommand strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 如果是新的 cURL 命令
		if strings.HasPrefix(strings.ToLower(line), "curl ") {
			// 保存之前的命令
			if currentCommand.Len() > 0 {
				commands = append(commands, currentCommand.String())
				currentCommand.Reset()
			}
			currentCommand.WriteString(line)
		} else if currentCommand.Len() > 0 {
			// 继续当前命令
			if strings.HasSuffix(line, "\\") {
				// 处理续行符
				currentCommand.WriteString(" " + strings.TrimSuffix(line, "\\"))
			} else {
				currentCommand.WriteString(" " + line)
			}
		}
	}

	// 添加最后一个命令
	if currentCommand.Len() > 0 {
		commands = append(commands, currentCommand.String())
	}

	return commands
}

// parseSingleCommand 解析单个 cURL 命令
func (p *Parser) parseSingleCommand(input string) (types.Request, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(strings.ToLower(input), "curl ") {
		return types.Request{}, fmt.Errorf("无效的 cURL 命令：必须以 'curl' 开头")
	}

	// 解析 URL
	urlPattern := regexp.MustCompile(`'([^']*)'|"([^"]*)"`)
	urlMatches := urlPattern.FindStringSubmatch(input)
	if len(urlMatches) < 3 {
		return types.Request{}, fmt.Errorf("无法解析 URL")
	}
	url := urlMatches[1]
	if url == "" {
		url = urlMatches[2]
	}

	// 解析请求方法
	method := "GET"
	if strings.Contains(input, "-X ") || strings.Contains(input, "--request ") {
		methodPattern := regexp.MustCompile(`-X\s+(\w+)|--request\s+(\w+)`)
		methodMatches := methodPattern.FindStringSubmatch(input)
		if len(methodMatches) > 1 {
			if methodMatches[1] != "" {
				method = methodMatches[1]
			} else {
				method = methodMatches[2]
			}
		}
	}

	// 解析请求头
	headers := make(map[string]string)
	headerPattern := regexp.MustCompile(`-H\s+'([^']*)'|-H\s+"([^"]*)"`)
	headerMatches := headerPattern.FindAllStringSubmatch(input, -1)
	for _, match := range headerMatches {
		header := match[1]
		if header == "" {
			header = match[2]
		}
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// 解析请求体
	body := make(map[string]interface{})
	dataPattern := regexp.MustCompile(`-d\s+'([^']*)'|-d\s+"([^"]*)"`)
	dataMatches := dataPattern.FindStringSubmatch(input)
	if len(dataMatches) > 1 {
		data := dataMatches[1]
		if data == "" {
			data = dataMatches[2]
		}
		if err := json.Unmarshal([]byte(data), &body); err != nil {
			// 如果不是 JSON，则作为普通字符串处理
			body["data"] = data
		}
	}

	// 解析查询参数
	params := make(map[string]interface{})
	if strings.Contains(url, "?") {
		parts := strings.SplitN(url, "?", 2)
		url = parts[0]
		queryParams := strings.Split(parts[1], "&")
		for _, param := range queryParams {
			kv := strings.SplitN(param, "=", 2)
			if len(kv) == 2 {
				params[kv[0]] = kv[1]
			}
		}
	}

	// 创建请求对象
	request := types.Request{
		RequestID:  fmt.Sprintf("req_%d", time.Now().UnixNano()),
		Method:     method,
		URL:        url,
		Headers:    headers,
		Body:       body,
		Params:     params,
		Timestamp:  time.Now(),
	}

	return request, nil
} 