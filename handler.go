package panicsync

import (
	"runtime"
	"sync"
)

type HandlerFunc func(Info)

type Handler struct {
	info    chan Info
	quit    chan bool
	handler HandlerFunc
	lock    *sync.Mutex
}

func NewHandler(fn HandlerFunc) *Handler {
	h := &Handler{
		info:    make(chan Info),
		handler: fn,
		lock:    &sync.Mutex{},
		quit:    make(chan bool),
	}
	go h.listen()
	return h
}

func (self *Handler) listen() {
	for info := range self.info {
		self.handle(info)
	}
}

func (self *Handler) handle(info Info) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.handler(info)
}

func (self *Handler) Done() {
	select {
	case info, ok := <-self.info:
		if ok {
			close(self.quit)
			close(self.info)
			self.lock.Lock()
			defer self.lock.Unlock()
			self.handle(info)
		}
	default:
		close(self.info)
	}
}

func (self *Handler) Sync() {
	err := recover()
	if err == nil {
		return
	}

	buf := make([]byte, 10000)
	traceSize := runtime.Stack(buf, false)
	buf = buf[:traceSize]
	info := Info{
		Error:      err,
		StackTrace: string(buf),
	}

	select {
	case self.info <- info:
	case <-self.quit:
	}
}
