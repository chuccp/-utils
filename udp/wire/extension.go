package wire

import (
	"bytes"
	"github.com/chuccp/utils/udp/util"
)

type ExtensionType uint16

const (



)

func Extansions() []byte {
	var buff =new(bytes.Buffer)
	statusRequest :=[]byte{0x00,0x05,0x00,0x05,0x01,0x00,0x00,0x00,0x00}
	buff.Write(statusRequest)
	supportedGroups :=[]byte{0x00,0x0a,0x00,0x0a,0x00,0x08,0x00,0x1d,0x00,0x17,0x00,0x18,0x00,0x19}
	buff.Write(supportedGroups)
	ecPointFormats :=[]byte{0x00,0x0b,0x00,0x02,0x01,0x00}
	buff.Write(ecPointFormats)
	signatureAlgorithms:=[]byte{0,0xd,0,0x1a,0,0x18,8,4,4,3,8,7,8,5,8,6,4,1,5,1,6,1,5,3,6,3,2,1,2,3}
	buff.Write(signatureAlgorithms)
	renegotiationInfo :=[]byte{0xff,0x01,0x00,0x01,0x00}
	buff.Write(renegotiationInfo)
	name:=[]byte("quic-echo-example")
	ns:= util.U16B(uint16(len(name)+3))
	ns2:= util.U16B(uint16(len(name)+1))
	alpn:=[]byte{0x00,0x10,ns[0],ns[1],ns2[0],ns2[1],byte(len(name))}
	applicationLayerProtocolNegotiation := append(alpn, name...)
	buff.Write(applicationLayerProtocolNegotiation)
	signedCertificateTimestamp :=[]byte{0x00,0x12,0x00,0x00}
	buff.Write(signedCertificateTimestamp)
	supportedVersions :=[]byte{0x00,0x2b,0x00,0x03,2,3,4}
	buff.Write(supportedVersions)
	data,_:=util.RandId(32)
	keyShare :=[]byte{0x00,0x33,0,0x26,0,0x24,0,0x1d,0x00,0x20}
	//keyShare = append(keyShare, data...)
	buff.Write(keyShare)
	buff.Write(data)

	tp:=NewTransportParameters()
	tp.Init()
	tps:=tp.Bytes()
	den :=util.U16B(uint16(len(tps)))
	transportParameter :=[]byte{0x00,0x39, den[0], den[1]}
	buff.Write(transportParameter)
	buff.Write(tps)
	return buff.Bytes()

}