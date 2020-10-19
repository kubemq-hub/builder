package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/bridges"
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
	//utils.Println(c.String())

	file, err := ioutil.ReadFile("./bridges.yaml")
	if err != nil {
		log.Fatal(err)
	}
	list := &bridges.Bindings{}
	err = yaml.Unmarshal(file, list)
	if err != nil {
		log.Fatal(err)
	}

	var bindingsYaml []byte
	if bindingsYaml, err = bridges.NewBridges("kubemq-bridges").
		SetBindings(list.Bindings).
		Render(); err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./bridges.yaml", bindingsYaml, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bindingsYaml))
}
