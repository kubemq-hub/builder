package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
)

type Queue struct {
	MaxReceiveMessagesRequest int `json:"max_receive_messages_request"`
	MaxWaitTimeoutSeconds     int `json:"max_wait_timeout_seconds"`
	MaxExpirationSeconds      int `json:"max_expiration_seconds"`
	MaxDelaySeconds           int `json:"max_delay_seconds"`
	MaxReQueues               int `json:"max_re_queues"`
	MaxVisibilitySeconds      int `json:"max_visibility_seconds"`
	DefaultVisibilitySeconds  int `json:"default_visibility_seconds"`
	DefaultWaitTimeoutSeconds int `json:"default_wait_timeout_seconds"`
}

func NewQueue() *Queue {
	return &Queue{}
}
func (q *Queue) Clone() *Queue {
	return &Queue{
		MaxReceiveMessagesRequest: q.MaxReceiveMessagesRequest,
		MaxWaitTimeoutSeconds:     q.MaxWaitTimeoutSeconds,
		MaxExpirationSeconds:      q.MaxExpirationSeconds,
		MaxDelaySeconds:           q.MaxDelaySeconds,
		MaxReQueues:               q.MaxReQueues,
		MaxVisibilitySeconds:      q.MaxVisibilitySeconds,
		DefaultVisibilitySeconds:  q.DefaultVisibilitySeconds,
		DefaultWaitTimeoutSeconds: q.DefaultWaitTimeoutSeconds,
	}
}
func (q *Queue) askMaxReceiveMessagesRequest() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("MaxReceiveMessagesRequest").
		SetMessage("Set max of sending / receiving batch of queue message").
		SetDefault("1024").
		SetHelp("Set max of sending / receiving batch of queue message").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&q.MaxReceiveMessagesRequest)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) askMaxWaitTimeoutSeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("MaxWaitTimeoutSeconds").
		SetMessage("Set max wait timeout allowed for message").
		SetDefault("3600").
		SetHelp("Set max wait timeout allowed for message").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&q.MaxWaitTimeoutSeconds)
	if err != nil {
		return err
	}
	return nil
}
func (q *Queue) askMaxExpirationSeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("MaxWaitTimeoutSeconds").
		SetMessage("Set max expiration allowed for message").
		SetDefault("43200").
		SetHelp("Set max expiration allowed for message").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&q.MaxExpirationSeconds)
	if err != nil {
		return err
	}
	return nil
}
func (q *Queue) asMaxDelaySeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("MaxDelaySeconds").
		SetMessage("Set max delay seconds allowed for message").
		SetDefault("43200").
		SetHelp("Set max delay seconds allowed for message").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&q.MaxDelaySeconds)
	if err != nil {
		return err
	}
	return nil
}
func (q *Queue) askMaxReQueues() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("MaxReQueues").
		SetMessage("Set max retires to receive message before discard").
		SetDefault("1024").
		SetHelp("Set max retires to receive message before discard").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&q.MaxReQueues)
	if err != nil {
		return err
	}
	return nil
}
func (q *Queue) askMaxVisibilitySeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("MaxVisibilitySeconds").
		SetMessage("Set max time of hold received message before returning to queue").
		SetDefault("43200").
		SetHelp("Set max time of hold received message before returning to queue").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&q.MaxVisibilitySeconds)
	if err != nil {
		return err
	}
	return nil
}
func (q *Queue) askDefaultVisibilitySeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("DefaultVisibilitySeconds").
		SetMessage("Set default time of hold received message before returning to queue").
		SetDefault("60").
		SetHelp("Set default time of hold received message before returning to queue").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&q.DefaultVisibilitySeconds)
	if err != nil {
		return err
	}
	return nil
}
func (q *Queue) askDefaultWaitTimeoutSeconds() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("DefaultWaitTimeoutSeconds").
		SetMessage("Set default time to wait for a message in a queue").
		SetDefault("1").
		SetHelp("Set default time to wait for a message in a queue").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&q.DefaultWaitTimeoutSeconds)
	if err != nil {
		return err
	}
	return nil
}
func (q *Queue) Validate() error {

	if q.MaxReceiveMessagesRequest < 0 {
		return fmt.Errorf("max receive messages request cannot be less than 0")
	}
	if q.MaxWaitTimeoutSeconds < 0 {
		return fmt.Errorf("max wait timeout seconds cannot be less than 0")
	}

	if q.MaxExpirationSeconds < 0 {
		return fmt.Errorf("max expiration seconds cannot be less than 0")
	}
	if q.MaxDelaySeconds < 0 {
		return fmt.Errorf("max delay seconds cannot be less than 0")
	}
	if q.MaxReQueues < 0 {
		return fmt.Errorf("max re-queues cannot be less than 0")
	}
	if q.MaxVisibilitySeconds < 0 {
		return fmt.Errorf("max visibility seconds cannot be less than 0")
	}
	if q.DefaultVisibilitySeconds < 0 {
		return fmt.Errorf("default visibility seconds cannot be less than 0")
	}
	if q.DefaultWaitTimeoutSeconds < 0 {
		return fmt.Errorf("default wait timeout seconds cannot be less than 0")
	}
	return nil
}

func (q *Queue) Render() (*Queue, error) {

	if err := q.askMaxReceiveMessagesRequest(); err != nil {
		return nil, err
	}
	if err := q.askMaxWaitTimeoutSeconds(); err != nil {
		return nil, err
	}
	if err := q.askMaxExpirationSeconds(); err != nil {
		return nil, err
	}
	if err := q.asMaxDelaySeconds(); err != nil {
		return nil, err
	}
	if err := q.askMaxReQueues(); err != nil {
		return nil, err
	}
	if err := q.askMaxVisibilitySeconds(); err != nil {
		return nil, err
	}
	if err := q.askDefaultVisibilitySeconds(); err != nil {
		return nil, err
	}
	if err := q.askDefaultWaitTimeoutSeconds(); err != nil {
		return nil, err
	}
	return q, nil
}

var _ Validator = NewQueue()
