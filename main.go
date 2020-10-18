package main

import (
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/bridges/binding"
	"github.com/kubemq-hub/builder/connector/bridges/source"
	"github.com/kubemq-hub/builder/connector/bridges/target"
	"github.com/kubemq-hub/builder/pkg/utils"
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
	//b, err := builder.NewBindings().
	//	SetManifest(m).
	//	SetOptions(options).Render()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(b))
	//m, err := common.LoadFromUrl("https://raw.githubusercontent.com/kubemq-hub/kubemq-sources/master/manifest.json")
	//if err != nil {
	//	log.Fatal(m)
	//}
	br := &connector.Bridges{
		Bindings: []*binding.Binding{
			&binding.Binding{
				Name: "b1",
				Sources: &source.Source{
					Name: "b1.s1",
					Kind: "s",
					Connections: []map[string]string{
						{
							"c11": "c11",
							"c12": "c12",
						},
						{
							"c21": "c21",
							"c22": "c22",
							"c23": "\nasdasdasdasdadasda\nasddasdasd",
						},
					},
					ConnectionSpec: "",
				},
				Targets: &target.Target{
					Name: "b1.t1",
					Kind: "t",
					Connections: []map[string]string{
						{
							"tc11": "tc11",
							"tc12": "tc12",
						},
						{
							"tc21": "tc21",
							"tc22": "tc22",
							"tc23": "\nasdasdasdasdadasda\nasddasdasd",
						},
					},
					ConnectionSpec: "",
				},
				Properties: map[string]string{
					"log_level":                     "debug",
					"rate_per_second":               "100",
					"retry_attempts":                "1",
					"retry_delay_milliseconds":      "100",
					"retry_delay_type":              "fixed",
					"retry_max_jitter_milliseconds": "100",
				},
				SourcesSpec:    "",
				TargetSpec:     "",
				PropertiesSpec: "",
			},
			&binding.Binding{
				Name: "b1",
				Sources: &source.Source{
					Name: "b1.s1",
					Kind: "s",
					Connections: []map[string]string{
						{
							"c11": "c11",
							"c12": "c12",
						},
						{
							"c21": "c21",
							"c22": "c22",
							"c23": "\nasdasdasdasdadasda\nasddasdasd",
						},
					},
					ConnectionSpec: "",
				},
				Targets: &target.Target{
					Name: "b1.t1",
					Kind: "t",
					Connections: []map[string]string{
						{
							"tc11": "tc11",
							"tc12": "tc12",
						},
						{
							"tc21": "tc21",
							"tc22": "tc22",
							"tc23": "\nasdasdasdasdadasda\nasddasdasd",
						},
					},
					ConnectionSpec: "",
				},
				Properties: map[string]string{
					"log_level":                     "debug",
					"rate_per_second":               "100",
					"retry_attempts":                "1",
					"retry_delay_milliseconds":      "100",
					"retry_delay_type":              "fixed",
					"retry_max_jitter_milliseconds": "100",
				},
				SourcesSpec:    "",
				TargetSpec:     "",
				PropertiesSpec: "",
			},
		},
	}
	cfg, err := br.Yaml()
	if err != nil {
		log.Fatal(err)
	}
	c := &connector.Connector{
		Name:        "kubemq-connector-1",
		Namespace:   "abc",
		Type:        "bridges",
		Replicas:    1,
		Config:      string(cfg),
		NodePort:    0,
		ServiceType: "",
		Image:       "",
	}
	utils.Println(c.String())
	//
	//c, err := connector.NewConnector().Render()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//	m1 := map[string]string{
	//		"key1": "value1",
	//		"aksdasda2": `value2
	//asdasd
	//asd
	//asdasdasd
	//asd`,
	//	}
	//	m2 := map[string]string{
	//		"key1":      "value1",
	//		"aksdasda2": "value2",
	//	}
	//
	//	utils.Println(MapArrayToYaml([]map[string]string{m1, m2}))
	//
	//}
}
