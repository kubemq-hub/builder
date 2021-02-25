package main

import (
	"encoding/json"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/web"
	"io/ioutil"
)

func main() {

	man, err := common.LoadManifestFromFile("./targets-manifest.json")
	if err != nil {
		panic(err)
	}
	schema, err := web.ConvertToJsonSchema(man.Targets)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(schema, "", "\t")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("schema.json", data, 0644)
	if err != nil {
		panic(err)
	}
}
