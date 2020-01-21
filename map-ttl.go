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
	value       interface{}
	CloseChan   chan comData
	TimeoutTime time.Time
}

type comData struct {
	Flag   int
	NewTtl time.Duration
}
type MapTtl struct {
	sync.RWMutex
	ttl          map[interface{}]data
	callbackChan *chan interface{}
}

func (slf *MapTtl) tll(key interface{}, ttl time.Duration, _chan chan comData) {

	for {
		var timeoutChan <-chan time.Time
		if ttl > 0 {
			timeoutChan = time.After(ttl)
		}
		select {
		case <-timeoutChan:
			slf.Lock()
			if slf.callbackChan != nil {
				if v, ok := slf.ttl[key]; ok {
					*slf.callbackChan <- v.value
				}
			}
			delete(slf.ttl, key)
			slf.Unlock()
			break
		case data := <-_chan:
			if data.Flag == Del {
				delete(slf.ttl, key)
			} else if data.Flag == Reset {
				ttl = data.NewTtl
			}
		}
	}

}
func (slf *MapTtl) SetCallback(callbackChan *chan interface{}) {
	slf.Lock()
	defer slf.Unlock()
	slf.callbackChan = callbackChan
}
func (slf *MapTtl) Init() {
	slf.ttl = make(map[interface{}]data)
}
func (slf *MapTtl) SetData(key, value interface{}) bool {
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
func (slf *MapTtl) UnsafeSetData(key, value interface{}) bool {
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			v.value = value
			slf.ttl[key] = v
			return true
		}
	}
	return false
}
func (slf *MapTtl) Set(key, value interface{}, ttl time.Duration) {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			v.CloseChan <- comData{
				Flag:   Reset,
				NewTtl: ttl,
			}
		}
		CloseChan := make(chan comData)
		slf.ttl[key] = data{
			TimeoutTime: time.Now().Add(ttl),
			CloseChan:   CloseChan,
			value:       value,
		}
		go slf.tll(key, ttl, CloseChan)
	}
}
func (slf *MapTtl) Get(key interface{}) interface{} {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			return v.value
		}
	}
	return nil
}
func (slf *MapTtl) Del(key interface{}) {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			v.CloseChan <- comData{
				Flag: Del,
			}
		}
	}
}
func (slf *MapTtl) GetTtl(key interface{}) time.Duration {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		if v, ok := slf.ttl[key]; ok {
			return v.TimeoutTime.Sub(time.Now())
		}
	}
	return -1
}
func (slf *MapTtl) Clear() {
	slf.Lock()
	defer slf.Unlock()
	for _, data := range slf.ttl {
		data.CloseChan <- comData{
			Flag: Del,
		}
	}
	for {
		if len(slf.ttl) == 0 {
			break
		}
	}
}
func (slf *MapTtl) Len() int {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl == nil {
		return 0
	} else {
		return len(slf.ttl)
	}
}
func (slf *MapTtl) Range(f func(key interface{}, value interface{})) {
	slf.Lock()
	defer slf.Unlock()
	if slf.ttl != nil {
		for k, v := range slf.ttl {
			f(k, v.value)
		}
	}
}
