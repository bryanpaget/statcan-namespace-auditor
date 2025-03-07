StatCan Namespace Auditor
=========================

Objective
---------
The StatCan Namespace Auditor is a Kubernetes controller written in Go designed to ensure that namespaces associated with Statcan.gc.ca email accounts remain up to date. The controller periodically audits namespaces to verify that the corresponding user exists in the StatCan EntraID tenant (via Microsoft Graph API). If a namespace is associated with an email that no longer exists in EntraID, the namespace is marked for deletion and eventually removed after a safety (grace) period.

Background & Logic Considerations
---------------------------------
- StatCan Email Requirement: All StatCan employees should have a Statcan.gc.ca email registered in the STC Tenant EntraID.
- Zone User Variability: Not all zone users will have a Statcan.gc.ca email (e.g., our team might use cloud.statcan.ca). Only users with Statcan.gc.ca addresses in their rolebindings are considered.
- Automation & Safety Window: The process includes:
  1. Querying namespaces for associated user accounts.
  2. Validating each account against the EntraID service (using Microsoft Graph API).
  3. Marking the namespace for deletion if the corresponding account is not found.
  4. Deleting the namespace only after a configurable grace period, providing a safety window to recover from potential errors.

Features
--------
- Automated Auditing: Periodically scans namespaces for user accounts with a "user-email" label (restricted to Statcan.gc.ca addresses).
- Graph API Integration: Uses the Microsoft Graph API to verify the existence of the user in EntraID. (Note: Only accounts with Statcan.gc.ca emails are checked.)
- Grace Period Enforcement: Applies a configurable grace period (e.g., 48 hours) before deletion, allowing time for review or rollback if needed.
- Logging & Audit Trail: Maintains detailed logs of all actions (marking, deletion, errors) for auditing and troubleshooting.

Project Structure
-----------------
The project follows a modular structure to clearly separate controller logic, EntraID integration, and Kubernetes deployment manifests.

Prerequisites
-------------
- Go 1.18+
- Access to a Kubernetes cluster
- kubectl installed and configured
- Azure AD credentials (Tenant ID, Client ID, Client Secret) with appropriate permissions to query the Microsoft Graph API for user details

Installation & Deployment
-------------------------
1. Clone the Repository:
   git clone https://github.com/bryanpaget/statcan-namespace-auditor.git
   cd statcan-namespace-auditor

2. Build the Controller:
   go build -o statcan-namespace-auditor .

3. Configure EntraID Credentials:
   Set the following environment variables or use Kubernetes Secrets:
     - TENANT_ID
     - CLIENT_ID
     - CLIENT_SECRET

4. Deploy to Kubernetes:
   Use the provided manifests in the "config/" directory. For example, deploy the manager using:
     kubectl apply -f config/manager/manager.yaml

Development
-----------
- Reconciliation Logic:
  * The controller watches namespaces with a "user-email" label.
  * It performs EntraID validation via Microsoft Graph API.
  * Namespaces flagged for deletion are annotated with a timestamp.
  * Periodically re-evaluates flagged namespaces and deletes them once the grace period has elapsed.
- Testing:
  * Run unit tests locally using: go test ./...
  * Consider a staging environment deployment to verify real-world behavior.

Contributing
------------
Contributions and feedback are welcome! Please fork the repository, create feature branches for your changes, and open a pull request with detailed descriptions. For any issues or enhancement ideas, open an issue on the GitHub repository.

License
-------
This project is licensed under the MIT License.

Further Resources
-----------------
- Kubebuilder Documentation: https://book.kubebuilder.io/
- Controller-runtime Library: https://github.com/kubernetes-sigs/controller-runtime
- Microsoft Graph API Documentation: https://learn.microsoft.com/en-us/graph/overview


Proposed Project File Structure
=================================

statcan-namespace-auditor/
├── README.txt                  # Project overview, setup instructions, and usage guidelines
├── go.mod                      # Go module file
├── go.sum                      # Dependency checksums
├── main.go                     # Entry point; initializes the controller manager and sets up logging
├── controllers/                # Contains the controller logic and integrations
│   ├── namespace_controller.go   # Reconciliation logic for auditing and cleaning namespaces
│   ├── entra.go               # Functions for obtaining tokens and querying Microsoft Graph API (EntraID)
│   └── types.go               # (Optional) Definitions for custom types or CRD schemas
├── config/                     # Kubernetes manifests for deploying the controller
│   ├── crd/                   # Custom Resource Definitions (if any)
│   ├── rbac/                  # RBAC roles and bindings needed for the controller
│   └── manager/               # Deployment manifests for the controller manager (e.g., manager.yaml)
└── pkg/                        # (Optional) Additional libraries or utility functions
