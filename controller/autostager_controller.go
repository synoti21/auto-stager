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

//+kubebuilder:rbac:groups=autostager.autostager.com,resources=autostagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autostager.autostager.com,resources=autostagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autostager.autostager.com,resources=autostagers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

func (r *AutostagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Autostager.AutostagerClient.Reconcile(ctx, req)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutostagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Autostager{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
