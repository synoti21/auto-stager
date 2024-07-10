/*
Copyright 2024.

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

package controller

import (
	"context"

	appv1alpha1 "github.com/synoti21/auto-stager/api/v1alpha1"
	autostager "github.com/synoti21/auto-stager/internal"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AutostagerReconciler reconciles a Autostager object
type AutostagerReconciler struct {
	client.Client
	Autostager *autostager.Manager
	Scheme     *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.synoti21,resources=autostagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.synoti21,resources=autostagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.synoti21,resources=autostagers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Autostager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *AutostagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Autostager.AutostagerClient.Reconcile(ctx, req)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutostagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Autostager{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
