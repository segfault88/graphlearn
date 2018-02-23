package graphlearn

// Edge graph edge
type Edge struct {
	linksTo string
	length  float32
}

// Node graph node
type Node struct {
	name  string
	edges []Edge
}
