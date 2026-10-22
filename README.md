# mrlm-net/simbrief

[![Go Reference](https://pkg.go.dev/badge/github.com/mrlm-net/simbrief.svg)](https://pkg.go.dev/github.com/mrlm-net/simbrief)
[![Docs](https://img.shields.io/badge/docs-simbrief.mrlm.net-3d7a4f)](https://simbrief.mrlm.net/)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-6c7480)](LICENSE)

A Go client for the SimBrief API: fetch a pilot's latest OFP (JSON v2 or XML), read its navlog, weights, fuel, weather and runway analysis, and build the URL that generates a new flight plan.

|  |  |
|--|--|
| Package name | github.com/mrlm-net/simbrief |
| Go version | 1.27+ |
| License | Apache 2.0 License |
| Platform | Cross-platform |

## Documentation

Full documentation is at **[simbrief.mrlm.net](https://simbrief.mrlm.net/)**: the client, fetching, the OFP and runway analysis, generating plans, helpers and errors. The sources are in [docs/](docs/getting-started.md); the API reference is on [pkg.go.dev](https://pkg.go.dev/github.com/mrlm-net/simbrief).

## Installation

```bash
go get github.com/mrlm-net/simbrief
```

Requirements:

- Go 1.27 or newer
- An internet connection (the SimBrief API)
- A web browser, for generating a new plan (SimBrief signs the user in there)

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/mrlm-net/simbrief/pkg/client"
    "github.com/mrlm-net/simbrief/pkg/types"
)

func main() {
    // Create a new SimBrief client
    simbrief := client.NewClient()

    // Create a basic flight plan request
    request := &types.FlightPlanRequest{
        Origin:      "KJFK",  // New York JFK
        Destination: "KLAX",  // Los Angeles LAX
        Aircraft:    "B38M",  // Boeing 737 MAX 8
    }

    // Generate flight plan URL (actual generation requires browser authentication)
    planURL := simbrief.GenerateFlightPlanURL(request)
    fmt.Printf("Flight plan URL: %s\n", planURL)

    // Get supported aircraft types
    options, err := simbrief.GetSupportedOptions()
    if err != nil {
        log.Fatal("Failed to get aircraft options:", err)
    }

    fmt.Printf("Found %d aircraft types available\n", len(options.Aircraft))
}
```

## Examples

The `examples/` directory contains practical demonstrations of SimBrief features:

| Example | Description | Features |
|---------|-------------|----------|
| [basic](examples/basic/) | Aircraft list, plan URL, latest OFP | Supported options, builder + validation, fetch by user ID, navlog |
| [advanced](examples/advanced/) | Custom aircraft plan URL, full OFP walk-through | Custom aircraft data, ATC, weights, fuel, weather, alternates, navlog, files |

Run any example:

```bash
cd examples/basic
go run main.go
```

## API Reference

### Core Components

#### Client Creation

```go
client := client.NewClient()                      // Create new client
client := client.NewClientWithConfig(...)        // Create with custom config
```

#### Flight Plan Generation

```go
request := &types.FlightPlanRequest{
    Origin:      "KJFK",     // Origin airport ICAO
    Destination: "KLAX",     // Destination airport ICAO
    Aircraft:    "B38M",     // Aircraft type
    Route:       "DCT",      // Flight route
}

// Generate URL for browser-based flight plan creation
planURL := client.GenerateFlightPlanURL(request)
```

#### Data Retrieval

```go
// Get supported aircraft and layouts
options, err := client.GetSupportedOptions()

// Fetch existing flight plan by user ID
flightPlan, err := client.GetFlightPlanByUserID("user_id")

// Get flight plan by username  
flightPlan, err := client.GetFlightPlanByUsername("username")

// Or choose the format: JSON v2 (json=v2) or XML; both decode the same way
flightPlan, err := client.GetFlightPlan(&types.FetchRequest{UserID: "user_id", JSON: true})
```

#### OFP Data

`types.FlightPlanResponse` decodes the OFP from JSON v2 and from XML
(root element `<OFP>`). Numeric values are `types.Number`: the raw text
SimBrief sends ("0365", "05000", ".57") with `Int()`, `Float()` and
`Bool()` helpers. Time values stay strings (JSON v2: ISO 8601 and HH:MM:SS,
XML: Unix seconds and seconds).

```go
fmt.Println(flightPlan.ATC.Callsign)                          // "CEF007"
fmt.Println(flightPlan.ATC.InitialAltUnit, flightPlan.ATC.InitialAlt) // "F" "150"
fmt.Println(flightPlan.General.STARIdent)                     // "BEKV1Q"
fmt.Println(flightPlan.General.InitialAltitude.Int())         // 15000
fmt.Println(flightPlan.Origin.ICAO, flightPlan.Origin.Runway) // "LKPR" "06"
fmt.Println(flightPlan.Origin.TransAlt.Int())                 // 5000
fmt.Println(flightPlan.Aircraft.ICAOCode)                     // "A319"

for _, altn := range flightPlan.Alternates {
    fmt.Println(altn.ICAO, altn.Runway, altn.Burn.Int())
}
for _, fix := range flightPlan.NavLog {
    fmt.Println(fix.Ident, fix.ViaAirway, fix.AltitudeFeet.Int(), fix.IsSIDSTAR.Bool())
}
```

#### Take-off and Landing Speeds (Runway Analysis)

When the plan is generated with runway analysis (`RunwayAnalysis` in the
request, for an aircraft whose `TLRData` is true), `flightPlan.TLR` holds
the take-off and landing report from both JSON v2 and XML. Without it, `TLR`
is the zero value and the helpers return `false`.

```go
tlr := flightPlan.TLR
if rwy, ok := tlr.TakeoffRunway(flightPlan.Origin.Runway); ok { // "06"; "6" works too
    fmt.Println(rwy.SpeedsV1.Int(), rwy.SpeedsVR.Int(), rwy.SpeedsV2.Int()) // 129 129 134
    fmt.Println(rwy.FlapSetting, rwy.ThrustSetting, rwy.FlexTemperature.Int()) // "1" "FLEX" 65
}
if vref, ok := tlr.LandingVref(false); ok { // false = dry runway, true = wet
    fmt.Println(vref) // 125
}
fmt.Println(tlr.Landing.DistanceDry.FactoredDistance.Int()) // 3546 (feet)
```

`TakeoffRunway` tries the exact identifier first, then ignores case and
leading zeros; a number without a side letter finds a lettered runway only
when exactly one matches.

#### Custom Configuration

```go
httpClient := &http.Client{Timeout: 60 * time.Second}
c := client.NewClientWithConfig(client.DefaultBaseURL, httpClient)
c.SetUserAgent("my-app/1.0")
```

## Troubleshooting

The examples read the user to fetch from `SIMBRIEF_USER_ID`:

```bash
export SIMBRIEF_USER_ID=your_id
```

Common issues:

- **Fetch fails with `types.APIError`**: the message is SimBrief's fetch status (unknown user, no plan generated yet)
- **Invalid aircraft code**: list valid codes with `GetAircraftTypes()`
- **Empty `TLR`**: the plan was generated without runway analysis, or the aircraft has no `TLRData`

**API Documentation**: [SimBrief API Docs](https://developers.navigraph.com/docs/simbrief/using-the-api)

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## Contributing

Contributions are welcome. They follow the [Code of Conduct](https://github.com/mrlm-net/.github/blob/main/docs/CODE_OF_CONDUCT.md) and the [contribution guidelines](https://github.com/mrlm-net/.github/blob/main/docs/CONTRIBUTING.md). Report bugs and request features in [issues](https://github.com/mrlm-net/simbrief/issues).

## Security

To report a vulnerability, follow the [security policy](https://github.com/mrlm-net/.github/blob/main/docs/SECURITY.md); please do not open a public issue for it.

## License

Apache License 2.0, see [LICENSE](LICENSE).

Copyright © Martin Hrášek [<@marley-ma>](https://github.com/marley-ma) and WANTED.solutions s.r.o. [<@wanted-solutions>](https://github.com/wanted-solutions)
