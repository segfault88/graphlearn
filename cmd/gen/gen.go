package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/segfault88/graphlearn"
)

const (
	graphX           = 16
	graphY           = 16
	wordLength       = 4
	linksMin         = 2
	linksMax         = 5
	linkMaxNodeRange = 3 // up to 3 rows/columns away
)

func getWords() []string {
	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatalf("Couldn't open wors %s", err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Couldn't read words file")
	}

	out := []string{}

	words := strings.Split(string(data), "\n")
	for _, word := range words {
		if len(word) > wordLength {
			out = append(out, word)
		}
	}

	return out
}

func randomWord(words []string, usedWords []string) string {
	j := int(rand.Int31()) % len(words)
	word := words[j]

	// check that word hasn't already been used
	for _, used := range usedWords {
		if word == used {
			// word already used, try again
			return randomWord(words, usedWords)
		}
	}

	// word not already used
	return word
}

func wrapToRange(n int, max int) int {
	// if less than 0, wrap around to max
	// if greater than max, wrap around to min
	if n < 0 {
		return max - 1 - (n * -1)
	}
	if n > max-1 {
		return n - max
	}

	return n
}

func main() {
	fmt.Printf("Generate a test graph\n")
	fmt.Printf("Makes a grid %d by %d (%d nodes), gives a random name, shifts around a little\n", graphX, graphY, graphX*graphY)
	fmt.Printf("randomly joins it up working out length, then marshals just the nodes\n")
	fmt.Printf("with names, connections and connection lengths.\n")
	fmt.Printf("Then finally saves into graph.json.\n")

	rand.Seed(time.Now().UnixNano())

	words := getWords()
	usedWords := []string{}

	fmt.Printf("Got %d words for node names\n", len(words))

	grid := make([][]graphlearn.GridNode, graphX)
	for x := 0; x < graphX; x++ {
		grid[x] = make([]graphlearn.GridNode, graphY)

		for y := 0; y < graphY; y++ {
			word := randomWord(words, usedWords)
			usedWords = append(usedWords, word)

			grid[x][y].Node.Name = word
			grid[x][y].Edges = []graphlearn.Edge{}

			// set x and y for the grid, shifting them around by a random amount
			grid[x][y].X = float32(x) + float32(math.Mod(rand.Float64(), 15.0)-7.5)
			grid[x][y].Y = float32(y) + float32(math.Mod(rand.Float64(), 15.0)-7.5)
		}
	}

	totalDistance := float64(0.0)

	// now build links
	for x := 0; x < graphX; x++ {
		for y := 0; y < graphY; y++ {
			thisNode := grid[x][y]
			nLinks := int(linksMin + rand.Int31()%(linksMax-linksMin))

			for i := 0; i < nLinks; i++ {
				// how far away is the node we're going to link to
				deltaX := rand.Intn(linkMaxNodeRange*2) - linkMaxNodeRange
				deltaY := rand.Intn(linkMaxNodeRange*2) - linkMaxNodeRange

				otherX := wrapToRange(x+deltaX, graphX)
				otherY := wrapToRange(y+deltaY, graphY)

				if otherX == 0 && otherY == 0 {
					// don't link to self
					nLinks++
					continue
				}

				otherNode := grid[otherX][otherY]

				dX := float64(otherNode.Y - thisNode.Y)
				dY := float64(otherNode.X - thisNode.Y)
				distance := float32(math.Sqrt(dX*dX + dY*dY))

				totalDistance += float64(distance)

				edge := graphlearn.Edge{LinksTo: otherNode.Name, Length: distance}
				grid[x][y].Edges = append(grid[x][y].Edges, edge)
				// add the edge in the other direction
				otherEdge := graphlearn.Edge{LinksTo: thisNode.Name, Length: distance}
				grid[otherX][otherY].Edges = append(grid[otherX][otherY].Edges, otherEdge)
			}
		}
	}

	graph := []graphlearn.Node{}
	for x := 0; x < graphX; x++ {
		for y := 0; y < graphY; y++ {
			graph = append(graph, grid[x][y].Node)
		}
	}

	out, err := json.Marshal(graph)
	if err != nil {
		panic("Couldn't marshal graph!")
	}

	ioutil.WriteFile("graph.json", out, 0644)
	fmt.Printf("Output length %d, total distance: %f\n", len(out), totalDistance)
}
