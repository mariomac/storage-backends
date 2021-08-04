


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