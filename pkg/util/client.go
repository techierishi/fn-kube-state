package util

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient() (kubernetes.Interface, error) {

	clientType := os.Getenv("KUBE_CLIENT")
	Println("clientType", clientType)

	var kubeConfig *rest.Config
	var err error
	if len(clientType) > 0 && clientType == "inCluster" {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		Println(clientcmd.RecommendedHomeFile)
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(kubeConfig)
}
