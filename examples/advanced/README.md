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
Route: KJFK → EGLL
Origin: John F Kennedy Intl (KJFK) - New York, United States
Destination: London Heathrow (EGLL) - London, United Kingdom

Aircraft: Boeing 777-200LR (B77L)
Registration: N787BA
Engine: GE90-110B1

Flight Planning:
Distance: 3459 nm
Route: HAPIE6 HAPIE N247A ALLRY DCT KANNI N866B BEXET UL9 BOGNA UL607 REDFA
Cruise Altitude: FL380
Cost Index: 100

Timing:
Flight Time: 07:45
Block Time: 08:15
Taxi Out: 15 min
Taxi In: 15 min

Weight & Balance (LBS):
Operating Empty Weight: 346500
Zero Fuel Weight: 438900
Takeoff Weight: 766000
Landing Weight: 445200
Payload: 92400 (Pax: 340 @ 190)

Fuel Planning (LBS):
Total Planned: 327100
Trip Fuel: 281300
Taxi Fuel: 2800
Alternate Fuel: 18500
Contingency: 8400
Reserve: 13100
Extra: 3000
Average Flow: 36280

Weather:
Origin METAR: KJFK 211251Z 25018KT 10SM FEW050 BKN120 BKN250 M03/M15 A3012 RMK AO2 SLP203 T10281150
Destination METAR: EGLL 211220Z 27015KT 9999 FEW035 SCT120 02/M05 Q1018 NOSIG
Average Wind: 285°/35 kts
Temperature Range: -8°C to 2°C (avg: -3°C)

Alternate: London Gatwick (EGKK)
Distance: 24 nm, Bearing: 180°
Fuel Required: 18500 LBS

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
