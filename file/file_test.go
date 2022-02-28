package file

import (
	"testing"
)


func TestExists(t *testing.T) {
	file,err:=NewFile("static")
	if err==nil{
		t.Log(file.Name(),"====",file.IsDir())
	}
}

func TestOpenOrCreateFile(t *testing.T) {
	file,err:=NewFile("static/sss.ini")
	t.Log("!!!========",err)
	if err==nil{
		data,flag,err3:=file.Read()
		if flag && err3==nil{
			t.Log(data)
		}else{
			t.Log(flag,err3)
		}
	}else{
		t.Log(err)
	}
}

func TestGetRootPath(t *testing.T) {
	files,err:=GetRootPath()
	if err==nil{
		for _, file := range files {
			t.Logf(file.Name())
		}
	}
}
func TestFile_Read(t *testing.T) {


}