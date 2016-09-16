package panicsync

import (
	"fmt"
	"sync"
	"testing"
)

func TestRootPanic(t *testing.T) {
	ph := NewHandler(func(info *Info) {
		fmt.Println(info.Error)
		info.Print()
	})
	defer ph.Sync()

	panic("Error")
}

func TestSinglePanic(t *testing.T) {
	quit := make(chan bool)

	ph := NewHandler(func(info *Info) {
		fmt.Println(info.Error)
		info.Print()
		close(quit)
	})

	go func() {
		defer ph.Sync()
		panic("Error")
	}()

	<-quit
}

func TestMultiplePanics(t *testing.T) {
	ph := NewHandler(func(info *Info) {
		fmt.Println(info.Error)
		info.Print()
	})

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer ph.Sync()
			defer wg.Done()

			panic("Error from goroutine")
		}()
	}

	wg.Wait()
}

func TestNestedPanics(t *testing.T) {
	ph := NewHandler(func(info *Info) {
		fmt.Println(info.Error)
		info.Print()
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
}
