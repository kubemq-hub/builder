package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"strconv"
	"strings"
)

const volumeTml = `
<red>volume:</>
  <yellow>size:</> <white>{{ .Size }}</>
  <yellow>storageClass:</> <white>{{ .StorageClass}}</>
`

type Volume struct {
	Size           string `json:"size"`
	StorageClass   string `json:"storage_class"`
	storageClasses []string
}

func NewVolume() *Volume {
	return &Volume{}
}
func (v *Volume) Clone() *Volume {
	return &Volume{
		Size:           v.Size,
		StorageClass:   v.StorageClass,
		storageClasses: v.storageClasses,
	}
}
func (v *Volume) Validate() error {
	if err := v.checkSize(v.Size); err != nil {
		return err
	}
	if v.StorageClass == "" {
		return fmt.Errorf("volume storage class must be defined")
	}
	return nil
}
func (v *Volume) SetStorageClasses(value []string) *Volume {
	v.storageClasses = value
	return v
}
func (v *Volume) checkSize(val interface{}) error {
	str, _ := val.(string)
	if !strings.HasSuffix(str, "Gi") {
		return fmt.Errorf("size should be in form of xxGi")
	}
	intVal, err := strconv.Atoi(strings.Trim(str, "Gi"))
	if err != nil {
		return fmt.Errorf("size should be an integer value")
	}
	if intVal < 1 {
		return fmt.Errorf("size should be at least 1Gi")
	}
	return nil
}
func (v *Volume) askSize() error {
	err := survey.NewString().
		SetKind("string").
		SetName("size").
		SetMessage("Set persistence volume size").
		SetDefault("30Gi").
		SetHelp("Set persistence volume size for each node").
		SetRequired(true).
		SetValidator(v.checkSize).
		Render(&v.Size)
	if err != nil {
		return err
	}
	return nil
}
func (v *Volume) askStorageClass() error {
	err := survey.NewString().
		SetKind("string").
		SetName("storage-class").
		SetMessage("Set persistence volume storage class").
		SetDefault("default").
		SetOptions(v.storageClasses).
		SetHelp("Set persistence volume size storage class").
		SetRequired(false).
		Render(&v.StorageClass)
	if err != nil {
		return err
	}
	if v.StorageClass == "" {
		v.StorageClass = "default"
	}
	return nil
}

func (v *Volume) Render() (*Volume, error) {
	if len(v.storageClasses) == 0 {
		v.storageClasses = append(v.storageClasses, "default", "Other")
	}
	if err := v.askSize(); err != nil {
		return nil, err
	}
	if err := v.askStorageClass(); err != nil {
		return nil, err
	}
	return v, nil
}
func (v *Volume) ColoredYaml() (string, error) {
	t := NewTemplate(volumeTml, v)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring volume spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewVolume()
