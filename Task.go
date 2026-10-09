package limbov1

type Task uint32

const INVALID_TASK = Task(^uint32(0))

func MakeTask(id uint32, gen uint8) Task {
	return Task((uint32(id) >> 8) << 8 | uint32(gen))
}

func (x Task) Id() uint32 {
	return uint32(x >> 8)
}

func (x Task) Gen() uint8 {
	return uint8(x)
}

func (x Task) SetId(v uint32) Task {
	return MakeTask(v, x.Gen())
}

func (x Task) SetGen(v uint8) Task {
	return MakeTask(x.Id(), v)
}

func (x Task) Closer() uint32 {
	return uint32(x)
}