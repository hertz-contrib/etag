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
	h.Use(etag.New(etag.WithGenerator(
		func(ctx context.Context, c *app.RequestContext) []byte {
			return []byte("my-custom-etag")
		},
	)))
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "pong")
	})
	h.Spin()
}
