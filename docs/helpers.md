---
title: Helpers
description: Route, fuel and time helpers in pkg/client for parsing routes and flight levels, converting fuel units, and parsing times and durations.
order: 6
section: reference
---

Small stateless helpers in `pkg/client`.

## RouteHelper

```go
rh := client.NewRouteHelper()
fmt.Println(rh.ParseRoute("VOZ7B  VOZ  DCT")) // [VOZ7B VOZ DCT]
fmt.Println(rh.ValidateICAOCode("LKPR"))      // true
fmt.Println(rh.FormatFlightLevel(36000))      // FL360
feet, err := rh.ParseFlightLevel("FL360")     // 36000 (plain "5000" works too)
if err != nil {
	return err
}
fmt.Println(feet)
```

`ValidateICAOCode` accepts upper-case letters and digits only.

## FuelHelper

```go
fh := client.NewFuelHelper()
fmt.Println(fh.ConvertLBSToKGS(1000), fh.ConvertKGSToLBS(453.592)) // 453.592 1000
pct, minutes, err := fh.ParseFuelValue("0.05/15")                  // 0.05 "15"
if err != nil {
	return err
}
fmt.Println(pct, minutes)
```

## TimeHelper

```go
th := client.NewTimeHelper()
h, m, err := th.ParseTimeString("14:30") // 14 30
if err != nil {
	return err
}
fmt.Println(h, m, th.FormatTimeString(h, m)) // 14 30 14:30
mins, err := th.ParseDuration("01:27")       // 87
if err != nil {
	return err
}
fmt.Println(mins, th.FormatDuration(mins)) // 87 01:27
```
