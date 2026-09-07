package limbov1

import (
	"time"
)

type NetMessageEntry struct {
	Id         string
	Typ        NetMessageType
	Payload    []byte
	Data       []byte
	Recipients []Recipient
	CreatedAt  time.Time
}
