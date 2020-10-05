package main

import (
	"fmt"
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/builder/connector"
	"log"
)

func main() {
	//options := map[string][]string{
	//	"target.command/address":      []string{"cluster-1:50000", "cluster-2:50000", "Other"},
	//	"target.query/address":        []string{"cluster-1:50000", "cluster-2:50000", "Other"},
	//	"target.queue/address":        []string{"cluster-1:50000", "cluster-2:50000", "Other"},
	//	"target.events/address":       []string{"cluster-1:50000", "cluster-2:50000", "Other"},
	//	"target.events-store/address": []string{"cluster-1:50000", "cluster-2:50000", "Other"},
	//}
	//manData, err := ioutil.ReadFile("sources.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//m, err := builder.LoadManifest(manData)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//b, err := builder.NewBuilder().
	//	SetManifest(m).
	//	SetOptions(options).Render()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(b))
	m, err := common.LoadFromUrl("https://raw.githubusercontent.com/kubemq-hub/kubemq-sources/master/manifest.json")
	if err != nil {
		log.Fatal(m)
	}
	t, err := connector.NewSource().SetManifest(m.Marshal()).Render()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(t))
}
