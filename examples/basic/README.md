# Basic example

Fetches the supported aircraft and layouts, builds a flight plan request with
the builder, and fetches the latest OFP of a user.

## Running

```bash
export SIMBRIEF_USER_ID=123456   # optional; without it the program stops before step 3
cd examples/basic
go run .
```

## What it does

### 1. Supported options

```go
simbrief := client.NewClient()

options, err := simbrief.GetSupportedOptions()
if err != nil {
	log.Fatalf("Failed to get supported options: %v", err)
}
fmt.Printf("Found %d aircraft types and %d layouts\n", len(options.Aircraft), len(options.Layouts))
```

It then prints five aircraft (`id`, `Name`, `Accuracy`) and three layouts
(`id`, `NameLong`). Map order is random, so the entries differ per run.

### 2. Plan request and URL

```go
request := client.NewFlightPlan("KJFK", "KLAX", "B738").
	Route("HUSKY6 PENNS CRP LAS1").
	Altitude("FL380").
	Registration("N123AB").
	Captain("JOHN DOE").
	Passengers(148).
	StaticID("EXAMPLE_FLIGHT_123").
	EnableNavLog().
	Units(types.UnitsLBS).
	Build()

if err := simbrief.ValidateFlightPlanRequest(request); err != nil {
	log.Fatalf("Flight plan validation failed: %v", err)
}
fmt.Printf("Flight plan generation URL: %s\n", simbrief.GenerateFlightPlanURL(request))
```

The URL opens SimBrief's dispatch page with the fields filled in; generating
the plan needs a browser logged in to SimBrief.

### 3. Latest OFP

```go
flightPlan, err := simbrief.GetFlightPlanByUserID(userID)
if err != nil {
	log.Printf("Failed to fetch flight plan: %v", err)
	return
}
fmt.Printf("Flight: %s → %s\n", flightPlan.Origin.ICAO, flightPlan.Destination.ICAO)
fmt.Printf("Distance: %d nm\n", flightPlan.General.RouteDistance.Int())
fmt.Printf("Planned Fuel: %d %s\n", flightPlan.Fuel.PlanRamp.Int(), flightPlan.Params.Units)
```

It also prints the callsign, aircraft, route, enroute time and the first five
navlog fixes (ident, airway, altitude).

## Sample output

The fetch part, rendered from the test OFP in `pkg/types/testdata` (a
CEF007 LKPR-LKPD plan):

```
=== Fetching Flight Plan Data ===
Flight: LKPR → LKPD
Callsign: CEF007
Aircraft: A319 (FENIX A319)
Route: DCT BEKVI BEKV1Q
Distance: 107 nm
Planned Fuel: 3979 kgs
Flight Time: 00:23:03
  BEKVI  DCT      13500 ft
  TOC    BEKV1Q   15000 ft
  GOLIN  BEKV1Q   15000 ft
  KAFIC  BEKV1Q   15000 ft
  PD905  BEKV1Q   15000 ft
```

## See also

- [Advanced example](../advanced/)
- [Usage guide](../../docs/usage.md)
