package util

type PacketNumber uint64

func (pn PacketNumber) GetPacketNumberLength() uint8 {
	if pn <= 0xFF {
		return 1
	}
	if pn <= 0xFF_FF {
		return 2
	}
	if pn <= 0xFF_FF_FF {
		return 3
	}
	return 4
}
