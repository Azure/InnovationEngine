package kube

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Obtains the Kubernetes clientset based on the environment
// this function is executing in.
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if _, err := rest.InClusterConfig(); err != nil {
		kubeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
