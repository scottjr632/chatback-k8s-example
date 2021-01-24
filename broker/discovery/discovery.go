package discovery

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/scottjr632/chatback/broker/config"
)

type Discovery struct {
	client kubernetes.Interface
	config *config.Config
}

func createClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// New creates and returns a new Discovery instance
func New(config *config.Config) (*Discovery, error) {
	discovery := &Discovery{config: config}
	if client, err := createClient(); err != nil {
		return nil, err
	} else {
		discovery.client = client
		return discovery, nil
	}
}

// GetBrokerPodsIPs returns a slice of ip addresses retrieved from
// k8s API using the LabelSelector
func (d *Discovery) GetBrokerPodsIPs() ([]string, error) {
	ips := []string{}
	pods, err := d.client.CoreV1().Pods(d.config.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: d.config.LabelSelector,
	})
	if err != nil {
		return ips, err
	}

	for _, pod := range pods.Items {
		ips = append(ips, pod.Status.PodIP)
	}
	return ips, nil
}

func (d *Discovery) getK8sVersion() (string, error) {
	version, err := d.client.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", version), nil
}
