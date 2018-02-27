package graphlearn

import (
	"log"
	"math"
)

// searchContext wrap up some stuff together with the index
type searchContext struct {
	nodes      []Node
	start      *Node
	end        *Node
	index      map[string]int
	visited    map[int]bool
	foundPaths []Path
}

// SearchForShortestPath try to find the shortest path from the first node to the last
func SearchForShortestPath(nodes []Node) {
	sc := searchContext{
		nodes:      nodes,
		start:      &nodes[0],
		end:        &nodes[len(nodes)-1],
		index:      map[string]int{},
		visited:    map[int]bool{},
		foundPaths: make([]Path, len(nodes)),
	}

	log.Printf("%d nodes found, going to search from '%s' to '%s', setting up\n",
		len(nodes),
		sc.start.Name,
		sc.end.Name,
	)

	// build an index of node name -> index
	vertices := 0
	for i, node := range nodes {
		sc.index[node.Name] = i
		vertices += len(node.Edges)
	}

	log.Printf("built index, %d total vertices", vertices)

	// starting node added with length 0
	sc.foundPaths[0].Length = 0.0
	sc.foundPaths[0].Path = []string{}
	// all the rest have infinite length (max float32 actually)
	for i := 1; i < len(nodes); i++ {
		sc.foundPaths[i].Length = math.MaxFloat32
	}

	searchNode(&sc, *sc.start, []string{sc.start.Name}, 0.0)

	log.Printf("End\n")
	log.Printf("%#v", sc.end.Name)
	pathToEnd := sc.foundPaths[sc.index[sc.end.Name]]
	log.Printf("Path from start to end, length: %f, path: %#v", pathToEnd.Length, pathToEnd.Path)
	// log.Printf("Final paths %#v", sc.foundPaths)
}

func searchNode(sc *searchContext, node Node, pathSoFar []string, lengthSoFar float32) {
	// mark node visited
	sc.visited[sc.index[node.Name]] = true

	for _, edge := range node.Edges {
		thisLength := lengthSoFar + edge.Length
		linksToIdx := sc.index[edge.LinksTo]
		pathTo := sc.foundPaths[linksToIdx]

		// check if we've found a path shorter than what's been found already
		if thisLength < pathTo.Length {
			log.Printf("Found shorter path to %s, was %f, now %f",
				edge.LinksTo,
				pathTo.Length,
				thisLength,
			)

			// add this to the found paths
			sc.foundPaths[linksToIdx].Length = thisLength
			sc.foundPaths[linksToIdx].Path = append(pathSoFar, edge.LinksTo)
		}
	}

	for _, edge := range node.Edges {
		linksToIdx := sc.index[edge.LinksTo]

		// if this edge has not been visited, visit it!
		if !sc.visited[linksToIdx] {
			log.Printf("Links to node %s has not been visited, searching", edge.LinksTo)
			linksToNode := sc.nodes[linksToIdx]
			searchNode(sc, linksToNode, append(pathSoFar, linksToNode.Name), lengthSoFar+edge.Length)
		} else {
			log.Printf("Links to node %s has already been visited, skipping", edge.LinksTo)
		}
	}
}
