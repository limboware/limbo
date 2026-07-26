package limbov1

import "testing"

func TestThatBuildingWorkCorrect(t *testing.T) {
	data := make([]byte, NETMSG_HEADER_SIZE)

	(*NetMessage)(&data[0]).Builder().
		PutFlags(FlagEncrypted).
		PutType(1).
		Build()

	if (*NetMessage)(&data[0]).Version() != NETMSG_VERSION {
		t.Errorf("expected version: %d, but got: %v", NETMSG_VERSION, data)
	}
}
