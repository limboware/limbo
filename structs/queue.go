package limbov1_structs

type Queue[T comparable] struct{}

func (x *Queue[T]) Pop(seq *[]T, buff *T) bool {
	if buff == nil || seq == nil || len(*seq) == 0 {
		return false
	}

	*buff = (*seq)[0] 

	*seq = (*seq)[1:]

	return true
}

func (x *Queue[T]) Peak(seq *[]T, buff *T) bool {
	if buff == nil || seq == nil || len(*seq) == 0 {
		return false
	}

	*buff = (*seq)[0]
	return true
}

func (x *Queue[T]) Push(buff *[]T, v T) {
	if buff == nil {
		return 
	}

	if len(*buff) < cap(*buff) {
		(*buff)[len(*buff)] = v
	} else {
		*buff = append(*buff, v)
	}
}