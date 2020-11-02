package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Manager struct {
	clusterManager      *ClustersManager
	connectorManager    *ConnectorsManager
	catalogManager      *CatalogManager
	loadedOptions       common.DefaultOptions
	integrationsManager *IntegrationsManager
	contextHandler      ContextHandler
}

func NewManager() *Manager {
	return &Manager{}
}
func (m *Manager) Init(loadedOptions common.DefaultOptions, connectorHandler connector.ConnectorsHandler, clusterHandler cluster.ClustersHandler, contextHandler ContextHandler) error {
	m.loadedOptions = loadedOptions
	m.catalogManager = NewCatalogManager()
	if err := m.catalogManager.Update(false); err != nil {
		return fmt.Errorf("error on connector catalog initialzation: %s", err.Error())
	}
	m.connectorManager = NewConnectorsManager(connectorHandler, m.catalogManager, loadedOptions)
	m.clusterManager = NewClustersManager(clusterHandler, m.connectorManager, loadedOptions)
	m.integrationsManager = NewIntegrationsManager(m.clusterManager, m.connectorManager, m.catalogManager)
	m.contextHandler = contextHandler
	current := m.contextHandler.Get()
	m.clusterManager.SetCurrentContext(current)
	m.connectorManager.SetCurrentContext(current)
	m.integrationsManager.SetCurrentContext(current)

	return nil
}
func (m *Manager) changeKubernetesContext() error {
	if err := m.contextHandler.Set(); err != nil {
		return err
	}
	current := m.contextHandler.Get()
	m.clusterManager.SetCurrentContext(current)
	m.connectorManager.SetCurrentContext(current)
	m.integrationsManager.SetCurrentContext(current)
	utils.Println("\n")
	return nil
}
func (m *Manager) Render() error {

	if err := survey.NewMenu("Select Manager Option:").
		AddItem("<c> Manage KubeMQ Clusters", func() error {
			return m.clusterManager.Render()
		}).
		AddItem("<o> Manage KubeMQ Connectors", func() error {
			return m.connectorManager.Render()
		}).
		AddItem("<i> Manage KubeMQ Integrations", func() error {
			return m.integrationsManager.Render()
		}).
		AddItem("<x> Change Kubernetes Context", func() error {
			return m.changeKubernetesContext()
		}).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		Render(); err != nil {
		return err
	}
	return nil
}
