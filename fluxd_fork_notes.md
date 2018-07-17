# Docker

```bash
docker run --network host --expose 8093 akurilov/fluxd
```

# Memory Consumption

To reproduce try the script `tsdb_fill.sh` to prepare the 100K of points with total cardinality
of 100K and time step of 5 minutes.

```bash
curl -XPOST --data-urlencode 'q=from(db: "test_db") |> range(start: 1970-01-01T00:00:00.0Z) |> window(start: 1970-01-01T00:00:00.0Z, every: 1m)' http://127.0.0.1:8093/v1/query?orgID=00
```

`window` function step value | Memory consumption | Request duration
-----------------------------|--------------------|------------------
1m | 63336 |