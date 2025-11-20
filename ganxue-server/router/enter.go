package router

import (
	"ganxue-server/handler/test"
	"ganxue-server/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

var (
	AuthGroup *route.RouterGroup
	OpenGroup *route.RouterGroup
)

func Init(h *server.Hertz) {
	// 使用自动刷新的JWT中间件
	AuthGroup = h.Group("/auth", middleware.TokenAutoRefreshMiddleware())
	OpenGroup = h.Group("/open")
	h.GET("/test", test.Test())

	OpenRouter()
	AuthRouter()
}
