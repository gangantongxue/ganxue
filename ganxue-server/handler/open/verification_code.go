package open

import (
	"context"
	"ganxue-server/utils/ver_code"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gangantongxue/ggl"
)

func VerCode() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取邮箱
		email := ctx.Query("email")
		if err := ver_code.GetVerCode(email); err != nil {
			ggl.Error("验证码获取失败", ggl.Err(err.ToError()))
			ctx.String(400, "获取验证码失败")
			return
		}
		ctx.String(200, "获取验证码成功")
	}
}
