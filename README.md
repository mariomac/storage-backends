# storage-backends
Testing backends for OVN-K8s network flows storage

https://github.com/jotak/loki-sandbox

## Install steps

1. Deploy a cluster with at least 6 worker nodes
2. Install loki

```
helm install loki -f install/helm-values-override.yaml grafana/loki-distributed  --set access_key=... --set secret_access_key=... --set bucketname=...
oc apply -f install/emitters.yaml
```