package cluster

import "github.com/kubemq-hub/builder/connector"

type ClustersHandler interface {
	Name() string
	Add(cluster *Cluster) error
	Edit(cluster *Cluster) error
	Delete(cluster *Cluster) error
	Get(namespace, name string) (*Cluster, error)
	List() ([]*Cluster, error)
	ConnectorsHandler() connector.ConnectorsHandler
}
