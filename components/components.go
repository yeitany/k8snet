package components

import (
	"context"

	"github.com/yeitany/k8s_net/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetWorkloads(clientset *kubernetes.Clientset) (map[string]graph.Node, error) {
	entities := make(map[string]graph.Node, 0)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
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
	return entities, nil
}

func GetServices(clientset *kubernetes.Clientset) (map[string]graph.Node, error) {
	entities := make(map[string]graph.Node, 0)
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
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
	return entities, nil
}

func GetComponents(clientset *kubernetes.Clientset) (map[string]graph.Node, error) {
	entities, err := GetWorkloads(clientset)
	if err != nil {
		return nil, err
	}
	services, err := GetServices(clientset)
	if err != nil {
		return nil, err
	}
	for k, svc := range services {
		entities[k] = svc
	}
	return entities, nil
}
