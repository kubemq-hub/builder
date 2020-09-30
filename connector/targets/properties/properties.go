package properties

import "github.com/kubemq-hub/builder/connector/common/survey"

type Properties struct {
	values map[string]string
}

func NewProperties() *Properties {
	return &Properties{
		values: map[string]string{},
	}
}

func (p *Properties) Render() (map[string]string, error) {
	confirmVal := false
	err := survey.NewConfirm().
		SetKind("confirm").
		SetName("add-middleware").
		SetMessage("Would you to add middlewares to this binding").
		SetDefault("false").
		SetHelp("Add a middleware properties").
		SetRequired(true).
		Render(&confirmVal)
	if err != nil {
		return nil, err
	}
	if !confirmVal {
		return nil, nil
	}
	if values, err := NewLog().Render(); err != nil {
		return nil, err
	} else {
		for key, val := range values {
			p.values[key] = val
		}
	}
	if values, err := NewRateLimiter().Render(); err != nil {
		return nil, err
	} else {
		for key, val := range values {
			p.values[key] = val
		}
	}
	if values, err := NewRetry().Render(); err != nil {
		return nil, err
	} else {
		for key, val := range values {
			p.values[key] = val
		}
	}
	return p.values, nil
}
