package properties

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/pkg/survey"
)

type RateLimiter struct {
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

func (r *RateLimiter) Render() (map[string]string, error) {
	confirmVal := false
	err := survey.NewConfirm().
		SetKind("confirm").
		SetName("add-rate-limiter-middleware").
		SetMessage("Would you like to set a rate limiting middleware").
		SetDefault("false").
		SetHelp("Add a rate limit middleware properties").
		SetRequired(true).
		Render(&confirmVal)
	if err != nil {
		return nil, err
	}
	if !confirmVal {
		return nil, nil
	}
	val := 0
	err = survey.NewInput().
		SetKind("int").
		SetName("rate-limiter").
		SetMessage("Set rate request per second limiting").
		SetDefault("100").
		SetHelp("Sets how many request per second to limit").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return nil, err
	}

	return map[string]string{"rate_per_second": fmt.Sprintf("%d", val)}, nil
}
