package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yeitany/k8s_net/graph"
	k8snet "github.com/yeitany/k8s_net/network"
	"github.com/yeitany/k8s_net/utils"
	"k8s.io/client-go/kubernetes"

	//
	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

var (
	entities      map[string]graph.Node
	conntrackMeta []k8snet.Meta
)

func main() {
	config, clientset := utils.GetKubeConfig()

	http.HandleFunc("/graph", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("syncNodes")
		entities = syncEnitities(clientset)
		fmt.Println("syncConntrack")
		conntrackMeta = k8snet.SyncConntracks(clientset, config)
		fmt.Println("parseConntrackMeta")
		conntrackMetaParsed := k8snet.ParseConntrackMeta(conntrackMeta)
		//fmt.Printf("%+v", entities)
		//fmt.Printf("%+v", conntrackMetaParsed)

		fmt.Println("graphviz")
		filename := graphviz(conntrackMetaParsed)
		fmt.Println("commad")
		ouput, err := exec.Command("circo", "-Tpng", filename).Output()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v", err)
		}
		fmt.Println("writing results")
		err = os.WriteFile(filename+".png", ouput, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v", err)
		}
		fmt.Println("done")
		w.WriteHeader(http.StatusOK)
		w.Write(ouput)
	})
	go http.ListenAndServe(":9001", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func syncEnitities(clientset *kubernetes.Clientset) map[string]graph.Node {
	entities := utils.GetWorkloads(clientset)
	for k, svc := range utils.GetServices(clientset) {
		entities[k] = svc
	}
	return entities
}

func graphviz(edges map[string]graph.Edge) string {
	var s string
	for k := range edges {
		var (
			src graph.Node
			dst graph.Node
			ok  bool
		)
		if src, ok = entities[edges[k].Src]; !ok {
			src = graph.Node{
				Name: "extrnal ip",
			}
		}
		if dst, ok = entities[edges[k].Dst]; !ok {
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
	s = fmt.Sprintf("digraph k8s_net \n{\n%v}", s)
	err := os.WriteFile("dot", []byte(s), 0644)
	if err != nil {
		panic(err.Error())
	}
	return "dot"
}
