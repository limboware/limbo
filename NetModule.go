package limbov1

import (
	errnov1 "github.com/rejchev/errno"
)

func InitNetModule() errnov1.Code {
	if err := NetConnTypes().Init(); errnov1.FAIL(err) {
		return err
	}

	if err := NetworkHandlers().Init(); errnov1.FAIL(err) {
		return err
	}

	return Networks().Init()
}