package limbov1

import "testing"

func TestAB(t *testing.T) {
	x := ComponentMask{}
	v := uint8(0) | 0x4

	x = x.With(Compotype(8+8+2))

	if y := x.Get(2); y != v {
		t.Errorf("expected %d, but got %d", v, y)
	} 
}
