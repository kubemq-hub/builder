package properties

import "github.com/kubemq-hub/builder/connector/common/survey"

type Log struct {
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Render() (map[string]string, error) {
	confirmVal := false
	err := survey.NewConfirm().
		SetKind("confirm").
		SetName("add-log-middleware").
		SetMessage("Would you to  set a logging middleware").
		SetDefault("true").
		SetHelp("Add logging middleware properties").
		SetRequired(true).
		Render(&confirmVal)
	if err != nil {
		return nil, err
	}
	if !confirmVal {
		return nil, nil
	}
	val := ""
	err = survey.NewInput().
		SetKind("string").
		SetName("log-level").
		SetMessage("Set Log level").
		SetOptions([]string{"debug", "info", "error"}).
		SetDefault("info").
		SetHelp("Sets Log level printing").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return nil, err
	}
	return map[string]string{"log_level": val}, nil
}
