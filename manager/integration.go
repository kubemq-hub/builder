package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
)

type Integration struct {
	Cluster   *cluster.Cluster
	Connector *connector.Connector
	Binding   *common.Binding
}

func NewIntegration() *Integration {
	return &Integration{}
}
func (i *Integration) Clone() *Integration {
	return &Integration{
		Cluster:   i.Cluster,
		Connector: i.Connector,
		Binding:   i.Binding.Clone(),
	}
}
func (i *Integration) SetCluster(value *cluster.Cluster) *Integration {
	i.Cluster = value
	return i
}
func (i *Integration) SetConnector(value *connector.Connector) *Integration {
	i.Connector = value
	return i
}
func (i *Integration) SetBinding(value *common.Binding) *Integration {
	i.Binding = value
	return i
}

func (i *Integration) Name() string {
	if i.Connector == nil {
		return ""
	}
	if i.Binding == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s/%s", i.Connector.Type, i.Connector.Namespace, i.Connector.Name, i.Binding.Name)
}

func EditIntegration(origin *Integration, connectorManager *ConnectorsManager) (*Integration, error) {
	cloned := origin.Clone()
	var manifest *common.Manifest
	bindings, err := common.Unmarshal([]byte(cloned.Connector.Config))
	if err != nil {
		return nil, err
	}
	switch cloned.Connector.Type {
	case "targets":
		manifest = connectorManager.catalog.ToTargetManifest()
		bindings.Side = "targets"
	case "sources":
		manifest = connectorManager.catalog.ToSourcesManifest()
		bindings.Side = "sources"
	}
	if manifest == nil {
		return nil, fmt.Errorf("no valid manifest found")
	}
	if cloned.Binding, err = cloned.Binding.
		SetEditMode(true).
		Render(); err != nil {
		return nil, err
	}
	if cloned.Binding.Name != origin.Binding.Name {
		for _, binding := range bindings.Bindings {
			if cloned.Binding.Name == binding.Name {
				return nil, fmt.Errorf("binding name %s is not unique, binding %s was not edited", cloned.Binding.Name, origin.Binding.Name)
			}
		}
	}

	if err := cloned.Binding.Validate(); err != nil {
		return nil, err
	}

	bindings.SwitchOrRemove(origin.Binding, cloned.Binding)
	data, err := bindings.Yaml()
	if err != nil {
		return nil, err
	}
	cloned.Connector.Config = string(data)
	if err := connectorManager.handler.Edit(cloned.Connector); err != nil {
		return nil, err
	}
	return cloned, nil
}
func DeleteIntegration(origin *Integration, connectorManager *ConnectorsManager) error {
	bindings, err := common.Unmarshal([]byte(origin.Connector.Config))
	if err != nil {
		return err
	}

	bindings.SwitchOrRemove(origin.Binding, nil)
	data, err := bindings.Yaml()
	if err != nil {
		return err
	}
	origin.Connector.Config = string(data)
	if err := connectorManager.handler.Edit(origin.Connector); err != nil {
		return err
	}
	return nil
}
func generateUniqueIntegrationName(takenNames []string) string {
	for i := len(takenNames) + 1; i < 10000000; i++ {
		name := fmt.Sprintf("integration-%d", i)
		found := false
		for _, taken := range takenNames {
			if taken == name {
				found = true
				break
			}
		}
		if !found {
			return name
		}
	}
	return ""
}
func CopyIntegration(origin *Integration, connectorManager *ConnectorsManager) error {
	bindings, err := common.Unmarshal([]byte(origin.Connector.Config))
	if err != nil {
		return err
	}
	cloned := origin.Clone()
	cloned.Binding.Name = generateUniqueIntegrationName(cloned.Connector.GetBindingNames())
	bindings.Bindings = append(bindings.Bindings, cloned.Binding)
	if err := bindings.Validate(); err != nil {
		return err
	}
	data, err := bindings.Yaml()
	if err != nil {
		return err
	}
	origin.Connector.Config = string(data)
	if err := connectorManager.handler.Edit(origin.Connector); err != nil {
		return err
	}
	utils.Println("<cyan>New Integration %s was added\n</>", cloned.Name())
	return nil
}
