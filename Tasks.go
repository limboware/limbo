package limbov1

import (
	"context"
	"sync"
	"time"

	errnov1 "github.com/rejchev/errno"
)

type Task_t struct {
	Id         Task
	TTL        time.Duration
	Worker     Worker
	CancelFn   context.CancelFunc
	DoneAt     int64
	ExecutedAt int64
}

type TaskExecuteFn = func(context.Context, Worker, Task) errnov1.Code

type TaskManager struct {
	tasks      []Task
	ttl        []time.Duration
	execFn     []TaskExecuteFn
	worker     []Worker
	cancel     []context.CancelFunc
	code       []errnov1.Code
	doneAt     []int64
	executedAt []int64

	gen  []uint8
	free [][3]byte

	router map[Task]int
	rwm    sync.RWMutex
}

var tasks = TaskManager{}

func Tasks() *TaskManager {
	return &tasks
}

func (x *TaskManager) Init() errnov1.Code {
	x.rwm = sync.RWMutex{}
	x.router = map[Task]int{}

	return errnov1.OK
}

func (x *TaskManager) route(v Task) int {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	if inner, ok := x.router[v]; ok {
		return inner
	}

	return -1
}

func (x *TaskManager) CancelFn(v Task) context.CancelFunc {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	return x.cancel[v.Id()]
}

func (x *TaskManager) Worker(v Task) Worker {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	return x.worker[v.Id()]
}

func (x *TaskManager) DoneAt(v Task) int64 {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	return x.doneAt[v.Id()]
}

func (x *TaskManager) Code(v Task) errnov1.Code {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	return x.code[v.Id()]
}

func (x *TaskManager) ExecutedAt(v Task) int64 {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	return x.executedAt[v.Id()]
}

func (x *TaskManager) TTL(v Task) time.Duration {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	return x.ttl[v.Id()]
}

func (x *TaskManager) IsAlive(v Task) bool {
	x.rwm.RLock()
	defer x.rwm.RUnlock()
	return x.gen[v.Id()] == v.Gen()
}

func (x *TaskManager) SetDone(w Worker, t Task, code errnov1.Code, at int64) {
	x.rwm.Lock()
	if x.gen[t.Id()] == t.Gen() && x.executedAt[t.Id()] != 0 {
		x.doneAt[t.Id()] = at
		x.code[t.Id()] = code
	}
	x.rwm.Unlock()
}

func (x *TaskManager) SetExecutedAt(w Worker, t Task, v int64) {
	x.rwm.Lock()
	if x.gen[t.Id()] == t.Gen() && x.worker[t.Id()] != w {
		x.executedAt[t.Id()] = v
	}
	x.rwm.Unlock()
}

func (x *TaskManager) SetCancelFn(w Worker, t Task, v context.CancelFunc) {
	x.rwm.Lock()
	if x.gen[t.Id()] == t.Gen() && x.worker[t.Id()] != w {
		x.cancel[t.Id()] = v
	}
	x.rwm.Unlock()
}

func (x *TaskManager) Execute(w Worker, t Task, cancelFn context.CancelFunc) TaskExecuteFn {
	x.rwm.Lock()
	defer x.rwm.Unlock()

	if x.gen[t.Id()] == t.Gen() && x.worker[t.Id()] == w && x.executedAt[t.Id()] == 0 {
		x.executedAt[t.Id()] = time.Now().Unix()
		x.cancel[t.Id()] = cancelFn
		return x.execFn[t.Id()]
	}

	return nil
}

func (x *TaskManager) Get(w Worker, buff *Task) bool {
	if buff == nil {
		return false
	}

	iter := x.Iterator()

	if iter == nil {
		return false
	}

	if !iter.First(buff, func(t Task) bool {
		return Tasks().ExecutedAt(t) == 0
	}) {
		return false
	}

	x.rwm.Lock()
	defer x.rwm.Unlock()

	if x.gen[(*buff).Id()] != (*buff).Gen() {
		return false
	}

	x.worker[(*buff).Id()] = w

	return true
}

func (x *TaskManager) Iterator() *Iterator[Task] {
	x.rwm.RLock()
	l := len(x.tasks)
	x.rwm.RUnlock()

	if l == 0 {
		return nil
	}

	buff := make([]Task, l)

	x.rwm.RLock()
	copy(buff, x.tasks)
	x.rwm.RUnlock()

	return NewIterator(buff)
}

func (x *TaskManager) IsInWork(v Task) bool {
	return x.ExecutedAt(v) != 0 && x.DoneAt(v) == 0
}

func (x *TaskManager) ExecDuration(v Task) int64 {
	if !x.IsAlive(v) || !x.IsInWork(v) {
		return 0
	}

	return time.Now().Unix() - x.ExecutedAt(v)
}

func (x *TaskManager) Remove(v Task) {
	if !x.IsAlive(v) || x.IsInWork(v) {
		return
	}

}
