package limbov1

import (
	"context"
	"math"
	"sync"
	"time"

	errnov1 "github.com/rejchev/errno"
)

const WORKERS_MAX = ^uint8(0)

type WorkerFn = func(ctx context.Context)

func WorkerLoop(ctx context.Context, w Worker, tickrate int, shutdownedAt *int64) {
	var taskCtx context.Context
	var taskCancelFn context.CancelFunc
	var task Task

	tickd := time.Second / time.Duration(tickrate)
	next := time.Now()

	for (*shutdownedAt) == 0 {
		select {
		case <-ctx.Done():
			return

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

			if Tasks().Get(w, &task) {
				taskCtx, taskCancelFn = context.WithDeadline(ctx, time.Now().Add(Tasks().TTL(task)))

				errno := Tasks().Execute(w, task, taskCancelFn)(taskCtx, w, task)

				Tasks().SetDone(w, task, errno, time.Now().Unix())

				taskCancelFn()
			}
		}
	}
}

type WorkerManager struct {
	workers []Worker

	sinces       []int64
	cancels      []context.CancelFunc
	shutdownedAt []int64

	gen  []uint8
	free []uint8

	router map[Worker]int

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	rwm    sync.RWMutex
}

var w = WorkerManager{
	workers: make([]Worker, 0, 8),

	sinces:       make([]int64, 0, 8),
	cancels:      make([]context.CancelFunc, 0, 8),
	shutdownedAt: make([]int64, 0, 8),

	gen:  make([]uint8, 0, 8),
	free: make([]uint8, 0, 8),

	router: map[Worker]int{},

	wg:     sync.WaitGroup{},
	rwm:    sync.RWMutex{},
	ctx:    nil,
	cancel: nil,
}

func Workers() *WorkerManager {
	return &w
}

func (x *WorkerManager) Init() errnov1.Code {
	if x.ctx == nil {
		x.ctx, x.cancel = context.WithCancel(context.Background())
	}

	return errnov1.OK
}

func (x *WorkerManager) Destroy() {
	x.cancel()
	x.wg.Wait()
}

func (x *WorkerManager) Run(buff *Worker) bool {
	if len(x.workers) >= int(WORKERS_MAX - 1) {
		return false
	}

	id := uint8(0)

	if len(x.free) > 0 {
		id = x.free[len(x.free)-1]
		x.free = x.free[:len(x.free)-1]
	} else {
		id = uint8(len(x.gen))
		x.gen = append(x.gen, 0)
	}

	if len(x.cancels) <= int(id) {
		x.cancels = append(x.cancels, nil)
	}

	if len(x.sinces) <= int(id) {
		x.sinces = append(x.sinces, time.Now().Unix())
	}

	if len(x.shutdownedAt) <= int(id) {
		x.shutdownedAt = append(x.shutdownedAt, 0)
	}

	ctx, cancel := context.WithCancel(x.ctx)

	x.cancels[id] = cancel
	
	*buff = MakeWorker(id, x.gen[id])

	x.wg.Go(func() {
		WorkerLoop(ctx, *buff, 64, &x.shutdownedAt[id])
	})

	innerIdx := len(x.workers)

	x.workers = append(x.workers, *buff)

	x.router[*buff] = innerIdx

	Events().Publish("workers.new", *buff)

	return true
}

func (x *WorkerManager) Since(v Worker) int64 {
	return x.sinces[v.Id()]
}

func (x *WorkerManager) Cancel(v Worker) context.CancelFunc {
	return x.cancels[v.Id()]
}

func (x *WorkerManager) IsAlive(v Worker) bool {
	return x.gen[v.Id()] == v.Gen()
}

func (x *WorkerManager) ShutdownedAt(v Worker) int64 {
	return x.shutdownedAt[v.Id()]
}

func (x *WorkerManager) Shutdown(v Worker) {
	x.shutdownedAt[v.Id()] = time.Now().Unix()
} 

func (x *WorkerManager) IsRemovable(v Worker) bool {
	return x.ShutdownedAt(v) != 0 && (time.Now().Unix() - x.ShutdownedAt(v)) > 10 
}

func (x *WorkerManager) Iterator() *Iterator[Worker] {
	buff := make([]Worker, len(x.workers))

	copy(buff, x.workers)

	return NewIterator(buff)
}

func (x *WorkerManager) route(v Worker) int {
	if innerIdx, ok := x.router[v]; ok {
		return innerIdx
	}

	return -1
}

func (x *WorkerManager) Remove(v Worker) {
	if !x.IsAlive(v) || !x.IsRemovable(v) {
		return
	}

	innerIdx := x.route(v)

	if innerIdx == -1 {
		return
	}

	Events().Publish("workers.remove", v)

	idx := v.Id()
	if (x.gen[idx] + 1) == math.MaxUint8 {
		x.gen[idx] = 0
	}

	x.cancels[idx]()

	x.gen[idx]++
	x.cancels[idx] = nil
	x.sinces[idx] = 0
	x.free = append(x.free, idx)

	len := len(x.workers)

	if len > 1 && len-1 != innerIdx {
		x.workers[innerIdx] = x.workers[len-1]
		x.router[x.workers[len-1]] = innerIdx
	}

	delete(x.router, v)
	x.workers = x.workers[:len-1]

	Events().Publish("workers.removed", v)
}
