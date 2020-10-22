package connector

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
	"strconv"
)

type Replicas struct {
	value int
}

func (r *Replicas) Validate() error {
	if r.value < 0 {
		return fmt.Errorf("number of replicase must be >= 0")
	}
	return nil
}
func (r *Replicas) checkValue(val interface{}) error {
	if str, ok := val.(string); ok {
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		if val < 0 {
			return fmt.Errorf("number of replicase must be >= 0")
		}
	}
	return nil
}
func NewReplicas() *Replicas {
	return &Replicas{}
}
func (r *Replicas) Render(defaultVal int) (int, error) {

	err := survey.NewInt().
		SetKind("int").
		SetName("replicas").
		SetMessage("Set connector's replicas").
		SetDefault(fmt.Sprintf("%d", defaultVal)).
		SetRange(0, math.MaxInt32).
		SetHelp("Set how many replicas for the connector").
		SetRequired(true).
		SetValidator(r.checkValue).
		Render(&r.value)
	if err != nil {
		return 0, err
	}
	return r.value, nil
}

var _ Validator = NewReplicas()
