package sse

import (
	"fmt"
	"testing"
)

func TestGetMyIPs(t *testing.T) {
	ips, err := getMachineIPs()
	if err != nil {
		t.Fatal(err)
	}

	if len(ips) == 0 {
		t.Fatalf("ips should be greater than 0")
	}
	fmt.Printf("len(ips) ->  %d\n", len(ips))
	for ip, _ := range ips {
		fmt.Println(ip)
	}
}
