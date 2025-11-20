package open

import (
	"context"
	"ganxue-server/utils/token"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gangantongxue/ggl"
)

// Refresh 刷新短token
func Refresh() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		longToken := ctx.Cookie("long_token")
		userID, err := token.ParseToken(string(longToken))
		if err != nil {
			ctx.AbortWithStatus(401)
			ggl.Debug("Parse long token error")
			return
		}

		// 获取当前短token（如果有的话）
		authHeader := ctx.Request.Header.Get("Authorization")
		oldShortToken := authHeader
		if len(oldShortToken) > 7 && oldShortToken[:7] == "Bearer " {
			oldShortToken = oldShortToken[7:]
		}

		var shortToken string
		if oldShortToken != "" && oldShortToken != authHeader {
			// 如果有旧短token，使用刷新功能
			shortToken, err = token.RefreshShortToken(oldShortToken, userID)
		} else {
			// 否则生成新的短token
			shortToken, err = token.GenerateShortToken(userID)
		}

		if err != nil {
			ctx.AbortWithStatus(401)
			ggl.Debug("Generate/Refresh short token error", ggl.Err(err.ToError()))
			return
		}

		ctx.JSON(200, map[string]string{
			"token": shortToken,
		})
	}
}
