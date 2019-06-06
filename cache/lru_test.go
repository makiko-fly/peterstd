package cache

import (
	"sync"
	"testing"
)

type Foobar struct {
	Foo  string
	Bar  string
	Hits int64
}

func Test_LRUCacheRace(t *testing.T) {
	c := NewLRUCache()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Set("aaa", Foobar{
					Foo: "foo",
					Bar: "bar",
				})
			}
		}()

		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				v, _ := c.Get("aaa")
				a := v.(Foobar)
				a.Bar = "bar2"
			}
		}()
	}
	wg.Wait()
}
