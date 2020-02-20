package map_ttl

import (
	"testing"
	"time"
)

var testmap MapTtl
var c chan interface{}

func init() {
	c = make(chan interface{})
	testmap = MapTtl{}
	testmap.Init(&c)
	time.Sleep(time.Second)
}

//func TestMap_tll_Set_FPS(t *testing.T) {
//	Count := 1000
//
//	go func() {
//		for i := 0; i < Count; i++ {
//			testmap.Set(i, i, time.Second * time.Duration(rand.Int31n(100)) , true)
//		}
//	}()
//
//	for {
//		select {
//			case v := <- c :
//				fmt.Println(v.(int))
//		}
//	}
//
//}

//func TestMap_tll_Set_callback(t *testing.T) {
//	println("callback")
//	testmap.Set("jiang", "zhu", time.Second*1, false)
//	timeout := time.Tick(time.Second * 5)
//FOR:
//	for {
//		select {
//		case v := <-c:
//			println(v.(string))
//		case <-timeout:
//			println("recv timeout")
//			break FOR
//		}
//	}
//}
func TestMap_tll_Set(t *testing.T) {

	println("Test Set start")
	testmap.Set("jiang", "zhu", 0, true)
	testmap.Set("jiang", "zhu1", 0, true)
	testmap.Set("jiang1", "zhu1", 0, true)

	time.Sleep(time.Second * 10)

	println("Test set end")
}

//func TestMap_tll_Get(t *testing.T) {
//	println("Test Get start")
//	if obj := testmap.Get("jiang"); obj != nil {
//		println(obj.(string))
//	}
//}

func TestMap_tll_Del(t *testing.T) {
	println("Test Del Start")

	testmap.Del("jiang")
	time.Sleep(time.Second)
	println("Test Del end")
}

//
//func TestMap_tll_Clear(t *testing.T) {
//	println("clear")
//	testmap.Set("jiang1", "zhu1", 0, true)
//	testmap.Set("jiang2", "zhu2", 0, true)
//	testmap.Set("jiang3", "zhu3", 0, true)
//
//	fmt.Printf("%+v \n ", testmap.data)
//	testmap.Clear()
//	fmt.Printf("%+v \n ", testmap.data)
//}
//func TestMap_tll_Len(t *testing.T) {
//	println("len")
//	testmap.Set("jiang1", "zhu1", 0, true)
//	testmap.Set("jiang2", "zhu2", 0, true)
//	testmap.Set("jiang3", "zhu3", 0, true)
//
//	fmt.Printf("%d \n ", testmap.Len())
//}
//

//func TestMap_tll_UnsafeSetData(t *testing.T) {
//	println("unsafesetdata")
//	testmap.Set("jiang1", "zhu1", 0, true)
//	testmap.Set("jiang2", "zhu2", 0, true)
//	testmap.Set("jiang3", "zhu3", 0, true)
//	testmap.Range(func(key interface{}, value interface{}) {
//		testmap.UnsafeSetData(key, 111)
//	})
//	fmt.Printf("%+v \n ", testmap.data)
//}
