package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/common"
	connectorSources "github.com/kubemq-hub/builder/connector/sources"
	"io/ioutil"
	"log"
)

func main() {
	//m, err := common.LoadManifestFromFile("./sources-manifest.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c, err := connector.NewConnector().
	//	SetSourcesManifest(m.Marshal()).
	//	Render()
	//if err != nil {
	//	log.Fatal(err)
	//}
	////c, err := connector.NewConnector().
	////	SetSourcesManifest(m.Marshal()).
	////	Render()
	////if err != nil {
	////	log.Fatal(err)
	////}
	//utils.Println(c.ColoredYaml())

	file, err := ioutil.ReadFile("./sources.yaml")
	if err != nil {
		log.Fatal(err)
	}
	list := &common.Bindings{}
	err = yaml.Unmarshal(file, list)
	if err != nil {
		log.Fatal(err)
	}

	var bindingsYaml []byte
	if bindingsYaml, err = connectorSources.NewSource("kubemq-bridges").
		SetBindings(list.Bindings).
		SetManifestFile("./sources-manifest.json").
		Render(); err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./sources.yaml", bindingsYaml, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bindingsYaml))
}
