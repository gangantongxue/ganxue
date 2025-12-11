package middleware

import (
	"context"
	"ganxue-server/utils/token"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func JwtMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 解析token
		authHeader := ctx.Request.Header.Get("Authorization")
		shortToken := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := token.ParseShortToken(shortToken)
		if err != nil {
			// 如果shortToken验证失败，尝试从cookie中获取自动登录token
			autoLoginToken := string(ctx.Cookie("auto_login_token"))
			if autoLoginToken != "" {
				// 验证自动登录token
				userID, err = token.ParseAutoLoginToken(autoLoginToken)
				if err != nil {
					// 如果自动登录token验证失败，返回401
					ctx.AbortWithStatus(401)
					return
				}
				// 自动登录token验证成功，生成新的shortToken并设置到响应头
				newShortToken, err := token.GenerateShortToken(userID)
				if err != nil {
					ctx.AbortWithStatus(500)
					return
				}
				// 设置新的shortToken到响应头
				ctx.Response.Header.Set("X-New-Token", newShortToken)
			} else {
				// 如果没有自动登录token，返回401
				ctx.AbortWithStatus(401)
				return
			}
		}

		c = context.WithValue(c, "userID", userID)
		ctx.Next(c)
	}
}
