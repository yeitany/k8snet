package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetKubeProxyPods(clientset *kubernetes.Clientset) *corev1.PodList {
	kubeProxyPods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: "component=kube-proxy",
	})
	if err != nil {
		panic(err.Error())
	}
	return kubeProxyPods
}
