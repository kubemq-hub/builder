package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

const clusterTmpl = `
<red>apiVersion:</> <white>core.k8s.kubemq.io/v1alpha1</>
<red>kind:</> <white>KubemqCluster</>
<red>metadata:</>
 <red>name:</> <white>{{.Name}}</>
 <red>namespace:</> <white>{{.Namespace}}</>
<red>spec:</>
  <red>replicas:</> <white>{{.Replicas}}</>
  {{- .VolumeSpec | indent 2 -}}
  {{- .ImageSpec | indent 2 -}}
  {{- .LicenseSpec | indent 2 -}}
  {{- .ApiServiceSpec | indent 2 -}}
  {{- .GrpcServiceSpec | indent 2 -}}
  {{- .RestServiceSpec | indent 2 -}}
  {{- .TlsSpec | indent 2 -}}
  {{- .QueueSpec | indent 2 -}}
  {{- .StoreSpec | indent 2 -}}
  {{- .AuthenticationSpec | indent 2 -}}
  {{- .AuthorizationSpec | indent 2 -}}
  {{- .RoutingSpec | indent 2 -}}
  {{- .HealthSpec | indent 2 -}}
  {{- .ResourcesSpec | indent 2 -}}
  {{- .NodeSelectorsSpec | indent 2 -}}
  {{- .NotificationSpec | indent 2 -}}
  {{- .LogSpec | indent 2 -}}
`

type Cluster struct {
	Name               string            `json:"name"`
	Namespace          string            `json:"namespace"`
	Replicas           int               `json:"replicas"`
	Authentication     *Authentication   `json:"authentication"`
	Authorization      *Authorization    `json:"authorization"`
	Health             *Health           `json:"health"`
	Image              *Image            `json:"image"`
	License            string            `json:"license"`
	Log                *Log              `json:"log"`
	NodeSelectors      map[string]string `json:"node_selectors"`
	Notification       *Notification     `json:"notification"`
	Queue              *Queue            `json:"queue"`
	Resource           *Resource         `json:"resource"`
	Api                *Service          `json:"api"`
	Grpc               *Service          `json:"grpc"`
	Rest               *Service          `json:"rest"`
	Routing            *Routing          `json:"routing"`
	Store              *Store            `json:"store"`
	Tls                *Tls              `json:"tls"`
	Volume             *Volume           `json:"volume"`
	takenNames         []string
	namespaces         []string
	questionsMap       map[string]func() error
	questionOptions    []string
	handler            ClustersHandler
	ImageSpec          string
	AuthenticationSpec string
	AuthorizationSpec  string
	HealthSpec         string
	NotificationSpec   string
	QueueSpec          string
	StoreSpec          string
	ApiServiceSpec     string
	GrpcServiceSpec    string
	RestServiceSpec    string
	VolumeSpec         string
	TlsSpec            string
	ResourcesSpec      string
	RoutingSpec        string
	LogSpec            string
	NodeSelectorsSpec  string
	LicenseSpec        string
}

func NewCluster(handler ClustersHandler) *Cluster {
	return &Cluster{
		handler: handler,
	}
}
func (c *Cluster) Clone(handler ClustersHandler) *Cluster {
	newCluster := &Cluster{
		Name:            c.Name,
		Namespace:       c.Namespace,
		License:         c.License,
		Replicas:        c.Replicas,
		NodeSelectors:   nil,
		takenNames:      c.takenNames,
		namespaces:      c.namespaces,
		questionsMap:    nil,
		questionOptions: nil,
		handler:         handler,
	}

	if c.Authentication != nil {
		newCluster.Authentication = c.Authentication.Clone()
	}

	if c.Authorization != nil {
		newCluster.Authorization = c.Authorization.Clone()
	}
	if c.Health != nil {
		newCluster.Health = c.Health.Clone()
	}
	if c.Image != nil {
		newCluster.Image = c.Image.Clone()
	}

	if c.Log != nil {
		newCluster.Log = c.Log.Clone()
	}

	if c.Notification != nil {
		newCluster.Notification = c.Notification.Clone()
	}
	if c.Queue != nil {
		newCluster.Queue = c.Queue.Clone()
	}
	if c.Resource != nil {
		newCluster.Resource = c.Resource.Clone()
	}
	if c.Api != nil {
		newCluster.Api = c.Api.Clone()
	}
	if c.Grpc != nil {
		newCluster.Grpc = c.Grpc.Clone()
	}
	if c.Rest != nil {
		newCluster.Rest = c.Rest.Clone()
	}
	if c.Routing != nil {
		newCluster.Routing = c.Routing.Clone()
	}
	if c.Store != nil {
		newCluster.Store = c.Store.Clone()
	}
	if c.Tls != nil {
		newCluster.Tls = c.Tls.Clone()
	}
	if c.Volume != nil {
		newCluster.Volume = c.Volume.Clone()
	}
	if c.NodeSelectors != nil {
		newCluster.NodeSelectors = map[string]string{}
		for key, val := range c.NodeSelectors {
			newCluster.NodeSelectors[key] = val
		}
	}
	return newCluster
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

func (c *Cluster) askName(defaultName string) error {
	if name, err := NewName().
		SetTakenNames(c.takenNames).
		Render(defaultName); err != nil {
		return err
	} else {
		c.Name = name.Name
	}
	return nil
}

func (c *Cluster) askNamespace(defaultNamespace string) error {
	if n, err := NewNamespace().
		SetNamespaces(c.namespaces).
		Render(defaultNamespace); err != nil {
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
	if c.Grpc, err = NewService().
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

func (c *Cluster) askSelection(withAskDefault bool) error {
	if withAskDefault {
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
	}
	var err error
	var selections []string
	err = survey.NewList().
		SetKind("list").
		SetName("selection").
		SetMessage("Select at least one option").
		SetOptions(c.questionOptions).
		SetRequired(true).
		SetPageSize(15).
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
	if err := c.askName("kubemq-cluster"); err != nil {
		return nil, err
	}
	if err := c.askNamespace("kubemq"); err != nil {
		return nil, err
	}
	if err := c.askReplicas(); err != nil {
		return nil, err
	}

	if err := c.askSelection(true); err != nil {
		return nil, err
	}
	return c, nil
}
func (c *Cluster) Key() string {
	return fmt.Sprintf("%s/%s", c.Namespace, c.Name)
}
func AddCluster(handler ClustersHandler) (*Cluster, error) {
	added, err := NewCluster(handler).Render()
	if err != nil {
		return nil, err
	}
	err = handler.Add(added)
	if err != nil {

	}
	return added, nil
}

func EditCluster(origin *Cluster, handler ClustersHandler, isCopyMode bool) (*Cluster, error) {
	var result *Cluster
	cloned := origin.Clone(handler)
	cloned.complete()
	ftReplicas := new(string)
	*ftReplicas = fmt.Sprintf("<r> Edit Cluster Replicas (%d)", cloned.Replicas)

	form := survey.NewForm(fmt.Sprintf("Select Edit %s Cluster Option:", cloned.Key()))
	form.AddItem(ftReplicas, func() error {
		if err := cloned.askReplicas(); err != nil {
			return err
		}
		*ftReplicas = fmt.Sprintf("<r> Edit Cluster Replicas (%d)", cloned.Replicas)
		return nil
	})

	form.AddItem("<c> Edit Cluster Configurations", func() error {
		if err := cloned.askSelection(false); err != nil {
			return err
		}
		return nil
	})

	form.AddItem("<s> Show Cluster Configuration", func() error {
		str, _ := cloned.ColoredYaml()
		utils.Println(str)
		return nil
	})

	form.SetOnSaveFn(func() error {
		if isCopyMode {
			result = cloned
			return nil
		}
		originStr, _ := origin.ColoredYaml()
		clonedStr, _ := cloned.ColoredYaml()
		if !(originStr == clonedStr) {
			if err := cloned.Validate(); err != nil {
				return err
			}
			err := cloned.handler.Edit(cloned)
			if err != nil {
				return err
			}
			result = cloned
		} else {

			result = origin
		}
		return nil
	})
	form.SetOnCancelFn(func() error {
		result = origin
		return nil
	})

	form.SetOnErrorFn(survey.FormShowErrorFn)
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
}

func CopyCluster(origin *Cluster, handler ClustersHandler) (*Cluster, error) {
	copied := origin.Clone(handler)

	if err := copied.askName(copied.Name); err != nil {
		return nil, err
	}
	if err := copied.askNamespace(copied.Namespace); err != nil {
		return nil, err
	}

	if copied.Name == origin.Name && copied.Namespace == origin.Namespace {
		return nil, fmt.Errorf("copied cluster must have different name or namespace\n")
	}
	checkEdit := false
	if err := survey.NewBool().
		SetKind("bool").
		SetMessage("Would you like to edit the copied cluster before saving").
		SetRequired(true).
		SetDefault("false").
		Render(&checkEdit); err != nil {
		return nil, err
	}
	if checkEdit {
		var err error
		copied, err = EditCluster(copied, handler, true)
		if err != nil {
			return nil, err
		}
	}
	if err := copied.Validate(); err != nil {
		return nil, err
	}
	err := copied.handler.Add(copied)
	if err != nil {
		return nil, err
	}
	return copied, nil
}

func (c *Cluster) ColoredYaml() (string, error) {
	if c.Image != nil {
		if spec, err := c.Image.ColoredYaml(); err != nil {
			return "", err
		} else {
			c.ImageSpec = spec
		}
	}
	if c.Authentication != nil {
		if spec, err := c.Authentication.ColoredYaml(); err != nil {
			return "", err
		} else {
			c.AuthenticationSpec = spec
		}
	}
	if c.Authorization != nil {
		if spec, err := c.Authorization.ColoredYaml(); err != nil {
			return "", err
		} else {
			c.AuthorizationSpec = spec
		}
	}
	//if c.Health != nil {
	//	if spec, err := c.Health.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.HealthSpec = spec
	//	}
	//}
	//if c.license != nil {
	//	if spec, err := c.license.Render(); err != nil {
	//		return "", err
	//	} else {
	//		c.LicenseSpec = spec
	//	}
	//}
	//if c.Notification != nil {
	//	if spec, err := c.Notification.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.NotificationSpec = spec
	//	}
	//}
	//if c.Queue != nil {
	//	if spec, err := c.Queue.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.QueueSpec = spec
	//	}
	//}
	//if c.Store != nil {
	//	if spec, err := c.Store.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.StoreSpec = spec
	//	}
	//}
	//if c.Api != nil {
	//	if spec, err := c.Api.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.ApiServiceSpec = spec
	//	}
	//}
	//if c.Grpc != nil {
	//	if spec, err := c.Grpc.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.GrpcServiceSpec = spec
	//	}
	//}
	//if c.Rest != nil {
	//	if spec, err := c.Rest.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.RestServiceSpec = spec
	//	}
	//}
	//if c.Tls != nil {
	//	if spec, err := c.Tls.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.TlsSpec = spec
	//	}
	//}
	//if c.Volume != nil {
	//	if spec, err := c.Volume.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.VolumeSpec = spec
	//	}
	//}
	//if c.Resource != nil {
	//	if spec, err := c.Resource.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.ResourcesSpec = spec
	//	}
	//}
	//if c.Routing != nil {
	//	if spec, err := c.Routing.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.RoutingSpec = spec
	//	}
	//}
	//if c.Log != nil {
	//	if spec, err := c.Log.ColoredYaml(); err != nil {
	//		return "", err
	//	} else {
	//		c.LogSpec = spec
	//	}
	//}
	//if c.nodeSelectors != nil {
	//	if spec, err := c.nodeSelectors.Render(); err != nil {
	//		return "", err
	//	} else {
	//		c.NodeSelectorsSpec = spec
	//	}
	//}
	t := NewTemplate(clusterTmpl, c)
	b, err := t.Get()
	if err != nil {
		return err.Error(), err
	}

	return string(b), nil
}
