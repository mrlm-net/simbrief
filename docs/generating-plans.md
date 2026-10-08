---
title: Generating a Plan
description: Build SimBrief dispatch URLs with FlightPlanBuilder, validate requests, and list aircraft types and plan formats.
order: 5
section: guides
---

The library does not create plans itself. It builds the URL of SimBrief's dispatch page with all parameters filled in; open it in a browser that is logged in to SimBrief. Give the plan a `StaticID` to fetch exactly that plan afterwards with `GetFlightPlanByStaticID`.

```go
req := client.NewFlightPlan("LKPR", "EGLL", "A320").
	Route("VOZ7B VOZ DCT ...").
	Airline("CSA").
	FlightNumber("123").
	AltitudeFromFlightLevel(360).
	DepartureTime(14, 30).
	Passengers(150).
	Units(types.UnitsKGS).
	StaticID("MYAPP_42").
	EnableNavLog().
	Build()

runwayAnalysis := true
req.RunwayAnalysis = &runwayAnalysis // TLR in the OFP (aircraft with TLRData)

if err := c.ValidateFlightPlanRequest(req); err != nil {
	return err
}
fmt.Println(c.GenerateFlightPlanURL(req)) // open in a browser logged in to SimBrief
```

## The builder

`FlightPlanBuilder` (from `client.NewFlightPlan(origin, destination, aircraft)`) also has `Alternate`, `Date`/`DateFromTime`, `Altitude`, `AltitudeFromFeet`, `Cargo`, `PlanFormat`, `Registration`, `CallSign`, `Captain`, `Dispatcher`, `DisableNavLog`, `EnableETOPS`, `EnableStepClimbs`, `CustomAircraftData` (`*types.AircraftData`), `TaxiTimes` and `Runways`.

Every other SimBrief input (fuel, profiles, alternates 1-4, NOTAMs, maps ...) is a field on `types.FlightPlanRequest`; set it on the built request or build the struct directly:

```go
req := types.NewFlightPlanRequest("KJFK", "KLAX", "B38M")
req.Route = "DCT"
req.CostIndex = "25"
fmt.Println(c.GenerateFlightPlanURL(req))
```

## Validation

`ValidateFlightPlanRequest` checks that origin, destination and aircraft are set, that both airports are 4-letter codes and that a departure time is in range. `req.Validate()` only checks the three required fields and returns the sentinel errors listed in [Errors](errors.md).

## Editing a plan

To send the user to edit a plan they already generated with a static ID:

```go
fmt.Println(c.GetDirectEditURL("MYAPP_42"))
// https://www.simbrief.com/system/dispatch.php?editflight=last&static_id=MYAPP_42
```

## Aircraft types and plan formats

```go
aircraft, err := c.GetAircraftTypes() // map[string]types.AircraftOption
if err != nil {
	return err
}
for id, a := range aircraft {
	fmt.Println(id, a.Name, a.Accuracy, a.TLRData)
}

layouts, err := c.GetPlanFormats() // map[string]types.LayoutOption
if err != nil {
	return err
}
for id, l := range layouts {
	fmt.Println(id, l.NameShort, l.NameLong)
}
```

Both read `/api/inputs.list.json`; `GetSupportedOptions()` returns the whole `types.SupportedOptions` (aircraft, layouts, last update) in one request.
