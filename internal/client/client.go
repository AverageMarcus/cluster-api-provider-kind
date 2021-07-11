package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/AverageMarcus/cluster-api-provider-kind/api/v1alpha4"
)

var client = http.Client{Timeout: 300 * time.Second}

// CreateCluster creates a new cluster in Kind
func CreateCluster(kindCluster *v1alpha4.KindCluster) error {
	payload, err := json.Marshal(*kindCluster)
	if err != nil {
		return err
	}

	resp, err := client.Post(getAPIEndpoint(), "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		fmt.Println(resp.Status)
		return fmt.Errorf("unexpected error returned from server")
	}

	return nil
}

// IsReady checks if the cluster is ready in Kind
func IsReady(clusterName string) (bool, error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s", getAPIEndpoint(), clusterName))
	if err != nil {
		return false, err
	}
	if resp.StatusCode >= 400 {
		fmt.Println(resp.Status)
		return false, fmt.Errorf("unexpected error returned from server")
	}

	isReady := false
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(body, &isReady); err != nil {
		return false, err
	}

	return isReady, nil
}

// GetKubeConfig returns the KubeConfig for the cluster in Kind matching the given name
func GetKubeConfig(clusterName string) (string, error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s/kubeconfig", getAPIEndpoint(), clusterName))
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		fmt.Println(resp.Status)
		return "", fmt.Errorf("unexpected error returned from server")
	}

	kubeconfig := ""
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &kubeconfig); err != nil {
		return "", err
	}

	return kubeconfig, nil
}

// DeleteCluster removes the cluster from Kind
func DeleteCluster(clusterName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", getAPIEndpoint(), clusterName), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		fmt.Println(resp.Status)
		return fmt.Errorf("unexpected error returned from server")
	}

	return nil
}

func getAPIEndpoint() string {
	return fmt.Sprintf("http://%s:%s", os.Getenv("KIND_SERVER_ENDPOINT"), os.Getenv("KIND_SERVER_PORT"))
}
