package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"strconv"
	"strings"
)

type Resource struct {
	LimitsCpu      string `json:"limits_cpu"`
	LimitsMemory   string `json:"limits_memory"`
	RequestsCpu    string `json:"requests_cpu"`
	RequestsMemory string `json:"requests_memory"`
}

func NewResource() *Resource {
	return &Resource{}
}
func (r *Resource) Validate() error {
	if err := r.checkCPU(r.LimitsCpu); err != nil {
		return err
	}
	if err := r.checkMem(r.LimitsMemory); err != nil {
		return err
	}
	if err := r.checkCPU(r.RequestsCpu); err != nil {
		return err
	}
	if err := r.checkMem(r.RequestsMemory); err != nil {
		return err
	}
	return nil
}

func (r *Resource) checkCPU(val interface{}) error {
	str := val.(string)
	intVal, err := strconv.Atoi(strings.Trim(str, "m"))
	if err != nil {
		return fmt.Errorf("cpu should be an integer value")
	}
	if intVal < 1 {
		return fmt.Errorf("size should be at least 1")
	}
	return nil
}
func (r *Resource) checkMem(val interface{}) error {
	str, _ := val.(string)
	if !strings.HasSuffix(str, "Gi") && !strings.HasSuffix(str, "Mi") {
		return fmt.Errorf("size should be in form of xxGi or xxMi")
	}
	strVal := strings.Trim(str, "Gi")
	strVal = strings.Trim(strVal, "Mi")
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return fmt.Errorf("memory size should be an integer value")
	}
	if intVal < 1 {
		return fmt.Errorf("memory size should be at least 1")
	}
	return nil
}
func (r *Resource) askLimitCpuSize() error {
	err := survey.NewString().
		SetKind("string").
		SetName("limit-cpu").
		SetMessage("Set Limit CPU resource").
		SetDefault("1").
		SetHelp("Set Limit CPU resource").
		SetRequired(false).
		SetValidator(r.checkCPU).
		Render(&r.LimitsCpu)
	if err != nil {
		return err
	}
	return nil
}
func (r *Resource) askRequestCpuSize() error {
	err := survey.NewString().
		SetKind("string").
		SetName("request-cpu").
		SetMessage("Set Request CPU resource").
		SetDefault("1").
		SetHelp("Set Request CPU resource").
		SetRequired(false).
		SetValidator(r.checkCPU).
		Render(&r.RequestsCpu)
	if err != nil {
		return err
	}
	return nil
}
func (r *Resource) askLimitMemSize() error {
	err := survey.NewString().
		SetKind("string").
		SetName("limit-mem").
		SetMessage("Set Limit Memory resource").
		SetDefault("1Gi").
		SetHelp("Set Limit Memory resource").
		SetRequired(false).
		SetValidator(r.checkMem).
		Render(&r.LimitsMemory)
	if err != nil {
		return err
	}
	return nil
}
func (r *Resource) askRequestMemSize() error {
	err := survey.NewString().
		SetKind("string").
		SetName("request-mem").
		SetMessage("Set Request Memory resource").
		SetDefault("1Gi").
		SetHelp("Set Request Memory resource").
		SetRequired(false).
		SetValidator(r.checkMem).
		Render(&r.RequestsMemory)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resource) Render() (*Resource, error) {
	if err := r.askLimitCpuSize(); err != nil {
		return nil, err
	}
	if err := r.askLimitMemSize(); err != nil {
		return nil, err
	}
	if err := r.askRequestCpuSize(); err != nil {
		return nil, err
	}
	if err := r.askRequestMemSize(); err != nil {
		return nil, err
	}
	return r, nil
}

var _ Validator = NewResource()
