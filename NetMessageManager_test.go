package limbov1

import (
	"math/rand"
	"testing"

	errnov1 "github.com/rejchev/errno"
)

func init() {
	NetMessages().Init()
}

func TestThatSendIsCorrect(t *testing.T) {
	const expectedType = 10

	if errno := NetMessages().Init(); errnov1.FAIL(errno) {
		t.Errorf("expected: OK, but got %s", errno.String())
		return
	}

	recipes := make([]Entity, 64)

	randomizeRecipes(recipes)

	id := NetMessages().Send(NetMessageType(expectedType), nil, nil, recipes)

	if id == "" {
		t.Errorf("expected uuid4, but got empty")
		return
	}

	rcp := NetMessages().Recipients(id)
	if !recipesEqual(recipes, rcp) {
		t.Errorf("expected %v, but got %v", recipes, rcp)
		return
	}

	payload := NetMessages().Payload(id)
	if payload != nil {
		t.Errorf("expected payload nil, but got %v", payload)
		return
	}

	data := NetMessages().Data(id)
	if data != nil {
		t.Errorf("expected data nil, but got %v", data)
		return
	}

	typ := NetMessages().Type(id)
	if typ != expectedType {
		t.Errorf("expected type %d, but got %d", expectedType, typ)
		return
	}
}

func TestThatPeakIsCorrect(t *testing.T) {
	const expectedType = 10

	recipes := make([]Entity, 64)

	randomizeRecipes(recipes)

	target := recipes[0]

	expectedId := NetMessages().Send(NetMessageType(expectedType), nil, nil, recipes)

	genMessages(recipes, 1)

	id := NetMessages().Peak(target)

	if expectedId != id {
		t.Errorf("expected %s, but got %s", expectedId, id)
	}
}

func genMessages(recipes []Entity, size int) {
	const bits = 64

	strategy := make([]uint64, size/bits+1)

	for i := range len(strategy) {
		strategy[i] = rand.Uint64()
	}

	for i := range len(strategy) {
		for j := range bits {
			if (1<<j)&strategy[i] == 0 {
				s := int(rand.Uint32()) % (len(recipes) / 2)
				e := (int(rand.Uint32()) % ((len(recipes) / 2) - 1)) + s

				NetMessages().Send(NetMessageType(rand.Uint32()), nil, nil, recipes[s:e])
			} else {
				s := int(rand.Uint32()) % (len(recipes) / 2)

				NetMessages().Send(NetMessageType(rand.Uint32()), nil, nil, recipes[s:s+1])
			}
		}
	}

}

func recipesEqual(ents []Entity, recipes []Recipient) bool {
	if len(ents) != len(recipes) {
		return false
	}

	for i := range len(ents) {
		if ents[i] != recipes[i].Target {
			return false
		}
	}

	return true
}

func randomizeRecipes(v []Entity) {
	for i := range len(v) {
		v[i] = CreateEntity(uint32(rand.Uint64()), uint16(rand.Uint64()))
	}
}
