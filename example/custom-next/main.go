package main

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/etag"
)

func main() {
	h := server.Default()
	h.Use(etag.New(etag.WithNext(
		func(ctx context.Context, c *app.RequestContext) bool {
			if string(c.Method()) == http.MethodPost {
				return true
			} else {
				return false
			}
		},
	)))
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "pong")
	})
	h.Spin()
}
