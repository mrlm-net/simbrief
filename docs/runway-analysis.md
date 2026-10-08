---
title: Runway Analysis
description: Take-off V-speeds, flaps, thrust and flex temperature, and the landing Vref from SimBrief's take-off and landing report (TLR).
order: 4
section: guides
---

`ofp.TLR` holds SimBrief's take-off and landing report. It is filled only when the plan was generated with runway analysis on (`FlightPlanRequest.RunwayAnalysis`) for an aircraft that supports it (`AircraftOption.TLRData`); otherwise it is the zero value and the helpers return `false`. Both JSON v2 and XML fill it. Weights are in `Params.Units`, lengths in feet, speeds in knots, temperatures in degrees Celsius.

```go
tlr := ofp.TLR
if rwy, ok := tlr.TakeoffRunway(ofp.Origin.Runway); ok { // "06"; "6" works too
	fmt.Println(rwy.SpeedsV1.Int(), rwy.SpeedsVR.Int(), rwy.SpeedsV2.Int())   // 129 129 134
	fmt.Println(rwy.FlapSetting, rwy.ThrustSetting, rwy.FlexTemperature.Int()) // "1" "FLEX" 65
}
if vref, ok := tlr.LandingVref(false); ok { // false = dry runway, true = wet
	fmt.Println(vref) // 125
}
fmt.Println(tlr.Landing.DistanceDry.FactoredDistance.Int()) // 3546 (feet)
```

## Finding a runway

`TakeoffRunway(ident)` tries the exact identifier first, then ignores case, a "RW"/"RWY" prefix and leading zeros ("6l" finds "06L"); a number without a side letter finds a lettered runway only when exactly one matches.

`LandingVref(wet)` returns the Vref of the dry or wet landing distance.

## The rest of the report

The rest of the report is plain data:

| Field | Type | Content |
|-------|------|---------|
| `TLR.Takeoff.Conditions`, `TLR.Landing.Conditions` | `TLRConditions` | planned runway and weight, wind, temperature, altimeter, surface |
| `TLR.Takeoff.Runways` | `[]TakeoffRunway` | flaps, thrust, flex temperature, V-speeds, limits and distances |
| `TLR.Landing.DistanceDry`, `DistanceWet` | `TLRLandingDistance` | landing distances and Vref |
| `TLR.Landing.Runways` | `[]LandingRunway` | per-runway landing data |

Both runway types embed `TLRRunway` (lengths, course, wind components, ILS frequency).

## Getting a plan with runway analysis

Turn it on in the request when you [generate the plan](generating-plans.md):

```go
req := client.NewFlightPlan("LKPR", "EGLL", "A320").Build()
runwayAnalysis := true
req.RunwayAnalysis = &runwayAnalysis // TLR in the OFP (aircraft with TLRData)
fmt.Println(c.GenerateFlightPlanURL(req))
```

An empty `TLR` means the plan was generated without runway analysis, or the aircraft has no `TLRData`.
