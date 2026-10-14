package limbov1

import (
	"context"
	"sync"
	"time"
	"unsafe"

	"github.com/limboware/limbo/debug"
	errnov1 "github.com/rejchev/errno"
)

type WorkerSystem struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func WorkerSystemFactory(buff *System) {
	v := new(WorkerSystem)

	*buff = System{
		Instance:   unsafe.Pointer(v),
		Init:       v.Load,
		Activate:   nil,
		Update:     nil,
		Deactivate: nil,
		Destroy:    v.Unload,
	}
}

func (x *WorkerSystem) Load() errnov1.Code {
	debug.Get().Printf("<WorkerSystem::Init> start")

	Events().Subscribe("workers.new", x.onWorker)
	Events().Subscribe("workers.running", x.onWorkerRunning)
	Events().Subscribe("workers.finished", x.onWorkerFinished)
	Events().Subscribe("workers.remove", x.onWorkerRemove)

	x.ctx, x.cancel = context.WithCancel(context.Background())

	if err := Workers().Init(); errnov1.FAIL(err) {
		return err
	}

	debug.Get().Printf("<WorkerSystem::Init> done")

	return errnov1.OK
}

func (x *WorkerSystem) onWorker(_ string, data any) {
	debug.Get().Printf("<WorkerSystem::onWorker> worker %v", data)

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

		Events().Publish("workers.created", w)

		debug.Get().Printf("<WorkerSystem::onWorker> worker %d created", w)
	}
}

func (x *WorkerSystem) onWorkerRunning(_ string, data any) {
	debug.Get().Printf("<WorkerSystem::onWorkerRunning> called with: %v", data)
	if w := data.(*OnWorkerRun); Workers().IsAlive(w.ID) {
		Workers().SetRunnedAt(w.ID, w.At)

		debug.Get().Printf("<WorkerSystem::onWorkerRunning> worker: %d; at: %d", w.ID, w.At)
	}
}

func (x *WorkerSystem) onWorkerFinished(_ string, data any) {
	debug.Get().Printf("<WorkerSystem::onWorkerFinished> called with: %v", data)
	if w, ok := data.(*OnWorkerFinished); ok {

		if Workers().IsAlive(w.ID) {
			Workers().SetShuthdownedAt(w.ID, w.At)
		}

		debug.Get().Printf("<WorkerSystem::onWorkerFinished> worker: %d; at: %d", w.ID, w.At)
	}
}

func (x *WorkerSystem) onWorkerRemove(_ string, data any) {
	debug.Get().Printf("<WorkerSystem::onWorkerRemove> called with: %v", data)
	if w := data.(Worker); Workers().IsAlive(w) {
		if cancelFn := Workers().CancelFn(w); cancelFn != nil {
			cancelFn()

			debug.Get().Printf("<WorkerSystem::onWorkerRemove> worker: %d canceled", w)
		}
	}
}

func (x *WorkerSystem) Unload() {
	debug.Get().Printf("<WorkerSystem::Unload> start")
	Workers().Destroy()

	x.cancel()
	x.wg.Wait()

	debug.Get().Printf("<WorkerSystem::Unload> done")
}
