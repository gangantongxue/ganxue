package middleware

import (
	"context"
	"log"

	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func LogMid() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		path := ctx.Request.URI().PathOriginal()

		ctx.Next(c)
		duration := time.Since(start)

		log.Println("Request: ", string(path), "Duration: ", duration.String())
	}
}
