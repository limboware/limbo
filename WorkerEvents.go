package limbov1

import errnov1 "github.com/rejchev/errno"

type OnWorkerFinished struct {
	ID    Worker
	Tasks []Task
	At    int64
	Code  errnov1.Code
}

type OnWorkerRun struct {
	ID Worker
	At int64
}
