package log

import (
	"log"
	"testing"
)




func TestCut1(t *testing.T) {
	var str = "log/${time:2006-01-02-15-04}${line:2000}${size:200mb}${level}.log"
	c,err:=parse(str)

	if err==nil{

		log.Println(c)
	}

}
