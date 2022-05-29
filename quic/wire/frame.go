package wire

import (
	"bytes"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/quic/util"
)



func ParseFrame(data []byte) ([]byte,error){
	stream := io.NewReadBytesStream(data)
	cryptoMap:=make(map[uint32]*Crypto)
	for{
		b,err:=stream.ReadByte()
		if err!=nil{
			break
		}
		if b==0x00 || b==0x01{
			continue
		}
		if b==0x06{
			readValue := util.NewReadValue(stream)
			offset,err:=readValue.ReadVariableValueLength()
			if err!=nil{
				return nil,err
			}
			len,err:=readValue.ReadVariableValueLength()
			if err!=nil{
				return nil,err
			}
			data2,err:=stream.ReadBytes(int(len))
			if err!=nil{
				return nil,err
			}
			cryptoMap[offset] = NewCrypto(offset,data2)
		}
	}
	buff:=new(bytes.Buffer)
	var offset uint32= 0
	for{
		crypto:=cryptoMap[(offset)]
		if crypto==nil{
			break
		}
		buff.Write(crypto.data)
		offset = crypto.Offset+crypto.Length
	}
	return buff.Bytes(),nil
}
