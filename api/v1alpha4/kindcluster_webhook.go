/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha4

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var kindclusterlog = logf.Log.WithName("kindcluster-resource")

func (r *KindCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha4-kindcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kindclusters,verbs=create;update,versions=v1alpha4,name=mkindcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &KindCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *KindCluster) Default() {
	kindclusterlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha4-kindcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kindclusters,verbs=create;update,versions=v1alpha4,name=vkindcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &KindCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KindCluster) ValidateCreate() error {
	kindclusterlog.Info("validate create", "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KindCluster) ValidateUpdate(old runtime.Object) error {
	kindclusterlog.Info("validate update", "name", r.Name)
	oldCluster := old.(*KindCluster)

	if oldCluster.Spec.Replicas != r.Spec.Replicas {
		return fmt.Errorf("Unable to modify replicas")
	}

	if oldCluster.Spec.Image != r.Spec.Image {
		return fmt.Errorf("Unable to modify image")
	}

	if oldCluster.Spec.Version != r.Spec.Version {
		return fmt.Errorf("Unable to modify version")
	}

	if !reflect.DeepEqual(oldCluster.Spec.FeatureGates, r.Spec.FeatureGates) {
		return fmt.Errorf("Unable to modify featureGates")
	}

	if !reflect.DeepEqual(oldCluster.Spec.RuntimeConfig, r.Spec.RuntimeConfig) {
		return fmt.Errorf("Unable to modify runtimeConfig")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KindCluster) ValidateDelete() error {
	kindclusterlog.Info("validate delete", "name", r.Name)
	return nil
}
