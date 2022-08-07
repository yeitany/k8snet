package graph

import "fmt"

type Node struct {
	Name      string
	Namespace string
	Type      string
	IP        string
}

func (n *Node) Format() string {
	return fmt.Sprintf("%v/%v(%v)", n.Name, n.Type, n.Namespace)
}

type Edge struct {
	Src string
	Dst string
}
