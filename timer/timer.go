package timer

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

var (
	ZeroTime = time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.Local)
)

func parseFunc(cb interface{}, args ...interface{}) (reflect.Value, []reflect.Value) {
	cbValue := reflect.ValueOf(cb)
	if cbValue.Kind() != reflect.Func {
		panic("cb must be a function")
	}
	values := make([]reflect.Value, len(args))
	for i, v := range args {
		values[i] = reflect.ValueOf(v)
	}
	return cbValue, values
}

func callback(cb reflect.Value, values []reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	cb.Call(values)
}

type Timer struct {
	t    *time.Timer
	stop chan bool
	mu   sync.Mutex
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.t.Stop() {
		t.stop <- true
	}
}

func After(delay time.Duration, cbFunc interface{}, args ...interface{}) *Timer {
	cb, ov := parseFunc(cbFunc, args...)
	t := &Timer{t: time.NewTimer(delay), stop: make(chan bool, 1)}

	go func() {
		select {
		case <-t.t.C:
			callback(cb, ov)
		case <-t.stop:
		}

		t.mu.Lock()
		t.t.Stop()
		t.mu.Unlock()
	}()

	return t
}

func Loop(delay, interval time.Duration, times int, cbFunc interface{}, args ...interface{}) *Timer {
	cb, ov := parseFunc(cbFunc, args...)
	t := &Timer{t: time.NewTimer(delay), stop: make(chan bool, 1)}

	go func() {
		var count int
	TAG_STOP_FOR:
		for {
			select {
			case <-t.t.C:
				if times > 0 {
					if count < times {
						count++
						t.t.Reset(interval)
					} else {
						break TAG_STOP_FOR
					}
				} else {
					t.t.Reset(interval)
				}

				go callback(cb, ov)
			case <-t.stop:
				break TAG_STOP_FOR
			}
		}

		t.mu.Lock()
		t.t.Stop()
		t.mu.Unlock()
	}()

	return t
}

func checkMomentHappened(t time.Time, moment string) bool {
	now := time.Now()
	trigger, _ := time.ParseInLocation("15:04:05", moment, time.Local)
	trigger = time.Date(t.Year(), t.Month(), t.Day(), trigger.Hour(), trigger.Minute(), trigger.Second(), 0, time.Local)
	if !trigger.Before(now) {
		return false
	}
	if trigger.After(t) {
		return true
	}
	trigger = trigger.AddDate(0, 0, 1)
	return trigger.Before(now)
}

func getMomentDelay(moment string) time.Duration {
	now := time.Now()
	trigger, _ := time.ParseInLocation("15:04:05", moment, time.Local)
	trigger = time.Date(now.Year(), now.Month(), now.Day(), trigger.Hour(), trigger.Minute(), trigger.Second(), 0, time.Local)
	if trigger.After(now) {
		return trigger.Sub(now)
	}

	return now.AddDate(0, 0, 1).Sub(trigger)
}

func DayLoop(moment string, cbFunc interface{}, args ...interface{}) *Timer {
	return Loop(getMomentDelay(moment), time.Hour*24, -1, cbFunc, args...)
}

func DayLoop2(last time.Time, moment string, cbFunc interface{}, args ...interface{}) *Timer {
	if checkMomentHappened(last, moment) {
		go callback(parseFunc(cbFunc, args...))
	}
	return Loop(getMomentDelay(moment), time.Hour*24, -1, cbFunc, args...)
}

func HourlyLoop(cbFunc interface{}, args ...interface{}) *Timer {
	now := time.Now()
	d := time.Minute*time.Duration(now.Minute()) + time.Second*time.Duration(now.Second()) + time.Nanosecond*time.Duration(now.Nanosecond())
	return Loop(time.Hour-d, time.Hour, -1, cbFunc, args...)
}

func MinuteLoop(cbFunc interface{}, args ...interface{}) *Timer {
	now := time.Now()
	d := time.Second*time.Duration(now.Second()) + time.Nanosecond*time.Duration(now.Nanosecond())
	return Loop(time.Minute-d, time.Minute, -1, cbFunc, args...)
}

func TenMinutesLoop(cbFunc interface{}, args ...interface{}) *Timer {
	now := time.Now()
	d := time.Minute*time.Duration(now.Minute()%10) + time.Second*time.Duration(now.Second()) + time.Nanosecond*time.Duration(now.Nanosecond())
	return Loop(time.Minute*10-d, time.Minute*10, -1, cbFunc, args...)
}

func HalfhourLoop(cbFunc interface{}, args ...interface{}) *Timer {
	now := time.Now()
	d := time.Minute*time.Duration(now.Minute()%30) + time.Second*time.Duration(now.Second()) + time.Nanosecond*time.Duration(now.Nanosecond())
	return Loop(time.Minute*30-d, time.Minute*30, -1, cbFunc, args...)
}
