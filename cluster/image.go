package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Image struct {
	Image      string `json:"image"`
	PullPolicy string `json:"pull_policy"`
}

func NewImage() *Image {
	return &Image{}
}

func (i *Image) askImage() error {
	err := survey.NewString().
		SetKind("string").
		SetName("image").
		SetMessage("Set docker image").
		SetDefault("docker.io/kubemq/kubemq:latest").
		SetHelp("Set cluster docker image").
		SetRequired(true).
		Render(&i.Image)
	if err != nil {
		return err
	}
	return nil
}
func (i *Image) askPullPolicy() error {
	err := survey.NewString().
		SetKind("string").
		SetName("pull-policy").
		SetMessage("Set image pull policy").
		SetDefault("Always").
		SetOptions([]string{"Always", "IfNotPresent", "Never"}).
		SetHelp("Set image pull policy").
		SetRequired(true).
		Render(&i.PullPolicy)
	if err != nil {
		return err
	}
	return nil
}
func (i *Image) Validate() error {
	if i.Image == "" {
		return fmt.Errorf("docker image name is required")
	}
	switch i.PullPolicy {
	case "Always", "IfNotPresent", "Never":
	default:
		return fmt.Errorf("invalid Image Pull Data value")
	}
	return nil
}
func (i *Image) Render() (*Image, error) {

	if err := i.askImage(); err != nil {
		return nil, err
	}
	if err := i.askPullPolicy(); err != nil {
		return nil, err
	}
	return i, nil
}

var _ Validator = NewImage()
