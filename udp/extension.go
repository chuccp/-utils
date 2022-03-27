package udp

import "bytes"

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
	signatureAlgorithms :=[]byte{0x00,0x0d,0x00,10,0,8,0x08,0x07,0x04,0x03,0x05,0x03,0x06,0x03}
	buff.Write(signatureAlgorithms)
	renegotiationInfo :=[]byte{0xff,0x01,0x00,0x01,0x00}
	buff.Write(renegotiationInfo)

	name:=[]byte("coke-name")
	ns:= U16B(uint16(len(name)+3))
	ns2:= U16B(uint16(len(name)+1))
	alpn:=[]byte{0x00,0x10,ns[0],ns[1],ns2[0],ns2[1],byte(len(name))}
	applicationLayerProtocolNegotiation := append(alpn, name...)
	buff.Write(applicationLayerProtocolNegotiation)
	signedCertificateTimestamp :=[]byte{0x00,0x12,0x00,0x00}
	buff.Write(signedCertificateTimestamp)
	supportedVersions :=[]byte{0x00,0x2b,0x00,0x03,2,3,4}
	buff.Write(supportedVersions)
	data,_:=RandId(32)
	keyShare :=[]byte{0x00,0x33,0,0x26,0,0x24,0,0x1d,0x00,0x20}
	//keyShare = append(keyShare, data...)
	buff.Write(keyShare)
	buff.Write(data)

	tp:=NewTransportParameters()
	tp.Init()
	tps:=tp.Bytes()
	den :=U16B(uint16(len(tps)))
	transportParameter :=[]byte{0x00,0x39, den[0], den[1]}
	buff.Write(transportParameter)
	buff.Write(tps)
	return buff.Bytes()

}