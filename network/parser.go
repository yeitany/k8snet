package network

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/yeitany/k8s_net/graph"
	"github.com/yeitany/k8s_net/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func ParseConntrackMeta(conntrackMeta []Meta) map[string]graph.Edge {
	edges := make(map[string]graph.Edge, 0)
	for i := range conntrackMeta {
		src := conntrackMeta[i].Layer3.Src
		dst := conntrackMeta[i].Layer3.Dst
		if src == "" || dst == "" {
			continue
		}
		key := fmt.Sprintf("%v:%v", src, dst)
		if _, ok := edges[key]; !ok {
			edges[key] = graph.Edge{
				Src: conntrackMeta[i].Layer3.Src,
				Dst: conntrackMeta[i].Layer3.Dst,
			}
		}
	}
	return edges
}

func SyncConntracks(clientset *kubernetes.Clientset, config *rest.Config) []Meta {
	kubeProxyPods := utils.GetKubeProxyPods(clientset)
	conntrackMeta := make([]Meta, 0)
	for _, pod := range kubeProxyPods.Items {
		req := clientset.CoreV1().
			RESTClient().
			Post().
			Namespace(pod.Namespace).
			Resource("pods").
			Name(pod.Name).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: pod.Spec.Containers[0].Name,
				Command:   []string{"conntrack", "-L", "-o", "xml"},
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       true,
			}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
		if err != nil {
			panic(err.Error())
		}
		var stdout, stderr bytes.Buffer
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: &stdout,
			Stderr: &stderr,
			Tty:    false,
		})
		if err != nil {
			panic(err.Error())
		}
		var x Conntrack
		err = xml.Unmarshal(stdout.Bytes(), &x)
		if err != nil {
			panic(err.Error())
		}
		conntrackMeta = append(conntrackMeta, x.Flow.Meta...)
	}
	return conntrackMeta
}
