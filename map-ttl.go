package map_ttl

import (
	"sync"
	"time"
)

const (
	Set   = 0
	Del   = 1
	Reset = 2
	Clear = 3
)

type data struct {
	changeChan chan comData
	data       interface{}
}
type newMapData struct {
	Flag       int
	key        interface{}
	value      interface{}
	ChangeChan chan comData
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
	sync.RWMutex
	data         map[interface{}]data
	mapChan      chan newMapData
	callbackChan *chan interface{}
	len          uint
}

func (slf *MapTtl) Init(callbackChan *chan interface{}) {
	slf.data = make(map[interface{}]data)
	slf.mapChan = make(chan newMapData, 1000)
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
				}
				if value, ok := slf.data[v.key]; ok {
					value.changeChan <- comData{
						Flag:              Reset,
						NewTtl:            v.Ttl,
						timeoutDeleteFlag: v.deleteFlag,
						value:             v.value,
					}
				} else {
					slf.len++
					slf.data[v.key] = data
					go slf.tll(v.key, v.value, v.Ttl, v.ChangeChan, v.deleteFlag)
				}
			} else if v.Flag == Del {
				if value, ok := slf.data[v.key]; ok {
					slf.len--
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
			}
		}
	}
}
func (slf *MapTtl) Set(key, value interface{}, ttl time.Duration, TimeOutDelete bool) {
	slf.mapChan <- newMapData{
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
	slf.mapChan <- newMapData{
		Flag: Del,
		key:  key,
	}

}
func (slf *MapTtl) Clear() {
	slf.mapChan <- newMapData{
		Flag: Clear,
	}
}
