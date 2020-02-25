package map_ttl

import (
	"sync/atomic"
	"time"
)

const (
	Set   = 0
	Del   = 1
	Reset = 2
	Get   = 3
	Clear = 4
)

type data struct {
	changeChan chan comData
	value      interface{}
}
type MapData struct {
	Flag       int
	key        interface{}
	value      interface{}
	ChangeChan chan comData
	getChan    chan interface{}
	Ttl        time.Duration
	deleteFlag bool
}

type comData struct {
	Flag              int
	NewTtl            time.Duration
	timeoutDeleteFlag bool
	value             interface{}
}

type MapTtl struct {
	data         map[interface{}]data
	mapChan      chan MapData
	callbackChan *chan interface{}
	len          int64
}

func (slf *MapTtl) Init(callbackChan *chan interface{}) {
	slf.data = make(map[interface{}]data)
	slf.mapChan = make(chan MapData, 1000)
	slf.callbackChan = callbackChan
	go slf.goMap()
	time.Sleep(time.Second)
}

func (slf *MapTtl) goMap() {
	for {
		select {
		case v := <-slf.mapChan:
			if v.Flag == Set {
				data := data{
					changeChan: v.ChangeChan,
					value:      v.value,
				}
				if value, ok := slf.data[v.key]; ok {
					value.changeChan <- comData{
						Flag:              Reset,
						NewTtl:            v.Ttl,
						timeoutDeleteFlag: v.deleteFlag,
						value:             v.value,
					}
				} else {
					atomic.AddInt64(&slf.len, 1)
					slf.data[v.key] = data
					go slf.tll(v.key, v.value, v.Ttl, v.ChangeChan, v.deleteFlag)
				}
			} else if v.Flag == Del {
				if value, ok := slf.data[v.key]; ok {
					atomic.AddInt64(&slf.len, -1)
					delete(slf.data, v.key)
					value.changeChan <- comData{
						Flag: Del,
					}
				}
			} else if v.Flag == Clear {
				for _, v := range slf.data {
					v.changeChan <- comData{
						Flag: Del,
					}
				}
				slf.data = make(map[interface{}]data)
			} else if v.Flag == Get {
				if value, ok := slf.data[v.key]; ok {
					v.getChan <- value.value
				} else {
					v.getChan <- nil
				}
			}
		}
	}
}

func (slf *MapTtl) Set(key, value interface{}, ttl time.Duration, TimeOutDelete bool) {
	slf.mapChan <- MapData{
		Flag:       Set,
		key:        key,
		value:      value,
		ChangeChan: make(chan comData),
		Ttl:        ttl,
		deleteFlag: TimeOutDelete,
	}
}

func (slf *MapTtl) tll(key interface{}, value interface{}, ttl time.Duration, _chan chan comData, DeleteFlag bool) {
	timeoutDeleteFlag := DeleteFlag
	for {
		var timeoutChan <-chan time.Time
		if ttl > 0 {
			timeoutChan = time.After(ttl)
		}
		select {
		case <-timeoutChan:
			if slf.callbackChan != nil {
				*slf.callbackChan <- value
			}
			if timeoutDeleteFlag == true {
				close(_chan)
				return
			}

		case data := <-_chan:
			if data.Flag == Del {
				close(_chan)
				return
			} else if data.Flag == Reset {
				ttl = data.NewTtl
				timeoutDeleteFlag = data.timeoutDeleteFlag
				value = data.value
			}
		}
	}
}

func (slf *MapTtl) Del(key interface{}) {
	slf.mapChan <- MapData{
		Flag: Del,
		key:  key,
	}

}

func (slf *MapTtl) Clear() {
	slf.mapChan <- MapData{
		Flag: Clear,
	}
}

//The value of length is not necessarily accurate,Because set() is asynchronous
func (slf *MapTtl) Len() int64 {
	return atomic.LoadInt64(&slf.len)
}

func (slf *MapTtl) Get(key interface{}) interface{} {
	GetChan := make(chan interface{})
	slf.mapChan <- MapData{
		Flag:    Get,
		key:     key,
		getChan: GetChan,
	}
	return <-GetChan
}
