package map_ttl

import (
	"sync"
	"time"
)

const (
	Del    = 1
	Re_set = 2
)

type ttl_data struct {
	Close_chan   chan int
	Timeout_time time.Time
}
type Map_ttl struct {
	sync.RWMutex
	data      map[interface{}]interface{}
	ttl       map[interface{}]ttl_data
	data_chan *chan interface{}
}

func (slf *Map_ttl) tll(key interface{}, ttl time.Duration, close_chan chan int) {
	defer close(close_chan)
	select {
	case <-time.After(ttl):
		slf.Lock()
		defer slf.Unlock()
		if slf.data_chan != nil {
			if v, ok := slf.data[key]; ok {
				*slf.data_chan <- v
			}
		}
		delete(slf.ttl, key)
		delete(slf.data, key)
	case flag := <-close_chan:
		if flag == Del {
			delete(slf.ttl, key)
			delete(slf.data, key)
		}
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
func (slf *Map_ttl) SetData(key, value interface{}) {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		slf.data[key] = value
	}
}
func (slf *Map_ttl) Set(key, value interface{}, ttl time.Duration) {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if v, ok := slf.ttl[key]; ok {
			v.Close_chan <- Re_set
		}
		if ttl == 0 {
			ttl = time.Minute * 65535
		}
		Close_chan := make(chan int)
		slf.ttl[key] = ttl_data{
			Timeout_time: time.Now().Add(ttl),
			Close_chan:   Close_chan,
		}
		go slf.tll(key, ttl, Close_chan)

		slf.data[key] = value
	}
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
		if len(slf.data) == 0 {
			break
		}
	}
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
