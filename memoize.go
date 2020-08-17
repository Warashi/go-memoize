package memoize

import (
	"reflect"
	"sync"
)

var (
	Default Memoize
)

// An InvalidCallError describes an invalid argument passed to Call.
// (The argument to Call must be a non-nil pointer.)
type InvalidCallError struct {
	Type reflect.Type
}

func (e *InvalidCallError) Error() string {
	if e.Type == nil {
		return "memoize: dst is (nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "memoize: dst is (non-pointer " + e.Type.String() + ")"
	}
	return "memoize: dst is (nil " + e.Type.String() + ")"
}

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
