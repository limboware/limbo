package limbov1

type Worker uint16

func MakeWorker(id, gen uint8) Worker {
	return Worker(uint16(id) << 8 | uint16(gen))
}

func (x Worker) Id() uint8 {
	return uint8(x >> 8)
}

func (x Worker) Gen() uint8 {
	return uint8(x)
}

func (x Worker) SetId(v uint8) Worker {
	return MakeWorker(v, x.Gen())
}

func (x Worker) SetGen(v uint8) Worker {
	return MakeWorker(x.Id(), v)
}

func (x Worker) Closer() uint16 {
	return uint16(x)
}