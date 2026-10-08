package limbov1

import errnov1 "github.com/rejchev/errno"


// TODO: impl
type jobSystem struct {}


func (x *jobSystem) Init() errnov1.Code {
	return errnov1.OK
}

func (x *jobSystem) onNewTask(_ string, data any) {}

func (x *jobSystem) onTaskRemove(_ string, data any) {}

func (x *jobSystem) onWorkerFinished(_ string, data any) {}