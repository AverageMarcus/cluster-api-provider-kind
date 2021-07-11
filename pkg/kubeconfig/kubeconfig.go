package kubeconfig

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type kubeConfig struct {
	Clusters []cluster `json:"clusters"`
}

type cluster struct {
	Cluster clusterDetails `json:"cluster"`
	Name    string         `json:"name"`
}

type clusterDetails struct {
	Server string `json:"server"`
}

// ClusterEndpoint contains the control plane endpoint parts
type ClusterEndpoint struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

// ExtractEndpoint parses the provided kubeconfig and attempts to pull out the
// cluster endpoint details matching the given cluster name
func ExtractEndpoint(kubeconfig string, clusterName string) (*ClusterEndpoint, error) {
	var config kubeConfig
	err := yaml.Unmarshal([]byte(kubeconfig), &config)
	if err != nil {
		return nil, err
	}

	for _, cluster := range config.Clusters {
		if cluster.Name == fmt.Sprintf("kind-%s", clusterName) {
			serverEndpoint := cluster.Cluster.Server
			serverEndpoint = strings.TrimPrefix(serverEndpoint, "https://")
			serverEndpoint = strings.TrimPrefix(serverEndpoint, "http://")
			parts := strings.Split(serverEndpoint, ":")

			if len(parts) != 2 {
				// Unexpected server endpoint URL, lets keep looking
				continue
			}

			port, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				// Unable to parse port, lets keep looking
				continue
			}

			endpoint := &ClusterEndpoint{
				Host: parts[0],
				Port: int32(port),
			}
			return endpoint, nil
		}
	}

	return nil, fmt.Errorf("Unable to find valid server details for %s", clusterName)
}
