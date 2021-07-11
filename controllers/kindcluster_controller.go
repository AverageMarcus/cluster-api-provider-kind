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

package controllers

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/AverageMarcus/cluster-api-provider-kind/api/v1alpha4"
	infrastructurev1alpha4 "github.com/AverageMarcus/cluster-api-provider-kind/api/v1alpha4"
	"github.com/AverageMarcus/cluster-api-provider-kind/pkg/kind"
	"github.com/AverageMarcus/cluster-api-provider-kind/pkg/kubeconfig"
	"github.com/AverageMarcus/cluster-api-provider-kind/pkg/utils"
)

// KindClusterReconciler reconciles a KindCluster object
type KindClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const finalizerName = "kindcluster.cluster.x-k8s.io/finalizer"

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KindClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("kindcluster", req.NamespacedName)

	// Fetch the KindCluster instance
	kindCluster := &infrastructurev1alpha4.KindCluster{}
	if err := r.Get(ctx, req.NamespacedName, kindCluster); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch KindCluster")
			return ctrl.Result{}, err
		}

		// Cluster no longer exists so lets stop now
		return ctrl.Result{}, nil
	}

	// Fetch the owner Cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, kindCluster.ObjectMeta)
	if err != nil {
		log.Error(err, "failed to get owner cluster")
		return ctrl.Result{}, err
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, kindCluster) {
		log.Info("KindCluster or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", kindCluster.Name)
	helper, err := patch.NewHelper(kindCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to init patch helper")
	}

	k := kind.New(log)

	// Ensure we always patch the resource with the latest changes when exiting function
	defer func() {
		helper.Patch(
			context.TODO(),
			kindCluster,
			patch.WithOwnedConditions{
				Conditions: []clusterv1.ConditionType{
					clusterv1.ReadyCondition,
				}},
		)
	}()

	if !kindCluster.ObjectMeta.DeletionTimestamp.IsZero() {
		// The KindCluster is being deleted
		if controllerutil.ContainsFinalizer(kindCluster, finalizerName) {
			log.Info("deleting cluster")

			kindCluster.Status.Phase = infrastructurev1alpha4.KindClusterPhaseDeleting
			kindCluster.Status.Ready = false
			if err := helper.Patch(ctx, kindCluster); err != nil {
				log.Error(err, "failed to update KindCluster status")
				return ctrl.Result{}, err
			}

			if err := k.DeleteCluster(kindCluster.NamespacedName()); err != nil {
				log.Error(err, "failed to delete cluster")
				kindCluster.Status.FailureReason = v1alpha4.FailureReasonDeleteFailed
				kindCluster.Status.FailureMessage = utils.StringPtr(err.Error())
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(kindCluster, finalizerName)
			log.Info("removed finalizer")

			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	// Ensure our finalizer is present
	controllerutil.AddFinalizer(kindCluster, finalizerName)
	if err := helper.Patch(ctx, kindCluster); err != nil {
		return ctrl.Result{}, err
	}

	if kindCluster.Status.Phase == "" {
		kindCluster.Status.Phase = infrastructurev1alpha4.KindClusterPhaseCreating
		if err := helper.Patch(ctx, kindCluster); err != nil {
			log.Error(err, "failed to update KindCluster status")
			return ctrl.Result{}, err
		}

		if err := k.CreateCluster(kindCluster); err != nil {
			log.Error(err, "failed to create cluster in kind")
			kindCluster.Status.FailureReason = v1alpha4.FailureReasonCreateFailed
			kindCluster.Status.FailureMessage = utils.StringPtr(err.Error())
			return ctrl.Result{}, err
		}

		kindCluster.Status.Ready = true
		kindCluster.Status.Phase = infrastructurev1alpha4.KindClusterPhaseReady
		if err := helper.Patch(ctx, kindCluster); err != nil {
			log.Error(err, "failed to update KindCluster status")
			return ctrl.Result{}, err
		}
	}

	// Ensure ready status is up-to-date
	isReady, err := k.IsReady(kindCluster.NamespacedName())
	if err != nil {
		log.Error(err, "failed to check status of cluster")
		kindCluster.Status.FailureReason = v1alpha4.FailureReasonClusterNotFound
		kindCluster.Status.FailureMessage = utils.StringPtr(err.Error())
		return ctrl.Result{}, err
	}
	kindCluster.Status.Ready = isReady
	if isReady {
		kindCluster.Status.Phase = infrastructurev1alpha4.KindClusterPhaseReady
	} else {
		kindCluster.Status.Phase = infrastructurev1alpha4.KindClusterPhaseCreating
	}

	// Ensure kubeconfig is up-to-date
	kc, err := k.GetKubeConfig(kindCluster.NamespacedName())
	if err != nil {
		log.Error(err, "failed to check status of cluster")
		kindCluster.Status.FailureReason = v1alpha4.FailureReasonKubeConfig
		kindCluster.Status.FailureMessage = utils.StringPtr(err.Error())
		return ctrl.Result{}, err
	}
	kindCluster.Status.KubeConfig = &kc

	// Populate the server endpoint details
	endpoint, err := kubeconfig.ExtractEndpoint(kc, kindCluster.NamespacedName())
	if err != nil {
		log.Error(err, "failed to get control plane endpoint")
		kindCluster.Status.FailureReason = v1alpha4.FailureReasonEndpoint
		kindCluster.Status.FailureMessage = utils.StringPtr(err.Error())
		return ctrl.Result{}, err
	}
	kindCluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
		Host: endpoint.Host,
		Port: endpoint.Port,
	}

	if err := helper.Patch(ctx, kindCluster); err != nil {
		log.Error(err, "failed to update KindCluster status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KindClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha4.KindCluster{}).
		Complete(r)
}
