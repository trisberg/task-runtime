/*

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

	buildv1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"

	"k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	taskv1alpha1 "github.com/projectriff/task/api/v1alpha1"
)

// TaskExecutionReconciler reconciles a TaskExecution object
type TaskExecutionReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=task.projectriff.io,resources=taskexecutions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=task.projectriff.io,resources=taskexecutions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=build.projectriff.io,resources=containers,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods/logs,verbs=get;list;watch

func (r *TaskExecutionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("taskexecution", req.NamespacedName)

	// your logic here
	var taskExec taskv1alpha1.TaskExecution
	if err := r.Get(ctx, req.NamespacedName, &taskExec); err != nil {
		log.Error(err, "unable to fetch TaskExecution")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}
	log.Info("TaskExecution", "name", taskExec.Name)
	log.Info("TaskExecution", "taskLauncher", taskExec.Spec.TaskLauncherRef)
	namespacedTaskLauncherRef := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      taskExec.Spec.TaskLauncherRef,
	}
	var taskLauncher taskv1alpha1.TaskLauncher
	if err := r.Get(ctx, namespacedTaskLauncherRef, &taskLauncher); err != nil {
		log.Error(err, "unable to find TaskLauncherRef")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}
	log.Info("TaskLauncherRef", "name", taskLauncher.Name)
	log.Info("TaskLauncherRef", "containerRef", taskLauncher.Spec.Build.ContainerRef)

	namespacedContainerRef := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      taskLauncher.Spec.Build.ContainerRef,
	}

	var container buildv1alpha1.Container
	if err := r.Get(ctx, namespacedContainerRef, &container); err != nil {
		log.Error(err, "unable to find ContainerRef")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}
	image := container.Spec.Image

	log.Info("ContainerRef", "image", image)

	taskPod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskExec.Name,
			Namespace: req.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "task",
					Image: image,
				},
			},
			RestartPolicy: "Never",
		},
	}
	err := r.Create(ctx, taskPod)
	if err != nil {
		log.Error(err, "CREATING POD!", "name", taskExec.Name)
	}

	return ctrl.Result{}, nil
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

func (r *TaskExecutionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1alpha1.TaskExecution{}).
		Complete(r)
}
