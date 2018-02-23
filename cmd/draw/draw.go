package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/segfault88/graphlearn"
)

func main() {
	fmt.Printf("Make a graphviz view of the data in graph.json\n")

	data, err := ioutil.ReadFile("graph.json")
	if err != nil {
		log.Fatalf("couldn't read graph %s\n", err)
	}

	var nodes []graphlearn.Node
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		log.Fatalf("couldn't unmarshal graph data %s\n", err)
	}

	fmt.Printf("Loaded %d nodes\n", len(nodes))

	viz := strings.Builder{}
	viz.WriteString("graph {\n")

	for _, node := range nodes {
		// for i := 0; i < 1000; i++ {
		// node := nodes[i]
		connections := ""
		for _, c := range node.Edges {
			connections += c.LinksTo
			connections += " "
		}
		connections = strings.TrimSpace(connections)
		viz.WriteString(fmt.Sprintf("\t%s -- { %s }\n", node.Name, connections))
	}

	viz.WriteString("}")
	err = ioutil.WriteFile("graph.dot", []byte(viz.String()), 0644)
	if err != nil {
		log.Fatalf("Couldn't write graph viz file %s", err)
	}
}
