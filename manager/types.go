package manager

import "github.com/kubemq-hub/builder/connector"

const (
	promptCatalogLoadingStarted   = `<yellow>Loading Connectors catalog...</>`
	promptCatalogLoadingCompleted = `<yellow>Loading Connectors catalog completed</>`
	promptCatalogLoadingError     = `<red>Loading Connectors catalog error:%s</>`
)

type ConnectorsHandler interface {
	Add(connector *connector.Connector) error
	Edit(connector *connector.Connector) error
	Delete(connector *connector.Connector) error
	Get(namespace, name string) (*connector.Connector, error)
	List() []*connector.Connector
}
