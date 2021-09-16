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
	"github.com/go-logr/logr"
	webappv1 "github.com/ivyxjc/webapp-operator/api/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WebDeploymentReconciler reconciles a WebDeployment object
type WebDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=webapp.webapp.ivyxjc.com,resources=webdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.webapp.ivyxjc.com,resources=webdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.webapp.ivyxjc.com,resources=webdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *WebDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("webDeployment", req.NamespacedName)

	var webDeployment webappv1.WebDeployment
	if err := r.Get(ctx, req.NamespacedName, &webDeployment); err != nil {
		log.Error(err, "unable to fetch CronJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployment, err := r.constructJobForWebDeployment(&webDeployment)
	log.V(1).Info("created Web Deployment")
	if err != nil {
		log.Error(err, "unable to construct deployment")
		return ctrl.Result{}, nil
	}
	if err := r.Create(ctx, deployment); err != nil {
		log.Error(err, "unable to create deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WebDeploymentReconciler) constructJobForWebDeployment(webDeployment *webappv1.WebDeployment) (*appsv1beta1.Deployment, error) {
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        webDeployment.Name,
			Namespace:   webDeployment.Namespace,
		},
		Spec: *webDeployment.Spec.Deployment.DeepCopy(),
	}

	if err := ctrl.SetControllerReference(webDeployment, deployment, r.Scheme); err != nil {
		return nil, err
	}
	return deployment, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.WebDeployment{}).
		Complete(r)
}
