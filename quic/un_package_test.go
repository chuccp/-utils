package quic

import (
	"os"
	"testing"
)

func TestUnPackage(t *testing.T) {
	data,err:=os.ReadFile("C:\\Users\\cao\\Documents\\shandshake.bin")
	if err!=nil{
		t.Log(err)
	}else{
		un_package(data)
	}
}
