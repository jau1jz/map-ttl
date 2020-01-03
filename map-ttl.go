package map_ttl

import (
	"sync"
	"time"
)

const (
	Del   = 1
	Reset = 2
)

type data struct {
	value        interface{}
	Close_chan   chan int
	Timeout_time time.Time
}
type Map_ttl struct {
	sync.RWMutex
	ttl       map[interface{}]data
	data_chan *chan interface{}
}

func (slf *Map_ttl) tll(key interface{}, ttl time.Duration, close_chan chan int) {
	select {
	case <-time.After(ttl):
		slf.Lock()
		defer slf.Unlock()
		if slf.data_chan != nil {
			if v, ok := slf.ttl[key]; ok {
				*slf.data_chan <- v.value
			}
		}
		delete(slf.ttl, key)
	case flag := <-close_chan:
		if flag == Del {
			delete(slf.ttl, key)
		}
	}

}
func (slf *Map_ttl) Set_callback(callback_chan *chan interface{}) {
	slf.Lock()
	defer slf.Unlock()
	slf.data_chan = callback_chan
}
func (slf *Map_ttl) Init() {
	slf.ttl = make(map[interface{}]data)
}
func (slf *Map_ttl) SetData(key, value interface{}) bool {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			v.value = value
			slf.ttl[key] = v
			return true
		}
	}
	return false
}
func (slf *Map_ttl) UnsafeSetData(key, value interface{}) bool {
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			v.value = value
			slf.ttl[key] = v
			return true
		}
	}
	return false
}
func (slf *Map_ttl) Set(key, value interface{}, ttl time.Duration) {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			v.Close_chan <- Reset
		}
		if ttl == 0 {
			ttl = time.Minute * 65535
		}
		Close_chan := make(chan int)
		slf.ttl[key] = data{
			Timeout_time: time.Now().Add(ttl),
			Close_chan:   Close_chan,
			value:        value,
		}
		go slf.tll(key, ttl, Close_chan)
	}
}
func (slf *Map_ttl) Get(key interface{}) interface{} {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			return v.value
		}
	}
	return nil
}
func (slf *Map_ttl) Del(key interface{}) {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			v.Close_chan <- Del
		}
	}
}
func (slf *Map_ttl) Get_ttl(key interface{}) time.Duration {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			return v.Timeout_time.Sub(time.Now())
		}
	}
	return -1
}
func (slf *Map_ttl) Clear() {
	slf.Lock()
	defer slf.Unlock()
	for _, data := range slf.ttl {
		data.Close_chan <- Del
	}
	for {
		if len(slf.ttl) == 0 {
			break
		}
	}
}
func (slf *Map_ttl) Len() int {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl == nil {
		return 0
	} else {
		return len(slf.ttl)
	}
}
func (slf *Map_ttl) Range(f func(key interface{}, value interface{})) {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		for k, v := range slf.ttl {
			f(k, v.value)
		}
	}
}
