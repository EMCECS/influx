from(bucket: "test")
    |> range(start: 2018-01-01T00:00:00Z)
    |> dedup()