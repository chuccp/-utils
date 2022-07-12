package util

import "encoding/binary"

type VersionNumber uint32

const (
	Version1 VersionNumber = 0xFF000001
)

func (vn VersionNumber) ToBytes()[]byte  {
	v:=[]byte{0,0,0,0}
	binary.BigEndian.PutUint32(v, uint32(vn))
	return v
}