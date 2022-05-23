package log

import (
	"os"
	"time"
)

type WriteFile struct {
	fileCut *cut
	t       *time.Time
	file *os.File
	path string

	line uint64
	size uint64


}

func NewWriteFile(fileCut *cut) *WriteFile {
	var wf WriteFile
	wf.fileCut = fileCut
	wf.line  = 0
	wf.size = 0
	return &wf
}

func (wf *WriteFile)WriteLog(entry *Entry)error{
	num,err:=wf.file.WriteString(entry.GetLog())
	wf.size = wf.size+uint64(num)
	wf.line = wf.line+1
	return err
}


func (wf *WriteFile) fileTo(t *time.Time, level Level) error {
	path,flag:=wf.fileCut.getPath(t, wf.line, wf.size, level)
	if flag{
		wf.size = 0
		wf.line = 0
	}
	file, err := wf.fileCut.getOut(path)
	if err != nil {
		return  err
	}
	wf.file = file
	return nil
}
