package graph

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/yeitany/k8s_net/utils"
)

func GenerateGraph(ctx context.Context, nodes map[string]Node, edges map[string]Edge) (*bytes.Buffer, error) {
	parameters := ctx.Value(utils.CtxKey("asd")).(utils.UrlParmeters)

	g := graphviz.New()
	g.SetLayout(parameters.Layout)
	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	unrgistedNode := &Node{
		Name: "external_ip",
	}

	for i := range edges {
		var (
			src Node
			dst Node
			ok  bool
		)
		if edges[i].Src == "" || edges[i].Dst == "" {
			continue
		}
		if src, ok = nodes[edges[i].Src]; !ok {
			src = *unrgistedNode
		}
		if dst, ok = nodes[edges[i].Dst]; !ok {
			dst = *unrgistedNode
		}
		if src.CNode == nil {
			src.CNode, err = graph.CreateNode(src.Format())
			if err != nil {
				log.Println(err.Error())
			}
		}
		if dst.CNode == nil {
			dst.CNode, err = graph.CreateNode(dst.Format())
			if err != nil {
				log.Println(err.Error())
			}
		}
		_, err = graph.CreateEdge(fmt.Sprintf("%v:%v", src.Format(), dst.Format()), src.CNode, dst.CNode)
		if err != nil {
			log.Println(err.Error())
		}

	}

	var buf bytes.Buffer
	if err := g.Render(graph, parameters.Format, &buf); err != nil {
		log.Fatal(err)
	}
	log.Println("done")
	return &buf, nil
}
