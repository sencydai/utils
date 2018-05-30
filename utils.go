package utils

import (
	"fmt"
	"runtime"
)

//FileLine 调用处的文件名和行号
func FileLine(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	i, count := len(file)-4, 0
	for ; i > 0; i-- {
		if file[i] == '/' {
			count++
			if count == 2 {
				break
			}
		}
	}
	return fmt.Sprintf("%s:%d", file[i+1:], line)
}

type Semaphore struct {
	ch chan bool
}

func NewSemaphore(count int) *Semaphore {
	if count <= 0 {
		return nil
	}
	return &Semaphore{ch: make(chan bool, count)}
}

func (s *Semaphore) Require() {
	s.ch <- true
}

func (s *Semaphore) Release() {
	<-s.ch
}
