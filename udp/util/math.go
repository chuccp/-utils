package util

func VariableLengthToBytes(length uint32)[]byte  {
	if length<=0x3f{
		return []byte{byte(length)}
	}
	if length<=0x3fff{
		return []byte{byte(length>>8)|0x40,byte(length)}
	}
	if length<=0x3fffff{
		return []byte{byte(length>>16)|0x80,byte(length>>8),byte(length)}
	}else{
		return []byte{byte(length>>24)|0xc0,byte(length>>16),byte(length>>8),byte(length)}
	}
}

func BTU16(ds []byte) uint16 {
	return uint16(ds[0])<<8|uint16(ds[1])
}