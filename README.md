# ETag (This is a community driven project)

ETag middleware for Hertz framework, inspired by [fiber-etag](https://github.com/gofiber/fiber/tree/master/middleware/etag).

The ETag (or entity tag) HTTP response header is an identifier for a specific version of a resource. 
It lets caches be more efficient and save bandwidth, as a web server does not need to resend a full response if the content was not changed. 
Additionally, etags help to prevent simultaneous updates of a resource from overwriting each other ("mid-air collisions").

## Install

```shell
go get github.com/hertz-contrib/etag
```

## Example

```go
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
	h.Use(etag.New())
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "pong")
	})
	h.Spin()
}
```

## Configuration

| Configuration   | Default | Description                                                    | Example                               |
|-----------------|---------|----------------------------------------------------------------|---------------------------------------|
| `WithWeak`      | `false` | Enable weak validator                                          | [example](./example/custom-weak)      |
| `WithNext`      | `nil`   | Defines a function to skip etag middleware when return is true | [example](./example/custom-next)      |
| `WithGenerator` | `nil`   | Custom etag generation logic                                   | [example](./example/custom-generator) |