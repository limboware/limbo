package limbov1

type Iterator[T comparable] struct {
	iterate []T
	pointer int
}

func NewIterator[T comparable](v []T) *Iterator[T] {
	return &Iterator[T]{
		iterate: v,
		pointer: 0,
	}
}

func (x *Iterator[T]) HasNext() bool {
	return len(x.iterate) > (x.pointer + 1)
}

func (x *Iterator[T]) ForEach(fn func(T)) {
	for _, y := range x.iterate {
		fn(y)
	}
}

func (x *Iterator[T]) First(buff *T, fn func(T) bool) bool {
	for _, y := range x.iterate {
		if fn(y) {
			*buff = y
			return true
		}
	}

	return false
}

func (x *Iterator[T]) Count() int {
	return len(x.iterate)
}

func (x *Iterator[T]) Reset() {
	x.pointer = 0
}

func (x Iterator[T]) Next() {
	x.pointer++
}

func (x *Iterator[T]) Get() T {
	return x.iterate[x.pointer]
}
