# plex-init

## Use Case
I use this as an init container when deploying plex to my home kubernetes cluster so I don't have to manually get a new token. The init container will create or update an existing kubernetes secret with the plex claim token as the value.

I then use this secret as an environment variable for my plex container to claim the server.

## Usage

Example:

This assumes there is an environment variable available named `PLEX_TOKEN` where the value is your plex token. This is required to get a new claim token.

Running outside of a container
```
plex-init claim token --secret-name=<name-of-secret-to-create> --namespace=<namespace-to-create-secret-in> --kube-config='/path/to/kubeconfig'
```

When running inside of a container in a kubernetes cluster, you do not need to specify the path to the kubeconfig. Instead, the in cluster config will be used. 

If a kube-config is specified then that will be used instead of the in cluster config.

```
plex-init claim token --secret-name=<name-of-secret-to-create> --namespace=<namespace-to-create-secret-in>
```

## Init container in a cluster

An example init container that sets the PLEX_TOKEN environment variable from a kubernetes secret:

I use rbac to create a service account for the init container
```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: plex-init-service-account
  namespace: media
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: plex-init-cluster-role
  namespace: media
rules:
  - apiGroups:
        - ""
    resources:
      - secrets
    verbs: ["create", "update"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: plex-init-cluster-role-binding
subjects:
- namespace: media 
  kind: ServiceAccount
  name: plex-init-service-account
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: plex-init-cluster-role
```

Then use that service account with the init container
```
serviceAccountName: plex-init-service-account
initContainers:
- name: init-plex-claim
    image: kdwilson/plex-init
    args:
    [
        "claim",
        "token",
        "--secret-name=plex-claim",
        "--namespace=media",
    ]
    env:
    - name: PLEX_TOKEN
        valueFrom:
        secretKeyRef:
            key: token
            name: plex
```
## Using the secret
The secret created will look something like this
```
apiVersion: v1
data:
  token: b64 encoded string
kind: Secret
metadata:
  name: plex-claim
  namespace: media
type: Opaque
```

Once the secret has been created, we can pass it as an environment variable to other containers.

For example, consider this as the env block for a container running in the same namespace the secret was created in.
```
env:
- name: PLEX_CLAIM
    valueFrom:
    secretKeyRef:
        key: token
        name: plex-claim
```