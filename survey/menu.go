package survey

import "fmt"

type Menu struct {
	title   string
	fnMap   map[string]func() error
	fnItems []string
}

func NewMenu(title string) *Menu {
	return &Menu{
		title:   title,
		fnMap:   map[string]func() error{},
		fnItems: []string{},
	}
}

func (m *Menu) AddItem(title string, fn func() error) *Menu {
	m.fnMap[title] = fn
	m.fnItems = append(m.fnItems, title)
	return m
}

func (m *Menu) Render() error {
	for {
		val := ""
		err := NewString().
			SetKind("string").
			SetName("menu").
			SetMessage(m.title).
			SetDefault(m.fnItems[0]).
			SetRequired(true).
			SetOptions(m.fnItems).
			Render(&val)
		if err != nil {
			return err
		}
		fn, ok := m.fnMap[val]
		if !ok {
			return fmt.Errorf("menu function for %s not found", val)
		}
		if fn == nil {
			return nil
		}
		if err := fn(); err != nil {
			return fn()
		}
	}
}
