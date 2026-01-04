package kube

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewClient() (*kubernetes.Clientset, error) {
	home := homedir.HomeDir()
	if home == "" {
		return nil, fmt.Errorf("cannot determine home directory")
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("loading kubeconfig: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	return client, nil
}
