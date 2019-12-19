package map_ttl

import (
	"context"
	"sync"
	"time"
)

type ttl_data struct {
	CancelFunc context.CancelFunc
	ttl        time.Time
}
type Map_ttl struct {
	sync.RWMutex
	data      map[interface{}]interface{}
	ttl       map[interface{}]ttl_data
	data_chan *chan interface{}
}

func (slf *Map_ttl) tll(key interface{}, ctx context.Context) {
	select {
	case <-ctx.Done():
		slf.Lock()
		defer slf.Unlock()
		if slf.data_chan != nil {
			if v, ok := slf.data[key]; ok {
				*slf.data_chan <- v
			}

		}
		delete(slf.ttl, key)
		delete(slf.data, key)
	}
}
func (slf *Map_ttl) Set_callback(callback_chan *chan interface{}) {
	slf.Lock()
	defer slf.Unlock()
	slf.data_chan = callback_chan
}
func (slf *Map_ttl) Init() {
	slf.data = make(map[interface{}]interface{})
	slf.ttl = make(map[interface{}]ttl_data)
}

func (slf *Map_ttl) Set(key, value interface{}, ttl time.Duration) {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if ttl > 0 {
			ctx := context.Background()
			timeout_time := time.Now().Add(ttl)
			if v, ok := slf.ttl[key]; ok {
				v.CancelFunc()
			} else {
				time_ctx, cancel_fun := context.WithDeadline(ctx, timeout_time)
				slf.ttl[key] = ttl_data{
					CancelFunc: cancel_fun,
					ttl:        timeout_time,
				}
				go slf.tll(key, time_ctx)
			}
		}
	}
	slf.data[key] = value
}
func (slf *Map_ttl) Get(key interface{}) interface{} {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if v, ok := slf.data[key]; ok {
			return v
		}
	}
	return nil
}
func (slf *Map_ttl) Del(key interface{}) {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			delete(slf.data, key)
			v.CancelFunc()
		}
	}
}
func (slf *Map_ttl) Get_ttl(key interface{}) time.Duration {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			return v.ttl.Sub(time.Now())
		}
	}
	return -1
}
func (slf *Map_ttl) Clear() {
	slf.Lock()
	defer slf.Unlock()
	slf.data = make(map[interface{}]interface{})
	slf.ttl = make(map[interface{}]ttl_data)
}
func (slf *Map_ttl) Len() int {
	slf.Lock()
	defer slf.Unlock()
	if slf.data == nil {
		return 0
	} else {
		return len(slf.data)
	}
}
func (slf *Map_ttl) Range(f func(key interface{}, value interface{})) {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		for k, v := range slf.data {
			f(k, v)
		}
	}
}
