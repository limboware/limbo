package debug

import "sync"

type Debug sync.Mutex

func Debugf(buff string, data ...any) {
	if Get() != nil {
		Get().Printf(buff, data...)
	}
}