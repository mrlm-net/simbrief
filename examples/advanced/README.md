# Advanced example

Builds a plan URL with custom aircraft data, then fetches the latest OFP of a
user and prints most of it.

## Running

```bash
export SIMBRIEF_USER_ID=123456   # required; the program exits at the start without it
cd examples/advanced
go run .
```

## What it does

### 1. Plan URL with custom aircraft data

```go
customAircraft := &types.AircraftData{
	ICAO:        "B39M",
	Name:        "737 MAX 9",
	Engines:     "CFM LEAP-1B",
	Category:    "M",
	Equipment:   "SDE3FGHIRWY",
	Transponder: "S",
	PBN:         "PBN/A1B1C1D1",
	MaxPax:      "220",
	OEW:         99.5, // thousands of lbs
	MZFW:        138.8,
	MTOW:        194.7,
	MLW:         155.0,
	MaxFuel:     46.0,
	// ... ExtraRemark, HexCode, Per, PaxWgt
}

request := client.NewFlightPlan("KJFK", "EGLL", "B39M").
	Route("HAPIE6 HAPIE N247A ALLRY DCT KANNI N866B BEXET UL9 BOGNA UL607 REDFA").
	Altitude("FL380").
	Airline("UAL").
	FlightNumber("918").
	// ... Registration, Captain, Dispatcher, Passengers, Cargo
	Alternate("EGKK").
	CustomAircraftData(customAircraft).
	EnableNavLog().
	Units(types.UnitsLBS).
	StaticID("ADVANCED_EXAMPLE").
	Build()

fmt.Printf("Advanced flight plan URL: %s\n", simbrief.GenerateFlightPlanURL(request))
```

The custom aircraft data goes into the URL as the `acdata` JSON parameter.

### 2. OFP walk-through

`GetFlightPlanByUserID(userID)` fetches the latest plan (JSON v2); the program
then prints, section by section:

| Section | Fields |
|---------|--------|
| Airports | `Origin`/`Destination`: `Name`, `ICAO`, `Runway`, `Elevation` |
| ATC | `ATC.Callsign`, `FlightRules`, `FlightType`, `InitialAltUnit` + `InitialAlt`, `Route` |
| Aircraft | `Aircraft.Name`, `ICAOCode`, `Registration`, `Engines` |
| Flight planning | `General.RouteDistance`, `Route`, `SIDIdent`, `STARIdent`, `InitialAltitude`, `CostIndex`, `AvgWindDir`/`AvgWindSpd` |
| Timing | `Times.EstTimeEnroute`, `EstBlock`, `TaxiOut`, `TaxiIn` (JSON v2: `HH:MM:SS`) |
| Weights | `Weights.OEW`, `EstZFW`, `EstTOW`, `EstLDW`, `Payload`, `PaxCount`, `PaxWeight`, in `Params.Units` |
| Fuel | `Fuel.PlanRamp`, `EnrouteBurn`, `Taxi`, `AlternateBurn`, `Contingency`, `Reserve`, `Extra`, `AvgFuelFlow` |
| Weather | `Weather.OrigMETAR`, `DestMETAR` |
| Alternates | each of `Alternates`: `Name`, `ICAO`, `Runway`, `Distance`, `TrackMag`, `Burn` |
| Navlog | every fix: `Ident`, `Type`, `ViaAirway`, `AltitudeFeet`, `Stage` |
| Files | `Files.Directory` + `Files.PDF.Link`, number of other `Files.Files` |

Numeric fields are `types.Number`: printed with `%s` they show SimBrief's raw
text, `.Int()` gives the number.

## See also

- [Basic example](../basic/)
- [Usage guide](../../docs/usage.md) - including the runway analysis (TLR), which this example does not print
