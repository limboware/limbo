package limbov1

import "context"

type OnJobExecute struct {
	ID       Job
	At       int64
	CancelFn context.CancelFunc
}

type OnJobDone struct {
	ID Job
	At int64
}

type OnJobCanceled struct {
	ID Job
	At int64
}

type OnJobRemoved struct {
	ID Job
	At int64
}
