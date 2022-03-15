package str

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strconv"
	"strings"
)

func Decode(data []byte, encoder encoding.Encoding) string {
	data, _ = ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(data)), encoder.NewDecoder()))
	return string(data)
}

func StringDecode(data []byte, encoder string) string {
	if encoder == "GBK" {
		data, _ = ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(data)), simplifiedchinese.GBK.NewDecoder()))
		return string(data)
	}
	return string(data)
}

func BytesToHex(data []byte) []string {
	 sdata :=make([]string,len(data))
	for _, dat := range data {
		sdata = append(sdata,strings.TrimSpace(strconv.FormatInt(int64(dat),16)) )
	}
	return sdata
}