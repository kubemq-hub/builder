package manager

const (
	promptCatalogLoadingStarted         = `<yellow>Loading Connectors catalog...</>`
	promptCatalogLoadingCompleted       = `<yellow>Loading Connectors catalog completed.</>`
	promptCatalogLoadingError           = `<red>Loading Connectors catalog error:%s</>`
	promptConnectorDelete               = "<cyan>Connector %s deleted successfully\n</>"
	promptConnectorAdded                = "<cyan>Connector %s added successfully\n</>"
	promptClusterAdded                  = "<cyan>Cluster %s added successfully\n</>"
	promptConnectorEdited               = "<cyan>Connector %s edited successfully\n</>"
	promptClusterEdited                 = "<cyan>Cluster %s edited successfully\n</>"
	promptConnectorCopied               = "<cyan>Connector %s copied successfully\n</>"
	promptClusterCopied                 = "<cyan>Cluster %s copied successfully\n</>"
	promptClusterDelete                 = "<cyan>Cluster %s deleted successfully\n</>"
	promptIntegrationEditedConfirmation = "<cyan>Integration %s edited successfully\n</>"
	promptIntegrationAddedConfirmation  = "<cyan>Integration %s added successfully\n</>"
	promptIntegrationDeleteConfirmation = "<cyan>Integration %s deleted successfully\n</>"
	promptIntegrationCopiedConfirmation = "<cyan>Integration %s copied successfully\n</>"
)

type ContextHandler interface {
	Set() error
	Get() string
}
