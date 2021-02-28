package main

import (
	"github.com/kubemq-hub/builder/web"
	"io/ioutil"
)

func main() {
	schema, err := web.NewIntegrationSchema().Load("./sources-manifest.json", "./targets-manifest.json")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("schema.json", schema.Marshal(), 0644)
	if err != nil {
		panic(err)
	}
}
