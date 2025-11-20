package auth

import (
	"context"
	"ganxue-server/utils/db/mysql"
	"ganxue-server/utils/password"
	"github.com/cloudwego/hertz/pkg/app"
)

type SecureDeleteUserRequest struct {
	Password string `json:"password"`
}

// SecureDeleteUser 安全删除用户（需要密码验证）
func SecureDeleteUser() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID := c.Value("userID")
		if userID == nil {
			ctx.JSON(401, map[string]string{"message": "未授权，请重新登录"})
			return
		}

		var req SecureDeleteUserRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, map[string]string{"message": "请求格式错误"})
			return
		}

		if req.Password == "" {
			ctx.JSON(400, map[string]string{"message": "密码不能为空"})
			return
		}

		// 查找用户
		user, err := mysql.FindUserByID(userID.(uint))
		if err != nil {
			ctx.JSON(500, map[string]string{"message": "服务器错误"})
			return
		}
		if user == nil {
			ctx.JSON(401, map[string]string{"message": "未授权，请重新登录"})
			return
		}

		// 验证密码
		if !password.ComparePasswords(user.Password, req.Password) {
			ctx.JSON(400, map[string]string{"message": "密码错误"})
			return
		}

		// 删除用户（级联删除用户信息）
		if err := mysql.Delete(user); err != nil {
			ctx.JSON(500, map[string]string{"message": "服务器错误"})
			return
		}

		ctx.JSON(200, map[string]string{"message": "账户删除成功"})
	}
}