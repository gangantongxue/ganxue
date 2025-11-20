package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"ganxue-server/global"
	"ganxue-server/model/answer_model"
	"ganxue-server/utils/db/mongodb"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gangantongxue/ggl"
	"go.mongodb.org/mongo-driver/bson"
)

func RunCode() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取请求体
		userCode := struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		}{}
		if err := ctx.Bind(&userCode); err != nil {
			ggl.Error("解析请求体失败", ggl.Err(err))
			ctx.JSON(400, map[string]string{"message": "解析请求体失败"})
			return
		}
		var ans answer_model.Answer
		if err := mongodb.Find(global.ANSWER, bson.M{"id": userCode.ID}, &ans); err != nil {
			ggl.Error("查询数据失败", ggl.Err(err.ToError()))
			ctx.JSON(400, map[string]string{"message": "查询数据失败"})
			return
		}
		ctx.JSON(100, map[string]string{
			"message": "等待运行中...",
		})

		userID := fmt.Sprintf("%d", c.Value("userID").(uint))

		err := global.RDB.Set(context.Background(), "userCode"+userID, userCode.Code, 15*time.Minute).Err()
		if err != nil {
			ggl.Error("将用户代码加入缓存失败", ggl.Err(err))
			ctx.JSON(500, map[string]string{"message": "将用户代码加入缓存失败"})
			return
		}
		err = global.RDB.Set(context.Background(), "input"+userID, ans.Input, 15*time.Minute).Err()
		if err != nil {
			ggl.Error("将用户输入加入缓存失败", ggl.Err(err))
			ctx.JSON(500, map[string]string{"message": "将用户输入加入缓存失败"})
			return
		}

		switch userCode.ID[0] {
		case '0':
			err := global.RDB.RPush(context.Background(), "runcode-go", userID).Err()
			if err != nil {
				ggl.Error("将用户ID加入队列失败", ggl.Err(err))
				ctx.JSON(500, map[string]string{"message": "将用户ID加入队列失败"})
				return
			}
		case '1':
			err := global.RDB.RPush(context.Background(), "runcode-c", userID).Err()
			if err != nil {
				ggl.Error("将用户ID加入队列失败", ggl.Err(err))
				ctx.JSON(500, map[string]string{"message": "将用户ID加入队列失败"})
				return
			}
		case '2':
			err := global.RDB.RPush(context.Background(), "runcode-cpp", userID).Err()
			if err != nil {
				ggl.Error("将用户ID加入队列失败", ggl.Err(err))
				ctx.JSON(500, map[string]string{"message": "将用户ID加入队列失败"})
				return
			}
		}

		r, err := global.RDB.BLPop(context.Background(), 15*time.Minute, "runcode-finish"+userID).Result()
		if err != nil {
			ggl.Error("从队列中获取结果失败", ggl.Err(err))
			ctx.JSON(500, map[string]string{"message": "从队列中获取结果失败"})
			return
		}
		result := struct {
			Status int    `json:"status"`
			Output string `json:"output"`
		}{}

		err = json.Unmarshal([]byte(r[1]), &result)
		if err != nil {
			ggl.Error("解析队列中的结果失败", ggl.Err(err))
			ctx.JSON(500, map[string]string{"message": "解析队列中的结果失败"})
			return
		}

		if result.Status == 1 {
			ctx.JSON(500, map[string]string{
				"message": "运行失败",
				"output":  result.Output,
			})
		} else if result.Status == 2 {
			ctx.JSON(206, map[string]string{
				"message": "运行失败",
				"output":  result.Output,
			})
		} else if result.Status == 0 {
			if result.Output == ans.Output {
				ctx.JSON(200, map[string]string{
					"message": "运行成功",
					"output":  result.Output,
				})
			} else {
				ctx.JSON(206, map[string]string{
					"message": "运行失败",
					"output":  result.Output,
				})
			}
		}
	}
}
