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
	Flag              int
	NewTtl            time.Duration
	timeoutDeleteFlag bool
}
type MapTtl struct {
	sync.RWMutex
	data         map[interface{}]data
	callbackChan *chan interface{}
}

func (slf *MapTtl) tll(key interface{}, ttl time.Duration, _chan chan comData, DeleteFlag bool) {
	timeoutDeleteFlag := DeleteFlag
	for {
		var timeoutChan <-chan time.Time
		if ttl > 0 {
			timeoutChan = time.After(ttl)
		}
		select {
		case <-timeoutChan:
			slf.Lock()
			if slf.callbackChan != nil {
				v, ok := slf.data[key]
				if ok {
					*slf.callbackChan <- v.value
				}
			}

			slf.Unlock()
			if timeoutDeleteFlag == true {
				delete(slf.data, key)
				break
			}
		case data := <-_chan:
			if data.Flag == Del {
				delete(slf.data, key)
			} else if data.Flag == Reset {
				ttl = data.NewTtl
				timeoutDeleteFlag = data.timeoutDeleteFlag
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
	slf.data = make(map[interface{}]data)
}
func (slf *MapTtl) SetData(key, value interface{}) bool {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if v, ok := slf.data[key]; ok {
			v.value = value
			slf.data[key] = v
			return true
		}
	}
	return false
}
func (slf *MapTtl) UnsafeSetData(key, value interface{}) bool {
	if slf.data != nil {
		if v, ok := slf.data[key]; ok {
			v.value = value
			slf.data[key] = v
			return true
		}
	}
	return false
}
func (slf *MapTtl) Set(key, value interface{}, ttl time.Duration, TimeOutDelete bool) {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if v, ok := slf.data[key]; ok {
			v.CloseChan <- comData{
				Flag:   Reset,
				NewTtl: ttl,
			}
		}
		CloseChan := make(chan comData)
		slf.data[key] = data{
			TimeoutTime: time.Now().Add(ttl),
			CloseChan:   CloseChan,
			value:       value,
		}
		go slf.tll(key, ttl, CloseChan, TimeOutDelete)
	}
}
func (slf *MapTtl) Get(key interface{}) interface{} {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if v, ok := slf.data[key]; ok {
			return v.value
		}
	}
	return nil
}
func (slf *MapTtl) Del(key interface{}) {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if v, ok := slf.data[key]; ok {
			v.CloseChan <- comData{
				Flag: Del,
			}
		}
	}
}
func (slf *MapTtl) GetTtl(key interface{}) time.Duration {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		if v, ok := slf.data[key]; ok {
			return v.TimeoutTime.Sub(time.Now())
		}
	}
	return -1
}
func (slf *MapTtl) Clear() {
	slf.Lock()
	defer slf.Unlock()
	for _, data := range slf.data {
		data.CloseChan <- comData{
			Flag: Del,
		}
	}
	for {
		if len(slf.data) == 0 {
			break
		}
	}
}
func (slf *MapTtl) Len() int {
	slf.Lock()
	defer slf.Unlock()
	if slf.data == nil {
		return 0
	} else {
		return len(slf.data)
	}
}
func (slf *MapTtl) Range(f func(key interface{}, value interface{})) {
	slf.Lock()
	defer slf.Unlock()
	if slf.data != nil {
		for k, v := range slf.data {
			f(k, v.value)
		}
	}
}
