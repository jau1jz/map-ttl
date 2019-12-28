# map-ttl
golang ttl map like redis
```
func main() {
	tmap := map_ttl.Map_ttl{}
	tmap.Init()
	c := make(chan interface{})
	tmap.Set_callback(&c) //set callback chan this chan will be receive data when timeout , but if you del key by func Del() this chan don't work
	tmap.Set("jone","taylar",time.Second * 60) // 60s ttl
	tmap.Set("jack","ma",0) //no ttl
	tmap.Set("jack","ma",time.Second * 60) //reset ttl to 60s
	tmap.SetData("jack","li") //if don't have key "jack" func will be return false
	tmap.Get("jone") //get value
	tmap.Del("jone")
	tmap.Clear()//clear map
	tmap.Range(func(key interface{}, value interface{}) {
        //in func Range just can use unsafe func
        tmap.UnsafeSetData("jack","Hu")
        tmap.UnsafeSetData("jack","Hu")
	})
	go func() {
		for {
			select {
			case <- c :
				//timeout do something

			}
		}
	}()
}
```
