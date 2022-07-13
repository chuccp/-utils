package udp

import "github.com/chuccp/utils/udp/util"

func UnPacket(data []byte) error {
	fistByte := data[0]
	if (fistByte & 0x80) != 0 {
		var longHeader LongHeader
		rb := util.NewReadBuffer(data)
		return longHeader.Read(rb)

	}
	return nil
}

