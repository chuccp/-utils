package string

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"

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
