package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/yeitany/k8s_net/graph"
	k8snet "github.com/yeitany/k8s_net/network"
	"github.com/yeitany/k8s_net/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type GraphHandler struct {
	Config        *rest.Config
	Clientset     *kubernetes.Clientset
	entities      map[string]graph.Node
	conntrackMeta []k8snet.Meta
}

func (h *GraphHandler) Handle(w http.ResponseWriter, req *http.Request) {
	log.Println("syncNodes")
	h.entities = h.syncEnitities(h.Clientset)
	log.Println("syncConntrack")
	h.conntrackMeta = k8snet.SyncConntracks(h.Clientset, h.Config)
	log.Println("parseConntrackMeta")
	conntrackMetaParsed := k8snet.ParseConntrackMeta(h.conntrackMeta)

	log.Println("graphviz")
	filename := h.graphviz(conntrackMetaParsed)
	defer func() {
		os.Remove(filename)
	}()
	log.Println("commad")
	ouput, err := exec.Command("circo", "-Tpng", filename).Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
	}
	log.Println("done")
	w.WriteHeader(http.StatusOK)
	w.Write(ouput)
}

func (h *GraphHandler) syncEnitities(clientset *kubernetes.Clientset) map[string]graph.Node {
	entities := utils.GetWorkloads(clientset)
	for k, svc := range utils.GetServices(clientset) {
		entities[k] = svc
	}
	return entities
}

func (h *GraphHandler) graphviz(edges map[string]graph.Edge) string {
	var s string
	for k := range edges {
		var (
			src graph.Node
			dst graph.Node
			ok  bool
		)
		if src, ok = h.entities[edges[k].Src]; !ok {
			src = graph.Node{
				Name: "extrnal ip",
			}
		}
		if dst, ok = h.entities[edges[k].Dst]; !ok {
			dst = graph.Node{
				Name: "extrnal ip",
			}
		}
		if src.Namespace == "kube-system" || dst.Namespace == "kube-system" {
			continue
		}
		if src.Namespace == "gmp-system" || dst.Namespace == "gmp-system" {
			continue
		}
		s += fmt.Sprintf("\"%v\" -> \"%v\";\n", src.Format(), dst.Format())
	}
	filename := strconv.FormatInt(int64(time.Now().Unix()), 10)
	s = fmt.Sprintf("digraph k8s_net \n{\n%v}", s)
	err := os.WriteFile(filename, []byte(s), 0644)
	if err != nil {
		panic(err.Error())
	}
	return filename
}
