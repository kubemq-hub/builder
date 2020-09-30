package main

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/bridges"
	"log"
)

func main() {
	spec, err := bridges.NewBuilder().
		SetAddress([]string{"cluster-1:50000", "cluster-2:50000", "Other"}).
		Render()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("%+v", spec))
}
