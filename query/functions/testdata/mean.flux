from(bucket: "test")
    |> range(start:2018-05-22T19:53:26Z)
    |> group(by:["_measurement", "_start"])
    |> mean(timeSrc:"_start")
    |> map(fn: (r) => {_time: r._time, mean: r._value})
    |> yield(name: "0")