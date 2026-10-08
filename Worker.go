package limbov1

import (
	"context"
	"sync"
	"time"

	errnov1 "github.com/rejchev/errno"
)

type Worker uint16

func MakeWorker(id, gen uint8) Worker {
	return Worker(uint16(id)<<8 | uint16(gen))
}

func (x Worker) Id() uint8 {
	return uint8(x >> 8)
}

func (x Worker) Gen() uint8 {
	return uint8(x)
}

func (x Worker) SetId(v uint8) Worker {
	return MakeWorker(v, x.Gen())
}

func (x Worker) SetGen(v uint8) Worker {
	return MakeWorker(x.Id(), v)
}

func (x Worker) Closer() uint16 {
	return uint16(x)
}

///////////////////////////

type Worker_t struct {
	ID   Worker
	Rate int

	alive bool

	tasks     []Task
	doneFn    []DoneFn
	executeFn []ExecuteFn
	ttl       []time.Duration
	ctx       []context.Context
	cancelFn  []context.CancelFunc

	router map[Task]int

	rw sync.RWMutex
}

func (x *Worker_t) Init() {
	x.alive = true

	x.tasks = make([]Task, 0, 16)
	x.doneFn = make([]DoneFn, 0, 16)
	x.executeFn = make([]ExecuteFn, 0, 16)
	x.cancelFn = make([]context.CancelFunc, 0, 16)
	x.ctx = make([]context.Context, 0, 16)

	x.router = map[Task]int{}

	x.rw = sync.RWMutex{}

	Events().Subscribe("jobs.new", x.onNew)
	Events().Subscribe("jobs.remove", x.onRemove)
	Events().Subscribe("jobs.cancel", x.onCancel)

	Events().Subscribe("workers.shutdown", x.onShutdown)
	// Events().Subscribe("jobs.")
}

func (x *Worker_t) Run(ctx context.Context) errnov1.Code {
	tickd := time.Second / time.Duration(x.Rate)
	next := time.Now()
	last := time.Now()

	Events().PublishAsync("workers.running", x.ID)

	for x.alive {
		select {
		case <-ctx.Done():
			return errnov1.OK

		default:
			now := time.Now()

			if now.Before(next) {
				remaining := next.Sub(now)

				if remaining > 2*time.Millisecond {
					time.Sleep(remaining - time.Millisecond)
				}

				continue
			}

			next = now.Add(tickd)

			dt := now.Sub(last)
			last = now

			x.rw.Lock()
			l := len(x.tasks)
			x.rw.Unlock()

			buff := make([]*Task_t, l)

			x.rw.RLock()
			for i := range l {
				if x.ctx[i] == nil {
					x.ctx[i], x.cancelFn[i] = context.WithCancel(ctx)
				}

				buff[i] = &Task_t{
					ID:        x.tasks[i],
					TTL:       x.ttl[i],
					DoneFn:    x.doneFn[i],
					ExecuteFn: x.executeFn[i],
					Context:   x.ctx[i],
					CancelFn:  x.cancelFn[i],
				}
			}
			x.rw.RUnlock()

			// Execute
			for i := range len(buff) {
				if doneFn := buff[i].DoneFn; doneFn != nil {
					if !doneFn() {
						if updateFn := buff[i].ExecuteFn; updateFn != nil {

							updateFn(buff[i].Context, dt)

							if doneFn() {
								Events().PublishAsync("jobs.done", &OnJobDone{
									ID: MakeJob(x.ID, buff[i].ID),
									At: time.Now().Unix(),
								})
							}
						}
					}
				}
			}
		}
	}

	return errnov1.OK
}

func (x *Worker_t) setContext(t Task, ctx context.Context, cancelFn context.CancelFunc) int {
	if idx := x.route(t); idx != -1 {
		x.cancelFn[idx] = cancelFn
		x.ctx[idx] = ctx

		return idx
	}

	return -1
}

func (x *Worker_t) route(v Task) int {
	if innerIdx, ok := x.router[v]; ok {
		return innerIdx
	}

	return -1
}

func (x *Worker_t) Tasks() []Task {
	return x.tasks
}

func (x *Worker_t) onNew(_ string, data any) {
	if job := data.(*Job_t); job != nil && job.Worker == x.ID && job.Task != nil {

		x.rw.Lock()
		idx := x.route(job.Task.ID)
		x.rw.Unlock()

		if idx != -1 {
			Events().PublishAsync("jobs.execute", &OnJobExecute{
				ID: MakeJob(x.ID, job.Task.ID),
				At: -1,
			})

			return
		}

		x.rw.Lock()
		idx = len(x.tasks)
		x.tasks = append(x.tasks, job.Task.ID)
		x.doneFn = append(x.doneFn, job.Task.DoneFn)
		x.executeFn = append(x.executeFn, job.Task.ExecuteFn)
		x.ttl = append(x.ttl, job.Task.TTL)
		x.ctx = append(x.ctx, nil)
		x.cancelFn = append(x.cancelFn, nil)
		x.router[job.Task.ID] = idx
		x.rw.Unlock()

		Events().PublishAsync("jobs.execute", &OnJobExecute{
			ID: MakeJob(x.ID, job.Task.ID),
			At: time.Now().Unix(),
		})
	}
}

func (x *Worker_t) onShutdown(_ string, data any) {
	if w, ok := data.(Worker); ok && w == x.ID {
		x.alive = false
	}
}

func (x *Worker_t) onRemove(_ string, data any) {
	if task, ok := data.(Task); ok {
		idx := -1
		x.rw.Lock()
		if idx = x.route(task); idx != -1 {
			l := len(x.tasks)

			if idx+1 != l {
				x.tasks[idx] = x.tasks[l-1]
				x.doneFn[idx] = x.doneFn[l-1]
				x.executeFn[idx] = x.executeFn[l-1]
				x.ttl[idx] = x.ttl[l-1]
				x.ctx[idx] = x.ctx[l-1]
				x.cancelFn[idx] = x.cancelFn[l-1]
			}

			x.tasks = x.tasks[:l-1]
			x.doneFn = x.doneFn[:l-1]
			x.executeFn = x.executeFn[:l-1]
			x.ttl = x.ttl[:l-1]
			x.ctx = x.ctx[:l-1]
			x.cancelFn = x.cancelFn[:l-1]
		}
		x.rw.Unlock()

		at := int64(0)
		if idx != -1 {
			at = time.Now().Unix()
		}

		Events().PublishAsync("jobs.removed", &OnJobRemoved{
			ID: MakeJob(x.ID, task),
			At: at,
		})
	}
}

func (x *Worker_t) onCancel(_ string, data any) {
	if task, ok := data.(Task); ok {
		cancelFn := (context.CancelFunc)(nil)

		x.rw.Lock()
		if idx := x.route(task); idx != -1 {
			cancelFn = x.cancelFn[idx]
		}
		x.rw.Unlock()

		at := int64(0)
		if cancelFn != nil {
			cancelFn()
			at = time.Now().Unix()
		}

		Events().PublishAsync("jobs.canceled", &OnJobCanceled{
			ID: MakeJob(x.ID, task),
			At: at,
		})
	}
}

func (x *Worker_t) Destroy() {
	x.tasks = x.tasks[:0]
	x.doneFn = x.doneFn[:0]
	x.executeFn = x.executeFn[:0]
	x.ttl = x.ttl[:0]
	x.ctx = x.ctx[:0]
	x.cancelFn = x.cancelFn[:0]
}
