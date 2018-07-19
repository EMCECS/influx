# Build

In this fork sources there are the internal imports which are pointing to the upstream source repo.
This may brake the build. It may be necessary to execute the additional steps before the build:

1. Get the upstream sources
    ```bash
    go get influxdata/platform
    ```
2. Go to the sources directory and set the Git origin URL to the fork
    ```bash
    cd $GOPATH/src/github.com/influxdata/platform
    git remote set-url https://github.com/akurilov/influx
    ```

3. Merge the changes from the fork, get the dependencies
    ```bash
    git pull
    dep ensure
    ```
    
4. Build Flux Daemon
    ```bash
    make bin/linux/fluxd
    ```

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

No `window` function step argument (`every`) value dependence was measured in the range of `1m` to `1h`.
The peak memory consumption was about 175MB and the residual memory consumption was 142MB.
  