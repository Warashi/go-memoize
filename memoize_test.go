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
		name    string
		args    args
		wantDst interface{}
		wantErr bool
	}{
		{
			name: "Normal Case",
			args: args{
				key: key("normal_case"),
				f: func() interface{} {
					return 1
				},
			},
			wantDst: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst interface{}
			err := Call(tt.args.key, &dst, tt.args.f)
			if !assert.Equal(t, tt.wantDst, dst) {
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestParallelCall(t *testing.T) {
	type key string
	var k key = "key"
	var count uint64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			var dst interface{}
			assert.NoError(t, Call(k, &dst, func() interface{} {
				atomic.AddUint64(&count, 1)
				time.Sleep(1 * time.Second)
				return 1
			}))
			assert.Equal(t, 1, dst)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, uint64(1), count)
}
