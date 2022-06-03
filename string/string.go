package str

import (
	"bytes"
	"encoding/hex"
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
	if encoder == "GBK" || encoder == "gbk" {
		data, _ = ioutil.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder()))
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
func removeSpaces(b []byte) []byte {
	for i := 0; i < len(b); {
		idx := bytes.IndexAny(b[i:], "\r\n\t ")
		if idx < 0 {
			break
		}
		i += idx
		copy(b[i:], b[i+1:])
		b = b[:len(b)-1]
	}
	return b
}
func DecodeHex(str string) []byte {
	data := removeSpaces([]byte(str))
	n, err := hex.Decode(data, data)
	if err != nil {
		panic(err)
	}
	return data[:n]
}