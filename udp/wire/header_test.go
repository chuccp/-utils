package wire

import "testing"

func TestParseHeader(t *testing.T) {



	header,err:=ParseFileHeader("C:\\Users\\cao\\Documents\\serverHello.bin")
	if err!=nil{
		t.Log(err)
	}else{
		t.Log(header)
	}

}
