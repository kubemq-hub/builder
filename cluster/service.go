package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
)

type Service struct {
	NodePort   int    `json:"node_port"`
	Expose     string `json:"expose"`
	BufferSize int    `json:"buffer_size"`
	BodyLimit  int    `json:"body_limit"`
	kind       string
}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Validate() error {
	switch s.Expose {
	case "ClusterIP", "LoadBalancer":
	case "NodePort":
		if s.NodePort == 0 {
			return nil
		}
		if s.NodePort < 30000 || s.NodePort > 32767 {
			return fmt.Errorf("node port value must be with the range 30000-32767")
		}
	default:
		return fmt.Errorf("invalid service type")
	}
	return nil
}
func (s *Service) SetKind(value string) *Service {
	s.kind = value
	return s
}
func (s *Service) askExpose() error {
	err := survey.NewString().
		SetKind("string").
		SetName("expose").
		SetMessage(fmt.Sprintf("Set cluster %s service type", s.kind)).
		SetDefault("ClusterIP").
		SetOptions([]string{"ClusterIP", "NodePort", "LoadBalancer"}).
		SetHelp(fmt.Sprintf("Set cluster %s service type", s.kind)).
		SetRequired(true).
		Render(&s.Expose)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) askNodePort() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("node-port").
		SetMessage(fmt.Sprintf("Set cluster %s service NodePort value", s.kind)).
		SetDefault("30000").
		SetHelp(fmt.Sprintf("Set cluster %s service NodePort value", s.kind)).
		SetRequired(false).
		SetRange(30000, 32767).
		Render(&s.NodePort)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) askBufferSize() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("buffer-size").
		SetMessage("Set subscribe message / requests buffer size to use on server").
		SetDefault("100").
		SetHelp("Set subscribe message / requests buffer size to use on server").
		SetRequired(false).
		SetRange(0, math.MaxInt32).
		Render(&s.BufferSize)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) askBodyLimit() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("body=limit").
		SetMessage("Set max size of payload in bytes (0 - no limit)").
		SetDefault("0").
		SetHelp("Set max size of payload in bytes").
		SetRequired(false).
		SetRange(0, math.MaxInt32).
		Render(&s.BodyLimit)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) Render() (*Service, error) {
	if s.kind == "" {
		return nil, fmt.Errorf("no cluster service kind was set")
	}
	if err := s.askExpose(); err != nil {
		return nil, err
	}
	if s.Expose == "NodePort" {
		if err := s.askNodePort(); err != nil {
			return nil, err
		}
	}
	switch s.kind {
	case "grpc", "rest":
		if err := s.askBufferSize(); err != nil {
			return nil, err
		}
		if err := s.askBodyLimit(); err != nil {
			return nil, err
		}
	}
	return s, nil
}

var _ Validator = NewService()
