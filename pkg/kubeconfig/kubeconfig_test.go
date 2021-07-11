package kubeconfig

import (
	"testing"
)

func TestExtractEndpoint(t *testing.T) {
	tests := []struct {
		kubeConfig  string
		clusterName string
		host        string
		port        int32
	}{
		{
			kubeConfig: `clusters:
- name: kind-test-cluster
  cluster:
    server: https://1.2.3.4:1000`,
			clusterName: "test-cluster",
			host:        "1.2.3.4",
			port:        1000,
		},
		{
			kubeConfig: `clusters:
- name: kind-test-cluster
  cluster:
    server: https://1.2.3.4:1000
- name: kind-other-cluster
  cluster:
    server: https://localhost:5000`,
			clusterName: "other-cluster",
			host:        "localhost",
			port:        5000,
		},
		{
			kubeConfig: `clusters:
- name: kind-http-endpoint
  cluster:
    server: http://1.2.3.4:1000`,
			clusterName: "http-endpoint",
			host:        "1.2.3.4",
			port:        1000,
		},
		{
			kubeConfig: `clusters:
- name: kind-duplicates
  cluster:
- name: kind-duplicates
  cluster:
    server: https://100.100.100.100:6000`,
			clusterName: "duplicates",
			host:        "100.100.100.100",
			port:        6000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.clusterName, func(t *testing.T) {
			endpoint, err := ExtractEndpoint(tt.kubeConfig, tt.clusterName)
			if err != nil {
				t.Errorf("unexpected error returned - %+v", err)
				return
			}
			if endpoint == nil {
				t.Errorf("endpoint not returned")
				return
			}
			if endpoint.Host != tt.host || endpoint.Port != tt.port {
				t.Errorf("returned endpoint details not as expected")
			}
		})
	}
}
