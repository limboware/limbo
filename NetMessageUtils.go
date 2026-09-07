package limbov1

import (
	"google.golang.org/protobuf/proto"
	pblimbov1 "limboware.com/pkg/proto/limbo/v1"
)

type netMessageUtils struct{}

var netmu = netMessageUtils{}

func NetMessageUtils() netMessageUtils {
	return netmu
}

func (x *netMessageUtils) Send(typ NetMessageType, payload proto.Message, data []byte, recipe ...Entity) string {
	return NetMessages().Send(typ, payload, data, recipe)
}

func (x *netMessageUtils) SendB(typ NetMessageType, data []byte, recipe ...Entity) string {
	return NetMessages().Send(typ, nil, data, recipe)
}

func (x *netMessageUtils) SendC(typ NetMessageType, recipe ...Entity) string {
	return NetMessages().Send(typ, nil, nil, recipe)
}

func (x *netMessageUtils) SendD(typ NetMessageType, payload proto.Message, recipe ...Entity) string {
	return NetMessages().Send(typ, payload, nil, recipe)
}

func (x *netMessageUtils) Ping(e Entity) string {
	return x.SendC(NetMessageType(pblimbov1.MsgType_Ping), e)
}

func (x *netMessageUtils) Pong(e Entity) string {
	return x.SendC(NetMessageType(pblimbov1.MsgType_Pong), e)
}

func (x *netMessageUtils) Event(name string, payload proto.Message, data []byte, recipe ...Entity) string {
	bEvent := ([]byte)(nil)

	if payload != nil {
		if bts, err := proto.Marshal(payload); true {
			if err != nil {
				return ""
			}

			bEvent = bts
		}
	}

	return NetMessages().Send(NetMessageType(pblimbov1.MsgType_Event), &pblimbov1.MsgEvent{
		Name: name,
		Data: bEvent,
	}, data, recipe)
}

func (x *netMessageUtils) EventB(name string, payload proto.Message, recipe ...Entity) string {
	bEvent := ([]byte)(nil)

	if payload != nil {
		if bts, err := proto.Marshal(payload); true {
			if err != nil {
				return ""
			}

			bEvent = bts
		}
	}

	return NetMessages().Send(NetMessageType(pblimbov1.MsgType_Event), &pblimbov1.MsgEvent{
		Name: name,
		Data: bEvent,
	}, nil, recipe)
}