package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yeitany/k8s_net/handlers"
	"github.com/yeitany/k8s_net/utils"

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

func main() {
	config, clientset := utils.GetKubeConfig()
	graphHandler := handlers.GraphHandler{
		Config:    config,
		Clientset: clientset,
	}

	port := flag.String("servePort", "8080", "serve port")
	metricPort := flag.String("metricPort", "9001", "metric port")
	enableMetrics := flag.Bool("enableMetrics", true, "enable metrics")
	flag.Parse()

	http.HandleFunc("/graph", graphHandler.ServeHttp)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	if *enableMetrics {
		go http.ListenAndServe(fmt.Sprintf(":%v", *metricPort), promhttp.Handler())
	}
	http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
}
