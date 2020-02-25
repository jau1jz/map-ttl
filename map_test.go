package map_ttl

import (
	"fmt"
	"testing"
)

var testmap MapTtl
var c chan interface{}

func init() {
	c = make(chan interface{})
	testmap = MapTtl{}
	testmap.Init(&c)
	//time.Sleep(time.Second)
}

//func TestMap_tll_Set_FPS(t *testing.T) {
//	Count := 1000000
//	c := make(chan int )
//	go func() {
//		for i := 0; i < Count; i++ {
//			testmap.Set(i, i, time.Second * time.Duration(rand.Int31n(100)) , true)
//		}
//		c <- 1
//	}()
//
//	<-c
//	//for {
//	//	select {
//	//		case v := <- c :
//	//			fmt.Println(v.(int))
//	//	}
//	//}
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
	testmap.Set("jiang1", "zhu1", 0, true)
	testmap.Set("jiang2", "zhu2", 0, true)
	testmap.Set("jiang3", "zhu3", 0, true)

	println("Test set End")
}
func TestMap_tll_Get(t *testing.T) {

	println("Test Get start")
	for {
		obj := testmap.Get("jiang")
		if obj != nil {
			println(obj.(string))
			break
		} else {
			println("Get Empty")
		}
	}
	println("Test Get End")
}
func TestMap_tll_Del(t *testing.T) {
	println("Test Del Start")

	testmap.Del("jiang")

	println("Test Del End")
}

func TestMap_tll_Len(t *testing.T) {
	println("Test Len Start")

	fmt.Printf("%d \n", testmap.Len())

	println("Test Len End")
}

func TestMap_tll_Clear(t *testing.T) {
	println("Test Clear Start")

	testmap.Clear()
	println("Test Clear End")
}

//func TestMap_tll_UnsafeSetData(t *testing.T) {
//	println("unsafesetdata")
//	testmap.Set("jiang1", "zhu1", 0, true)
//	testmap.Set("jiang2", "zhu2", 0, true)
//	testmap.Set("jiang3", "zhu3", 0, true)
//	testmap.Range(func(key interface{}, value interface{}) {
//		testmap.UnsafeSetData(key, 111)
//	})
//	fmt.Printf("%+v \n ", testmap.value)
//}
