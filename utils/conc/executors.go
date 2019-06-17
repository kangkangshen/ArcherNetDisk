package conc

import (
	"github.com/kangkangshen/ArcherNetDisk/config"
	"sync"
)

//java style executors framwork

type Executor interface {
	Execute(task func())
}

type Task interface {
	CallBack() interface{}
}
type ExecutorService interface {
	Executor
	Submit(task func() interface{}) Future
	SubmitTask(task *Task) Future
	ShutDown()
	IsShutDown()
	State() int
}

type baseExecutorService struct {
	lock           sync.Mutex
	state          int
	errHandler     map[string]func(err error)
	shutdownChan   chan string //关闭信号
	tasks          int64       //当前接收到的任务数
	completedTasks int64       //已完成任务数
}

func (es *baseExecutorService) Execute(task func()) {
	go task()
}
func (es *baseExecutorService) Submit(task func() interface{}) *Future {
	go func() {
		es.tasks++
		es.completedTasks++
	}()
	return nil
}
func (es *baseExecutorService) ShutDown() {
	es.lock.Lock()
	es.shutdownChan <- "SHUTDOWN"
	es.state = config.SHUTDOWN
	es.lock.Unlock()
}
func (es *baseExecutorService) IsShutDown() bool {
	return es.state == config.SHUTDOWN
}

//func (es *baseExecutorService){}
