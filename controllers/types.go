package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceStatus represents the status of a namespace in the cleanup process
type NamespaceStatus string

const (
	// StatusActive indicates the namespace is currently active
	StatusActive NamespaceStatus = "Active"

	// StatusMarkedForDeletion indicates the namespace is marked for deletion
	StatusMarkedForDeletion NamespaceStatus = "MarkedForDeletion"

	// StatusDeleted indicates the namespace has been deleted
	StatusDeleted NamespaceStatus = "Deleted"
)

// CleanupAnnotation is the annotation key used to mark namespaces for cleanup
const CleanupAnnotation = "statcan.gc.ca/cleanup"

// GracePeriodAnnotation is the annotation key for storing grace period start time
const GracePeriodAnnotation = "statcan.gc.ca/grace-period-start"

// CleanupGracePeriod defines the default grace period for namespace deletion (in days)
const CleanupGracePeriod = 30

// NamespaceInfo stores information about a namespace and its status
type NamespaceInfo struct {
	Name             string
	OwnerEmail       string
	Status           NamespaceStatus
	GracePeriodStart *metav1.Time
}

// EmailVerificationResult represents the result of checking an email against Entra ID
type EmailVerificationResult struct {
	EmailExists bool
	Error       error
}

// NamespaceList is a custom type to store a list of namespaces with metadata
type NamespaceList struct {
	Items []NamespaceInfo
}

// UserRoleBindingInfo represents information about a RoleBinding and associated email
type UserRoleBindingInfo struct {
	Namespace string
	Email     string
	Role      string
}

// UserRoleBindingList is a list of UserRoleBindingInfo objects
type UserRoleBindingList struct {
	Items []UserRoleBindingInfo
}
