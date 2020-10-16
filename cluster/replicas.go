package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"strconv"
)

type Replicas struct {
	value int
}

func (r *Replicas) Validate() error {
	if r.value < 3 {
		return fmt.Errorf("number of nodes must be >= 3")
	}
	if r.value%2 == 0 {
		return fmt.Errorf("number of nodes must be an odd number")
	}
	return nil
}
func (r *Replicas) checkValue(val interface{}) error {
	if str, ok := val.(string); ok {
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		if val < 3 {
			return fmt.Errorf("number of nodes must be >= 3")
		}
		if val%2 == 0 {
			return fmt.Errorf("number of nodes must be an odd number")
		}
	}
	return nil
}
func NewReplicas() *Replicas {
	return &Replicas{}
}
func (r *Replicas) Render() (int, error) {
	err := survey.NewInt().
		SetKind("int").
		SetName("replicas").
		SetMessage("Set Cluster replicas (3,5,7,9)").
		SetDefault("3").
		SetRange(3, 9).
		SetHelp("Set how many cluster nodes to deploy").
		SetRequired(true).
		SetValidator(r.checkValue).
		Render(&r.value)
	if err != nil {
		return 0, err
	}
	return r.value, nil
}

var _ Validator = NewReplicas()
