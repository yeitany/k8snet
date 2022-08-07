package main

import (
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

	http.HandleFunc("/graph", graphHandler.Handle)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(":9001", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
