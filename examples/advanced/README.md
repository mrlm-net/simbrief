# Advanced SimBrief Example

This example demonstrates advanced features of the SimBrief SDK including custom aircraft configurations, comprehensive flight plan analysis, and detailed data extraction.

## Features Demonstrated

- **Custom Aircraft Data**: Create flight plans with custom aircraft specifications
- **Comprehensive Flight Planning**: Build detailed flight plans with all parameters
- **Advanced Data Analysis**: Extract and analyze detailed flight plan information
- **Weather Integration**: Display weather data and analysis
- **Weight & Balance**: Detailed weight and balance calculations
- **Fuel Planning**: Comprehensive fuel planning with all components
- **Navigation Analysis**: Route and navigation data extraction
- **File Generation**: Access to generated PDF, XML, and flight simulator files

## Prerequisites

- Go 1.21 or higher
- Internet connection for SimBrief API access
- **Required**: SimBrief user ID (for fetching existing flight plans)

## Environment Variables

```bash
# Required: Set your SimBrief user ID
export SIMBRIEF_USER_ID=857341

# Optional: Enable debug logging
export SIMBRIEF_DEBUG=true
```

> **Note**: This example requires a valid SimBrief user ID to demonstrate flight plan analysis. You can find your user ID in your SimBrief account settings.

## Running the Example

```bash
cd examples/advanced
go run main.go
```

## Code Overview

### 1. Custom Aircraft Configuration

```go
// Define custom aircraft with detailed specifications
customAircraft := &types.AircraftData{
    ICAO:        "B39M", // Custom aircraft code
    Name:        "737 MAX 9",
    Engines:     "CFM LEAP-1B",
    Category:    "M", // Medium category
    Equipment:   "SDE3FGHIRWY",
    Transponder: "S",
    PBN:         "PBN/A1B1C1D1",
    ExtraRemark: "RMK/CUSTOM AIRCRAFT CONFIG",
    MaxPax:      "220",
    OEW:         99.5,  // Operating Empty Weight (thousands of lbs)
    MZFW:        138.8, // Max Zero Fuel Weight
    MTOW:        194.7, // Max Takeoff Weight
    MLW:         155.0, // Max Landing Weight
    MaxFuel:     46.0,  // Max Fuel Capacity
    HexCode:     "A1B2C3",
    Per:         "D",
    PaxWgt:      190, // Average passenger weight (lbs)
}
```

### 2. Comprehensive Flight Plan Builder

```go
// Build a detailed flight plan with all parameters
request := client.NewFlightPlan("KJFK", "EGLL", "B39M").
    Route("HAPIE6 HAPIE N247A ALLRY DCT KANNI N866B BEXET UL9 BOGNA UL607 REDFA").
    Altitude("FL380").
    Airline("UAL").
    FlightNumber("918").
    Registration("N39MAX").
    Captain("JANE SMITH").
    Dispatcher("JOHN DOE").
    Passengers(180).
    Cargo(12.5).
    Alternate("EGKK").
    CustomAircraftData(customAircraft).
    EnableNavLog().
    Units(types.UnitsLBS).
    StaticID("ADVANCED_EXAMPLE").
    Build()
```

### 3. Detailed Flight Plan Analysis

```go
// Fetch and analyze existing flight plan
flightPlan, err := simbrief.GetFlightPlanByUserID(userID)
if err != nil {
    log.Fatalf("Failed to fetch flight plan: %v", err)
}

// Analyze various aspects of the flight plan
fmt.Printf("Route: %s → %s\n", flightPlan.Origin.ICAO, flightPlan.Destination.ICAO)
fmt.Printf("Aircraft: %s (%s)\n", flightPlan.Aircraft.Name, flightPlan.Aircraft.ICAO)
fmt.Printf("Distance: %s nm\n", flightPlan.General.Distance)
fmt.Printf("Flight Time: %s\n", flightPlan.Times.FlightTime)
```

### 4. Weight & Balance Analysis

```go
// Display comprehensive weight and balance information
fmt.Printf("Weight & Balance (%s):\n", flightPlan.General.Units)
fmt.Printf("Operating Empty Weight: %s\n", flightPlan.Weights.OEW)
fmt.Printf("Zero Fuel Weight: %s\n", flightPlan.Weights.ZFW)
fmt.Printf("Takeoff Weight: %s\n", flightPlan.Weights.TakeoffWt)
fmt.Printf("Landing Weight: %s\n", flightPlan.Weights.LandingWt)
fmt.Printf("Payload: %s (Pax: %s @ %s)\n",
    flightPlan.Weights.Payload, flightPlan.Weights.PaxCount, flightPlan.Weights.PaxWeight)
```

### 5. Fuel Planning Analysis

```go
// Detailed fuel planning breakdown
fmt.Printf("Fuel Planning (%s):\n", flightPlan.General.Units)
fmt.Printf("Total Planned: %s\n", flightPlan.Fuel.Plan)
fmt.Printf("Trip Fuel: %s\n", flightPlan.Fuel.Trip)
fmt.Printf("Taxi Fuel: %s\n", flightPlan.Fuel.Taxi)
fmt.Printf("Alternate Fuel: %s\n", flightPlan.Fuel.Alternate)
fmt.Printf("Contingency: %s\n", flightPlan.Fuel.Contingency)
fmt.Printf("Reserve: %s\n", flightPlan.Fuel.Reserve)
fmt.Printf("Extra: %s\n", flightPlan.Fuel.Extra)
```

### 6. Weather Integration

```go
// Weather analysis and display
if flightPlan.Weather.OriginMETAR != "" {
    fmt.Printf("Origin METAR: %s\n", flightPlan.Weather.OriginMETAR)
}
if flightPlan.Weather.DestMETAR != "" {
    fmt.Printf("Destination METAR: %s\n", flightPlan.Weather.DestMETAR)
}
fmt.Printf("Average Wind: %s°/%s kts\n",
    flightPlan.Weather.AvgWindDir, flightPlan.Weather.AvgWindSpd)
```

## Sample Output

```
=== Creating Flight Plan with Custom Aircraft Data ===
Advanced flight plan URL: https://www.simbrief.com/system/dispatch.php?origin=KJFK&dest=EGLL&aircraft=B39M&route=HAPIE6%20HAPIE%20N247A%20ALLRY%20DCT%20KANNI%20N866B%20BEXET%20UL9%20BOGNA%20UL607%20REDFA&altitude=FL380&airline=UAL&flight=918&reg=N39MAX&captain=JANE%20SMITH&dispatcher=JOHN%20DOE&pax=180&cargo=12.5&altn=EGKK&units=LBS&navlog=1&static_id=ADVANCED_EXAMPLE

=== Detailed Flight Plan Analysis ===
Flight Plan Analysis for 857341
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Route: LIRF → LFBO
Origin: FIUMICINO (LIRF) - Rome, Italy
Destination: BLAGNAC (LFBO) - Toulouse, France

Aircraft: A320-251N (A20N)
Registration: G-INEO
Engine: CFM LEAP-1A

Flight Planning:
Distance: 532 nm
Route: POD9E PODOX T246 MIRSA DCT OMARD DCT STP DCT PADKO UL127 FJR DCT AFRIC AFRI8N
Cruise Altitude: FL320
Cost Index: 5

Timing:
Flight Time: 01:27
Block Time: 01:37
Taxi Out: 5 min
Taxi In: 5 min

Weight & Balance (KGS):
Operating Empty Weight: 42500
Zero Fuel Weight: 57298
Takeoff Weight: 62442
Landing Weight: 59577
Payload: 14798 (Pax: 151 @ 85)

Fuel Planning (KGS):
Total Planned: 5374
Trip Fuel: 2865
Taxi Fuel: 230
Alternate Fuel: 1087
Contingency: 287
Reserve: 905
Extra: 0
Average Flow: 1956

Weather:
Origin METAR: LIRF 201520Z 26010KT 9999 FEW030TCU 30/17 Q1015 NOSIG
Destination METAR: LFBO 201530Z AUTO VRB05KT CAVOK 36/12 Q1018 NOSIG
Average Wind: Variable/Light kts
Temperature Range: 17°C to 36°C (avg: 26°C)

Alternate: MERIGNAC (LFBD)
Distance: 132 nm, Bearing: 246°
Fuel Required: 1087 KGS

Navigation Log: Available (structure may vary)

Generated Files:
PDF: Available
XML: Available
FSX/P3D: Available

✅ Advanced analysis completed!
```

## Advanced Features Explained

### Custom Aircraft Data

This example shows how to define custom aircraft with detailed specifications including:

- **Weight Limits**: OEW, MZFW, MTOW, MLW
- **Fuel Capacity**: Maximum fuel capacity
- **Equipment Codes**: Navigation and communication equipment
- **Performance Data**: Category, engine type, passenger capacity

### Flight Plan Analysis

The comprehensive analysis includes:

- **Route Analysis**: Detailed route information and waypoints
- **Performance Calculations**: Weight, balance, and fuel planning
- **Weather Integration**: Current conditions and forecasts
- **Timing**: Flight time, block time, and taxi times
- **Alternate Planning**: Alternate airport and fuel requirements

### Data Extraction

Learn how to extract various types of data:

- **Airport Information**: Names, codes, cities, countries
- **Aircraft Details**: Type, registration, engine specifications
- **Operational Data**: Fuel flows, weights, distances
- **Weather Data**: METAR reports, wind, temperature
- **File Links**: Access to generated PDF, XML, and flight simulator files

## Use Cases

This advanced example is perfect for:

1. **Flight Operations Centers**: Comprehensive flight planning and analysis
2. **Training Applications**: Teaching flight planning concepts
3. **Analytics Tools**: Extracting data for further analysis
4. **Custom Integrations**: Building specialized flight planning tools
5. **Performance Monitoring**: Analyzing fuel efficiency and performance

## Integration Patterns

### Data Pipeline Integration

```go
// Extract key performance metrics
metrics := map[string]interface{}{
    "flight_time":    flightPlan.Times.FlightTime,
    "fuel_required":  flightPlan.Fuel.Plan,
    "distance":       flightPlan.General.Distance,
    "avg_fuel_flow":  flightPlan.Fuel.AvgFuelFlow,
    "takeoff_weight": flightPlan.Weights.TakeoffWt,
}

// Send to analytics system
sendToAnalytics(metrics)
```

### Custom Reporting

```go
// Generate custom reports
report := generateFlightReport(flightPlan)
saveToFile("flight_report.json", report)
emailReport(report, "ops@airline.com")
```

## Error Handling Best Practices

```go
flightPlan, err := simbrief.GetFlightPlanByUserID(userID)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "not found"):
        fmt.Println("No flight plan found for user ID")
        return
    case strings.Contains(err.Error(), "unauthorized"):
        fmt.Println("Invalid user ID or access denied")
        return
    case strings.Contains(err.Error(), "timeout"):
        fmt.Println("Request timeout - SimBrief may be busy")
        return
    default:
        log.Fatalf("Unexpected error: %v", err)
    }
}
```

## What You'll Learn

1. **Advanced API Usage**: Complex request building and parameter management
2. **Data Structure Navigation**: Working with nested response structures
3. **Custom Aircraft Configuration**: Defining and using custom aircraft data
4. **Performance Analysis**: Understanding flight planning calculations
5. **Weather Integration**: Working with meteorological data
6. **File Management**: Accessing generated flight plan files
7. **Error Handling**: Robust error handling for production use
8. **Integration Patterns**: Building flight planning into larger systems

## Next Steps

After running this example, consider:

1. **Build a Flight Tracker**: Combine with real-time flight data
2. **Create Custom Reports**: Generate PDF or Excel reports
3. **Integrate with Databases**: Store flight plans for historical analysis
4. **Add Performance Monitoring**: Track fuel efficiency over time
5. **Build a Web Interface**: Create a web-based flight planning tool

## Performance Considerations

- **Caching**: Cache aircraft data and route information
- **Rate Limiting**: Respect SimBrief API rate limits
- **Error Recovery**: Implement retry logic for transient failures
- **Data Validation**: Validate input parameters before API calls

## Related Documentation

- [Basic Example](../basic/) - Getting started with the SDK
- [Main README](../../README.md) - Full SDK documentation
- [Flight Planning Guide](../../docs/flight-planning.md) - Flight planning concepts
- [Performance Documentation](../../docs/performance.md) - Performance optimization
- [API Integration Guide](../../docs/api-integration.md) - Best practices

## Troubleshooting

### Common Issues

1. **Missing User ID**: Ensure `SIMBRIEF_USER_ID` environment variable is set
2. **No Flight Plans**: Create a flight plan on SimBrief website first
3. **API Timeouts**: Check internet connection and SimBrief status
4. **Invalid Aircraft**: Verify custom aircraft data is properly formatted

### Debug Mode

Enable debug logging to see detailed API interactions:

```bash
export SIMBRIEF_DEBUG=true
go run main.go
```

This will output detailed information about API requests and responses.
