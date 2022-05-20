package wire

import (
	"github.com/chuccp/utils/io"
	"os"
)



func newHeader(TypeByte byte) {

}

type Header struct {
	IsLongHeader    bool
	Type            PacketType
	TypeByte        byte
	Version         VersionNumber
	desConnIdLen    uint8
	DesConnId       ConnectionID
	sourceConnIdLen uint8
	SourceConnId    ConnectionID
	ParsedLen       ByteCount
}



func (header *Header) HandleLongHeader(read *io.ReadStream) error {
	switch (header.TypeByte & 0x30) >> 4 {
	case 0x0:
		header.Type = PacketTypeInitial
	case 0x1:
		header.Type = PacketType0RTT
	case 0x2:
		header.Type = PacketTypeHandshake
	case 0x3:
		header.Type = PacketTypeRetry
	}
	if data, err := read.ReadBytes(4); err != nil {
		return err
	} else {
		header.Version = ParseVersion(data)
	}
	if err := header.ReadConnectionId(read); err != nil {
		return err
	}
	if header.Type == PacketTypeRetry {
		header.ParsedLen = ByteCount(header.sourceConnIdLen + 1 + header.desConnIdLen + 1 + 5)
	}
	return nil
}





func (header *Header) ReadConnectionId(read *io.ReadStream) (err error) {

	if header.desConnIdLen, err = read.ReadByte(); err != nil {
		return
	}
	if header.desConnIdLen != 0 {
		if header.DesConnId, err = read.ReadBytes(int(header.desConnIdLen)); err != nil {
			return
		}
	}else{
		header.DesConnId = []byte{}
	}
	if header.sourceConnIdLen, err = read.ReadByte(); err != nil {
		return
	}
	if header.sourceConnIdLen != 0 {
		if header.SourceConnId, err = read.ReadBytes(int(header.sourceConnIdLen)); err != nil {
			return
		}
	}else{
		header.SourceConnId = []byte{}
	}
	return nil

}

func ParseFileHeader(path string)(header *Header, err error)  {
	file,err:=os.Open(path)
	if err!=nil{
		return nil, err
	}
	return ParseHeader(io.NewReadStream(file))
}

func ParseHeader(read *io.ReadStream) (header *Header, err error) {
	header = &Header{}
	header.TypeByte, err = read.ReadByte()
	if err != nil {
		return nil, err
	}
	header.IsLongHeader = !(header.TypeByte&0x80 == 0)
	if header.IsLongHeader{
		err:=header.HandleLongHeader(read)
		if err!=nil{
			return nil, err
		}
	}
	return header, nil
}
func ParsePacket(data []byte) (*Header,error) {
	read := io.NewReadBytesStream(data)
	header, err := ParseHeader(read)
	if err != nil {
		return nil,err
	}
	return header,nil
}
