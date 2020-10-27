package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
)

const healthTml = `
<red>health:</>
  <yellow>enabled:</> <white>{{ .Enable }}</>  
  <yellow>initialDelaySeconds:</> <white>{{ .InitialDelaySeconds}}</>
  <yellow>periodSeconds:</> <white>{{ .PeriodSeconds}}</>
  <yellow>timeoutSeconds:</> <white>{{ .TimeoutSeconds}}</>
  <yellow>successThreshold:</> <white>{{ .SuccessThreshold}}</>
  <yellow>failureThreshold:</> <white>{{ .FailureThreshold}}</>
`

type Health struct {
	Enabled             bool `json:"enabled"`
	InitialDelaySeconds int  `json:"initial_delay_seconds"`
	PeriodSeconds       int  `json:"period_seconds"`
	TimeoutSeconds      int  `json:"timeout_seconds"`
	SuccessThreshold    int  `json:"success_threshold"`
	FailureThreshold    int  `json:"failure_threshold"`
}

func NewHealth() *Health {
	return &Health{}
}
func (h *Health) Clone() *Health {
	return &Health{
		Enabled:             h.Enabled,
		InitialDelaySeconds: h.InitialDelaySeconds,
		PeriodSeconds:       h.PeriodSeconds,
		TimeoutSeconds:      h.TimeoutSeconds,
		SuccessThreshold:    h.SuccessThreshold,
		FailureThreshold:    h.FailureThreshold,
	}
}
func (h *Health) askInitialDelaySeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("initial-delay-seconds").
		SetMessage("Set initial delay in seconds").
		SetDefault("10").
		SetHelp("Set initial delay in seconds health point").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&h.InitialDelaySeconds)
	if err != nil {
		return err
	}
	return nil
}
func (h *Health) askPeriodSeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("period-seconds").
		SetMessage("Set health check period in seconds").
		SetDefault("5").
		SetHelp("Set health check period in seconds").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&h.PeriodSeconds)
	if err != nil {
		return err
	}
	return nil
}
func (h *Health) askTimeoutSeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("timeout-seconds").
		SetMessage("Set health check timeout in seconds").
		SetDefault("5").
		SetHelp("Set health check timeout in seconds").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&h.TimeoutSeconds)
	if err != nil {
		return err
	}
	return nil
}

func (h *Health) askSuccessThreshold() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("success-threshold").
		SetMessage("Set health check success threshold").
		SetDefault("1").
		SetHelp("Set health check health check success threshold").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&h.SuccessThreshold)
	if err != nil {
		return err
	}
	return nil
}
func (h *Health) askFailureThreshold() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("failure threshold").
		SetMessage("Set health check failure threshold").
		SetDefault("3").
		SetHelp("Set health check health check success threshold").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&h.FailureThreshold)
	if err != nil {
		return err
	}
	return nil
}
func (h *Health) Validate() error {
	if h.InitialDelaySeconds < 0 {
		return fmt.Errorf("initial delay seconds cannot be less than 0")
	}
	if h.PeriodSeconds < 0 {
		return fmt.Errorf("period seconds cannot be less than 0")
	}

	if h.TimeoutSeconds < 0 {
		return fmt.Errorf("timeout seconds cannot be less than 0")
	}
	if h.SuccessThreshold < 0 {
		return fmt.Errorf("success threshold cannot be less than 0")
	}
	if h.FailureThreshold < 0 {
		return fmt.Errorf("failure threshold cannot be less than 0")
	}
	return nil
}

func (h *Health) Render() (*Health, error) {
	h.Enabled = true

	if err := h.askInitialDelaySeconds(); err != nil {
		return nil, err
	}

	if err := h.askPeriodSeconds(); err != nil {
		return nil, err
	}

	if err := h.askTimeoutSeconds(); err != nil {
		return nil, err
	}

	if err := h.askSuccessThreshold(); err != nil {
		return nil, err
	}

	if err := h.askFailureThreshold(); err != nil {
		return nil, err
	}
	return h, nil
}
func (h *Health) ColoredYaml() (string, error) {
	t := NewTemplate(healthTml, h)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rending health spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewHealth()
