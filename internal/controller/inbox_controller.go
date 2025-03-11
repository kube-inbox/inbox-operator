/*
Copyright 2025.

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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	inboxv1 "kubeinbox.com/inbox-operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// InboxReconciler reconciles a Inbox object
type InboxReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=api.kubeinbox.com,resources=inboxes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.kubeinbox.com,resources=inboxes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.kubeinbox.com,resources=inboxes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Inbox object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *InboxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Begin reconciling Inbox", "namespace", req.Namespace, "name", req.Name)

	// Fetch the Inbox instance
	var inbox inboxv1.Inbox
	if err := r.Get(ctx, req.NamespacedName, &inbox); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "Unable to fetch Inbox")
			return ctrl.Result{}, err
		}
		// Object not found, return
		return ctrl.Result{}, nil
	}

	// Initialize the status if it's not set
	if inbox.Status.Phase == "" {
		inbox.Status.Phase = "Pending"
		if err := r.Status().Update(ctx, &inbox); err != nil {
			logger.Error(err, "Failed to update Inbox status")
			return ctrl.Result{}, err
		}
	}

	// Add finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(&inbox, "inbox.kubeinbox.com/finalizer") {
		controllerutil.AddFinalizer(&inbox, "inbox.kubeinbox.com/finalizer")
		if err := r.Update(ctx, &inbox); err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
	}

	// Handle deletion
	if !inbox.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, &inbox)
	}

	// Update status based on the inbox type
	inbox.Status.Phase = "Ready"
	now := metav1.NewTime(time.Now())
	inbox.Status.LastUpdated = &now

	if err := r.Status().Update(ctx, &inbox); err != nil {
		logger.Error(err, "Failed to update Inbox status")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled Inbox",
		"namespace", req.Namespace,
		"name", req.Name,
		"type", inbox.Spec.InboxType,
		"phase", inbox.Status.Phase)

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// handleDeletion handles the cleanup when an Inbox is being deleted
func (r *InboxReconciler) handleDeletion(ctx context.Context, inbox *inboxv1.Inbox) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Update the status to show it's being deleted
	inbox.Status.Phase = "Deleting"
	if err := r.Status().Update(ctx, inbox); err != nil {
		logger.Error(err, "Failed to update Inbox status during deletion")
		return ctrl.Result{}, err
	}

	// Perform cleanup logic here if needed
	// For example, you might want to clean up any resources created by this Inbox

	// Remove our finalizer from the list and update it
	controllerutil.RemoveFinalizer(inbox, "inbox.kubeinbox.com/finalizer")
	if err := r.Update(ctx, inbox); err != nil {
		logger.Error(err, "Failed to remove finalizer")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully deleted Inbox", "namespace", inbox.Namespace, "name", inbox.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inboxv1.Inbox{}).
		Named("inbox").
		Complete(r)
}
