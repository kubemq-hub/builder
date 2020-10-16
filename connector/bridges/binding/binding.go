package binding

import (
	"fmt"
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/builder/connector/bridges/source"
	"github.com/kubemq-hub/builder/connector/bridges/target"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

const (
	promptSource = `<cyan>In the next steps, we will configure the Source connection.
We will set:</>
<yellow>Name -</> A unique name for the Source's binding
<yellow>Kind -</> A Source connection type 
<yellow>Connections -</> A list of connections properties based on the selected kind`

	promptTarget = `<cyan>In the next steps, we will configure the Target connection.
We will set:</>
<yellow>Name -</> A unique name for the Source's binding
<yellow>Kind -</> A Source connection type 
<yellow>Connections -</> A list of connections properties based on the selected kind`
)

const bindingsTml = `
<red>name:</> {{.Name}}
{{- .SourcesSpec -}}
{{- .TargetSpec -}}
{{- .PropertiesSpec -}}
`

type Binding struct {
	Name              string            `json:"name"`
	Sources           *source.Source    `json:"sources"`
	Targets           *target.Target    `json:"targets"`
	Properties        map[string]string `json:"properties"`
	SourcesSpec       string
	TargetSpec        string
	PropertiesSpec    string
	addressOptions    []string
	takenSourceNames  []string
	takenTargetsNames []string
	takenBindingNames []string
}

func NewBinding() *Binding {
	return &Binding{}
}
func (b *Binding) SetAddress(value []string) *Binding {
	b.addressOptions = value
	return b
}
func (b *Binding) SetTakenSourceNames(value []string) *Binding {
	b.takenSourceNames = value
	return b
}
func (b *Binding) SetTakenTargetsNames(value []string) *Binding {
	b.takenTargetsNames = value
	return b
}
func (b *Binding) SetTakenBindingNames(value []string) *Binding {
	b.takenBindingNames = value
	return b
}
func (b *Binding) SourceName() string {
	if b.Sources != nil {
		return b.Sources.Name
	}
	return ""
}
func (b *Binding) TargetName() string {
	if b.Targets != nil {
		return b.Targets.Name
	}
	return ""
}
func (b *Binding) BindingName() string {
	return b.Name
}
func (b *Binding) confirmSource() bool {
	utils.Println(fmt.Sprintf("<cyan>Here is Binding Source configuration:</>%s", b.Sources.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println("<cyan>Lets reconfigure Binding Source:</>")
	}
	return val
}
func (b *Binding) confirmTarget() bool {
	utils.Println(fmt.Sprintf("<cyan>Here is Binding Target configuration:</>%s", b.Targets.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println("<cyan>Lets reconfigure Binding Target:</>")
	}
	return val
}
func (b *Binding) confirmProperties(p *common.Properties) bool {
	utils.Println(fmt.Sprintf("<cyan>Here is Binding Middleware Properties configuration:</>%s", p.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println("<cyan>Lets reconfigure Middleware Properties:</>")
	}
	return val
}
func (b *Binding) Render() (*Binding, error) {
	var err error
	if b.Name, err = NewName().
		SetTakenNames(b.takenBindingNames).
		Render(); err != nil {
		return nil, err
	}
	utils.Println(promptSource)
	utils.Println("<cyan>Lets Set Source Configuration:</>")
	for {
		if b.Sources, err = source.NewSource().
			SetAddress(b.addressOptions).
			SetTakenNames(b.takenSourceNames).
			Render(); err != nil {
			return nil, err
		}
		ok := b.confirmSource()
		if ok {
			b.SourcesSpec = b.Sources.String()
			break
		}
	}
	utils.Println(promptTarget)
	utils.Println("<cyan>Lets Set Target Configuration:</>")
	for {
		if b.Targets, err = target.NewTarget().
			SetAddress(b.addressOptions).
			SetTakenNames(b.takenTargetsNames).
			Render(); err != nil {
			return nil, err
		}
		ok := b.confirmTarget()
		if ok {
			b.TargetSpec = b.Targets.String()
			break
		}
	}
	utils.Println("<cyan>We have completed Source and Target Configuration</>")
	for {
		p := common.NewProperties()
		if b.Properties, err = p.
			Render(); err != nil {
			return nil, err
		}
		ok := b.confirmProperties(p)
		if ok {
			b.PropertiesSpec = p.String()
			break
		}

	}
	return b, nil
}

func (b *Binding) String() string {
	tpl := utils.NewTemplate(bindingsTml, b)
	bnd, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring binding spec,%s", err.Error())
	}
	return string(bnd)
}
