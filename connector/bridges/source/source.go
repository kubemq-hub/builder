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
	WasEdited      bool                `json:"-" yaml:"-"`
	defaultName    string
	isEdit         bool
}

func NewSource(defaultName string) *Source {
	return &Source{
		defaultName: defaultName,
	}
}
func (s *Source) Clone() *Source {
	newSrc := &Source{
		Name:           s.Name,
		Kind:           s.Kind,
		Connections:    []map[string]string{},
		ConnectionSpec: s.ConnectionSpec,
		WasEdited:      s.WasEdited,
		defaultName:    s.Name,
		isEdit:         s.isEdit,
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
		Render(s.Name, s.Kind); err != nil {
		return err
	} else {
		s.Connections = append(s.Connections, connection)
	}
	return nil
}
func (s *Source) add() (*Source, error) {
	var err error
	if s.Name, err = NewName(s.defaultName).
		Render(); err != nil {
		return nil, err
	}
	if s.Kind, err = NewKind("").
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

func (s *Source) editKind() (bool, error) {
	var err error
	current := s.Kind
	if s.Kind, err = NewKind(s.Kind).
		Render(); err != nil {
		return false, err
	}
	return s.Kind != current, nil
}
func (s *Source) editConnections() error {
	s.Connections = []map[string]string{}
	utils.Println(promptSourceFirstConnection, s.Kind)
	err := s.addConnection()
	if err != nil {
		return err
	}
	for {
		addMore, err := s.askAddConnection()
		if err != nil {
			return nil
		}
		if addMore {
			err = s.addConnection()
			if err != nil {
				return err
			}
		} else {
			goto done
		}
	}
done:
	s.WasEdited = true
	return nil
}

func (s *Source) edit() (*Source, error) {
	var result *Source
	edited := s.Clone()
	form := survey.NewForm(fmt.Sprintf("Select Edit %s Source Option", edited.Name))

	ftName := new(string)
	*ftName = fmt.Sprintf("<n> Edit Source Name (%s)", edited.Name)
	form.AddItem(ftName, func() error {
		var err error
		if edited.Name, err = NewName(edited.Name).
			Render(); err != nil {
			return err
		}
		*ftName = fmt.Sprintf("<n> Edit Source Name (%s)", edited.Name)
		return nil
	})

	ftKind := new(string)
	*ftKind = fmt.Sprintf("<k> Edit Source Kind (%s)", edited.Kind)
	ftConnections := new(string)
	*ftConnections = fmt.Sprintf("<c> Edit Source Connections (%s)", edited.Kind)

	form.AddItem(ftKind, func() error {
		if changed, err := s.editKind(); err != nil {
			return err
		} else {
			if changed {
				if err := s.editConnections(); err != nil {
					return err
				}
				s.WasEdited = true
			}
		}
		*ftKind = fmt.Sprintf("<k> Edit Source Kind (%s)", edited.Kind)
		*ftConnections = fmt.Sprintf("<c> Edit Source Connections (%s)", edited.Kind)
		return nil
	})

	form.AddItem(ftConnections, func() error {
		if err := s.editConnections(); err != nil {
			return err
		}
		*ftConnections = fmt.Sprintf("<c> Edit Source Connections (%s)", edited.Kind)
		return nil
	})

	form.AddItem("Show Source Configuration", func() error {
		utils.Println(promptShowSource, edited.Name)
		utils.Println("%s\n", edited.ColoredYaml())
		return nil
	})
	form.SetOnSaveFn(func() error {
		if err := edited.Validate(); err != nil {
			return err
		}
		result = edited
		return nil
	})

	form.SetOnCancelFn(func() error {
		result = s
		return nil
	})
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
}
func (s *Source) Render() (*Source, error) {
	if s.isEdit {
		return s.edit()
	}
	return s.add()
}

func (s *Source) ColoredYaml() string {
	s.ConnectionSpec = utils.MapArrayToYaml(s.Connections)
	t := utils.NewTemplate(sourceTemplate, s)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring source  spec,%s", err.Error())
	}
	return string(b)
}
func (s *Source) TableItemShort() string {
	return fmt.Sprintf("%s/%s/%d", s.Name, s.Kind, len(s.Connections))
}

func (s *Source) Validate() error {
	return nil
}
