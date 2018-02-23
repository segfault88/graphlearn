# graphlearn

Just messing around with graph algorithm stuff for learning purposes.

## gen

Generate a random graph of nodes and connections (with lengths)

```bash
go run ./cmd/gen/gen.go # generates graph.json randomly
go run ./cmd/draw/draw.go # generates graph.dot - graphviz format for viewing
dot -Tsvg graph.dot -o graph.svg # generate svg of graph, takes about 3 mins to layout
```

