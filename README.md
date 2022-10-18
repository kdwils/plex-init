# plex-init

## Use Case
A quick CLI to claim a plex token for your server.

## Usage

Example:

This assumes there is an environment variable available named `PLEX_TOKEN` where the value is your plex API token. This is required to get a new claim token.

You can prefix the command like `PLEX_TOKEN='123' plex-init ...`

Running outside of a container
```
plex-init claim token --secret-name=<name-of-secret-to-create> --namespace=<namespace-to-create-secret-in> --kube-config='/path/to/kubeconfig'
```

When running inside of a container in a kubernetes cluster, you do not need to specify the path to the kubeconfig. Instead, the in cluster config will be used. 

If a kube-config is specified then that will be used instead of the in cluster config.

```
plex-init claim token --secret-name=<name-of-secret-to-create> --namespace=<namespace-to-create-secret-in>
```