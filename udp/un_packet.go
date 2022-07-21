package udp

import (
	"github.com/chuccp/utils/udp/wire"
)

func ParseHeader(data []byte, header *wire.Header)error {
	fistByte := data[0]
	header.IsLongHeader = fistByte&0x80 > 0
	if header.IsLongHeader {
		err := header.ParseLongHeader(data)
		if err != nil {
			return err
		}
	}
	return nil
}





