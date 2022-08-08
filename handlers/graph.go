package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/goccy/go-graphviz"
	"github.com/yeitany/k8s_net/graph"
	k8snet_graph "github.com/yeitany/k8s_net/graph"
	k8snet "github.com/yeitany/k8s_net/network"
	"github.com/yeitany/k8s_net/workloads"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type GraphHandler struct {
	Config    *rest.Config
	Clientset *kubernetes.Clientset
}

func (h *GraphHandler) ServeHttp(w http.ResponseWriter, req *http.Request) {
	log.Println("syncNodes")
	nodes := h.syncEnitities(h.Clientset)
	log.Println("syncConntrack")
	conntrackMeta := k8snet.SyncConntracks(h.Clientset, h.Config)
	log.Println("parseConntrackMeta")
	edges := k8snet.ParseConntrackMeta(conntrackMeta)

	log.Println("graphviz")

	buf := h.generateGraph(nodes, edges)
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (h *GraphHandler) generateGraph(nodes map[string]k8snet_graph.Node, edges map[string]k8snet_graph.Edge) bytes.Buffer {
	g := graphviz.New()
	g.SetLayout(graphviz.CIRCO)
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	unrgistedNode := &k8snet_graph.Node{
		Name: "external_ip",
	}
	for i := range edges {
		var (
			src k8snet_graph.Node
			dst k8snet_graph.Node
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
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		log.Fatal(err)
	}
	log.Println("done")
	return buf
}

func (h *GraphHandler) syncEnitities(clientset *kubernetes.Clientset) map[string]graph.Node {
	entities := workloads.GetWorkloads(clientset)
	for k, svc := range workloads.GetServices(clientset) {
		entities[k] = svc
	}
	return entities
}
