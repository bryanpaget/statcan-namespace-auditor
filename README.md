# Namespace Cleanup Controller

This controller audits namespaces in a Kubernetes cluster based on associated emails in Entra ID (Microsoft Azure AD) and deletes namespaces if the associated email no longer exists.

## How It Works

1. **Configuration:**  
   Uses a ConfigMap (`namespace-cleanup-config`) for settings:
   - `emailDomain` (e.g., `statcan.gc.ca`)
   - `gracePeriodDays` (default: 30 days)
   - `schedule` (Cron expression)

2. **Process Flow:**
   - Queries all namespaces.
   - Checks if the associated email exists in Entra ID using Graph API.
   - Marks namespaces for deletion if the email does not exist.
   - Deletes namespaces after the grace period.

3. **RBAC Requirements:**
   - Needs permissions for namespaces and ConfigMaps.

## Configuration Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: namespace-cleanup-config
  namespace: kube-system
data:
  emailDomain: "statcan.gc.ca"
  gracePeriodDays: "30"
  schedule: "0 0 * * 0"
```

## Deployment

Apply the following manifests:

```bash
kubectl apply -f config/configmap.yaml
kubectl apply -f config/rbac.yaml
kubectl apply -f config/deployment.yaml
```
