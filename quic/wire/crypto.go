package wire

type Crypto struct {
	Offset uint32
	Length uint32
	data []byte
}

func  NewCrypto(Offset uint32,data []byte)*Crypto  {
	return &Crypto{Offset:Offset,Length: uint32(len(data)),data:data}
}
