package limbov1

type Job uint64

func MakeJob(w Worker, t Task) Job {
	return Job(uint64(w)<<32 | uint64(t))
}

func (x Job) Worker() Worker {
	return Worker(x >> 32)
}

func (x Job) Task() Task {
	return Task(x)
}

type Job_t struct {
	Worker Worker
	Task   *Task_t
}
