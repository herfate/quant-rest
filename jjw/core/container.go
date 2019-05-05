package core

import (
	"github.com/findthefirst/quant-rest/jjw/utils/log"
	"strconv"
)

var taskIns []Task

func GoTask(t Task) {
	taskIns = append(taskIns, t)
	go t.Run()
	log.WarningAndWrap(" a task is running ..." + t.GetDesc() + " \n ready task " + strconv.Itoa(len(taskIns)))
}


func IsExit() bool {
	readySize := len(taskIns)
	exitSize := 0
	for _, it := range taskIns {
		if it.IsExit() {
			exitSize++
		}
	}
	return readySize == exitSize
}

type Task interface {
	Run()
	GetDesc() string
	IsExit() bool
}
