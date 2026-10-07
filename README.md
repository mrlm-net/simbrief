# mrlm-net/simbrief

GoLang wrapper for the SimBrief API, allowing Gophers to create flight planning tools, calculate flight plans, and integrate with SimBrief's comprehensive flight planning services.

|  |  |
|--|--|
| Package name | github.com/mrlm-net/simbrief |
| Go version | 1.21+ |
| License | Apache 2.0 License |
| Platform | Cross-platform |

## Table of contents

• [Installation](#installation) • [Usage](#usage) • [Examples](#examples) • [API Reference](#api-reference) • [Debugging](#debugging) • [Advanced Usage](#advanced-usage) • [Contributing](#contributing)

## Installation

I'm using `go mod` so examples will be using it, you can install this package via Go modules.

```bash
go get github.com/mrlm-net/simbrief
```

**Requirements:**

• Go 1.21 or higher
• Internet connection (for SimBrief API access)
• Web browser (for flight plan generation authentication)

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
| [basic](examples/basic/) | Basic flight plan operations and data retrieval | Get supported aircraft, Fetch existing plans, Basic API usage |
| [advanced](examples/advanced/) | Complete flight planning with custom parameters | Custom aircraft data, Advanced routing, Detailed fuel planning |

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
client := client.NewClientWithConfig("https://api.simbrief.com", httpClient)
```

## Debugging

Enable detailed HTTP logging and error information:

**Environment Variables:**
```bash
export SIMBRIEF_DEBUG=true          # Enable debug logging
export SIMBRIEF_USER_ID=your_id     # Set user ID for examples
```

**Common issues:**

• **Network connectivity**: Ensure stable internet connection for API access
• **Invalid aircraft code**: Verify aircraft ICAO codes using `GetSupportedOptions()`
• **Route validation**: Check route format against SimBrief documentation
• **Rate limiting**: SimBrief may rate limit requests; implement appropriate delays

**API Documentation**: [SimBrief API Docs](https://developers.navigraph.com/docs/simbrief/using-the-api)

## Advanced Usage

For complex scenarios, see [Advanced Documentation](https://github.com/mrlm-net/simbrief/blob/main/docs):

• [Flight Planning](https://github.com/mrlm-net/simbrief/blob/main/docs/flight-planning.md) - Advanced flight plan customization and routing
• [Aircraft Management](https://github.com/mrlm-net/simbrief/blob/main/docs/aircraft-management.md) - Custom aircraft configurations and performance data
• [API Integration](https://github.com/mrlm-net/simbrief/blob/main/docs/api-integration.md) - Best practices for API usage and error handling
• [Performance](https://github.com/mrlm-net/simbrief/blob/main/docs/performance.md) - Optimization and caching strategies

## Contributing

Contributions are welcomed and must follow Code of Conduct and common [Contributions guidelines](https://github.com/mrlm-net/.github/blob/main/docs/CONTRIBUTING.md).

> **Security**: If you'd like to report security issue please follow security guidelines.

All rights reserved © Martin Hrášek [<@marley-ma>](https://github.com/marley-ma) and WANTED.solutions s.r.o. [<@wanted-solutions>](https://github.com/wanted-solutions)
