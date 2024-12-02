package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupRouter 配置路由
func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// 设置 HTML 渲染器
	router.SetHTMLTemplate(handler.templates)

	// 首页路由
	router.GET("/", handler.RenderIndex)

	// API 路由组
	v1 := router.Group("/api/v1")
	{
		v1.POST("/test-requests", handler.CreateTestRequest)
		// ... 其他 API 路由
	}

	return router
} 