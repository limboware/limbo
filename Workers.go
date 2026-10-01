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

func WorkerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:

		}
	}
}

type WorkerManager struct {
	workers []Worker

	sinces  []int64
	cancels []context.CancelFunc

	gen  []uint8
	free []uint8

	router map[Worker]int

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	rwm    sync.RWMutex

	total int8
}

var w = WorkerManager{
	workers: make([]Worker, 0, 8),

	sinces:  make([]int64, 0, 8),
	cancels: make([]context.CancelFunc, 0, 8),

	gen:  make([]uint8, 0, 8),
	free: make([]uint8, 0, 8),

	router: map[Worker]int{},

	wg:     sync.WaitGroup{},
	rwm:    sync.RWMutex{},
	ctx:    nil,
	cancel: nil,
	total:  8,
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
	if len(x.workers) >= int(x.total) {
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

	ctx, cancel := context.WithCancel(x.ctx)

	x.cancels[id] = cancel

	x.wg.Go(func() {
		WorkerLoop(ctx)
	})

	*buff = MakeWorker(id, x.gen[id])

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
	if !x.IsAlive(v) {
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
