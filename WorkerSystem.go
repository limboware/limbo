package limbov1

import (
	"context"
	"sync"
	"time"
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

type workerSystem struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func WorkerSystem(buff *System) {
	v := new(workerSystem)

	*buff = System{
		Instance:   unsafe.Pointer(v),
		Init:       v.Load,
		Activate:   nil,
		Update:     nil,
		Deactivate: nil,
		Destroy:    v.Unload,
	}
}

func (x *workerSystem) Load() errnov1.Code {

	Events().Subscribe("workers.new", x.onWorker)
	Events().Subscribe("workers.running", x.onWorkerRunning)
	Events().Subscribe("workers.finished", x.onWorkerFinished)
	Events().Subscribe("workers.remove", x.onWorkerRemove)

	x.ctx, x.cancel = context.WithCancel(context.Background())

	return Workers().Init()
}

func (x *workerSystem) onWorker(_ string, data any) {
	if w := data.(Worker); Workers().IsAlive(w) {
		ctx, cancel := context.WithCancel(x.ctx)

		x.wg.Go(func() {
			worker := Worker_t{
				ID:   w,
				Rate: 64,
			}

			worker.Init()

			errno := worker.Run(ctx)

			Events().PublishAsync("workers.finished", &OnWorkerFinished{
				ID:    worker.ID,
				Tasks: worker.Tasks(),
				At:    time.Now().Unix(),
				Code:  errno,
			})

			worker.Destroy()
		})

		Workers().SetCancelFn(w, cancel)
	}
}

func (x *workerSystem) onWorkerRunning(_ string, data any) {
	if w := data.(OnWorkerRun); Workers().IsAlive(w.ID) {
		Workers().SetRunnedAt(w.ID, w.At)
	}
}

func (x *workerSystem) onWorkerFinished(_ string, data any) {
	if w := data.(OnWorkerFinished); Workers().IsAlive(w.ID) {
		Workers().SetShuthdownedAt(w.ID, w.At)
	}
}

func (x *workerSystem) onWorkerRemove(_ string, data any) {
	if w := data.(Worker); Workers().IsAlive(w) {
		if cancelFn := Workers().CancelFn(w); cancelFn != nil && Workers().ShutdownedAt(w) == 0 {
			cancelFn()
		}
	}
}

func (x *workerSystem) Unload() {
	x.cancel()
	x.wg.Wait()
}
