package open

import (
	"context"
	"ganxue-server/utils/db/mysql"
	"ganxue-server/utils/password"
	"ganxue-server/utils/token"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/gangantongxue/ggl"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AutoLogin bool   `json:"autoLogin"`
}

// SignIn 登录
func SignIn() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		loginReq := LoginRequest{}
		// 解析请求体
		if err := ctx.Bind(&loginReq); err != nil {
			ctx.JSON(400, map[string]string{"message": "解析请求体失败"})
			return
		}
		// 查找用户
		_user, err := mysql.FindUserByEmail(loginReq.Email)
		// 用户不存在
		if err != nil {
			ctx.JSON(400, map[string]string{"message": "未注册"})
			return
		}

		// 密码校验
		if !password.ComparePasswords(_user.Password, loginReq.Password) {
			ctx.JSON(400, map[string]string{"message": "密码错误"})
			return
		}

		// 生成token
		var shortToken, longToken string
		if shortToken, err = token.GenerateShortToken(_user.ID); err != nil {
			ggl.Error("生成短token失败", ggl.Err(err.ToError()))
			ctx.JSON(500, map[string]string{"message": "服务器错误"})
			return
		}
		if longToken, err = token.GenerateLongToken(_user.ID); err != nil {
			ggl.Error("生成长token失败", ggl.Err(err.ToError()))
			ctx.JSON(500, map[string]string{"message": "服务器错误"})
			return
		}

		// 设置长token cookie用于刷新短token
		ctx.SetCookie(
			"long_token",
			longToken,
			60*60*24*7,
			"/api/open/refresh",
			"",
			protocol.CookieSameSiteNoneMode,
			false,
			true,
		)

		// 如果勾选了自动登录，生成自动登录token
		if loginReq.AutoLogin {
			var autoLoginToken string
			if autoLoginToken, err = token.GenerateAutoLoginToken(_user.ID); err != nil {
				ggl.Error("生成自动登录token失败", ggl.Err(err.ToError()))
				ctx.JSON(500, map[string]string{"message": "服务器错误"})
				return
			}
			// 设置自动登录token cookie
			ctx.SetCookie(
				"auto_login_token",
				autoLoginToken,
				60*60*24*7,
				"/",
				"",
				protocol.CookieSameSiteNoneMode,
				false,
				true,
			)
		}
		ctx.JSON(200, struct {
			Message string `json:"message"`
			Data    struct {
				Email    string `json:"email"`
				UserName string `json:"user_name"`
				Token    string `json:"token"`
			} `json:"data"`
		}{
			Message: "登录成功",
			Data: struct {
				Email    string `json:"email"`
				UserName string `json:"user_name"`
				Token    string `json:"token"`
			}{
				Email:    _user.Email,
				UserName: _user.UserName,
				Token:    shortToken,
			},
		})

	}
}

// CheckAutoLogin 检查自动登录token是否有效
func CheckAutoLogin() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从cookie中获取自动登录token
		autoLoginToken := string(ctx.Cookie("auto_login_token"))
		if autoLoginToken == "" {
			// 如果没有自动登录token，返回未登录状态
			ctx.JSON(401, map[string]string{"message": "未登录"})
			return
		}

		// 验证自动登录token
		userID, err := token.ParseAutoLoginToken(autoLoginToken)
		if err != nil {
			// 如果自动登录token无效，返回未登录状态
			ctx.JSON(401, map[string]string{"message": "自动登录token无效"})
			return
		}

		// 自动登录token有效，生成新的shortToken
		newShortToken, err := token.GenerateShortToken(userID)
		if err != nil {
			// 如果生成短token失败，返回服务器错误
			ctx.JSON(500, map[string]string{"message": "生成短token失败"})
			return
		}

		// 在响应头中返回新的shortToken
		ctx.Response.Header.Set("New-Access-Token", newShortToken)

		// 返回登录成功状态
		ctx.JSON(200, map[string]string{"message": "自动登录成功"})
	}
}
