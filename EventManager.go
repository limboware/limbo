package limbov1

import (
	"slices"
	"sync"

	errnov1 "github.com/rejchev/errno"
)

type ListenerFn func(string, any)

type ListenerEntry struct {
	Id  uint64
	Lis ListenerFn
}

type Event struct {
	Name string
	Data any
}

type EventManager struct {
	handlers        map[string][]ListenerEntry
	handlersCounter uint64
	postEvents      []Event

	rw sync.RWMutex
}

var m = EventManager{}

func Events() *EventManager {
	return &m
}

// PublishAsync implements [IManager].
func (x *EventManager) PublishAsync(n string, v any) {
	x.rw.Lock()
	x.postEvents = append(x.postEvents, Event{Name: n, Data: v})
	x.rw.Unlock()
}

func (x *EventManager) Init() errnov1.Code {
	x.handlers = map[string][]ListenerEntry{}
	x.handlersCounter = 0
	x.rw = sync.RWMutex{}
	x.postEvents = make([]Event, 0, 8)
	return errnov1.OK
}

func (x *EventManager) Destroy() {}

func (x *EventManager) PostPublisherRun() {
	if x.Count() == 0 {
		return
	}

	evs := x.postEvents
	x.postEvents = make([]Event, 0, 8)

	for _, v := range evs {
		x.publish(v)
	}
}

func (x *EventManager) Count() int {
	return len(x.postEvents)
}

func (x *EventManager) Subscribe(k string, v ListenerFn) uint64 {
	x.rw.Lock()
	id := x.handlersCounter

	if _, ok := x.handlers[k]; !ok {
		x.handlers[k] = make([]ListenerEntry, 0)
	}

	x.handlers[k] = append(x.handlers[k], ListenerEntry{
		Id:  id,
		Lis: v,
	})

	x.handlersCounter++
	x.rw.Unlock()

	return id
}

func (x *EventManager) Unsubscribe(id uint64) {
	x.rw.Lock()
	defer x.rw.Unlock()

	copy := x.handlers

	for k, v := range copy {
		x.handlers[k] = slices.DeleteFunc(v, func(e ListenerEntry) bool {
			return e.Id == id
		})
	}
}

func (x *EventManager) Publish(n string, v any) {
	x.publish(Event{Name: n, Data: v})
}

func (x *EventManager) publish(ev Event) {
	var h []ListenerFn

	x.rw.RLock()
	for _, y := range x.handlers[ev.Name] {
		h = append(h, y.Lis)
	}
	x.rw.RUnlock()

	for _, y := range h {
		y(ev.Name, ev.Data)
	}
}
