# Build

There are the internal imports which are pointing to the upstream source repo in this fork sources.
This may brake the build. It may be necessary to execute the additional steps before the build:

1. Get the upstream sources
    ```bash
    go get github.com/influxdata/platform
    ```
2. Go to the sources directory and set the Git origin URL to the fork
    ```bash
    cd $GOPATH/src/github.com/influxdata/platform
    git remote remove origin
    git remote add origin https://github.com/EMCECS/influx
    ```

3. Merge the changes from the fork, get the dependencies
    ```bash
    git pull
    dep init
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

# Changes

## Connection recovery

https://github.com/influxdata/platform/issues/171

## Deduplication function

https://github.com/influxdata/platform/issues/179

The dedup function behaves like previously existing "unique" function but compares the time series using all the columns
available. Example demonstrating the difference for the unique/dedup functions on the same data:

* unique:
```bash
curl -XPOST --data-urlencode 'q=from(db: "test_dedup_1") |> range(start: -100h) |> unique()' http://127.0.0.1:8093/v1/query?orgID=00
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string
#partition,false,false,true,true,false,false,true,true
#default,_result,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement
,,0,2018-07-01T04:02:57.742527076Z,2018-07-05T08:02:57.742527076Z,2018-07-05T06:41:52.015548952Z,1,value,x
```

* dedup:
```bash
curl -XPOST --data-urlencode 'q=from(db: "test_dedup_1") |> range(start: -100h) |> dedup()' http://127.0.0.1:8093/v1/query?orgID=00      
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string
#partition,false,false,true,true,false,false,true,true
#default,_result,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement
,,0,2018-07-01T04:03:06.486571927Z,2018-07-05T08:03:06.486571927Z,2018-07-05T06:41:52.015548952Z,1,value,x
,,0,2018-07-01T04:03:06.486571927Z,2018-07-05T08:03:06.486571927Z,2018-07-05T06:41:55.969062966Z,1,value,x
```