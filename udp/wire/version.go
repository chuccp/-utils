package wire

import (
	"encoding/binary"
	"math"
)

type VersionNumber uint32

const (
	VersionTLS      VersionNumber = 0x1
	VersionWhatever VersionNumber = math.MaxUint32 - 1 // for when the version doesn't matter
	VersionUnknown  VersionNumber = math.MaxUint32
	VersionDraft29  VersionNumber = 0xff00001d
	Version1        VersionNumber = 0x1
)

func ParseVersion(data []byte) VersionNumber {
	vn:=binary.BigEndian.Uint32(data)
	switch VersionNumber(vn) {
	case VersionTLS:
		return VersionTLS
	case VersionDraft29:
		return VersionDraft29
	case VersionWhatever:
		return VersionWhatever
	default:
		return VersionUnknown
	}
	return VersionUnknown
}