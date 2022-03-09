package udp

import "math/rand"

type ConnectionID []byte
func GenerateConnectionID(len int) (ConnectionID, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return ConnectionID(b), nil
}
