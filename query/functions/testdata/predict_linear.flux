from(db: "test")
    |> range(start:-5m)
    |> predictLinear(wantedValue: 10.0)
