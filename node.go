package graphlearn

// Edge graph edge
type Edge struct {
	LinksTo string
	Length  float32
}

// Node graph node
type Node struct {
	Name  string
	Edges []Edge
}

// GridNode grid node, wraps a node with a co-ordinate
type GridNode struct {
	Node
	X, Y float32
}

// Path path found to node
type Path struct {
	Length float32
	Path   []string
}
