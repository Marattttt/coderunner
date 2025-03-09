package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
)

type InMemoryStateProvider struct {
	states map[string]struct{}
	mu     sync.Mutex
}

func NewInMemoryStateProvider() *InMemoryStateProvider {
	return &InMemoryStateProvider{
		states: make(map[string]struct{}),
	}

}

func (i *InMemoryStateProvider) GenerateState(OAuthResourceServer) string {
	const stateEncodedLen = 128

	// A byte can have 256 distinct values, and a base64 character,
	// has 64 values, which means the byte buffer should be
	// smaller than the resulting string
	const bufLen = stateEncodedLen / (256 / 64)

	buf := make([]byte, bufLen)
	_, _ = rand.Read(buf)

	encoded := base64.StdEncoding.EncodeToString(buf)

	i.mu.Lock()
	i.states[encoded] = struct{}{}
	i.mu.Unlock()

	return encoded
}

var (
	ErrStateNotFound = fmt.Errorf("State not found")
)

func (i *InMemoryStateProvider) InvalidateState(state string) error {
	_, ok := i.states[state]

	if !ok {
		return ErrStateNotFound
	}

	i.mu.Lock()
	delete(i.states, state)
	i.mu.Unlock()

	return nil
}
