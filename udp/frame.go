package udp

import (
	"bytes"
	"github.com/chuccp/utils/io"
	io2 "io"
)

type Crypto struct {
	Offset uint32
	Length uint32
	data []byte
}

func  NewCrypto(Offset uint32,Length uint32,data []byte)*Crypto  {

	return &Crypto{Offset:Offset,Length:Length,data:data}
}

func ParseFrame(data []byte) (*bytes.Buffer,error){
	r := io.NewReadStream(bytes.NewReader(data))
	cryptoMap:=make(map[uint32]*Crypto)
	for{
		b,err:=r.ReadByte()
		if err!=nil{
			break
		}
		if b==0x00 || b==0x01{
			continue
		}
		if b==0x06{
			offset,num:=ReadVariableLength(r)
			if num==0{
				return nil,io2.EOF
			}
			len,num:=ReadVariableLength(r)
			if num==0{
				return nil,io2.EOF
			}

			data,err:=r.ReadUintBytes(len)
			if err!=nil{
				return nil,err
			}
			cryptoMap[offset] = NewCrypto(offset,len,data)
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
	return buff,nil
}
