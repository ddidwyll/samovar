package device

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"fmt"
	"os"
	"time"
)

type changeLog struct {
	act.Actor
	logFile *os.File
}

func newChangeLog() gen.ProcessBehavior {
	return &changeLog{}
}

func (l *changeLog) Init(_ ...any) error {
	l.Log().Debug("device.changeLog started (%s)", l.Name())
	return l.prepareFile()
}

func buildFileName() string {
	d := time.Now()
	return fmt.Sprintf("device_%d_%d.log.tsv", d.Year(), d.YearDay())
}

func (l *changeLog) prepareFile() error {
	fileName := buildFileName()
	flags := os.O_APPEND|os.O_WRONLY|os.O_CREATE
	file, err := os.OpenFile(fileName, flags, 0644)
	if err != nil {
  	return err
	}
	stat, err := file.Stat()
	if err != nil {
  	return err
	}
	if stat.Size() == 0 {
  	_, err = file.WriteString("key\tval\tts\n")
  	if err != nil {
    	return err
  	}
	}
	l.logFile = file
  return nil
}

func (l *changeLog) Terminate(_ error) {
  l.logFile.Close()
}
