package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

const logTml = `
<red>log:</>
  <yellow>level:</> <white>{{ .Level}}</>
`

type Log struct {
	Level       int `json:"level"`
	levelString string
}

func NewLog() *Log {
	return &Log{}
}
func (l *Log) Clone() *Log {
	return &Log{
		Level:       l.Level,
		levelString: l.levelString,
	}
}
func (l *Log) askLog() error {
	err := survey.NewString().
		SetKind("string").
		SetName("level string").
		SetMessage("Set log level printing").
		SetOptions([]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}).
		SetDefault("Info").
		SetHelp("Set log level printing").
		SetRequired(true).
		Render(&l.levelString)
	if err != nil {
		return err
	}
	return nil
}

func (l *Log) Validate() error {
	if l.Level < 0 || l.Level > 5 {
		return fmt.Errorf("invalid log level")
	}
	return nil
}
func (l *Log) Render() (*Log, error) {
	if err := l.askLog(); err != nil {
		return nil, err
	}
	switch l.levelString {
	case "Trace":
		l.Level = 0
	case "Debug":
		l.Level = 1
	case "Info":
		l.Level = 2
	case "Warn":
		l.Level = 3
	case "Error":
		l.Level = 4
	case "Fatal":
		l.Level = 5
	default:
		l.Level = -1
	}
	return l, nil
}
func (l *Log) ColoredYaml() (string, error) {
	t := NewTemplate(logTml, l)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring log spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewLog()
