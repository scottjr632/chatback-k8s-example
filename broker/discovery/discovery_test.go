package discovery

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/scottjr632/chatback/broker/config"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func newSimpleK8s(labelSelector string) *Discovery {
	client := Discovery{}
	client.client = fake.NewSimpleClientset()
	client.config = &config.Config{}
	client.config.LabelSelector = labelSelector
	return &client
}

func getPodObject() *core.Pod {
	id, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Println(err)
		return nil
	}
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("my-test-pod-%s", id),
			Namespace: "default",
			Labels: map[string]string{
				"teir": "broker",
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "busybox",
					Image:           "busybox",
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
					},
				},
			},
		},
	}
}

func createPod(d *Discovery) error {
	pod := getPodObject()
	_, err := d.client.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	return err
}

func TestCanCreatePod(t *testing.T) {
	client := newSimpleK8s("")
	pod := getPodObject()
	if pod == nil {
		t.Fatal("pod should not be nil")
	}
	if err := createPod(client); err != nil {
		t.Fatal(err)
	}
}

func TestCanListPodsByLabelSelector(t *testing.T) {
	client := newSimpleK8s("teir=broker")
	ips, err := client.GetBrokerPodsIPs()
	if err != nil {
		t.Fatal(err)
	}
	if len(ips) != 0 {
		t.Fatalf("No pods were created but received %d pods\n", len(ips))
	}

	if err := createPod(client); err != nil {
		t.Fatal(err)
	}

	ips, err = client.GetBrokerPodsIPs()
	if err != nil {
		t.Fatal(err)
	}

	if len(ips) < 1 {
		t.Fatalf("ips should be larger than 0: got %d\n", len(ips))
	}
}

func TestGetVersion(t *testing.T) {
	client := newSimpleK8s("")
	version, err := fmt.Println(client.getK8sVersion())
	if err != nil {
		t.Fatal("get version should not throw an error")
	}
	fmt.Println(version)
}

func TestCreateUUID(t *testing.T) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(out))
}
