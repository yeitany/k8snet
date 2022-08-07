package utils

import (
	"context"

	"github.com/yeitany/k8s_net/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetWorkloads(clientset *kubernetes.Clientset) map[string]graph.Node {
	entities := make(map[string]graph.Node, 0)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		n := graph.Node{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Type:      "pod",
			IP:        pod.Status.PodIP,
		}
		entities[pod.Status.PodIP] = n
	}
	return entities
}
