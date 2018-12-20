package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type CleanHandler func(*log.Logger)
type CleanUpHandlerArray []CleanHandler

func (c *CleanUpHandlerArray) Add(handler CleanHandler) {
	*c = append(*c, handler)
}
func (c *CleanUpHandlerArray) Wait() {
	select {
	case <-func() <-chan int {
		ch := make(chan int)
		go func() {
			for {
				if len(*c) == 0 {
					ch <- 0
				}
			}
		}()
		return ch
	}():
		return
	}
}
func NewCleanUpHandlerArray(l *log.Logger) *CleanUpHandlerArray {
	obj := &CleanUpHandlerArray{}

	go func() {
		ch := make(chan os.Signal)
		signal.Ignore(syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGQUIT)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGQUIT)
		sig := <-ch
		l.Println("catch signal :", sig, ", exit")
		l.Println("cleaning up begin,number of function to call :", len(*obj))

		for _, v := range *obj {
			v(l)
		}

		l.Println("cleaning up finished")
		for {
			*obj = make(CleanUpHandlerArray, 0)
		}
	}()
	return obj
}
