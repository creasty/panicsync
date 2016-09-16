panicsync
=========

Synchronize panic-recovery mechanism with multiple goroutines.


Example
-------

```go
package main

import (
  "fmt"
  "time"

  "github.com/creasty/panicsync"
)

func main() {
	ph := panicsync.NewHandler(func(info panicsync.Info) {
		info.Print()  // Handle error: just print and ignore
	})

	go func() {
		defer ph.Sync()
		panic("Error")
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("DONE")
}
```

```
panic: Error [panicsync]
goroutine 18 [running]:
github.com/creasty/panicsync.(*Handler).Sync(0xc82006c000)
        /Users/ykiwng/go/src/github.com/creasty/panicsync/handler.go:62 +0xad
panic(0xbc480, 0xc8200a4020)
        /Users/ykiwng/.anyenv/envs/goenv/versions/1.6/src/runtime/panic.go:426 +0x4e9
main.main.func2(0xc82006c000)
        /Users/ykiwng/go/src/github.com/creasty/panicsync/tmp/main.go:19 +0x90
created by main.main
        /Users/ykiwng/go/src/github.com/creasty/panicsync/tmp/main.go:20 +0x56

DONE
```


Advanced usage
--------------

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/creasty/panicsync"
)

func main() {
	ph := panicsync.NewHandler(func(info panicsync.Info) {
		fmt.Println(info.Error)
	})

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer ph.Sync()
			defer wg.Done()

			c := make(chan bool)
			go func(i int) {
				defer ph.Sync()
				defer close(c)

				if i%3 == 0 {
					panic("Error from goroutine in goroutine")
				}
			}(i)

			if i%2 == 0 {
				panic("Error from goroutine")
			}

			<-c
		}(i)
	}

	wg.Wait()

	time.Sleep(1 * time.Second)
	fmt.Println("DONE")
}
```

```
Error from goroutine
Error from goroutine
Error from goroutine in goroutine
Error from goroutine in goroutine
Error from goroutine
DONE
```
