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
	ph := panicsync.NewHandler(func(info *panicsync.Info) {
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
        main.main.func2 [/Users/ykiwng/go/src/github.com/creasty/panicsync/tmp/main.go:17]
        runtime.goexit [/Users/ykiwng/.anyenv/envs/goenv/versions/1.6/src/runtime/asm_amd64.s:1999]
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
	ph := panicsync.NewHandler(func(info *panicsync.Info) {
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
