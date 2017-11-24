package utils

import (
	"fmt"
	"runtime"
	"time"
)

//TimeFormat 时间格式化
func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FileLine 调用处的文件名和行号
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

type MutexGroup struct {
	ch chan bool
}

func NewMutexGroup(count int) *MutexGroup {
	if count <= 0 {
		return nil
	}
	return &MutexGroup{ch: make(chan bool, count)}
}

func (s *MutexGroup) Lock() {
	s.ch <- true
}

func (s *MutexGroup) UnLock() {
	<-s.ch
}
