package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
)

type Store struct {
	Clean                    bool   `json:"clean"`
	Path                     string `json:"path"`
	MaxChannels              int    `json:"max_channels"`
	MaxSubscribers           int    `json:"max_subscribers"`
	MaxMessages              int    `json:"max_messages"`
	MaxChannelSize           int    `json:"max_channel_size"`
	MessagesRetentionMinutes int    `json:"messages_retention_minutes"`
	PurgeInactiveMinutes     int    `json:"purge_inactive_minutes"`
}

func NewStore() *Store {
	return &Store{}
}
func (s *Store) Clone() *Store {
	return &Store{
		Clean:                    s.Clean,
		Path:                     s.Path,
		MaxChannels:              s.MaxChannels,
		MaxSubscribers:           s.MaxSubscribers,
		MaxMessages:              s.MaxMessages,
		MaxChannelSize:           s.MaxChannels,
		MessagesRetentionMinutes: s.MessagesRetentionMinutes,
		PurgeInactiveMinutes:     s.PurgeInactiveMinutes,
	}
}
func (s *Store) askClean() error {
	err := survey.NewBool().
		SetKind("bool").
		SetName("Clean").
		SetMessage("Set clear persistence data on start-up").
		SetDefault("false").
		SetHelp("Set clear persistence data on start-up").
		SetRequired(true).
		Render(&s.Clean)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) checkNonPathEmpty(val interface{}) error {
	str, _ := val.(string)
	if str == "" {
		return fmt.Errorf("path cannot be empty")
	}
	return nil
}
func (s *Store) askPath() error {
	err := survey.NewString().
		SetKind("string").
		SetName("path").
		SetMessage("Set persistence file path").
		SetDefault("./store").
		SetHelp("Set persistence file path").
		SetRequired(true).
		SetValidator(s.checkNonPathEmpty).
		Render(&s.Path)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) askMaxChannels() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("max channels").
		SetMessage("Set limit number of persistence channels (0 - no limit)").
		SetDefault("0").
		SetHelp("Set limit number of persistence channels").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&s.MaxChannels)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) askMaxSubscribers() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("max subscribers").
		SetMessage("Set limit of subscribers per channel (0 - no limit)").
		SetDefault("0").
		SetHelp("Set limit of subscribers per channel").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&s.MaxSubscribers)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) askMaxMessages() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("max messages").
		SetMessage("Set limit of messages per channel (0 - no limit)").
		SetDefault("0").
		SetHelp("Set limit of messages per channel").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&s.MaxMessages)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) asMaxChannelSize() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("max channel size").
		SetMessage("Set limit size of channel in bytes (0 - no limit)").
		SetDefault("0").
		SetHelp("Set limit size of channel in bytes").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&s.MaxChannelSize)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) askMessagesRetentionMinutes() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("message retention").
		SetMessage("Set message retention time in minutes (0 - no limit)").
		SetDefault("1440").
		SetHelp("Set message retention time in minutes").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&s.MessagesRetentionMinutes)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) askPurgeInactiveMinutes() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("purge inactive minutes").
		SetMessage("Set time in minutes of channel inactivity to delete (0 - no limit)").
		SetDefault("1440").
		SetHelp("Set health check health check success threshold").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&s.PurgeInactiveMinutes)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) Validate() error {
	if s.Path == "" {
		return fmt.Errorf("store path cannot be empty")
	}
	if s.MaxChannels < 0 {
		return fmt.Errorf("max channels cannot be less than 0")
	}
	if s.MaxSubscribers < 0 {
		return fmt.Errorf("max subscribers cannot be less than 0")
	}

	if s.MaxMessages < 0 {
		return fmt.Errorf("max messages cannot be less than 0")
	}
	if s.MaxChannelSize < 0 {
		return fmt.Errorf("max channel size cannot be less than 0")
	}
	if s.MessagesRetentionMinutes < 0 {
		return fmt.Errorf("messages retention minutes cannot be less than 0")
	}
	if s.PurgeInactiveMinutes < 0 {
		return fmt.Errorf("purge inactive minutes cannot be less than 0")
	}
	return nil
}

func (s *Store) Render() (*Store, error) {

	if err := s.askClean(); err != nil {
		return nil, err
	}
	if err := s.askPath(); err != nil {
		return nil, err
	}
	if err := s.askMaxChannels(); err != nil {
		return nil, err
	}
	if err := s.askMaxSubscribers(); err != nil {
		return nil, err
	}
	if err := s.askMaxMessages(); err != nil {
		return nil, err
	}
	if err := s.asMaxChannelSize(); err != nil {
		return nil, err
	}
	if err := s.askMessagesRetentionMinutes(); err != nil {
		return nil, err
	}
	if err := s.askPurgeInactiveMinutes(); err != nil {
		return nil, err
	}
	return s, nil
}

var _ Validator = NewStore()
