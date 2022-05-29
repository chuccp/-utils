package wire

import "github.com/chuccp/utils/io"

type ExtensionType uint16

const (
	StatusRequest                       ExtensionType = 5
	SupportedGroups                                   = 10
	ECPointFormats                                    = 11
	SignatureAlgorithms                               = 13
	RenegotiationInfo                                 = 65281
	ApplicationLayerProtocolNegotiation               = 16
	SignedCertificateTimestamp                        = 18
	SupportedVersions                                 = 43
	KeyShare                                          = 51
	QuicTransportParameters                           = 57
)

type Extensions struct {
}

type Extension struct {
	Type   ExtensionType
	Length uint16
	Value  []byte
}

func ParseExtension(data []byte) (map[ExtensionType]*Extension, error) {
	var dataMap = make(map[ExtensionType]*Extension)
	readStream := io.NewReadBytesStream(data)
	ln := uint16(len(data))
	for {
		type_, err := readStream.Read2Uint16()
		if err != nil {
			return nil, err
		}
		var extension Extension
		extension.Type = ExtensionType(type_)
		extension.Length, err = readStream.Read2Uint16()
		if err != nil {
			return nil, err
		}
		extension.Value, err = readStream.ReadBytes(int(extension.Length))
		if err != nil {
			return nil, err
		}
		dataMap[extension.Type] = &extension
		if readStream.Offset() == ln {
			break
		}
	}
	return dataMap, nil
}
