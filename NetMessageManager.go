package limbov1

import (
	"slices"
	"time"

	"github.com/google/uuid"
	errnov1 "github.com/rejchev/errno"
	"google.golang.org/protobuf/proto"
)

type Recipient struct {
	Target Entity
	IsDone bool
}

type NetMessageManager struct {
	id         [][36]byte
	typ        []NetMessageType
	payload    [][]byte
	data       [][]byte
	recipients [][]Recipient
	createdAt  []time.Time

	router map[[36]byte]int
}

var nmessages = NetMessageManager{}

func NetMessages() *NetMessageManager {
	return &nmessages
}

func (x *NetMessageManager) Init() errnov1.Code {
	x.id = make([][36]byte, 0, 64)
	x.typ = make([]NetMessageType, 0, 64)
	x.payload = make([][]byte, 0, 64)
	x.data = make([][]byte, 0, 64)
	x.createdAt = make([]time.Time, 0, 64)
	x.router = map[[36]byte]int{}

	return errnov1.OK
}

func (x *NetMessageManager) Send(typ NetMessageType, payload proto.Message, data []byte, recipes []Entity) string {
	var err error
	buff := ([]byte)(nil)
	dbuff := ([]byte)(nil)

	if len(recipes) == 0 {
		return ""
	}

	if payload != nil {
		if buff, err = proto.Marshal(payload); err != nil {
			return ""
		}
	}

	if len(data) != 0 {
		dbuff = make([]byte, len(data))
		copy(dbuff, data)
	}

	x.id = append(x.id, [36]byte{})
	x.typ = append(x.typ, 0)
	x.payload = append(x.payload, nil)
	x.data = append(x.data, nil)
	x.recipients = append(x.recipients, make([]Recipient, len(recipes)))
	x.createdAt = append(x.createdAt, time.Now())

	idx := len(x.id) - 1

	copy(x.id[idx][:], uuid.NewString())
	x.typ[idx] = typ
	x.payload[idx] = buff
	x.data[idx] = dbuff

	for i := range len(recipes) {
		x.recipients[idx][i].Target = recipes[i]
	}

	x.router[x.id[idx]] = idx

	return string(x.id[idx][:])
}

func (x *NetMessageManager) Len() int {
	return len(x.id)
}

// Id is get first message available for entity (this guarantees the synchronous transmission of listed messages and end‑to‑end transmission for chained messages)
func (x *NetMessageManager) Peak(v Entity) string {
	return x.GetFirst(func(id string) bool {
		recipes := x.Recipients(id)

		for i := range len(recipes) {
			if recipes[i].Target == v {
				return true
			}
		}

		return false
	})
}

func (x *NetMessageManager) PeakA() string {
	if x.Len() != 0 {
		return string(x.id[0][:])
	}
	return ""
}

func (x *NetMessageManager) GetFirst(cond func(string) bool) string {
	for i := range len(x.id) {
		if cond(string(x.id[i][:])) {
			return string(x.id[i][:])
		}
	}

	return ""
}

func (x *NetMessageManager) Payload(v string) []byte {
	return x.payload[x.route(v)]
}

func (x *NetMessageManager) IsAlive(v string) bool {
	return x.route(v) != -1
}

func (x *NetMessageManager) Data(v string) []byte {
	return x.data[x.route(v)]
}

func (x *NetMessageManager) Recipients(v string) []Recipient {
	return x.recipients[x.route(v)]
}

func (x *NetMessageManager) CreatedAt(v string) time.Time {
	return x.createdAt[x.route(v)]
}

func (x *NetMessageManager) Type(v string) NetMessageType {
	return x.typ[x.route(v)]
}

func (x *NetMessageManager) RecipientsCount(v string) int {
	return len(x.recipients[x.route(v)])
}

// IsDone is get done state of message for entity (ret false if entity is missing)
func (x *NetMessageManager) IsDone(v string, e Entity) bool {
	if x.IsAlive(v) {
		if recipes := x.Recipients(v); len(recipes) != 0 {
			if j := slices.IndexFunc(recipes, func(r Recipient) bool {
				return r.Target == e
			}); j != -1 {
				return recipes[j].IsDone
			}
		}
	}

	return false
}

func (x *NetMessageManager) Done(v string, e Entity) {
	if x.IsAlive(v) {
		if recipes := x.Recipients(v); len(recipes) != 0 {
			if j := slices.IndexFunc(recipes, func(r Recipient) bool {
				return r.Target == e
			}); j != -1 {
				recipes[j].IsDone = true
			}
		}
	}
}

func (x *NetMessageManager) Message(v string, buff *NetMessageEntry) bool {
	if buff == nil {
		return false
	}

	if i := x.route(v); i != -1 {
		*buff = NetMessageEntry{
			Id:         v,
			Typ:        x.Type(v),
			Payload:    x.Payload(v),
			Data:       x.Data(v),
			Recipients: x.Recipients(v),
			CreatedAt:  x.CreatedAt(v),
		}

		return true
	}

	return false
}

func (x *NetMessageManager) route(v string) int {
	if i, ok := x.router[([36]byte)(([]byte)(v))]; ok {
		return i
	}

	return -1
}

func (x *NetMessageManager) CleanDone() {
	if id := x.PeakA(); id != "" {
		recipes := x.Recipients(id)

		for i := range len(recipes) {
			if Entities().IsAlive(recipes[i].Target) && !recipes[i].IsDone {
				return
			}
		}

		x.Remove()
	}
}

func (x *NetMessageManager) Remove() {
	if len(x.id) == 0 {
		return
	}

	delete(x.router, x.id[0])

	x.id = x.id[1:]
	x.typ = x.typ[1:]
	x.payload = x.payload[1:]
	x.data = x.data[1:]
	x.recipients = x.recipients[1:]
	x.createdAt = x.createdAt[1:]

	defer func() {
		for k, v := range x.router {
			x.router[k] = v - 1
		}
	}()
}

func (x NetMessageManager) Destroy() {}
