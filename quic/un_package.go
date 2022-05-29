package quic

import (
	"github.com/chuccp/utils/quic/wire"
	"log"
)

func un_package(data []byte) (uint16, error) {
	isLongHeader := data[0]&80 != 0
	if isLongHeader {
		longPackage, err := wire.ParseLongPackage(data)
		if err != nil {
			return 0, err
		}

		log.Print(longPackage)

	}
	return 0, nil
}
