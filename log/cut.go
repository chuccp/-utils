package log

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TimeCut int

const (
	Second TimeCut = iota
	Minute
	Hour
	Day
	Month
	Year
)

type fileCut []string

func (fileCut *fileCut) parse(filePattern string) {
	end := len(filePattern)
	cutIndex := 0
	for i := 0; i < end; {
		start := i
		var flag = false
		for j := i; j < end; {
			if !flag {
				if j+7 < end {
					if filePattern[j:j+7] == "${time:" {
						start = j
						j = j + 7
						flag = true
					} else if filePattern[j:j+7] == "${line:" {
						start = j
						j = j + 7
						flag = true
					} else if filePattern[j:j+7] == "${size:" {
						start = j
						j = j + 7
						flag = true
					} else if filePattern[j:j+7] == "${level" {
						start = j
						j = j + 7
						flag = true
					} else {
						break
					}
				} else {
					break
				}
			} else {
				if filePattern[j] == '}' {
					if cutIndex != start {
						*fileCut = append(*fileCut, filePattern[cutIndex:start])
					}
					*fileCut = append(*fileCut, filePattern[start:j+1])
					cutIndex = j + 1
					i = j
					flag = false
					break
				}
				j++
			}
		}
		i++
	}
	if cutIndex < end-1 {
		*fileCut = append(*fileCut, filePattern[cutIndex:end])
	}
}

type cut struct {
	//时间切割
	timeCut  TimeCut
	line     uint64
	size     uint64
	hasLevel bool
	fileCut  fileCut
	io       io.Writer
	path     string
}

func (cut *cut) getPath(time *time.Time) string {

	return "log/info.log"
}
func (cut *cut) getOut(time *time.Time) (io io.Writer, err error) {
	path := cut.getPath(time)
	if cut.path == path {
		if io != nil {
			return io, nil
		}
	}
	ii := filepath.Dir(path)
	err = os.MkdirAll(ii, 0777)
	if err != nil {
		return nil, err
	}
	io, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	cut.io = io
	cut.path = path
	return
}
func getTimeCut(t string) TimeCut {

	if strings.Index(t, "05") > -1 {
		return Second
	}
	if strings.Index(t, "04") > -1 {
		return Minute
	}
	if strings.Index(t, "15") > -1 {
		return Hour
	}
	if strings.Index(t, "02") > -1 {
		return Day
	}
	if strings.Index(t, "01") > -1 {
		return Month
	}
	if strings.Index(t, "2006") > -1 {
		return Year
	}
	return -1
}

func parse(filePattern string) (*cut, error) {
	var cut cut
	cut.timeCut = -1
	cut.fileCut.parse(filePattern)
	for _, v := range cut.fileCut {
		if strings.HasPrefix(v, "${time:") {
			end := len(v)
			time := v[7 : end-1]
			tc := getTimeCut(time)
			if tc > cut.timeCut {
				cut.timeCut = tc
			}
		} else if strings.HasPrefix(v, "${line:") {
			end := len(v)
			line := v[7 : end-1]
			num, err := strconv.ParseUint(line, 10, 64)
			if err != nil {
				return nil, err
			}
			cut.line = num
		} else if strings.HasPrefix(v, "${size:") {
			end := len(v)
			size := v[7 : end-1]
			num, err := toNumber(size)
			if err != nil {
				return nil, err
			}
			cut.size = num
		} else if strings.HasPrefix(v, "${level") {
			cut.hasLevel = true
		}
	}
	return &cut, nil
}
