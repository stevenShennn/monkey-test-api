package api

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"monkey-test-api/internal/parser/curl"
	"monkey-test-api/internal/store"
	"monkey-test-api/internal/task"
)

type Handler struct {
	templates     *template.Template
	parser        *curl.Parser
	store         store.Store
	taskGenerator *task.Generator
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境中应该更严格地检查来源
	},
}

func NewHandler(parser *curl.Parser, store store.Store, taskGenerator *task.Generator) *Handler {
	// 解析所有模板
	tmpl, err := template.ParseGlob(filepath.Join("internal", "api", "templates", "*.html"))
	if err != nil {
		panic(err)
	}
	
	return &Handler{
		templates:     tmpl,
		parser:        parser,
		store:         store,
		taskGenerator: taskGenerator,
	}
}

// RenderIndex 渲染首页
func (h *Handler) RenderIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "layout", nil)
}

// RenderExecute 渲染执行页面
func (h *Handler) RenderExecute(c *gin.Context) {
	c.HTML(http.StatusOK, "layout", nil)
}

// CreateTestRequest 处理创建测试请求的API
func (h *Handler) CreateTestRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// HandleWebSocket 处理 WebSocket 连接
func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	// 处理 WebSocket 消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// 处理接收到的消息
		// TODO: 实现具体的测试执行逻辑

		// 发送响应
		if err := conn.WriteMessage(messageType, message); err != nil {
			break
		}
	}
}
