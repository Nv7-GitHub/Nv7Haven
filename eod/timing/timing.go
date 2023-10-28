package timing

import (
	"os"
	"strconv"
	"sync"
	"time"
)

var lock = &sync.RWMutex{}
var timing = make(map[string]*os.File)

type Timer struct {
	file  *os.File
	start time.Time
}

func GetTimer(commandName string) *Timer {
	lock.RLock()
	file, ok := timing[commandName]
	lock.RUnlock()
	if !ok {
		var err error
		file, err = os.OpenFile("data/eod/timing/"+commandName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		lock.Lock()
		timing[commandName] = file
		lock.Unlock()
	}

	return &Timer{file: file, start: time.Now()}
}

func (t *Timer) Stop() {
	_, err := t.file.WriteString(strconv.Itoa(int(time.Since(t.start).Microseconds())) + "\n")
	if err != nil {
		panic(err)
	}
}

func Init() {
	err := os.MkdirAll("data/eod/timing", os.ModePerm)
	if err != nil {
		panic(err)
	}
}
