


To locally test:

```
docker run -d --name=loki --mount source=loki-data,target=/loki -p 3100:3100 grafana/loki
```

Example query:

```
GET http://localhost:3100/loki/api/v1/query_range?query={source="fluentd"}|="\"DstAddr\":\"172.10.6.2\""
```

to locally build:

```
docker build --tag=quay.io/mmaciasl/fake-flow-loki-emitter:latest . 
```

CONFIG options

* `LOKI_HOST` (default: `http://localhost:3100`)
* `PODS` (default: `20`)
* `HASH_PODS_BASE` whether has to generate a base pods address hashed from the hostname (default: `false`)
* `FLOWS_PER_SECOND` (default: `2000`)
* `CONCURRENT` whether flows must be pushed in background while generating others (default: `true`)