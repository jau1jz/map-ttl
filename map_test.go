package map_ttl

import (
	"fmt"
	"testing"
	"time"
)

var test_map Map_ttl

func init() {
	test_map.Init()
}
func TestMap_tll_Set(t *testing.T) {
	test_map.Set("jiang", "zhu", time.Minute)
	test_map.Set("jiang", "zhu1", time.Minute)
	test_map.Set("jiang1", "zhu1", time.Minute)

	fmt.Printf("%+v \n ", test_map.data)
}
func TestMap_tll_Get(t *testing.T) {
	test_map.Set("jiang", "zhu", time.Minute)
	if obj := test_map.Get("jiang"); obj != nil {
		println(obj.(string))
	}
}
func TestMap_tll_Del(t *testing.T) {
	test_map.Set("jiang1", "zhu1", 0)
	test_map.Set("jiang2", "zhu2", 0)
	test_map.Set("jiang3", "zhu3", 0)
	fmt.Printf("%+v \n ", test_map.data)
	test_map.Del("jiang3")
	time.Sleep(time.Second)
	fmt.Printf("%+v \n ", test_map.data)
	test_map.SetData("jiang2", "zhu4")
	fmt.Printf("%+v \n ", test_map.data)
}
func TestMap_tll_Set_callback(t *testing.T) {
	c := make(chan interface{})
	test_map.Set_callback(&c)
	test_map.Set("jiang", "zhu", time.Second*5)
	v := <-c
	println(v.(string))
}
func TestMap_tll_Clear(t *testing.T) {
	test_map.Set("jiang1", "zhu1", 0)
	test_map.Set("jiang2", "zhu2", 0)
	test_map.Set("jiang3", "zhu3", 0)

	fmt.Printf("%+v \n ", test_map.data)
	test_map.Clear()
	fmt.Printf("%+v \n ", test_map.data)
}
func TestMap_tll_Len(t *testing.T) {
	test_map.Set("jiang1", "zhu1", 0)
	test_map.Set("jiang2", "zhu2", 0)
	test_map.Set("jiang3", "zhu3", 0)

	fmt.Printf("%d \n ", test_map.Len())
}

func TestMap_tll_Range(t *testing.T) {
	test_map.Set("jiang1", "zhu1", 0)
	test_map.Set("jiang2", "zhu2", 0)
	test_map.Set("jiang3", "zhu3", 0)
	test_map.Range(func(key interface{}, value interface{}) {
		fmt.Printf("%s %s \n", key, value)
	})
}
