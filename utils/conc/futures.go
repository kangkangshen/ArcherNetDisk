package conc

import (
	"time"
)

type Future interface {
	Cancel()
	IsDone() bool
	Pause()
	Resume()
	Get() interface{}
	GetWithTimeOut(timeout time.Duration) (error, interface{})
	GetOrDefault(timeout time.Duration, defa interface{}) (error, interface{})
}
type baseFuture struct {
	state      int
	resultChan chan interface{}
}

func (f *baseFuture) Cancel() {
	//

}

func (f *baseFuture) Pause() {

}
func (f *baseFuture) Resume() {

}
func (f *baseFuture) IsDone() bool {
	return false
}
func (f *baseFuture) Get() interface{} {
	return <-f.resultChan
}
func (f *baseFuture) GetWithTimeOut(timeout time.Duration) (error, interface{}) {
	return nil, nil
}
func (f *baseFuture) GetOrDefault(timeout time.Duration, defa interface{}) (error, interface{}) {
	return nil, nil
}
func newFuture() Future {
	return new(baseFuture)

}
