package handlers

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/goccy/go-graphviz"
	"github.com/yeitany/k8s_net/components"
	k8snet_graph "github.com/yeitany/k8s_net/graph"
	k8snet "github.com/yeitany/k8s_net/network"
	"github.com/yeitany/k8s_net/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type GraphHandler struct {
	Config    *rest.Config
	Clientset *kubernetes.Clientset
}

func (h *GraphHandler) ServeHttp(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	ctx := req.Context()

	ctx, cancelFunc := context.WithCancel(ctx)
	defer func() {
		cancelFunc()
		log.Printf("execution time in seconds:%v\n", time.Since(start).Seconds())
	}()

	urlParsed, err := parseUrl(req)
	if err != nil {
		log.Print("unable to parse url parmeters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("%+v", urlParsed)

	ctx = context.WithValue(ctx, utils.CtxKey("asd"), *urlParsed)
	resultCh := make(chan *bytes.Buffer, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		wg.Add(1)
		componentChan := make(chan map[string]k8snet_graph.Node, 1)
		go func() {
			defer wg.Done()
			nodes, err := components.GetComponents(h.Clientset)
			if err != nil {
				cancelFunc()
			}
			componentChan <- nodes
		}()

		wg.Add(1)
		edgesChan := make(chan map[string]k8snet_graph.Edge, 1)
		go func() {
			defer wg.Done()
			conntrackMeta, err := k8snet.SyncConntracks(h.Clientset, h.Config)
			if err != nil {
				cancelFunc()
			}
			edgesChan <- k8snet.ParseConntrackMeta(conntrackMeta)

		}()

		log.Println("graphviz")
		buf, err := k8snet_graph.GenerateGraph(ctx, <-componentChan, <-edgesChan)
		if err != nil {
			cancelFunc()
		}
		resultCh <- buf
	}()

	wg.Wait()
	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			log.Printf("context canceled : %v", err)
		}
		close(resultCh)
		return
	case result := <-resultCh:
		w.WriteHeader(http.StatusOK)
		w.Write(result.Bytes())
	}

}

func parseUrl(req *http.Request) (*utils.UrlParmeters, error) {
	qMap, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	layout := convertor(qMap, "layout", []string{string(graphviz.CIRCO)}).([]string)[0]
	format := convertor(qMap, "format", []string{string(graphviz.PNG)}).([]string)[0]
	targets := convertor(qMap, "targets", []string{}).([]string)
	blacklist := convertor(qMap, "blacklist", []string{}).([]string)
	whitelist := convertor(qMap, "whitelist", []string{}).([]string)
	var u utils.UrlParmeters = utils.UrlParmeters{
		Layout:              graphviz.Layout(layout),
		Format:              graphviz.Format(format),
		Targets:             targets,
		BlacklistNamespaces: blacklist,
		WhitelistNamespaces: whitelist,
	}
	if !u.IsValid() {
		return nil, errors.New("not valid")
	}
	return &u, nil
}

func convertor(y map[string][]string, key string, defaultVal interface{}) interface{} {
	v := defaultVal
	if formatArr, ok := y[key]; ok {
		if len(formatArr) > 0 {
			v = formatArr
		}
	}
	return v
}
