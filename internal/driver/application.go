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

package driver

import (
	"context"
	"time"

	"log"

	appv1alpha1 "github.com/synoti21/auto-stager/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AutostagerReconciler reconciles a Autostager object
type AutostagerClient struct {
	Kubernetes client.Client
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
func (a *AutostagerClient) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Println("Reconcile start at %s", time.Now())
	app := &appv1alpha1.Autostager{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, app)
	if err != nil { // No Kubernetes Cluster is found
		log.Println(err)
		if errors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (a *AutostagerClient) UpsertDeployment(ctx context.Context, req ctrl.Request, app *appv1alpha1.Autostager) error {
	deployment := &v1.Deployment{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Println("Making a Deployment")
			newDeployment := a.CreateNewDeployment(ctx, req, app)
			_ = ctrl.SetControllerReference(app, newDeployment, a.Scheme)
			return a.Kubernetes.Create(ctx, newDeployment)
		}
	}

	if app.Spec.Replicas != *&deployment.Spec.Replicas {
		deployment.Spec.Replicas = app.Spec.Replicas
		return a.Kubernetes.Update(ctx, deployment)
	}
	return err
}

func (a *AutostagerClient) CreateNewDeployment(ctx context.Context, req ctrl.Request, app *appv1alpha1.Autostager) *v1.Deployment {
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": req.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": req.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  req.Name,
							Image: app.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: app.Spec.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}

}
