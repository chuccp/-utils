package udp

import "encoding/binary"

type VersionNumber uint32

const (
	Version1        VersionNumber = 0xFF00001
)

func (vn VersionNumber) ToBytes()[]byte  {
	v:=[]byte{0,0,0,0}
	binary.LittleEndian.PutUint32(v, uint32(vn))
	return v
}