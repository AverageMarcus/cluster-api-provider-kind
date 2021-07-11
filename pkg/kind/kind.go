package kind

import (
	"fmt"
	"os"
	"path"
	"time"

	kindcluster "github.com/AverageMarcus/cluster-api-provider-kind/api/v1alpha4"
	"github.com/go-logr/logr"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
)

const createWaitTime = 60 * time.Second

// Kind provides function for interacting with Kind clusters
type Kind struct {
	provider *cluster.Provider
}

// New create a new instance of Kind
func New(log logr.Logger) *Kind {
	return &Kind{
		provider: cluster.NewProvider(cluster.ProviderWithLogger(Logger{log})),
	}
}

// CreateCluster creates a new cluster in Kind
func (k *Kind) CreateCluster(kindCluster *kindcluster.KindCluster) error {
	return k.provider.Create(
		kindCluster.NamespacedName(),
		cluster.CreateWithV1Alpha4Config(kindClusterToKindConfig(kindCluster)),
		cluster.CreateWithWaitForReady(createWaitTime),
		cluster.CreateWithKubeconfigPath(path.Join(os.TempDir(), "kubeconfig")),
		cluster.CreateWithRetain(false),
		cluster.CreateWithDisplayUsage(false),
		cluster.CreateWithDisplaySalutation(false),
	)
}

// GetKubeConfig returns the KubeConfig for the cluster in Kind matching the given name
func (k *Kind) GetKubeConfig(clusterName string) (string, error) {
	return k.provider.KubeConfig(clusterName, false)
}

// IsReady checks if the cluster is ready in Kind
func (k *Kind) IsReady(clusterName string) (bool, error) {
	readyClusters, err := k.provider.List()
	if err != nil {
		return false, err
	}
	for _, readyCluster := range readyClusters {
		if readyCluster == clusterName {
			return true, nil
		}
	}
	return false, nil
}

// DeleteCluster removes the cluster from Kind
func (k *Kind) DeleteCluster(clusterName string) error {
	return k.provider.Delete(clusterName, path.Join(os.TempDir(), "kubeconfig"))
}

func kindClusterToKindConfig(kindCluster *kindcluster.KindCluster) *v1alpha4.Cluster {
	replicas := 1
	featureGates := map[string]bool{}
	runtimeConfig := map[string]string{}
	image := "kindest/node"
	version := "v1.21.2"

	if kindCluster.Spec.Replicas > 0 {
		replicas = int(kindCluster.Spec.Replicas)
	}

	if kindCluster.Spec.FeatureGates != nil {
		featureGates = kindCluster.Spec.FeatureGates
	}

	if kindCluster.Spec.RuntimeConfig != nil {
		runtimeConfig = kindCluster.Spec.RuntimeConfig
	}

	if kindCluster.Spec.Image != "" {
		image = kindCluster.Spec.Image
	}

	if kindCluster.Spec.Version != "" {
		version = kindCluster.Spec.Version
	}

	nodes := []v1alpha4.Node{}
	for i := 0; i < replicas; i++ {
		nodes = append(nodes, v1alpha4.Node{
			Role:  v1alpha4.ControlPlaneRole,
			Image: fmt.Sprintf("%s:%s", image, version),
		})
	}

	return &v1alpha4.Cluster{
		FeatureGates:  featureGates,
		RuntimeConfig: runtimeConfig,
		Nodes:         nodes,
	}
}
