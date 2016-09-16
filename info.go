package panicsync

import (
	"fmt"
)

type Info struct {
	Error      interface{}
	StackTrace string
}

func (self Info) Print() {
	fmt.Printf("panic: %v [panicsync]\n", self.Error)
	fmt.Println(self.StackTrace)
}
