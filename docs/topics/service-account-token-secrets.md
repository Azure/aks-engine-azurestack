# Service account token secrets

A `kubernetes.io/service-account-token` type of Secret is used to store a token credential that identifies a ServiceAccount. This is a legacy mechanism that provides long-lived ServiceAccount credentials to Pods.

You can mount service account tokens as Kubernetes Secrets in Pods. These tokens don't expire and don't rotate. This method is not recommended, especially at scale, because of the risks associated with static, long-lived credentials.

## Autogeneration

Previously, these Secret-based service account tokens were auto generated for ServiceAccounts. In Kubernetes v1.27 and later, Kubernetes forced the `LegacyServiceAccountTokenNoAutoGeneration` feature gate value to true, preventing Kubernetes from automatically creating these service account token secrets for ServiceAccounts.

As a result, if using AKS Engine >= v0.79 with Kubernetes >= v1.27, Kubernetes will not automatically create these tokens.

## Switch to alternative methods in k8s v1.27+

The `LegacyServiceAccountTokenNoAutoGeneration` feature gate value is [forced by Kubernetes to true](https://github.com/kubernetes/kubernetes/pull/114522) in Kubernetes v1.27+. When enabled, Secret API objects containing service account tokens are no longer auto-generated for every ServiceAccount.

Alternatives:

- Use the [TokenRequest API](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-request-v1/) or an API client like kubectl (`kubectl create token`) to acquire service account tokens.
- Request a mounted token in a [projected volume](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#bound-service-account-token-volume) in your Pod manifest. Kubernetes creates the token and mounts it in the Pod. The token is automatically invalidated when the Pod that it's mounted in is deleted. For details, see [Launch a Pod using service account token projection](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#launch-a-pod-using-service-account-token-projection).
- If a non-expiring token is required, manually create a Secret API object for the token controller to populate with a service account token by following this [guide](https://kubernetes.io/docs/concepts/configuration/secret/#service-account-token-secrets). You should only create a ServiceAccount token Secret if you can't use the TokenRequest API to obtain a token, and the security exposure of persisting a non-expiring token credential in a readable API object is acceptable to you.
