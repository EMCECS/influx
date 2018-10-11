from(bucket: "test")
    |> range(start: 2018-08-10T09:30:00.00Z)
    |> predictLinear(wantedValue: 10.0)
