package source

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Source struct {
	Name           string              `json:"name"`
	Kind           string              `json:"kind"`
	Connections    []map[string]string `json:"connections"`
	ConnectionSpec string              `json:"-" yaml:"-"`
	addressOptions []string
	takenNames     []string
	defaultName    string
	isEdit         bool
}

func NewSource(defaultName string) *Source {
	return &Source{
		addressOptions: nil,
		defaultName:    defaultName,
	}
}
func (s *Source) Clone() *Source {
	newSrc := &Source{
		Name:           s.Name,
		Kind:           s.Kind,
		Connections:    []map[string]string{},
		ConnectionSpec: "",
		addressOptions: nil,
		takenNames:     nil,
		defaultName:    "",
	}
	for _, connection := range s.Connections {
		newConnection := map[string]string{}
		for Key, val := range connection {
			newConnection[Key] = val
		}
		newSrc.Connections = append(newSrc.Connections, newConnection)
	}
	return newSrc
}

func (s *Source) SetAddress(value []string) *Source {
	s.addressOptions = value
	return s
}
func (s *Source) SetTakenNames(value []string) *Source {
	s.takenNames = value
	return s
}
func (s *Source) SetIsEdit(value bool) *Source {
	s.isEdit = value
	return s
}
func (s *Source) askAddConnection() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-connection").
		SetMessage("Would you like to add another source connection").
		SetDefault("false").
		SetHelp("Add new source connection").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false, err
	}
	return val, nil
}

func (s *Source) addConnection() error {
	if connection, err := NewConnection().
		SetAddress(s.addressOptions).
		Render(s.Name, s.Kind); err != nil {
		return err
	} else {
		s.Connections = append(s.Connections, connection)
	}
	return nil
}

func (s *Source) Render() (*Source, error) {
	defaultName := ""
	if s.isEdit {
		defaultName = s.Name
	} else {
		defaultName = s.defaultName
	}
	var err error
	if s.Name, err = NewName(defaultName).
		SetTakenNames(s.takenNames).
		Render(); err != nil {
		return nil, err
	}
	defaultKind := ""
	if s.isEdit {
		defaultKind = s.Kind
	} else {
		defaultKind = ""
	}
	if s.Kind, err = NewKind(defaultKind).
		Render(); err != nil {
		return nil, err
	}
	utils.Println(promptSourceFirstConnection, s.Kind)
	err = s.addConnection()
	if err != nil {
		return nil, err
	}
	for {
		addMore, err := s.askAddConnection()
		if err != nil {
			return s, nil
		}
		if addMore {
			err = s.addConnection()
			if err != nil {
				return nil, err
			}
		} else {
			goto done
		}
	}
done:
	return s, nil

}

func (s *Source) String() string {
	s.ConnectionSpec = utils.MapArrayToYaml(s.Connections)
	t := utils.NewTemplate(sourceTemplate, s)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring source  spec,%s", err.Error())
	}
	return string(b)
}
