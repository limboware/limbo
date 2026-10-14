//go:build !limbo.debug

package debug

func Get() *Debug {
	return nil
}

func (x *Debug) Printf(v string, data ...any) {}
