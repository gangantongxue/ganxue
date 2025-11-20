package middleware

import (
	"context"
	"ganxue-server/utils/token"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gangantongxue/ggl"
)

// TokenRefreshMiddleware 自动刷新短token的中间件
func TokenRefreshMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取当前短token
		shortToken := string(ctx.Request.Header.Get("Authorization"))
		if len(shortToken) > 7 && shortToken[:7] == "Bearer " {
			shortToken = shortToken[7:]
		}

		if shortToken == "" {
			ctx.Next(c)
			return
		}

		// 验证当前短token是否有效
		userID, err := token.ParseShortToken(shortToken)
		if err != nil {
			ctx.Next(c)
			return
		}

		// 生成新的短token
		newShortToken, err := token.RefreshShortToken(shortToken, userID)
		if err != nil {
			ggl.Debug("刷新短token失败", ggl.Err(err.ToError()))
			ctx.Next(c)
			return
		}

		// 在响应头中返回新的短token
		ctx.Response.Header.Set("New-Access-Token", newShortToken)

		ctx.Next(c)
	}
}

// TokenAutoRefreshMiddleware 带自动刷新功能的JWT中间件
func TokenAutoRefreshMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 解析token
		shortToken := string(ctx.Request.Header.Get("Authorization"))
		if len(shortToken) > 7 && shortToken[:7] == "Bearer " {
			shortToken = shortToken[7:]
		}

		userID, err := token.ParseShortToken(shortToken)
		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		// 自动刷新token - 生成新的短token
		newShortToken, refreshErr := token.RefreshShortToken(shortToken, userID)
		if refreshErr != nil {
			ggl.Debug("自动刷新短token失败", ggl.Err(refreshErr.ToError()))
		} else {
			// 在响应头中返回新的短token
			ctx.Response.Header.Set("New-Access-Token", newShortToken)
		}

		c = context.WithValue(c, "userID", userID)
		ctx.Next(c)
	}
}