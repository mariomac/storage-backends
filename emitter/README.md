


To locally test:

```
docker run -d --name=loki --mount source=loki-data,target=/loki -p 3100:3100 grafana/loki
```


```
curl GET http://localhost:3100/loki/api/v1/query_range?query={baz="bar"}&step=600
```