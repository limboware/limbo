package limbov1

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	errnov1 "github.com/rejchev/errno"
)

type WorkerManager struct {
	cancels []context.CancelFunc
	// jobCreatedAt    []int64

	router map[string]int

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

var w = WorkerManager{
	cancels: make([]context.CancelFunc, 0, 8),
	router:          map[string]int{},
	wg:              sync.WaitGroup{},
	ctx:             nil,
	cancel:          nil,
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

func (x *WorkerManager) Close() {
	x.cancel()
	x.wg.Wait()
}

func (x *WorkerManager) Run(ttl time.Duration, jobFn func(context.Context)) string {
	if jobFn == nil {
		return ""
	}

	id := uuid.NewString()
	causeCtx, cancel := context.WithDeadline(x.ctx, time.Now().Add(ttl+time.Second))

	idx := len(x.cancels)
	x.cancels = append(x.cancels, cancel)

	x.router[id] = idx

	x.wg.Go(func() {
		jobFn(causeCtx)
	})

	return id
}

func (x *WorkerManager) Cancel(v string) bool {
	if idx := x.navigate(v); idx != -1 {
		x.cancels[idx]()
		return true
	}

	return false
}

func (x *WorkerManager) navigate(v string) int {
	if id, ok := x.router[v]; ok {
		return id
	}

	return -1
}

func (x *WorkerManager) Destroy(v string) {
	idx := -1

	if idx = x.navigate(v); idx < 0 {
		return
	}

	x.cancels[idx]()

	len := len(x.cancels)

	last := len - 1

	if last != idx {
		lastKey := ""

		for k, v := range x.router {
			if v == last {
				lastKey = k
				break
			}
		}

		x.router[lastKey] = idx
		x.cancels[idx] = x.cancels[last]
	}

	x.cancels = x.cancels[:last]
	delete(x.router, v)
}
