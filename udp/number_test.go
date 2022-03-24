package udp

import (
	"encoding/binary"
	"testing"
)

func TestReadVariableLength(t *testing.T) {



	println(0xf_ff_ff_ff)

	t.Log(ReadBytesVariableLength(VariableLengthBytes(0xf_ff_ff_ff)))

}
func TestAAAA(t *testing.T) {
	var packetNumber uint64=1
	data:=make([]byte,8)
	binary.BigEndian.PutUint64(data,packetNumber)
	t.Log(data)
}