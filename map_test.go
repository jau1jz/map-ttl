package map_ttl

import (
	"fmt"
	"testing"
	"time"
)

var testmap Map_ttl

func init() {
	testmap.Init()
}
func TestMap_tll_Set(t *testing.T) {
	testmap := Map_ttl{}
	testmap.Init()
	println("Set")
	testmap.Set("jiang", "zhu", time.Minute)
	testmap.Set("jiang", "zhu1", time.Minute)
	testmap.Set("jiang1", "zhu1", time.Minute)

	fmt.Printf("%+v \n ", testmap.ttl)
}
func TestMap_tll_Get(t *testing.T) {
	println("Get")
	testmap.Set("jiang", "zhu", time.Minute)
	if obj := testmap.Get("jiang"); obj != nil {
		println(obj.(string))
	}
}
func TestMap_tll_Del(t *testing.T) {
	println("Del")
	testmap.Set("jiang1", "zhu1", 0)
	testmap.Set("jiang2", "zhu2", 0)
	testmap.Set("jiang3", "zhu3", 0)
	fmt.Printf("%+v \n ", testmap.ttl)
	testmap.Del("jiang3")
	time.Sleep(time.Second)
	fmt.Printf("%+v \n ", testmap.ttl)
	testmap.SetData("jiang2", "zhu4")
	fmt.Printf("%+v \n ", testmap.ttl)
}
func TestMap_tll_Set_callback(t *testing.T) {
	println("callback")
	c := make(chan interface{})
	testmap.Set_callback(&c)
	testmap.Set("jiang", "zhu", 0)
	//v := <-c
	//println(v.(string))
}
func TestMap_tll_Clear(t *testing.T) {
	println("clear")
	testmap.Set("jiang1", "zhu1", 0)
	testmap.Set("jiang2", "zhu2", 0)
	testmap.Set("jiang3", "zhu3", 0)

	fmt.Printf("%+v \n ", testmap.ttl)
	testmap.Clear()
	fmt.Printf("%+v \n ", testmap.ttl)
}
func TestMap_tll_Len(t *testing.T) {
	println("len")
	testmap.Set("jiang1", "zhu1", 0)
	testmap.Set("jiang2", "zhu2", 0)
	testmap.Set("jiang3", "zhu3", 0)

	fmt.Printf("%d \n ", testmap.Len())
}

func TestMap_tll_Range(t *testing.T) {
	println("range")
	testmap.Set("jiang1", "zhu1", 0)
	testmap.Set("jiang2", "zhu2", 0)
	testmap.Set("jiang3", "zhu3", 0)
	testmap.Range(func(key interface{}, value interface{}) {
		fmt.Printf("%s %s \n", key, value)
	})
}
func TestMap_tll_UnsafeSetData(t *testing.T) {
	println("unsafesetdata")
	testmap.Set("jiang1", "zhu1", 0)
	testmap.Set("jiang2", "zhu2", 0)
	testmap.Set("jiang3", "zhu3", 0)
	testmap.Range(func(key interface{}, value interface{}) {
		testmap.UnsafeSetData(key, 111)
	})
	fmt.Printf("%+v \n ", testmap.ttl)
}
