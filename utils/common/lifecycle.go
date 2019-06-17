package common

import (
	"errors"
	"log"
	"sync"
)

type Lifecycle interface {
	Init() error
	Start() error
	Stop() error
	Destroy() error
	CurrentState() string
	ListLifecycleListeners() []LifecycleListener
	RemoveLifecycleListener(listener LifecycleListener)
	AddLifecycleListener(listener LifecycleListener)
}

type LifecycleBase struct {
	log                *log.Logger
	appid              string
	state              string
	lock               *sync.Mutex
	lifecycleListeners []LifecycleListener
	preStateTable      map[string][]string //for example ,current state is STOPED, this pre state is INITIALIZED or STARTED , so map["STARTED"] = {"INITIALIZED","STARTED"}
	wg                 *sync.WaitGroup
}

func (this *LifecycleBase) ListLifecycleListeners() []LifecycleListener {
	//return a copy
	return []LifecycleListener(this.lifecycleListeners)
}

func (this *LifecycleBase) RemoveLifecycleListener(listener LifecycleListener) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for idx, lis := range this.lifecycleListeners {
		if lis == listener {
			this.lifecycleListeners = append(this.lifecycleListeners[0:idx], this.lifecycleListeners[idx:]...)
			return
		}
	}
}

func (this *LifecycleBase) AddLifecycleListener(listener LifecycleListener) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.lifecycleListeners = append(this.lifecycleListeners, listener)
}

func (this *LifecycleBase) Init() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	var (
		err   error
		event *LifecycleEvent
	)
	this.log.Println("---------------------init----------------------- ")
	if !this.checkState(INITIALIZED) {
		return errors.New("want to change state to +" + INITIALIZED + " , but current state " + this.state + " not be allow ")
	}
	this.state = INITIALIZING
	event = NewLifecycleEvent(this, BEFORE_INIT_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	event = NewLifecycleEvent(this, AFTER_INIT_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	this.state = INITIALIZED
	return nil
}

func (this *LifecycleBase) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	var (
		err   error
		event *LifecycleEvent
	)
	this.log.Println("---------------------STARTED----------------------- ")
	if !this.checkState(STARTED) {
		return errors.New("want to change state to" + STARTED + " , but current state " + this.state + " not be allow ")
	}
	this.state = STARTING
	event = NewLifecycleEvent(this, BEFORE_START_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	event = NewLifecycleEvent(this, AFTER_START_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	this.state = STARTED
	return nil
}

func (this *LifecycleBase) Stop() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	var (
		err   error
		event *LifecycleEvent
	)
	this.log.Println("---------------------init----------------------- ")
	if !this.checkState(STOPPED) {
		return errors.New("want to change state to" + STOPPED + " , but current state " + this.state + " not be allow ")
	}
	this.state = STOPPING
	event = NewLifecycleEvent(this, BEFORE_STOP_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	event = NewLifecycleEvent(this, AFTER_STOP_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	this.state = STOPPED
	return nil
}

func (this *LifecycleBase) Destroy() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	var (
		err   error
		event *LifecycleEvent
	)
	this.log.Println("---------------------init----------------------- ")
	if !this.checkState(DESTROYED) {
		return errors.New("want to change state to" + DESTROYED + " , but current state " + this.state + " not be allow ")
	}
	this.state = DESTROYING
	event = NewLifecycleEvent(this, BEFORE_DESTROY_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	event = NewLifecycleEvent(this, AFTER_DESTROY_EVENT)
	if err = this.fireEventAndWaitDone(event); err != nil {
		this.log.Println(err)
		return err
	}
	this.state = DESTROYED
	return nil
}

func (this *LifecycleBase) CurrentState() string {
	return this.state
}

func (this *LifecycleBase) fireEvent(event *LifecycleEvent) {
	for _, listener := range this.lifecycleListeners {
		listener.AcceptEvent(event)
	}
}

func (this *LifecycleBase) fireEventAndWaitDone(event *LifecycleEvent) error {
	this.wg.Add(len(this.lifecycleListeners))
	this.fireEvent(event)
	this.wg.Wait()
	return nil
}

func (this *LifecycleBase) checkState(wantTo string) bool {
	pre := this.preStateTable[wantTo]
	for _, state := range pre {
		if state == wantTo {
			return true
		}
	}
	return false
}

type LifecycleListener interface {
	OnEvent(event *LifecycleEvent)
	AcceptEvent(event *LifecycleEvent)
}

type LifecycleListenerBase struct {
	EventChan chan *LifecycleEvent
}

func (this *LifecycleListenerBase) OnEvent(event *LifecycleEvent) {
	//need sub class override
	//insert some code
	event.Done()
}

func (this *LifecycleListenerBase) AcceptEvent(event *LifecycleEvent) {
	this.EventChan <- event
}

//lifecycle event
type LifecycleEvent struct {
	Host      *LifecycleBase
	EventType string
}

func NewLifecycleEvent(host *LifecycleBase, eventType string) *LifecycleEvent {
	//check host type require LifecycleBase
	return &LifecycleEvent{host, eventType}
}

func (this *LifecycleEvent) Done() {
	//notify host the event has done
	this.Host.wg.Done()
}

//event type
const (
	BEFORE_INIT_EVENT     = "before_init"
	AFTER_INIT_EVENT      = "after_init"
	START_EVENT           = "start"
	BEFORE_START_EVENT    = "before_start"
	AFTER_START_EVENT     = "after_start"
	STOP_EVENT            = "stop"
	BEFORE_STOP_EVENT     = "before_stop"
	AFTER_STOP_EVENT      = "after_stop"
	AFTER_DESTROY_EVENT   = "after_destroy"
	BEFORE_DESTROY_EVENT  = "before_destroy"
	CONFIGURE_START_EVENT = "configure_start"
	CONFIGURE_STOP_EVENT  = "configure_stop"
)

//life cycle state
const (
	NEW          = "NEW"
	INITIALIZING = "INITIALIZING"
	INITIALIZED  = "INITIALIZED"
	STARTING     = "STARTING"
	STARTED      = "STARTED"
	SHUTDOWN     = "SHUTDOWN"
	STOPPING     = "STOPPING"
	STOPPED      = "STOPPED"
	DESTROYING   = "DESTROYING"
	DESTROYED    = "DESTROYED"
	FAILED       = "FAILED" //means init failed generally
)

// life cycle failed error
/*
const(
	INIT_FAILED_ERR=errors.New("init failed err")
	START_FAILED_ERR=errors.New("start failed err")
	STOP_FAILED_ERR=errors.New("stop failed err")
	DESTROY_FAILED_ERR=errors.New("destory failed err")
)
*/
