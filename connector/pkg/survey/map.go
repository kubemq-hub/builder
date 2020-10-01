package survey

import "fmt"

type Map struct {
	*KindMeta
	*ObjectMeta
	Fields []*Input `json:"fields"`
}

func (m *Map) complete() error {
	return nil
}
func (m *Map) Render() (map[string]string, error) {
	results := make(map[string]string)
	for i, field := range m.Fields {
		if err := field.Complete(); err != nil {
			return nil, fmt.Errorf("error on field input %d: %s", i, err.Error())
		}
		switch field.Kind {
		case "string":
			result := ""
			err := field.Render(&result)
			if err != nil {
				return nil, err
			}
			results[field.Name] = result
		case "int":
			result := 0
			err := field.Render(&result)
			if err != nil {
				return nil, err
			}
			results[field.Name] = fmt.Sprintf("%d", result)
		}
	}
	return results, nil
}
