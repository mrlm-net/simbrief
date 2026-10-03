# Flight Planning

This document covers advanced flight planning capabilities of the SimBrief SDK.

## Table of Contents

- [Basic Flight Plan Creation](#basic-flight-plan-creation)
- [Advanced Routing](#advanced-routing)
- [Fuel Planning](#fuel-planning)
- [Weather Integration](#weather-integration)
- [Alternate Airports](#alternate-airports)
- [Custom Aircraft Configuration](#custom-aircraft-configuration)

## Basic Flight Plan Creation

### Simple Flight Plan

```go
request := &types.FlightPlanRequest{
    Origin:      "KJFK",
    Destination: "KLAX", 
    Aircraft:    "B38M",
}

flightPlan, err := client.GenerateFlightPlan(request)
```

### With Basic Parameters

```go
request := &types.FlightPlanRequest{
    Origin:          "KJFK",
    Destination:     "KLAX",
    Aircraft:        "B38M",
    Airline:         "ABC",
    FlightNumber:    "1234",
    DepartureHour:   14,
    DepartureMinute: 30,
    Passengers:      150,
}
```

## Advanced Routing

### Custom Route

```go
request := &types.FlightPlanRequest{
    Origin:      "KJFK",
    Destination: "KLAX",
    Aircraft:    "B38M",
    Route:       "HAPIE6 HAPIE J80 WILMINGTON J82 LYNCH J134 TUS EAGUL4",
    Altitude:    "FL370",
}
```

### Route Planning Options

```go
request := &types.FlightPlanRequest{
    // ... basic parameters
    RouteGenerator: "PFPX",     // Route planning system
    NavSystem:     "GPS",       // Navigation system
    SIDType:       "RNAV",      // SID type preference
    STARType:      "RNAV",      // STAR type preference
}
```

## Fuel Planning

### Custom Fuel Configuration

```go
request := &types.FlightPlanRequest{
    // ... basic parameters
    
    // Fuel planning
    ExtraFuelPercent: 5,        // Extra fuel percentage
    ExtraFuelTime:    30,       // Extra fuel in minutes
    TaxiFuel:         500,      // Taxi fuel in pounds
    ContingencyFuel:  5,        // Contingency fuel percentage
    
    // Payload
    Passengers:       180,      // Number of passengers
    Baggage:          2000,     // Baggage weight
    Cargo:            5000,     // Cargo weight
}
```

### Reserve Fuel Options

```go
request := &types.FlightPlanRequest{
    // ... basic parameters
    ReserveFuel:      3000,     // Reserve fuel in pounds
    MinimumFuel:     2000,      // Minimum fuel requirement
    AlternateFuel:   1500,      // Fuel for alternate airport
}
```

## Weather Integration

### Weather Planning

```go
request := &types.FlightPlanRequest{
    // ... basic parameters
    
    // Weather settings
    WeatherSource:   "NOAA",    // Weather data source
    TurbulenceData: true,       // Include turbulence data
    IcingData:      true,       // Include icing information
    
    // Wind planning
    WindData:       "HISTORICAL", // Use historical wind data
    WindOptimize:   true,         // Optimize for winds
}
```

## Alternate Airports

### Single Alternate

```go
request := &types.FlightPlanRequest{
    // ... basic parameters
    Alternate: "KLAS",  // Las Vegas as alternate
}
```

### Multiple Alternates

```go
request := &types.FlightPlanRequest{
    // ... basic parameters
    
    AltnCount: 2,
    Altn1ID:   "KLAS",
    Altn2ID:   "KPHX",
    
    // Avoid specific airports as alternates
    AltnAvoid: "KJFK KLGA KEWR",
}
```

### Alternate-Specific Routing

```go
request := &types.FlightPlanRequest{
    // ... basic parameters
    
    Altn1ID:     "KLAS",
    Altn1Route:  "DCT",           // Direct routing to alternate
    Altn1Runway: "08L",           // Preferred runway
    
    Altn2ID:     "KPHX",
    Altn2Route:  "BAYST1",        // Specific arrival
    Altn2Runway: "08",
}
```

## Custom Aircraft Configuration

### Basic Aircraft Data

```go
aircraftData := &types.AircraftData{
    ICAO:        "B39M",
    Name:        "737 MAX 9",
    Engines:     "CFM LEAP-1B",
    Category:    "M",
    WakeClass:   "M",
}

request := &types.FlightPlanRequest{
    // ... basic parameters
    AircraftData: aircraftData,
}
```

### Performance Parameters

```go
aircraftData := &types.AircraftData{
    ICAO:            "B39M",
    Name:            "737 MAX 9",
    
    // Performance data
    MaxAltitude:     41000,      // Service ceiling
    CruiseSpeed:     0.785,      // Mach number
    MaxSpeed:        0.82,       // Maximum speed
    
    // Weight limits
    MaxTakeoffWeight: 194700,    // MTOW in pounds
    MaxLandingWeight: 155500,    // MLW in pounds
    MaxZeroFuelWeight: 138500,   // MZFW in pounds
    
    // Fuel capacity
    MaxFuelCapacity: 25820,      // Maximum fuel in pounds
    
    // Operating costs
    FuelBurnRate:    5500,       // Fuel burn per hour
    OperatingCost:   12000,      // Operating cost per hour
}
```

### Engine-Specific Configuration

```go
aircraftData := &types.AircraftData{
    // ... basic data
    
    Engines:         "CFM LEAP-1B",
    EngineCount:     2,
    ThrustPerEngine: 28000,      // Thrust in pounds
    
    // Fuel consumption by phase
    TaxiFuelBurn:    200,        // Fuel burn during taxi
    ClimbFuelBurn:   8000,       // Fuel burn during climb
    CruiseFuelBurn:  5500,       // Fuel burn during cruise
    DescentFuelBurn: 3000,       // Fuel burn during descent
}
```

## Response Data Structure

### Flight Plan Response

`types.FlightPlanResponse` (pkg/types/response.go) decodes the OFP from
`xml.fetcher.php` both as JSON v2 (`json=v2`) and as XML (root `<OFP>`).

```go
type FlightPlanResponse struct {
    Fetch            FetchInfo       // fetch status ("Success" or an error)
    Params           FlightParams    // request id, AIRAC, units ("kgs"/"lbs")
    General          GeneralInfo     // route, route_ifps, SID/STAR, initial_altitude, costindex, passengers
    Origin           AirportInfo     // icao_code, plan_rwy, elevation, trans_alt, trans_level, metar
    Destination      AirportInfo
    Alternates       []AlternateInfo // AirportInfo + diversion leg (distance, burn, route)
    NavLog           NavLog          // []NavLogFix
    AlternateNavLogs []NavLog
    ATC              ATCInfo         // callsign, flight_rules, flight_type, initial_alt(+unit), route, FPL text
    Aircraft         AircraftInfo    // icaocode, icao_code, reg, name, engines
    Fuel             FuelInfo
    Times            TimeInfo
    Weights          WeightInfo
    Weather          WeatherInfo
    Files            FilesInfo
    Links            LinksInfo
}
```

Conventions:

- Text values (idents, routes, runways like `"06"`) are `string`.
- Numbers are `types.Number`, the raw SimBrief text (`"0365"`, `"05000"`)
  with `Int()`, `Float()`, `Bool()` and `IsSet()`. It decodes JSON strings,
  JSON numbers and XML text.
- Times stay `string`: JSON v2 sends ISO 8601 and HH:MM:SS, XML sends Unix
  seconds and seconds.

## Error Handling

### Common Flight Planning Errors

```go
flightPlan, err := client.GenerateFlightPlan(request)
if err != nil {
    switch {
    case errors.Is(err, types.ErrInvalidAircraft):
        // Handle invalid aircraft type
        fmt.Printf("Invalid aircraft: %s\n", request.Aircraft)
        
    case errors.Is(err, types.ErrInvalidRoute):
        // Handle routing errors
        fmt.Printf("Invalid route: %s\n", request.Route)
        
    case errors.Is(err, types.ErrAPIKeyRequired):
        // Handle authentication errors
        fmt.Println("API key required for flight plan generation")
        
    default:
        // Handle other errors
        fmt.Printf("Flight planning error: %v\n", err)
    }
}
```

### Validation

```go
// Validate request before sending
if err := request.Validate(); err != nil {
    fmt.Printf("Request validation failed: %v\n", err)
    return
}

// Generate flight plan
flightPlan, err := client.GenerateFlightPlan(request)
```

## Best Practices

1. **Always validate input**: Use the validation methods before making API calls
2. **Handle rate limits**: Implement backoff strategies for API calls
3. **Cache aircraft data**: Store supported aircraft types locally to reduce API calls
4. **Use appropriate timeouts**: Set reasonable timeouts for long-running operations
5. **Error handling**: Always handle specific error types appropriately
6. **Resource cleanup**: Ensure proper cleanup of HTTP resources
