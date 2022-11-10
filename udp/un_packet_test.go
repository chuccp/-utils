package udp

import (
	"github.com/chuccp/utils/file"
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/wire"
	"testing"
)

func TestUnInitial(t *testing.T) {

	rc := config.NewReceiveConfig()
	newFile, err := file.NewFile("data001.bb")
	if err == nil {
		data, b, err := newFile.Read()
		if err != nil || !b {
			return
		}
		wire.ParseInitialLongHeader(data, rc)
	}
}
