package wire

import (
	"bytes"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/quic/util"
)

const (
	ACKFrameType      byte = 2
	CRYPTOFrameType   byte = 6
	PADDINGFrameType0 byte = 0
	PADDINGFrameType1 byte = 1
)

type Frame struct {
	FrameType byte
	Data      []byte
}

func ParseFrame(data []byte) ([]byte, error) {
	stream := io.NewReadBytesStream(data)
	cryptoMap := make(map[uint32]*Crypto)
	for {
		b, err := stream.ReadByte()
		if err != nil {
			break
		}
		if PADDINGFrameType0 == b || b == PADDINGFrameType1 {
			continue
		}
		if b == ACKFrameType {
			_, err = ParseACK(stream)

		}
		if b == CRYPTOFrameType {
			readValue := util.NewReadValue(stream)
			offset, err := readValue.ReadVariableValueLength()
			if err != nil {
				return nil, err
			}
			len, err := readValue.ReadVariableValueLength()
			if err != nil {
				return nil, err
			}
			data2, err := stream.ReadBytes(int(len))
			if err != nil {
				return nil, err
			}
			cryptoMap[offset] = NewCrypto(offset, data2)
		}
	}
	buff := new(bytes.Buffer)
	var offset uint32 = 0
	for {
		crypto := cryptoMap[(offset)]
		if crypto == nil {
			break
		}
		buff.Write(crypto.data)
		offset = crypto.Offset + crypto.Length
	}
	return buff.Bytes(), nil
}
