package march

import (
	"fmt"
	"sync"
)

func commu() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for v := range ch {
			println(v)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("over")
}
