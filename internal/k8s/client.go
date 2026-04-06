package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client wraps Kubernetes clientset
type Client struct {
	Clientset *kubernetes.Clientset
	Config    *rest.Config
	Namespace string
}

// NewClient creates a new Kubernetes client
// Priority: in-cluster config > kubeconfig file
func NewClient(namespace string) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		// Running in-cluster
		return newClientFromConfig(config, namespace)
	}

	// Fallback to kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	if kubeconfig == "" {
		return nil, fmt.Errorf("no kubeconfig found")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	return newClientFromConfig(config, namespace)
}

func newClientFromConfig(config *rest.Config, namespace string) (*Client, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	if namespace == "" {
		namespace = "default"
	}

	return &Client{
		Clientset: clientset,
		Config:    config,
		Namespace: namespace,
	}, nil
}
