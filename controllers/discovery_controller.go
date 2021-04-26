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
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	discoveryv1alpha1 "github.com/liuqianhong6007/discovery-operator/api/v1alpha1"
)

// DiscoveryReconciler reconciles a Discovery object
type DiscoveryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=discovery.lqh.xzy,resources=discoveries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=discovery.lqh.xzy,resources=discoveries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=discovery.lqh.xzy,resources=discoveries/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Discovery object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *DiscoveryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("discovery", req.NamespacedName)

	// your logic here
	instance := discoveryv1alpha1.Discovery{}
	if err := r.Get(context.TODO(), req.NamespacedName, &instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	deploy := appsv1.Deployment{}
	if err := r.Get(context.TODO(), req.NamespacedName, &deploy); err != nil && errors.IsNotFound(err) {
		// create deploy
		newDeploy := appsv1.Deployment{
			TypeMeta:   instance.TypeMeta,
			ObjectMeta: instance.ObjectMeta,
			Spec:       instance.Spec.DeploymentSpec,
		}
		if err = r.Client.Create(context.TODO(), &newDeploy); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DiscoveryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&discoveryv1alpha1.Discovery{}).
		Complete(r)
}
