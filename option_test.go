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
	"fmt"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
)

func TestOptions(t *testing.T) {
	next := func(ctx context.Context, c *app.RequestContext) bool {
		return true
	}
	gen := func(ctx context.Context, c *app.RequestContext) []byte {
		return []byte("hello world")
	}
	options := newOptions(
		WithWeak(),
		WithNext(next),
		WithGenerator(gen),
	)
	assert.True(t, options.weak)
	assert.DeepEqual(t, fmt.Sprintf("%p", next), fmt.Sprintf("%p", options.next))
	assert.DeepEqual(t, fmt.Sprintf("%p", gen), fmt.Sprintf("%p", options.generator))
}

func TestDefaultOptions(t *testing.T) {
	options := newOptions()
	assert.False(t, options.weak)
	assert.Nil(t, options.next)
	assert.Nil(t, options.generator)
}
