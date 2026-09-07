package limbov1

import (
	"google.golang.org/protobuf/proto"
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
