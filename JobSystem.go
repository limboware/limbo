package limbov1

import (
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

type JobSystem struct {
	last Worker
}

func JobSystemFactory(buff *System) {
	i := new(JobSystem)

	*buff = System{
		Instance: unsafe.Pointer(i),
		Init: i.Init,

	}
}

func (x *JobSystem) Init() errnov1.Code {

	Events().Subscribe("tasks.new", x.onNewTask)
	Events().Subscribe("tasks.remove", x.onTaskRemove)
	
	Events().Subscribe("jobs.done", x.onDone)

	Events().Subscribe("workers.finished", x.onWorkerFinished)

	return errnov1.OK
}

func (x *JobSystem) onNewTask(_ string, data any) {
	if task, ok := data.(Task); ok {

		// TODO: balancer (no RR)

		iter := Workers().Iterator() 
		
		for iter.HasNext() {
			if iter.Get() == x.last {
				break
			}

			iter.Next()
		}

		if !iter.HasNext() {
			iter.Reset()
		}

		x.last = iter.Get()

		Events().PublishAsync("jobs.new", &Job_t{
			Worker: x.last,
			Task: Tasks().Struct(task),
		})
	}
}

func (x *JobSystem) onDone(_ string, data any) {
	if value, ok := data.(OnJobDone); ok {
		
		Events().Publish("tasks.done", value.ID)

		Tasks().Remove(value.ID.Task())
	} 
}

func (x *JobSystem) onTaskRemove(_ string, data any) {
	if task, ok := data.(Task); ok {
		Events().PublishAsync("jobs.remove", task)
	}
}

// TODO: impl
func (x *JobSystem) onWorkerFinished(_ string, data any) {}