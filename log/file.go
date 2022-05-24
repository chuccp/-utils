package log

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type WriteFile struct {
	fileCut *cut
	t       *time.Time
	file    *os.File
	path    string
	line    uint64
	size    uint64
}

func NewWriteFile(fileCut *cut) *WriteFile {
	var wf WriteFile
	wf.fileCut = fileCut
	wf.line = 0
	wf.size = 0
	return &wf
}
func NewWrite(path string) (*WriteFile,error) {
	var wf WriteFile
	wf.line = 0
	wf.size = 0
	_,err:=wf.getOut(path)
	return &wf,err
}

func (wf *WriteFile) WriteLog(entry *Entry) error {
	num, err := wf.file.WriteString(entry.GetLog())
	wf.size = wf.size + uint64(num)
	wf.line = wf.line + 1
	return err
}

func (wf *WriteFile) getPath(time *time.Time, level Level) (path string, flag bool) {
	path = wf.fileCut.path
	if wf.fileCut.timeCut != nil {
		ti := wf.fileCut.timeCut.parse(time)
		if ti == wf.fileCut.ctime {
			flag = false
		} else {
			wf.fileCut.ctime = ti
			flag = true
		}
		path = strings.ReplaceAll(path, "${TIME}", ti)
	}
	if wf.fileCut.line > 0 {
		line := (wf.line / wf.fileCut.line) * wf.fileCut.line
		path = strings.ReplaceAll(path, "${LINE}", strconv.FormatUint(line, 10))
	}
	if wf.fileCut.size > 0 {
		size := (wf.size / wf.fileCut.size) * wf.fileCut.size
		path = strings.ReplaceAll(path, "${SIZE}", strconv.FormatUint(size, 10))
	}
	if wf.fileCut.hasLevel {
		path = strings.ReplaceAll(path, "${LEVEL}", level.Level())
	}
	return path, flag
}
func (wf *WriteFile) getOut(path string) (file *os.File, err error) {
	if wf.path == path {
		if file != nil {
			return file, nil
		}
	} else {
		if len(wf.path) > 0 {
			wf.file.Close()
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
	wf.file = file
	wf.path = path
	return
}

func (wf *WriteFile) fileTo(t *time.Time, level Level) error {
	path, flag := wf.getPath(t, level)
	if flag {
		wf.size = 0
		wf.line = 0
	}
	file, err := wf.getOut(path)
	if err != nil {
		return err
	}
	wf.file = file
	return nil
}
