package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/yeitany/k8s_net/components"
	k8snet_graph "github.com/yeitany/k8s_net/graph"
	k8snet "github.com/yeitany/k8s_net/network"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type GraphHandler struct {
	Config    *rest.Config
	Clientset *kubernetes.Clientset
}

type UrlParmeters struct {
	BlacklistNamespaces []string `json:"blacklist_namespaces"`
	WhitelistNamespaces []string `json:"whitelist_namespaces"`
	Targets             []string `json:"targets"`
	Format              string   `json:"format"`
	Layout              string   `json:"layout"`
}

func (h *GraphHandler) ServeHttp(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	urlParsed, err := parseUrl(req, w)
	if err != nil {
		log.Print("unable to parse url parmeters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("%+v", urlParsed)

	nodes, err := components.GetComponents(h.Clientset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conntrackMeta, err := k8snet.SyncConntracks(h.Clientset, h.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	edges := k8snet.ParseConntrackMeta(conntrackMeta)

	log.Println("graphviz")
	buf, err := k8snet_graph.GenerateGraph(nodes, edges)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
	log.Printf("execution time in seconds:%v\n", time.Since(start).Seconds())
}

func parseUrl(req *http.Request, w http.ResponseWriter) (*UrlParmeters, error) {
	qMap, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	var u UrlParmeters
	qMarshal, err := json.Marshal(qMap)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(qMarshal, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
