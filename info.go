package panicsync

import (
	"fmt"
	"runtime"
)

const (
	MAX_STACKS   = 20
	STACK_OFFSET = 3
)

type Info struct {
	Error     interface{}
	Backtrace []*Caller
}

type Caller struct {
	File   string
	Line   int
	Method string
}

func NewInfo(err interface{}) *Info {
	info := &Info{
		Error: err,
	}
	info.buildBacktrace()
	return info
}

func (self *Info) buildBacktrace() {
	stack := make([]uintptr, MAX_STACKS)
	length := runtime.Callers(STACK_OFFSET, stack[:])

	backtrace := make([]*Caller, length)
	record := false
	i := 0

	for _, pc := range stack[:length] {
		f := runtime.FuncForPC(pc)
		if f == nil {
			continue
		}

		if !record {
			if f.Name() == "runtime.gopanic" {
				record = true
			}
			continue
		}

		file, line := f.FileLine(pc)

		backtrace[i] = &Caller{
			File:   file,
			Line:   line,
			Method: f.Name(),
		}
		i++
	}

	self.Backtrace = backtrace[:i]
}

func (self *Info) Print() {
	fmt.Printf("panic: %v [panicsync]\n", self.Error)
	for _, c := range self.Backtrace {
		fmt.Printf("\t%s [%s:%d]\n", c.Method, c.File, c.Line)
	}
}
