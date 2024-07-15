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
	"fmt"
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
	log.Printf("Reconcile start at %s", time.Now())
	app := &appv1alpha1.Autostager{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, app)
	if err != nil {
		log.Println(err)
		if errors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		return ctrl.Result{}, err
	}

	if err := a.UpsertDeployment(ctx, req, app); err != nil {
		log.Println(err)
		return ctrl.Result{}, err
	}

	log.Printf("Reconcile end at %s", time.Now())
	return ctrl.Result{}, nil
}

func NewAutostagerClient(kube client.Client, scheme *runtime.Scheme) (*AutostagerClient, error) {
	return &AutostagerClient{
		Kubernetes: kube,
		Scheme:     scheme,
	}, nil
}

func (a *AutostagerClient) UpsertDeployment(ctx context.Context, req ctrl.Request, app *appv1alpha1.Autostager) error {
	deployment := &v1.Deployment{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Println("Creating a new Deployment")
			newDeployment := a.CreateNewDeployment(ctx, req, app)
			if a.Scheme == nil {
				log.Println("a.Scheme is nil")
				return fmt.Errorf("a.Scheme is nil")
			}
			if err := ctrl.SetControllerReference(app, newDeployment, a.Scheme); err != nil {
				return err
			}
			return a.Kubernetes.Create(ctx, newDeployment)
		}
		return err
	}

	if deployment.Spec.Replicas == nil {
		log.Println("deployment.Spec.Replicas is nil")
		return fmt.Errorf("deployment.Spec.Replicas is nil")
	}

	if !app.Spec.Helm.UseHelm && app.Spec.Manifest.Replicas != deployment.Spec.Replicas {
		log.Printf("Updating Deployment replicas from %d to %d", *deployment.Spec.Replicas, app.Spec.Manifest.Replicas)
		deployment.Spec.Replicas = app.Spec.Manifest.Replicas
		return a.Kubernetes.Update(ctx, deployment)
	}
	return nil
}

func (a *AutostagerClient) CreateNewDeployment(ctx context.Context, req ctrl.Request, app *appv1alpha1.Autostager) *v1.Deployment {
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: app.Spec.Manifest.Replicas,
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
							Image: *app.Spec.Manifest.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: *app.Spec.Manifest.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}

}
