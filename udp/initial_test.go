package udp

import (
	"github.com/chuccp/utils/file"
	"github.com/chuccp/utils/io"
	"testing"
)


func TestRaw(t *testing.T) {

	fi, err := file.NewFile("C:\\Users\\cooge\\Documents\\quic\\Initial4.bin")
	if err == nil {
		file, err1 := fi.ToRawFile()
		if err1 == nil {
			read := io.NewReadStream(file)
			un_package_stream(read)
		}
	}

}
