// Copyright 2023 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package etag

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
)

func newTestEngine() *route.Engine {
	opt := config.NewOptions([]config.Option{})
	return route.NewEngine(opt)
}

func TestETagNext(t *testing.T) {
	t.Parallel()
	e := newTestEngine()
	e.Use(New(WithNext(
		func(ctx context.Context, c *app.RequestContext) bool {
			return true
		},
	)))
	resp := ut.PerformRequest(e, http.MethodGet, "/", nil)
	assert.DeepEqual(t, http.StatusNotFound, resp.Code)
}

func TestETagNotStatusOK(t *testing.T) {
	t.Parallel()
	e := newTestEngine()
	e.Use(New())
	e.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.SetStatusCode(http.StatusCreated)
	})
	resp := ut.PerformRequest(e, http.MethodGet, "/", nil)
	assert.DeepEqual(t, http.StatusCreated, resp.Code)
}

func TestETagNoBody(t *testing.T) {
	t.Parallel()
	e := newTestEngine()
	e.Use(New())
	e.GET("/", func(ctx context.Context, c *app.RequestContext) {})
	resp := ut.PerformRequest(e, http.MethodGet, "/", nil)
	assert.DeepEqual(t, http.StatusOK, resp.Code)
}

func TestETagNewETag(t *testing.T) {
	newETag := func(t *testing.T, headerIfNoneMatch, matched bool) {
		bodyString := "hello world"
		t.Helper()
		e := newTestEngine()
		e.Use(New())
		e.GET("/", func(ctx context.Context, c *app.RequestContext) {
			c.SetBodyString(bodyString)
		})
		resp := ut.PerformRequest(e, http.MethodGet, "/", nil)
		if headerIfNoneMatch {
			etag := `"non-match"`
			if matched {
				etag = `"11-222957957"`
			}
			resp = ut.PerformRequest(e, http.MethodGet, "/", nil, ut.Header{
				Key:   HeaderIfNoneMatch,
				Value: etag,
			})
		}
		if !headerIfNoneMatch || !matched {
			assert.DeepEqual(t, http.StatusOK, resp.Code)
			assert.DeepEqual(t, `"11-222957957"`, resp.Header().Get("Etag"))
			return
		}
		if matched {
			assert.DeepEqual(t, http.StatusNotModified, resp.Code)
			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.DeepEqual(t, 0, len(body))
		}
	}
	t.Parallel()
	t.Run("without HeaderIfNoneMatch", func(t *testing.T) {
		t.Parallel()
		newETag(t, false, false)
	})
	t.Run("with HeaderIfNoneMatch and not matched", func(t *testing.T) {
		t.Parallel()
		newETag(t, true, false)
	})
	t.Run("with HeaderIfNoneMatch and matched", func(t *testing.T) {
		t.Parallel()
		newETag(t, true, true)
	})
}

func TestETagWeakETag(t *testing.T) {
	weakEtag := func(t *testing.T, headerIfNoneMatch, matched bool) {
		bodyString := "hello world"
		t.Helper()
		e := newTestEngine()
		e.Use(New(WithWeak()))
		e.GET("/", func(ctx context.Context, c *app.RequestContext) {
			c.SetBodyString(bodyString)
		})
		resp := ut.PerformRequest(e, http.MethodGet, "/", nil)
		if headerIfNoneMatch {
			etag := `"W/non-match"`
			if matched {
				etag = `W/"11-222957957"`
			}
			resp = ut.PerformRequest(e, http.MethodGet, "/", nil, ut.Header{
				Key:   HeaderIfNoneMatch,
				Value: etag,
			})
		}
		if !headerIfNoneMatch || !matched {
			assert.DeepEqual(t, http.StatusOK, resp.Code)
			assert.DeepEqual(t, `W/"11-222957957"`, resp.Header().Get("Etag"))
			return
		}
		if matched {
			assert.DeepEqual(t, http.StatusNotModified, resp.Code)
			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.DeepEqual(t, 0, len(body))
		}
	}
	t.Parallel()
	t.Run("without HeaderIfNoneMatch", func(t *testing.T) {
		t.Parallel()
		weakEtag(t, false, false)
	})
	t.Run("with HeaderIfNoneMatch and not matched", func(t *testing.T) {
		t.Parallel()
		weakEtag(t, true, false)
	})
	t.Run("with HeaderIfNoneMatch and matched", func(t *testing.T) {
		t.Parallel()
		weakEtag(t, true, true)
	})
}

func TestETagGenerator(t *testing.T) {
	etagGenerator := func(t *testing.T, headerIfNoneMatch, matched bool) {
		bodyString := "hello world"
		t.Helper()
		e := newTestEngine()
		e.Use(New(WithGenerator(
			func(ctx context.Context, c *app.RequestContext) []byte {
				return []byte("my-custom-etag")
			})),
		)
		e.GET("/", func(ctx context.Context, c *app.RequestContext) {
			c.SetBodyString(bodyString)
		})
		resp := ut.PerformRequest(e, http.MethodGet, "/", nil)
		if headerIfNoneMatch {
			etag := "non-match"
			if matched {
				etag = "my-custom-etag"
			}
			resp = ut.PerformRequest(e, http.MethodGet, "/", nil, ut.Header{
				Key:   HeaderIfNoneMatch,
				Value: etag,
			})
		}
		if !headerIfNoneMatch || !matched {
			assert.DeepEqual(t, http.StatusOK, resp.Code)
			assert.DeepEqual(t, "my-custom-etag", resp.Header().Get("Etag"))
			return
		}
		if matched {
			assert.DeepEqual(t, http.StatusNotModified, resp.Code)
			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.DeepEqual(t, 0, len(body))
		}
	}
	t.Parallel()
	t.Run("without HeaderIfNoneMatch", func(t *testing.T) {
		t.Parallel()
		etagGenerator(t, false, false)
	})
	t.Run("with HeaderIfNoneMatch and not matched", func(t *testing.T) {
		t.Parallel()
		etagGenerator(t, true, false)
	})
	t.Run("with HeaderIfNoneMatch and matched", func(t *testing.T) {
		t.Parallel()
		etagGenerator(t, true, true)
	})
}
