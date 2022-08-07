package utils

import (
	"context"

	"github.com/yeitany/k8s_net/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetServices(clientset *kubernetes.Clientset) map[string]graph.Node {
	entities := make(map[string]graph.Node, 0)
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, svc := range services.Items {
		n := graph.Node{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      "svc",
			IP:        svc.Spec.ClusterIP,
		}
		entities[svc.Spec.ClusterIP] = n
	}
	return entities
}
