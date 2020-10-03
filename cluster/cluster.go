package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Cluster struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Replicas        int               `json:"replicas"`
	Authentication  *Authentication   `json:"authentication"`
	Authorization   *Authorization    `json:"authorization"`
	Health          *Health           `json:"health"`
	Image           *Image            `json:"image"`
	License         string            `json:"license"`
	Log             *Log              `json:"log"`
	NodeSelectors   map[string]string `json:"node_selectors"`
	Notification    *Notification     `json:"notification"`
	Queue           *Queue            `json:"queue"`
	Resource        *Resource         `json:"resource"`
	Api             *Service          `json:"api"`
	Grpc            *Service          `json:"grpc"`
	Rest            *Service          `json:"rest"`
	Routing         *Routing          `json:"routing"`
	Store           *Store            `json:"store"`
	Tls             *Tls              `json:"tls"`
	Volume          *Volume           `json:"volume"`
	takenNames      []string
	namespaces      []string
	questionsMap    map[string]func() error
	questionOptions []string
}

func NewCluster() *Cluster {
	return &Cluster{}
}

func (c *Cluster) complete() {
	c.questionOptions = append(c.questionOptions,

		"License",
		"Docker Image",
		"Persistence Volume",
		"GRPC Service",
		"REST Service",
		"Api Service",
		"TLS Configuration",
		"Persistence Store",
		"Queues Configuration",
		"Authentication",
		"Authorization",
		"Notification",
		"Smart Routing",
		"Resources Settings",
		"Health Prob Settings",
		"Node Selectors",
		"Logging",
	)
	c.questionsMap = map[string]func() error{}
	c.questionsMap["Authentication"] = c.askAuthentication
	c.questionsMap["Authorization"] = c.askAuthorization
	c.questionsMap["Health Prob Settings"] = c.askHealth
	c.questionsMap["Docker Image"] = c.askImage
	c.questionsMap["License"] = c.askLicense
	c.questionsMap["Logging"] = c.askLog
	c.questionsMap["Node Selectors"] = c.askNodeSelectors
	c.questionsMap["Notification"] = c.askNotification
	c.questionsMap["Queues Configuration"] = c.askQueue
	c.questionsMap["Resources Settings"] = c.askResource
	c.questionsMap["Api Service"] = c.askApi
	c.questionsMap["GRPC Service"] = c.askGrpc
	c.questionsMap["REST Service"] = c.askRest
	c.questionsMap["Smart Routing"] = c.askRouting
	c.questionsMap["Persistence Store"] = c.askStore
	c.questionsMap["TLS Configuration"] = c.askTls
	c.questionsMap["Persistence Volume"] = c.askVolume

}
func (c *Cluster) SetTakenClusterNames(value []string) *Cluster {
	c.takenNames = value
	return c
}
func (c *Cluster) SetNamespaces(value []string) *Cluster {
	c.namespaces = value
	return c
}

func (c *Cluster) askName() error {
	if name, err := NewName().
		SetTakenNames(c.takenNames).
		Render(); err != nil {
		return err
	} else {
		c.Name = name.Name
	}
	return nil
}

func (c *Cluster) askNamespace() error {
	if n, err := NewNamespace().
		SetNamespaces(c.namespaces).
		Render(); err != nil {
		return err
	} else {
		c.Namespace = n.Namespace
	}
	return nil
}
func (c *Cluster) askLicense() error {
	var err error
	if c.License, err = NewLicense().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askReplicas() error {
	var err error
	if c.Replicas, err = NewReplicas().
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) askAuthentication() error {
	var err error
	if c.Authentication, err = NewAuthentication().
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) askAuthorization() error {
	var err error
	if c.Authorization, err = NewAuthorization().
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) askHealth() error {
	var err error
	if c.Health, err = NewHealth().
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) askImage() error {
	var err error
	if c.Image, err = NewImage().
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) askLog() error {
	var err error
	if c.Log, err = NewLog().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askNodeSelectors() error {
	var err error
	if c.NodeSelectors, err = NewNodeSelector().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askNotification() error {
	var err error
	if c.Notification, err = NewNotification().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askQueue() error {
	var err error
	if c.Queue, err = NewQueue().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askResource() error {
	var err error
	if c.Resource, err = NewResource().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askApi() error {
	var err error
	if c.Api, err = NewService().
		SetKind("api").
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) askRest() error {
	var err error
	if c.Rest, err = NewService().
		SetKind("rest").
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askGrpc() error {
	var err error
	if c.Rest, err = NewService().
		SetKind("grpc").
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askRouting() error {
	var err error
	if c.Routing, err = NewRouting().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askStore() error {
	var err error
	if c.Store, err = NewStore().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askTls() error {
	var err error
	if c.Tls, err = NewTls().
		Render(); err != nil {
		return err
	}
	return nil
}
func (c *Cluster) askVolume() error {
	var err error
	if c.Volume, err = NewVolume().
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) askSelection() error {
	selection := ""
	err := survey.NewString().
		SetKind("string").
		SetName("set default").
		SetMessage("Select cluster configuration").
		SetOptions([]string{"Default", "Select configurations"}).
		SetDefault("Default").
		SetHelp("Select cluster configuration").
		SetRequired(true).
		Render(&selection)
	if err != nil {
		return err
	}
	if selection == "Default" {
		return nil
	}
	var selections []string
	err = survey.NewList().
		SetKind("list").
		SetName("selection").
		SetMessage("Select at least one option").
		SetOptions(c.questionOptions).
		SetRequired(true).
		Render(&selections)
	if err != nil {
		return err
	}
	var questionsList []func() error
	for _, option := range selections {
		q, ok := c.questionsMap[option]
		if ok {
			questionsList = append(questionsList, q)
		}
	}
	for _, questionFunc := range questionsList {
		if err := questionFunc(); err != nil {
			return err
		}
	}
	return nil
}
func (c *Cluster) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("cluster name cannot be empty")
	}
	if c.Namespace == "" {
		return fmt.Errorf("cluster namespace cannot be empty")
	}

	if c.Replicas < 3 {
		return fmt.Errorf("number of replicas must be >= 3")
	}
	if c.Replicas%2 == 0 {
		return fmt.Errorf("number of replicas must be an odd number")
	}

	if c.Authentication != nil {
		if err := c.Authentication.Validate(); err != nil {
			return err
		}
	}

	if c.Authorization != nil {
		if err := c.Authorization.Validate(); err != nil {
			return err
		}
	}
	if c.Health != nil {
		if err := c.Health.Validate(); err != nil {
			return err
		}
	}
	if c.Image != nil {
		if err := c.Image.Validate(); err != nil {
			return err
		}
	}

	if c.Log != nil {
		if err := c.Log.Validate(); err != nil {
			return err
		}
	}

	if c.Log != nil {
		if err := c.Log.Validate(); err != nil {
			return err
		}
	}

	if c.Notification != nil {
		if err := c.Notification.Validate(); err != nil {
			return err
		}
	}
	if c.Queue != nil {
		if err := c.Queue.Validate(); err != nil {
			return err
		}
	}
	if c.Resource != nil {
		if err := c.Resource.Validate(); err != nil {
			return err
		}
	}
	if c.Api != nil {
		if err := c.Api.Validate(); err != nil {
			return err
		}
	}
	if c.Grpc != nil {
		if err := c.Grpc.Validate(); err != nil {
			return err
		}
	}
	if c.Rest != nil {
		if err := c.Rest.Validate(); err != nil {
			return err
		}
	}
	if c.Routing != nil {
		if err := c.Routing.Validate(); err != nil {
			return err
		}
	}
	if c.Store != nil {
		if err := c.Store.Validate(); err != nil {
			return err
		}
	}
	if c.Tls != nil {
		if err := c.Tls.Validate(); err != nil {
			return err
		}
	}
	if c.Volume != nil {
		if err := c.Volume.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) Render() (*Cluster, error) {
	c.complete()
	if err := c.askName(); err != nil {
		return nil, err
	}
	if err := c.askNamespace(); err != nil {
		return nil, err
	}
	if err := c.askReplicas(); err != nil {
		return nil, err
	}

	if err := c.askSelection(); err != nil {
		return nil, err
	}
	return c, nil
}
