# tides: go library & cli

A tide prediction calculator that uses harmonic constituent data and astronomy math to provide tide predictions in go.

Please do not use this package if it's mission-critical for the tides to be accurate, such as in marine navigation. For all other purposes, it seems to do a pretty good job.

## Setup

Place the harmonic data for your desired station(s) in the data directory (as specified with the `dataDir/data-dir` parameter). See [Required Station Data](#required-station-data) for details.

## CLI
```bash
# install the binary
go install github.com/ryan-lang/tides/cmd@latest

# download the station metadata from NOAA
tides download noaaStation --station 9445719

# run a prediction
tides predict --station 9445719
```

## Library
```go
// Load harmonics from file
har, err := tides.LoadHarmonicsFromFile("./data", "9447130")
if err != nil {
    panic(err)
}

// Create a new prediction for a date range
start := time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)
end := start.Add(time.Hour * 1)
prediction := har.NewRangePrediction(start, end, tides.WithInterval(time.Minute*10))

// Get the prediction results
results := prediction.Predict()
for _, result := range results {
    fmt.Printf("%f @ %s\n", result.Level, result.Time)
}
```

## Required Station Data
Tides are calculated using harmonic constituent data, which can be found in several places online, or you can calculate your own through tide observations (which is outside the scope of this package).

See the [neaps tide database](https://github.com/neaps/tide-database) for a good repository of constituent data, or (for US stations only), use the CLI to download data from NOAA as shown below.

All values should be provided in meters.

#### Reference Stations vs Subordinate Stations

There are relatively few tide stations which actually use their own harmonic data, and these are called *reference stations*. All other stations are *subordinate stations* meaning they are pegged to a nearby reference station, and simply apply offsets to account for local differences.

This package supports both types of stations, but if you want to do calculations for a subordinate station, you need to provide the reference station data too. If downloading from NOAA, the CLI handles this for you.

#### Datum conversion

Results are relative to the MTL (mean tide level) datum. If a datum conversion is requested, then the datum metadata must be provided in the station json.

### Data structure
```json
// ./data/9447130.json (reference station Seattle, WA)
{
    "harmonic_constituents": [
        {
            "name": "M2",
            "phase_UTC": 10.6,
            "phase_local": 138.7,
            "amplitude": 1.072,
            "speed": 28.984104
        },
        ...
    ],
    "datums":[
        {
            "name": "MHHW",
            "value": 5.882
        },
        ...
    ]  
}
```
```json
// ./data/9445719.json (subordinate station Poulsbo, WA)
{
   "tide_pred_offsets": {
        "ref_station_id": "9447130",
        "height_offset_high_tide": 1.03,
        "height_offset_low_tide": 1.01,
        "time_offset_high_tide": 5,
        "time_offset_low_tide": 12
    },
    "datums":[
        {
            "name": "MHHW",
            "value": 5.3
        },
        ...
    ]  
}
```

## "Why don't my predictions perfectly match NOAA?"
It's a good question, and one that I don't know the answer to. The tests run a comparsion against NOAA values, and I can get them to pass only with a tolerance of +/- 0.15m, and +/- 10 min, which seems not great on the surface. 

However:

1) I am not aware of any other more accurate tide software. Most seem to get somewhere in the same ballpark.
2) Even if this package can't match NOAA predictions, we seem to be within the error bars when comparing against actual water levels, which is the whole point anyway.
3) It's very possible that there's nothing wrong with the math or the methodology, and NOAA just has some sort of secret sauce they use.
4) Open source makes things better... maybe you can help?

## Kudos
To the inimitable [Xtide](https://flaterco.com/xtide/), for helping to break down the difficult math, and to [Pytides](https://github.com/sam-cox/pytides) and [tide-predictor](https://github.com/neaps/tide-predictor) for implementation inspiration.
