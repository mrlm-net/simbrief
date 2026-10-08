---
title: The OFP
description: What types.FlightPlanResponse holds (ATC, aircraft, fuel, weights, times, weather, alternates, navlog) and how numbers and times are decoded.
order: 3
section: guides
---

`types.FlightPlanResponse` decodes the OFP from JSON v2 and from XML (root element `<OFP>`) and mirrors its blocks:

| Field | Type | Content |
|-------|------|---------|
| `Fetch` | `FetchInfo` | fetch outcome: `Status` ("Success" or an error), user ID, static ID |
| `Params` | `FlightParams` | request/sequence IDs, `TimeGen`, `AIRAC`, `OFPLayout`, `Units` ("kgs"/"lbs") |
| `General` | `GeneralInfo` | flight number, route, SID/STAR, cost index, initial altitude, distances, burn, average wind |
| `Origin`, `Destination` | `AirportInfo` | ICAO/IATA, position, elevation, planned `Runway`, `TransAlt`/`TransLevel`, METAR/TAF |
| `Alternates` | `[]AlternateInfo` | each alternate airport plus its diversion leg (distance, track, burn, route) |
| `NavLog` | `NavLog` | the route fixes |
| `AlternateNavLogs` | `[]NavLog` | one navlog per alternate, in `Alternates` order |
| `ATC` | `ATCInfo` | filed ICAO flight plan: callsign, rules, type, initial speed/level and units, route, FPL text |
| `Aircraft` | `AircraftInfo` | ICAO/IATA type, name, engines, registration, SELCAL, equipment |
| `Fuel` | `FuelInfo` | taxi, trip (`EnrouteBurn`), contingency, alternate, reserve, extra, `PlanRamp` (block) ... in `Params.Units` |
| `Weights` | `WeightInfo` | OEW, pax/bags, payload, ZFW/TOW/LDW estimates and limits, in `Params.Units` |
| `Times` | `TimeInfo` | scheduled/estimated out/off/on/in, block and enroute times |
| `Weather` | `WeatherInfo` | origin/destination METAR and TAF, one per alternate |
| `TLR` | `TLR` | runway analysis (take-off and landing); zero when the plan has none, see [Runway analysis](runway-analysis.md) |
| `Files`, `Links` | `FilesInfo`, `LinksInfo` | PDF and exported plan files (relative to `Files.Directory`), SkyVector link |

```go
fmt.Println(ofp.ATC.Callsign)                                 // "CEF007"
fmt.Println(ofp.ATC.InitialAltUnit, ofp.ATC.InitialAlt)       // "F" "150"
fmt.Println(ofp.General.Route, ofp.General.STARIdent)         // route, "BEKV1Q"
fmt.Println(ofp.General.InitialAltitude.Int())                // 15000
fmt.Println(ofp.Origin.ICAO, ofp.Origin.Runway)               // "LKPR" "06"
fmt.Println(ofp.Origin.TransAlt.Int())                        // 5000
fmt.Println(ofp.Aircraft.ICAOCode, ofp.Aircraft.Registration) // "A319" ...
fmt.Println(ofp.Fuel.PlanRamp.Int(), ofp.Params.Units)        // block fuel, "kgs"
fmt.Println(ofp.Weights.EstTOW.Int(), ofp.Times.EstBlock)     // weight, duration text
```

## Number

SimBrief sends numbers as text, often zero-padded ("0365", "05000"), without a leading zero (".57") or empty. Numeric fields are `types.Number`: the raw text, decoded from a JSON string, JSON number, boolean, null or XML text.

```go
n := types.Number("0365")
fmt.Println(n.Int(), n.Float(), n.IsSet(), n.String()) // 365 365 true "0365"
fmt.Println(types.Number(".57").Float())               // 0.57
fmt.Println(types.Number("").Int())                    // 0
fmt.Println(types.Number("1").Bool())                  // true
```

`Int` truncates fractions; `Int`, `Float` and `Bool` return zero/false for an empty or non-numeric value, so use `IsSet` when "missing" and "0" differ.

## NavLog

`NavLog` is a `[]NavLogFix`. JSON v2 sends a plain array, XML sends `<navlog><fix>...</fix></navlog>`; both decode the same way. A fix carries its ident, name, `Type` (wpt, vor, ndb, apt, ltlg for TOC/TOD ...), `Stage` (CLB, CRZ, DSC), `ViaAirway`, position, altitude, speeds, wind, leg and total time and fuel, FIR and MORA.

```go
for _, fix := range ofp.NavLog {
	if fix.IsSIDSTAR.Bool() {
		continue
	}
	fmt.Println(fix.Ident, fix.Type, fix.ViaAirway, fix.Stage,
		fix.AltitudeFeet.Int(), fix.Latitude.Float(), fix.Longitude.Float())
}
```

## Alternates

`Alternate()` returns the first alternate or nil. `AlternateNavLogs[i]` belongs to `Alternates[i]`.

```go
if altn := ofp.Alternate(); altn != nil {
	fmt.Println(altn.ICAO, altn.Runway, altn.Distance.Int(), altn.Burn.Int())
}
for i, altn := range ofp.Alternates {
	if i < len(ofp.AlternateNavLogs) {
		fmt.Println(altn.ICAO, len(ofp.AlternateNavLogs[i]), "fixes")
	}
}
```

## Times

Time values are strings because the formats differ: JSON v2 sends ISO 8601 timestamps and `HH:MM:SS` durations, XML sends Unix seconds and durations in seconds. This applies to `Times`, `Params.TimeGen`, `AlternateInfo.ETE` and the navlog's `TimeLeg`/`TimeTotal`.
