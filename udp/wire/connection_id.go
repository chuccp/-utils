package wire

import (
	"bytes"
	"fmt"
	"github.com/chuccp/utils/udp/util"
	"io"
)

type ConnectionID []byte

const maxConnectionIDLen = 20

func (c ConnectionID) Equal(other ConnectionID) bool {
	return bytes.Equal(c, other)
}
func ReadFullConnectionID(r io.Reader) (ConnectionID, uint32, error) {
	return nil, 0, nil
}
func ReadConnectionID(r io.Reader, len int) (ConnectionID, error) {
	if len == 0 {
		return nil, nil
	}
	c := make(ConnectionID, len)
	_, err := io.ReadFull(r, c)
	if err == io.ErrUnexpectedEOF {
		return nil, io.EOF
	}
	return c, err
}

func (c ConnectionID) Len() int {
	return len(c)
}
func (c ConnectionID) Bytes() []byte {
	return []byte(c)
}
func (c ConnectionID) String() string {
	if c.Len() == 0 {
		return "(empty)"
	}
	return fmt.Sprintf("%x", c.Bytes())
}

func GenerateConnectionID(len int) (ConnectionID, error) {
	data,err:=util.RandId(len)
	return data, err
}
