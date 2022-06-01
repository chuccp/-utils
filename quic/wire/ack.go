package wire

import "github.com/chuccp/utils/io"

type ACK struct {
	LargestAcknowledged uint8
	AckDelay            uint8
	AckRangeCount       uint8
	FirstAckRange       uint8
}

func ParseACK(readStream *io.ReadStream) (*ACK, error) {
	var ack ACK
	data, err := readStream.ReadBytes(4)
	if err == nil {
		ack.LargestAcknowledged = data[0]
		ack.AckDelay = data[1]
		ack.AckRangeCount = data[2]
		ack.FirstAckRange = data[3]
		return &ack, err
	}
	return nil, err
}
