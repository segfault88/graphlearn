package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/segfault88/graphlearn"
)

func main() {
	data, err := ioutil.ReadFile("./graph.json")
	if err != nil {
		panic(err)
	}

	nodes := []graphlearn.Node{}
	err = json.Unmarshal(data, &nodes)

	if err != nil {
		panic(err)
	}

	graphlearn.SearchForShortestPath(nodes)
}
