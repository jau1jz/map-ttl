package map_ttl

import (
	"fmt"
	"testing"
	"time"
)

var testmap MapTtl

func init() {
	testmap = MapTtl{}
	testmap.Init()

}
func TestMap_tll_Set_callback(t *testing.T) {
	println("callback")
	c := make(chan interface{})
	testmap.SetCallback(&c)
	testmap.Set("jiang", "zhu", time.Second*1, false)
	timeout := time.Tick(time.Second * 5)
FOR:
	for {
		select {
		case v := <-c:
			println(v.(string))
		case <-timeout:
			println("recv timeout")
			break FOR
		}
	}
}
func TestMap_tll_Set(t *testing.T) {

	println("Set")
	testmap.Set("jiang", "zhu", time.Minute, true)
	testmap.Set("jiang", "zhu1", time.Minute, true)
	testmap.Set("jiang1", "zhu1", time.Minute, true)

	fmt.Printf("%+v \n ", testmap.data)
}

func TestMap_tll_Get(t *testing.T) {
	println("Get")
	testmap.Set("jiang", "zhu", time.Minute, true)
	if obj := testmap.Get("jiang"); obj != nil {
		println(obj.(string))
	}
}

func TestMap_tll_Del(t *testing.T) {
	println("Del")
	testmap.Set("jiang1", "zhu1", 0, true)
	testmap.Set("jiang2", "zhu2", 0, true)
	testmap.Set("jiang3", "zhu3", 0, true)
	fmt.Printf("%+v \n ", testmap.data)
	testmap.Del("jiang3")
	time.Sleep(time.Second)
	fmt.Printf("%+v \n ", testmap.data)
	testmap.SetData("jiang2", "zhu4")
	fmt.Printf("%+v \n ", testmap.data)
}

func TestMap_tll_Clear(t *testing.T) {
	println("clear")
	testmap.Set("jiang1", "zhu1", 0, true)
	testmap.Set("jiang2", "zhu2", 0, true)
	testmap.Set("jiang3", "zhu3", 0, true)

	fmt.Printf("%+v \n ", testmap.data)
	testmap.Clear()
	fmt.Printf("%+v \n ", testmap.data)
}
func TestMap_tll_Len(t *testing.T) {
	println("len")
	testmap.Set("jiang1", "zhu1", 0, true)
	testmap.Set("jiang2", "zhu2", 0, true)
	testmap.Set("jiang3", "zhu3", 0, true)

	fmt.Printf("%d \n ", testmap.Len())
}

func TestMap_tll_Range(t *testing.T) {
	println("range")
	testmap.Set("jiang1", "zhu1", 0, true)
	testmap.Set("jiang2", "zhu2", 0, true)
	testmap.Set("jiang3", "zhu3", 0, true)
	testmap.Range(func(key interface{}, value interface{}) {
		fmt.Printf("%s %s \n", key, value)
	})
}
func TestMap_tll_UnsafeSetData(t *testing.T) {
	println("unsafesetdata")
	testmap.Set("jiang1", "zhu1", 0, true)
	testmap.Set("jiang2", "zhu2", 0, true)
	testmap.Set("jiang3", "zhu3", 0, true)
	testmap.Range(func(key interface{}, value interface{}) {
		testmap.UnsafeSetData(key, 111)
	})
	fmt.Printf("%+v \n ", testmap.data)
}
