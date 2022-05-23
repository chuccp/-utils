package log

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TimeCut struct {
	timeCutFlag timeCutFlag
	timeFormat  string
}

func (tc *TimeCut) parse(ti *time.Time) string {
	return ti.Format(tc.timeFormat)
}

type timeCutFlag int

const (
	Second timeCutFlag = iota
	Minute
	Hour
	Day
	Month
	Year
)

func SplitPath(filePattern string) []string {
	var paths = make([]string, 0)
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
						paths = append(paths, filePattern[cutIndex:start])
					}
					paths = append(paths, filePattern[start:j+1])
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
		paths = append(paths, filePattern[cutIndex:end])
	}
	return paths
}

type cut struct {
	//时间切割
	timeCut  *TimeCut
	line     uint64
	size     uint64
	hasLevel bool
	file     *os.File
	path     string
	ctime    string
}

func (cut *cut) getPath(time *time.Time, line uint64, size uint64, level Level) (path string, flag bool) {
	ti := cut.timeCut.parse(time)
	if ti == cut.ctime {
		flag = false
	} else {
		cut.ctime = ti
		flag = true
	}
	path = strings.ReplaceAll(cut.path, "${TIME}", ti)
	if cut.line > 0 {
		line = (line / cut.line) * cut.line
		path = strings.ReplaceAll(path, "${LINE}", strconv.FormatUint(line, 10))
	}
	if cut.size > 0 {
		size = (size / cut.size) * cut.size
		path = strings.ReplaceAll(path, "${SIZE}", strconv.FormatUint(size, 10))
	}
	if cut.hasLevel {
		path = strings.ReplaceAll(path, "${LEVEL}", level.Level())
	}
	return path,flag
}
func (cut *cut) getOut(path string) (file *os.File, err error) {
	if cut.path == path {
		if file != nil {
			return file, nil
		}
	}
	ii := filepath.Dir(path)
	err = os.MkdirAll(ii, 0777)
	if err != nil {
		return nil, err
	}
	file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	cut.file = file
	cut.path = path
	return
}
func getTimeCut(t string) *TimeCut {
	var timeCut TimeCut
	timeCut.timeFormat = t
	if strings.Index(t, "05") > -1 {
		timeCut.timeCutFlag = Second
	}
	if strings.Index(t, "04") > -1 {
		timeCut.timeCutFlag = Minute
	}
	if strings.Index(t, "15") > -1 {
		timeCut.timeCutFlag = Hour
	}
	if strings.Index(t, "02") > -1 {
		timeCut.timeCutFlag = Day
	}
	if strings.Index(t, "01") > -1 {
		timeCut.timeCutFlag = Month
	}
	if strings.Index(t, "2006") > -1 {
		timeCut.timeCutFlag = Year
	}
	return &timeCut
}

func parse(filePattern string) (*cut, error) {
	var cut cut
	fileCut := SplitPath(filePattern)
	var buffer = new(bytes.Buffer)
	for _, v := range fileCut {
		if strings.HasPrefix(v, "${time:") {
			end := len(v)
			time := v[7 : end-1]
			cut.timeCut = getTimeCut(time)
			buffer.WriteString("${TIME}")
		} else if strings.HasPrefix(v, "${line:") {
			end := len(v)
			line := v[7 : end-1]
			num, err := strconv.ParseUint(line, 10, 64)
			if err != nil {
				return nil, err
			}
			cut.line = num
			buffer.WriteString("${LINE}")
		} else if strings.HasPrefix(v, "${size:") {
			end := len(v)
			size := v[7 : end-1]
			num, err := toNumber(size)
			if err != nil {
				return nil, err
			}
			cut.size = num
			buffer.WriteString("${SIZE}")
		} else if strings.HasPrefix(v, "${level") {
			cut.hasLevel = true
			buffer.WriteString("${LEVEL}")
		} else {
			buffer.WriteString(v)
		}
	}
	cut.path = buffer.String()
	return &cut, nil
}
