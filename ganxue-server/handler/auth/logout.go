package auth

import (
	"context"
	"ganxue-server/utils/token"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/gangantongxue/ggl"
)

// Logout 用户登出
func Logout() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从context中获取用户ID
		userID, ok := c.Value("userID").(uint)
		if !ok {
			ctx.JSON(400, map[string]string{"message": "获取用户信息失败"})
			return
		}

		// 删除用户的所有短token
		err := token.LogoutUser(userID)
		if err != nil {
			ggl.Error("删除用户短token失败", ggl.Err(err.ToError()))
			ctx.JSON(500, map[string]string{"message": "登出失败"})
			return
		}

		// 删除自动登录token cookie
		ctx.SetCookie(
			"auto_login_token",
			"",
			-1, // 设置过期时间为过去
			"/",
			"",
			protocol.CookieSameSiteNoneMode,
			false, // 不设置Secure
			true,  // 设置HttpOnly
		)

		ctx.JSON(200, map[string]string{"message": "登出成功"})
	}
}
