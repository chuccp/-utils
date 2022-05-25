package log

import (
	"errors"
	"strconv"
	"strings"
)

func toNumber(size string) (uint64,error) {
	size = strings.ToLower(size)
	var pre = ""
	end:=len(size)
	var v uint64= 1
	var index = 2
	if strings.HasSuffix(size,"kb"){
		v = 1<<10
	}else if strings.HasSuffix(size,"mb"){
		v = 1<<20
	}else if strings.HasSuffix(size,"gb"){
		v = 1<<30
	}else if strings.HasSuffix(size,"tb"){
		v = 1<<40
	}else if strings.HasSuffix(size,"pb"){
		v = 1<<50
	}else if strings.HasSuffix(size,"eb"){
		v = 1<<60
	}else if strings.HasSuffix(size,"zb"){
		return 0,errors.New("overflows uint64")
	}else{
		index = 1
	}
	pre = size[0:end-index]
	num,err:=strconv.ParseUint(pre,10,64)
	if err!=nil{
		return 0,err
	}

	return uint64(num * v),nil
}
func StringToLevel(level string)Level  {
	if level=="trace" || level=="TRACE"{
		return TraceLevel
	}
	if level=="debug" || level=="DEBUG"{
		return DebugLevel
	}
	if level=="info" || level=="INFO"{
		return InfoLevel
	}
	if level=="error" || level=="ERROR"{
		return ErrorLevel
	}
	if level=="PANIC" || level=="panic"{
		return PanicLevel
	}
	return InfoLevel
}