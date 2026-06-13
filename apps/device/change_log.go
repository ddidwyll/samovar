package device

import (
  "samovar/lib/inter"
	"samovar/lib/change"
	"samovar/lib/i"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
	"fmt"
	"os"
	"strings"
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
  inter.RegisterActor(l, "[device.change_log]")
	return l.prepareFile()
}

func buildFileName() string {
	d := time.Now()
	return fmt.Sprintf("device_%d_%d.log.tsv", d.Year(), d.YearDay())
}

func (l *changeLog) prepareFile() error {
	fileName := buildFileName()
	flags := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	file, err := os.OpenFile(fileName, flags, 0644)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	l.logFile = file
	if stat.Size() == 0 {
		return l.writeLine("field", "value", "time")
	}
	return nil
}

func (l *changeLog) writeLine(cols ...string) error {
	line := strings.Join(cols, "\t")
	_, err := l.logFile.WriteString(line + "\n")
	return err
}

func (l *changeLog) logChange(report change.Report) error {
	field := report.FormatField()
	value := report.FormatValue()
	time := report.FormatTime()
	return l.writeLine(field, value, time)
}

func (l *changeLog) HandleMessage(_ gen.PID, msg any) error {
	if report, ok := msg.(change.Report); ok {
		if i.N(
			report.Key,
			"term_d",
			"term_c",
			"term_k",
			"power_m",
			"flag_otb",
			"work",
			"otbor",
			"term_c_max",
			"term_c_min",
			"otbor_t",
			"num_error",
			"count_vent",
			"time_stop",
			"otbor_minus",
		) {
			return l.logChange(report)
		} else {
			return nil
		}
	} else {
		err := fmt.Sprintf("device.changeLog reveice unexpected message: %#v", msg)
		return errors.New(err)
	}
}

func (l *changeLog) Terminate(_ error) {
	l.logFile.Close()
}
