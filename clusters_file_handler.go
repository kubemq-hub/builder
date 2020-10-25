package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"

	"io/ioutil"
	"sort"
	"time"
)

type ClustersFileHandler struct {
	filename          string
	clusters          map[string]*cluster.Cluster
	connectorsHandler connector.ConnectorsHandler
}

func NewClustersFileHandler(filename string, clusterHandler connector.ConnectorsHandler) (*ClustersFileHandler, error) {
	c := &ClustersFileHandler{
		filename:          filename,
		clusters:          map[string]*cluster.Cluster{},
		connectorsHandler: clusterHandler,
	}
	if err := c.load(); err != nil {
		if err := c.save(); err != nil {
			return nil, err
		}
	}
	return c, nil
}
func (c *ClustersFileHandler) Name() string {
	return fmt.Sprintf("file://%s", c.filename)
}
func (c *ClustersFileHandler) load() error {
	data, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &c.clusters)
	if err != nil {
		return err
	}
	return nil
}
func (c *ClustersFileHandler) save() error {
	data, err := yaml.Marshal(c.clusters)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
func (c *ClustersFileHandler) Add(cluster *cluster.Cluster) error {
	fmt.Println("add cluster")
	c.clusters[cluster.Key()] = cluster
	time.Sleep(1 * time.Second)
	return c.save()
}

func (c *ClustersFileHandler) Edit(cluster *cluster.Cluster) error {
	fmt.Println("edit cluster")
	c.clusters[cluster.Key()] = cluster
	time.Sleep(1 * time.Second)
	return c.save()
}

func (c *ClustersFileHandler) Delete(cluster *cluster.Cluster) error {
	fmt.Println("delete cluster")
	delete(c.clusters, cluster.Key())
	time.Sleep(1 * time.Second)
	return c.save()
}

func (c *ClustersFileHandler) Get(namespace, name string) (*cluster.Cluster, error) {
	fmt.Println("get cluster")
	key := fmt.Sprintf("%s/%s", namespace, name)
	con, ok := c.clusters[key]
	if !ok {
		return nil, fmt.Errorf("cluster not found")
	}
	time.Sleep(1 * time.Second)
	return con, nil

}
func (c *ClustersFileHandler) ConnectorsHandler() connector.ConnectorsHandler {
	return c.connectorsHandler
}

func (c *ClustersFileHandler) List() ([]*cluster.Cluster, error) {
	fmt.Println("get list")
	var list []*cluster.Cluster
	for _, con := range c.clusters {
		list = append(list, con)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Key() < list[j].Key()
	})

	return list, nil
}
