package log

import (
	"testing"
	"time"
)




func TestCut1(t *testing.T) {
	var str = "log/${time:2006-01-02_15:04}_${line:2000}_${size:200mb}_${level}.log"
	c,err:=parse(str)
	if err==nil{
		t:=time.Now()
		path,_:=c.getPath(&t,100,10,DebugLevel)
		println(path)
	}

}
func TestName(t *testing.T) {

	getTimeCut("${time:2006-01-02-15-04}")

}
