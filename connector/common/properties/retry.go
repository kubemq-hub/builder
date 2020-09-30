package properties

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/common/survey"
	"math"
)

type Retry struct {
	values map[string]string
}

func NewRetry() *Retry {
	return &Retry{
		values: map[string]string{},
	}
}
func (r *Retry) askDelayType() error {
	val := ""
	err := survey.NewInput().
		SetKind("string").
		SetName("retry_delay_type").
		SetMessage("Sets retry delay type").
		SetOptions([]string{"fixed", "back-off", "random"}).
		SetDefault("fixed").
		SetHelp("Sets retry delay type").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_delay_type"] = val
	return nil
}
func (r *Retry) askAttempts() error {
	val := 0
	err := survey.NewInput().
		SetKind("int").
		SetName("retry_attempts").
		SetMessage("Sets retry max attempts type").
		SetDefault("1").
		SetHelp("Sets retry max attempts type").
		SetRequired(true).
		SetMin(1).
		SetMax(1024).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_attempts"] = fmt.Sprintf("%d", val)
	return nil
}
func (r *Retry) askDelayMillisecond() error {
	val := 0
	err := survey.NewInput().
		SetKind("int").
		SetName("retry_delay_milliseconds").
		SetMessage("Sets retry delay milliseconds").
		SetDefault("100").
		SetHelp("Sets retry delay milliseconds").
		SetRequired(true).
		SetMin(0).
		SetMax(math.MaxInt32).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_delay_milliseconds"] = fmt.Sprintf("%d", val)
	return nil
}
func (r *Retry) askDelayJitter() error {
	val := 0
	err := survey.NewInput().
		SetKind("int").
		SetName("retry_max_jitter_milliseconds").
		SetMessage("Sets retry delay milliseconds jitter").
		SetDefault("100").
		SetHelp("Sets retry delay milliseconds jitter").
		SetRequired(true).
		SetMin(1).
		SetMax(math.MaxInt32).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_max_jitter_milliseconds"] = fmt.Sprintf("%d", val)
	return nil
}
func (r *Retry) Render() (map[string]string, error) {
	confirmVal := false
	err := survey.NewConfirm().
		SetKind("confirm").
		SetName("add-retry-middleware").
		SetMessage("Would you like to set a request retries middleware").
		SetDefault("false").
		SetHelp("Add a request retries middleware properties").
		SetRequired(true).
		Render(&confirmVal)
	if err != nil {
		return nil, err
	}
	if !confirmVal {
		return nil, nil
	}
	if err := r.askDelayType(); err != nil {
		return nil, err
	}
	if err := r.askAttempts(); err != nil {
		return nil, err
	}
	if err := r.askDelayMillisecond(); err != nil {
		return nil, err
	}
	if err := r.askDelayJitter(); err != nil {
		return nil, err
	}
	return r.values, nil
}
