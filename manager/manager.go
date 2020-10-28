package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/survey"
)

type Manager struct {
	connectorHandler connector.ConnectorsHandler
	clusterHandler   cluster.ClustersHandler
	clusterManager   *ClustersManager
	connectorManager *ConnectorsManager
	catalogManager   *CatalogManager
	loadedOptions    common.DefaultOptions
}

func NewManager() *Manager {
	return &Manager{}
}
func (m *Manager) Init(loadedOptions common.DefaultOptions, connectorHandler connector.ConnectorsHandler, clusterHandler cluster.ClustersHandler) error {
	m.loadedOptions = loadedOptions
	m.catalogManager = NewCatalogManager()
	if err := m.catalogManager.Init(); err != nil {
		return fmt.Errorf("error on connector catalog initialzation: %s", err.Error())
	}
	m.connectorManager = NewConnectorsManager(connectorHandler, m.catalogManager, loadedOptions)
	m.clusterManager = NewClustersManager(clusterHandler, m.connectorManager, loadedOptions)
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
		AddItem("<i> Manage Integrations Catalog", func() error {
			return m.catalogManager.Render()
		}).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		Render(); err != nil {
		return err
	}
	return nil
}
