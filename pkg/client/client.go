package client

import (
	"fmt"
	"os"

	"github.com/caicloud/cyclone/pkg/k8s/clientset"
	"github.com/cd1989/cycli/pkg/console"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var K8sClient clientset.Interface

func init() {
	homepath := os.Getenv("HOME")
	if homepath == "" {
		console.Error("Environment variable $HOME not set")
		os.Exit(1)
	}

	var err error
	K8sClient, err = getClient(fmt.Sprintf("%s/.kube/config", homepath))
	if err != nil {
		console.Error("Create k8s client error: ", err)
		os.Exit(1)
	}
}

// getClient creates a client for k8s cluster
func getClient(kubeConfigPath string) (clientset.Interface, error) {
	var config *rest.Config
	var err error
	if kubeConfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return clientset.NewForConfig(config)
}
