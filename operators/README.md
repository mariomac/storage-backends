# Loki as operator howto

Clone [Loki Operator](https://github.com/ViaQ/loki-operator) (under development)

Install it locally by running, from the operator repo:
```
make olm-deploy REGISTRY_ORG=<whatever> VERSION=v0.0.1
```

Add your S3 secrets:

```
cp secrets.yml.template secrets.yml
vim secrets.yml
oc apply -f secrets.yml
```

Deploy the lokistack operator:

```
oc apply -f lokistack.yml
```

# Feedback

- If replicationFactor: 1, deploy loki-standalone (all in a single pod).
- Error deploying in small cluster: `0/3 nodes are available: 3 pod has unbound immediate PersistentVolumeClaims.`
- 