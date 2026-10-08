package limbov1

import (
	"context"
	"math"
	"time"
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

type Task_t struct {
	ID        Task
	Instance  unsafe.Pointer
	ExecuteFn ExecuteFn
	DoneFn    DoneFn
	TTL       time.Duration
	Context   context.Context
	CancelFn  context.CancelFunc
}

type DoneFn = func() bool

type ExecuteFn = func(ctx context.Context, dt time.Duration)

type TaskManager struct {
	tasks []Task

	instances  []unsafe.Pointer
	doneFns    []DoneFn
	executeFns []ExecuteFn
	ttls       []time.Duration

	gen  []uint8
	free []uint32

	router map[Task]int
}

var tasks = TaskManager{}

func Tasks() *TaskManager {
	return &tasks
}

func (x *TaskManager) Init() errnov1.Code {
	x.tasks = make([]Task, 0, 16)

	x.gen = make([]uint8, 0, 16)
	x.free = make([]uint32, 0, 16)
	x.router = map[Task]int{}

	x.instances = make([]unsafe.Pointer, 0, 16)
	x.executeFns = make([]ExecuteFn, 0, 16)
	x.doneFns = make([]DoneFn, 0, 16)
	x.ttls = make([]time.Duration, 0, 16)

	return errnov1.OK
}

func (x *TaskManager) route(v Task) int {
	if inner, ok := x.router[v]; ok {
		return inner
	}

	return -1
}

func (x *TaskManager) Instance(v Task) unsafe.Pointer {
	return x.instances[v.Id()]
}

func (x *TaskManager) ExecuteFn(v Task) ExecuteFn {
	return x.executeFns[v.Id()]
}

func (x *TaskManager) DoneFn(v Task) DoneFn {
	return x.doneFns[v.Id()]
}

func (x *TaskManager) TTL(v Task) time.Duration {
	return x.ttls[v.Id()]
}

func (x *TaskManager) Struct(v Task) *Task_t {
	return &Task_t{
		ID:        v,
		Instance:  x.Instance(v),
		ExecuteFn: x.ExecuteFn(v),
		DoneFn:    x.DoneFn(v),
		TTL:       x.TTL(v),
	}
}

func (x *TaskManager) Iterator() *Iterator[Task] {
	buff := make([]Task, len(x.tasks))

	copy(buff, x.tasks)

	return NewIterator(buff)
}

func (x *TaskManager) New(initFn func(*Task_t), buff *Task) bool {
	if initFn == nil || buff == nil {
		return false
	}

	pTask := new(Task_t)

	initFn(pTask)

	idx := uint32(0)

	if len(x.free) > 0 {
		idx = x.free[len(x.free)-1]
		x.free = x.free[:len(x.free)-1]
	} else {
		idx = uint32(len(x.gen))
		x.gen = append(x.gen, 0)
	}

	if len(x.instances) <= int(idx) {
		x.instances = append(x.instances, nil)
	}

	if len(x.executeFns) <= int(idx) {
		x.executeFns = append(x.executeFns, nil)
	}

	if len(x.doneFns) <= int(idx) {
		x.doneFns = append(x.doneFns, nil)
	}

	if len(x.ttls) <= int(idx) {
		x.ttls = append(x.ttls, 0)
	}

	x.instances[idx] = pTask.Instance
	x.executeFns[idx] = pTask.ExecuteFn
	x.doneFns[idx] = pTask.DoneFn
	x.ttls[idx] = pTask.TTL

	*buff = MakeTask(idx, x.gen[idx])

	innerIdx := len(x.tasks)

	x.tasks = append(x.tasks, *buff)

	x.router[*buff] = innerIdx

	Events().Publish("tasks.new", *buff)

	return true
}

func (x *TaskManager) IsAlive(v Task) bool {
	return x.route(v) != -1 && x.gen[v.Id()] == v.Gen()
}

func (x *TaskManager) Remove(v Task) {
	if !x.IsAlive(v) {
		return
	}

	Events().Publish("tasks.remove", v)

	idx := v.Id()
	if x.gen[idx]+1 == math.MaxUint8 {
		x.gen[idx] = 0
	}

	x.gen[idx]++

	x.instances[idx] = nil
	x.executeFns[idx] = nil
	x.doneFns[idx] = nil
	x.ttls[idx] = 0

	l := len(x.tasks)

	if innerIdx := x.route(v); l > 1 && l-1 != innerIdx {
		x.tasks[innerIdx] = x.tasks[l-1]
		x.router[x.tasks[l-1]] = innerIdx
	}

	delete(x.router, v)

	x.tasks = x.tasks[:l-1]

	Events().Publish("tasks.removed", v)
}
