package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/AverageMarcus/cluster-api-provider-kind/api/v1alpha4"
)

var (
	ts       *httptest.Server
	response = ""
)

func init() {
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))

	url := strings.TrimPrefix(ts.URL, "http://")
	parts := strings.Split(url, ":")

	os.Setenv("KIND_SERVER_ENDPOINT", parts[0])
	os.Setenv("KIND_SERVER_PORT", parts[1])
}

func TestCreateCluster(t *testing.T) {
	cluster := &v1alpha4.KindCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "default",
		},
	}

	if err := CreateCluster(cluster); err != nil {
		t.Errorf("unexpected error when creating cluster - %+v", err)
	}
}

func TestIsReady(t *testing.T) {
	response = "true"
	result, err := IsReady("test-cluster")
	if err != nil {
		t.Errorf("unexpected error when getting status - %+v", err)
	}
	if !result {
		t.Errorf("was expecting the cluster to be marked as ready")
	}

	response = "false"
	result, err = IsReady("test-cluster")
	if err != nil {
		t.Errorf("unexpected error when getting status - %+v", err)
	}
	if result {
		t.Errorf("was expecting the cluster to be marked as not ready")
	}

	response = ""
	result, err = IsReady("test-cluster")
	if err == nil {
		t.Errorf("was expecting an error when getting status")
	}
}

func TestGetKubeConfig(t *testing.T) {
	response = "\"test\""
	result, err := GetKubeConfig("test-cluster")
	if err != nil {
		t.Errorf("unexpected error when getting status - %+v", err)
	}
	if result != "test" {
		t.Errorf("unexpected value returned")
	}
}
