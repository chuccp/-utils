package wire

import "math"

type VersionNumber uint32

const(
	Version1        VersionNumber = 0x1
	VersionUnknown  VersionNumber = math.MaxUint32
)