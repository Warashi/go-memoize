package memoize

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	type key string
	type args struct {
		key interface{}
		f   func() interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Normal Case",
			args: args{
				key: key("normal_case"),
				f: func() interface{} {
					return 1
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Call(tt.args.key, tt.args.f)
			if !assert.Equal(t, tt.want, actual) {
			}
		})
	}
}

func TestParallelCall(t *testing.T) {
	type key string
	var k key = "key"
	var count int64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			actual := Call(k, func() interface{} {
				atomic.AddInt64(&count, 1)
				time.Sleep(1 * time.Second)
				return 1
			})
			assert.Equal(t, 1, actual)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(1), count)
}

func TestParallelCall2(t *testing.T) {
	type key int
	expc := 100
	var count int64
	var wg sync.WaitGroup
	for i := 0; i < expc; i++ {
		i := i
		wg.Add(1)
		go func() {
			actual := Call(i, func() interface{} {
				atomic.AddInt64(&count, 1)
				time.Sleep(1 * time.Second)
				return i
			})
			assert.Equal(t, i, actual)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(expc), count)
}
