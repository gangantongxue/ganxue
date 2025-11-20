package auth

import (
	"context"
	"ganxue-server/utils/db/mysql"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateUsernameRequest struct {
	NewUsername string `json:"new_username"`
}

// UpdateUsername 更新用户名
func UpdateUsername() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID := c.Value("userID")
		if userID == nil {
			ctx.JSON(401, map[string]string{"message": "未授权，请重新登录"})
			return
		}

		var req UpdateUsernameRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, map[string]string{"message": "请求格式错误"})
			return
		}

		if req.NewUsername == "" {
			ctx.JSON(400, map[string]string{"message": "用户名格式不正确"})
			return
		}

		// 验证用户名长度
		if len(req.NewUsername) < 3 || len(req.NewUsername) > 20 {
			ctx.JSON(400, map[string]string{"message": "用户名格式不正确"})
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

		// 更新用户名
		user.UserName = req.NewUsername
		if err := mysql.Update(user); err != nil {
			ctx.JSON(500, map[string]string{"message": "服务器错误"})
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"message": "用户名更新成功",
			"data": map[string]string{
				"user_name": req.NewUsername,
			},
		})
	}
}