//go:build limbo.debug

package debug

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var dbg = Debug{}

func Get() *Debug {
	return &dbg
}

func (x *Debug) Printf(buff string, data ...any) {
	y := fmt.Sprintf(header()+" "+buff+"\n", data...)

	x.closer().Lock()
	_, _ = os.Stderr.Write([]byte(y))
	x.closer().Unlock()
}

func header() string {
	t := time.Now()

	return fmt.Sprintf("[%d/%s/%d %d:%d:%d (%d) µ]",
		t.Local().Day(),
		t.Local().Month().String(),
		t.Local().Year(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond()/1000,
	)
}

func (x *Debug) closer() *sync.Mutex {
	return (*sync.Mutex)(x)
}
