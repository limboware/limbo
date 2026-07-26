package limbov1

import "fmt"

type Compotype uint8

func (x Compotype) Closer() uint8 {
	return (uint8)(x)
}

func (x Compotype) Int() int {
	return (int)(x)
}

func (x Compotype) String() string {
	return fmt.Sprintf("%d", x.Closer())
}
