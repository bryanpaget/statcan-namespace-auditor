package controllers

import (
	"context"
	"flag"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Grace period before deleting a namespace marked for deletion
const gracePeriod = 7 * 24 * time.Hour

// Dry-run flag
var dryRun bool

func init() {
	// Define the dry-run flag
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run without deleting namespaces")
	flag.Parse()
}

// NamespaceReconciler reconciles Namespace objects
type NamespaceReconciler struct {
	client.Client
	// Add any dependencies like a logger or Entra ID client here
}

// Reconcile is the main reconciliation loop for the controller
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Namespace instance
	var namespace corev1.Namespace
	if err := r.Get(ctx, req.NamespacedName, &namespace); err != nil {
		logger.Error(err, "unable to fetch Namespace")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Skip system namespaces or namespaces that are terminating
	if namespace.Status.Phase == corev1.NamespaceTerminating || isSystemNamespace(namespace.Name) {
		return ctrl.Result{}, nil
	}

	// Extract email from namespace annotations
	email := namespace.Annotations["user-email"]
	if email == "" || !isStatCanEmail(email) {
		return ctrl.Result{}, nil
	}

	// Check if the email exists in Entra ID
	exists, err := checkEmailInEntraID(email)
	if err != nil {
		logger.Error(err, "failed to check email in Entra ID")
		return ctrl.Result{}, err
	}

	if !exists {
		// Mark namespace for deletion if email doesn't exist
		if namespace.Annotations["marked-for-deletion"] == "" {
			if dryRun {
				logger.Info("[Dry Run] Would mark namespace for deletion", "namespace", namespace.Name)
			} else {
				patch := client.MergeFrom(namespace.DeepCopy())
				namespace.Annotations["marked-for-deletion"] = time.Now().Format(time.RFC3339)
				if err := r.Patch(ctx, &namespace, patch); err != nil {
					logger.Error(err, "failed to mark namespace for deletion")
					return ctrl.Result{}, err
				}
				logger.Info("marked namespace for deletion", "namespace", namespace.Name)
			}
			return ctrl.Result{RequeueAfter: gracePeriod}, nil
		}

		// Check if grace period has passed
		markedTime, err := time.Parse(time.RFC3339, namespace.Annotations["marked-for-deletion"])
		if err != nil {
			logger.Error(err, "invalid marked-for-deletion timestamp")
			return ctrl.Result{}, err
		}
		if time.Since(markedTime) > gracePeriod {
			if dryRun {
				logger.Info("[Dry Run] Would delete namespace", "namespace", namespace.Name)
			} else {
				if err := r.Delete(ctx, &namespace); err != nil {
					logger.Error(err, "failed to delete namespace")
					return ctrl.Result{}, err
				}
				logger.Info("deleted namespace", "namespace", namespace.Name)
			}
			return ctrl.Result{}, nil
		}
	}

	// Requeue to recheck later
	return ctrl.Result{RequeueAfter: 24 * time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}

// isSystemNamespace checks if a namespace is a system namespace
func isSystemNamespace(namespace string) bool {
	systemNamespaces := []string{"kube-system", "kube-public", "default"}
	for _, ns := range systemNamespaces {
		if namespace == ns {
			return true
		}
	}
	return false
}

// isStatCanEmail checks if an email belongs to the StatCan domain
func isStatCanEmail(email string) bool {
	return strings.HasSuffix(email, "@statcan.gc.ca")
}

// checkEmailInEntraID simulates a call to Entra ID to check if an email exists
func checkEmailInEntraID(email string) (bool, error) {
	// Placeholder for actual Entra ID API call
	return false, nil
}
