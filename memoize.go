package memoize

import (
	"sync"
)

var (
	Default Memoize
)

type Memoize struct {
	memo sync.Map
	once sync.Map
	done sync.Map
}

func (m *Memoize) Call(key interface{}, f func() (ret interface{})) interface{} {
	ret, ok := m.memo.Load(key)
	if !ok {
		ret = m.call(key, f)
	}
	return ret
}

func (m *Memoize) call(key interface{}, f func() (ret interface{})) (ret interface{}) {
	oi, _ := m.once.LoadOrStore(key, &sync.Once{})
	o := oi.(*sync.Once)
	di, _ := m.done.LoadOrStore(key, make(chan struct{}))
	done := di.(chan struct{})
	o.Do(func() {
		m.memo.Store(key, f())
		close(done)
	})
	<-done
	ret, _ = m.memo.Load(key)
	return ret
}

// Call memoize and call function f, then store return value to dst
func Call(key interface{}, f func() interface{}) interface{} {
	return Default.Call(key, f)
}
