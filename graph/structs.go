package graph

import (
	"fmt"

	"github.com/goccy/go-graphviz/cgraph"
)

type Node struct {
	Name      string
	Namespace string
	Type      string
	IP        string
	CNode     *cgraph.Node
}

func (n *Node) Format() string {
	return fmt.Sprintf("%v/%v(%v)", n.Name, n.Type, n.Namespace)
}

type Edge struct {
	Src string
	Dst string
}
