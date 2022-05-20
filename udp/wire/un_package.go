package wire




//func un_package_stream(read *io.ReadStream) {
//
//	buff := NewBuffer()
//	err1 := buff.readPack(read)
//	if err1 == nil {
//		header := parseLongHeader(buff.data)
//		if header.packetType==PacketTypeInitial{
//			origPNBytes := make([]byte, 4)
//			copy(origPNBytes, buff.data[header.parsedLen:header.parsedLen+4])
//			sample :=buff.data[header.parsedLen+4:header.parsedLen+4+16]
//			param2 := &buff.data[0]
//			param3 := buff.data[header.parsedLen:header.parsedLen+4]
//			aead:=NewInitialAEAD(header.desConnId)
//			mask:=make([]byte, aead.block.BlockSize())
//			log.Info("sample2",sample)
//			aead.block.Encrypt(mask,sample)
//			*param2 ^= mask[0] & 0xf
//			for i := range param3 {
//				param3[i] ^= mask[i+1]
//			}
//			pageNum:=int(buff.data[0]&0x3)+1
//			exlen:=header.parsedLen+pageNum
//			nonceBuf:=make([]byte, aead.aead.NonceSize())
//			binary.BigEndian.PutUint64(nonceBuf[len(nonceBuf)-8:], 0)
//			log.Info("aead.iv=====3:",aead.iv)
//			for i, b := range nonceBuf[len(nonceBuf)-8:] {
//				aead.iv[4+i] ^= b
//			}
//			log.Info("aead.iv=====4:",aead.iv)
//			copy(buff.data[header.parsedLen+pageNum:header.parsedLen+4],origPNBytes[pageNum:4] )
//			ciphertext:=buff.data[exlen:buff.len]
//			additionalData:=buff.data[:exlen]
//			log.Info("additionalData=====2:",additionalData)
//			data,err:=aead.aead.Open([]byte{},aead.iv,ciphertext,additionalData)
//			log.Info(ConnectionID(data),err)
//			if err==nil{
//				buff,_:=ParseFrame(data)
//				ParseCrypto(buff)
//			}
//		}
//	}
//
//
//}
