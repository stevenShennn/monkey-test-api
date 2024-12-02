package api

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"monkey-test-api/internal/input"
	"monkey-test-api/internal/logger"
	"monkey-test-api/internal/store"
	"monkey-test-api/internal/task"
)

//go:embed templates/*
var templatesFS embed.FS

type Handler struct {
	parser    input.Parser
	store     store.Store
	generator task.TaskGenerator
	templates *template.Template
}

// NewHandler 创建新的 API 处理器
func NewHandler(parser input.Parser, store store.Store, generator task.TaskGenerator) *Handler {
	// 加载模板
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
	}).ParseFS(templatesFS, "templates/*.html"))

	return &Handler{
		parser:    parser,
		store:     store,
		generator: generator,
		templates: tmpl,
	}
}

// RenderIndex 渲染首页
func (h *Handler) RenderIndex(c *gin.Context) {
	// 获取历史记录
	history, err := h.store.GetRequestsByTime(c.Request.Context(), 10)
	if err != nil {
		logger.Errorf("获取历史记录失败: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取历史记录失败",
		})
		return
	}

	// 读取样式和脚本文件
	styles, err := templatesFS.ReadFile("templates/styles.css")
	if err != nil {
		logger.Errorf("读取样式文件失败: %v", err)
		styles = []byte{}
	}

	scripts, err := templatesFS.ReadFile("templates/scripts.js")
	if err != nil {
		logger.Errorf("读取脚本文件失败: %v", err)
		scripts = []byte{}
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"History": history,
		"Styles":  string(styles),
		"Scripts": string(scripts),
	})
}

// CreateTestRequest 处理创建测试请求的 API
func (h *Handler) CreateTestRequest(c *gin.Context) {
	// 1. 获取 cURL 命令
	var req struct {
		Curl string `json:"curl" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("解析请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	// 2. 解析 cURL 命令
	requests, err := h.parser.Parse(req.Curl)
	if err != nil {
		logger.Errorf("解析 cURL 命令失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 cURL 命令"})
		return
	}

	var results []map[string]interface{}
	
	// 3. 处理每个解析出的请求
	for _, request := range requests {
		// 3.1 存储父请求
		if err := h.store.InsertRequest(c.Request.Context(), &store.Request{
			RequestID:  request.RequestID,
			Method:    request.Method,
			URL:       request.URL,
			Headers:   request.Headers,
			Body:      request.Body,
			Params:    request.Params,
			Timestamp: request.Timestamp,
		}); err != nil {
			logger.Errorf("存储父请求失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "存储请求失败"})
			return
		}

		// 3.2 生成测试对象
		testObjects, err := h.generator.GenerateTestObjects(&task.Request{
			RequestID:  request.RequestID,
			Method:    request.Method,
			URL:       request.URL,
			Headers:   request.Headers,
			Body:      request.Body,
			Params:    request.Params,
			Timestamp: request.Timestamp,
		})
		if err != nil {
			logger.Errorf("生成测试对象失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成测试用例失败"})
			return
		}

		// 3.3 存储测试对象
		storeTestObjects := make([]store.TestObject, len(testObjects))
		for i, obj := range testObjects {
			storeTestObjects[i] = store.TestObject{
				TestID:         obj.TestID,
				ParentRequestID: obj.ParentRequestID,
				Method:         obj.Method,
				URL:            obj.URL,
				Headers:        obj.Headers,
				Body:           obj.Body,
				Params:         obj.Params,
				Reason:         obj.Reason,
				Status:         obj.Status,
				Timestamp:      obj.Timestamp,
			}
		}

		if err := h.store.InsertTestObjects(c.Request.Context(), storeTestObjects); err != nil {
			logger.Errorf("存储测试对象失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "存储测试用例失败"})
			return
		}

		// 3.4 添加结果
		results = append(results, map[string]interface{}{
			"requestID":     request.RequestID,
			"testCount":     len(testObjects),
			"testObjects":   testObjects,
		})
	}

	// 4. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "成功创建测试请求",
		"results": results,
	})
} 