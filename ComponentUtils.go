package limbov1

import errnov1 "github.com/rejchev/errno"

func ComponentPtr[T any](e Entity, t Compotype) *T {
	components := Components()

	if c := ComponentFor(e, t); components.Contains(c) {
		return (*T)(components.Get(c))
	}

	return nil
}

func NewComponent[T any](e Entity, t Compotype, buff *Component) *T {
	components := Components()
	if errnov1.FAIL(components.New(e, t, buff)) {
		return nil
	}

	return (*T)(components.Get(*buff))
}

func NewComponentB[T any](e Entity, t Compotype) *T {
	var component Component

	return NewComponent[T](e, t, &component)
}

func ContainsComponent(e Entity, t Compotype) bool {
	return Components().Contains(ComponentFor(e, t))
}

func DestroyComponent(e Entity, t Compotype) {
	Components().Destroy(ComponentFor(e, t))
}
