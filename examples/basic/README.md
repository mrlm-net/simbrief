# Basic SimBrief Example

This example demonstrates basic usage of the SimBrief SDK for common operations like getting supported aircraft types and fetching existing flight plans.

## Features Demonstrated

- **Get Supported Options**: Retrieve all available aircraft types and layouts
- **Fetch Flight Plans**: Get existing flight plans by user ID
- **Basic Error Handling**: Handle common API errors
- **Response Parsing**: Extract useful information from API responses

## Prerequisites

- Go 1.21 or higher
- Internet connection for SimBrief API access
- Optional: SimBrief user ID for fetching existing flight plans

## Environment Variables

```bash
# Optional: Set user ID to fetch existing flight plans
export SIMBRIEF_USER_ID=857341

# Optional: Enable debug logging
export SIMBRIEF_DEBUG=true
```

## Running the Example

```bash
cd examples/basic
go run main.go
```

## Code Overview

### 1. Getting Supported Aircraft Types

```go
// Initialize client (no API key needed for public data)
simbrief := client.NewClient()

// Get all supported aircraft and layouts
options, err := simbrief.GetSupportedOptions()
if err != nil {
    log.Fatalf("Failed to get supported options: %v", err)
}

// Display aircraft information
for id, aircraft := range options.Aircraft {
    fmt.Printf("ICAO: %s, Name: %s, Engine: %s\n", 
        id, aircraft.Name, aircraft.Engine)
}
```

### 2. Fetching Existing Flight Plans

```go
userID := os.Getenv("SIMBRIEF_USER_ID")
if userID != "" {
    // Fetch flight plan by user ID
    flightPlan, err := simbrief.GetFlightPlanByUserID(userID)
    if err != nil {
        log.Printf("Error fetching flight plan: %v", err)
    } else {
        fmt.Printf("Flight: %s to %s\n", 
            flightPlan.Origin.ICAO, flightPlan.Destination.ICAO)
        fmt.Printf("Distance: %.0f nm\n", flightPlan.General.DistanceNM)
        fmt.Printf("Aircraft: %s\n", flightPlan.Aircraft.Name)
    }
}
```

### 3. Error Handling

```go
options, err := simbrief.GetSupportedOptions()
if err != nil {
    // Handle different types of errors
    switch {
    case strings.Contains(err.Error(), "network"):
        fmt.Println("Network error - check internet connection")
    case strings.Contains(err.Error(), "timeout"):
        fmt.Println("Request timeout - try again later")
    default:
        fmt.Printf("API error: %v\n", err)
    }
    return
}
```

## Sample Output

```
=== Getting Supported Options ===
Found 847 aircraft types and 12 layouts

Sample Aircraft Types:
ICAO: A20N, Name: Airbus A320neo, Engine: CFM LEAP / PW GTF
ICAO: A21N, Name: Airbus A321neo, Engine: CFM LEAP / PW GTF
ICAO: A319, Name: Airbus A319, Engine: CFM56 / V2500
ICAO: A320, Name: Airbus A320, Engine: CFM56 / V2500
ICAO: A321, Name: Airbus A321, Engine: CFM56 / V2500

Sample Layouts:
ID: default, Name: SimBrief Default
ID: lido, Name: Lido mPilot
ID: jeppesen, Name: Jeppesen
ID: pfpx, Name: PFPX
ID: aivlasoft, Name: AivlaSoft EFB

=== Fetching Flight Plan ===
Flight plan found: KJFK to KLAX
Distance: 2475 nm
Aircraft: Boeing 737-800
Flight time: 05:23
Route: HAPIE6 HAPIE J80 WILMINGTON J82 LYNCH J134 TUS EAGUL4
```

## What You'll Learn

1. **Basic API Usage**: How to initialize the client and make API calls
2. **Data Structures**: Understanding the response formats
3. **Error Handling**: Proper error handling patterns
4. **Environment Configuration**: Using environment variables for configuration
5. **Response Parsing**: Extracting useful information from API responses

## Next Steps

After running this example, try:

1. **Modify the code** to filter aircraft by category (Light, Medium, Heavy)
2. **Add more error handling** for specific error types
3. **Cache the aircraft data** to avoid repeated API calls
4. **Explore the advanced example** for flight plan generation

## Related Documentation

- [Main README](../../README.md) - Full SDK documentation
- [Advanced Example](../advanced/) - Flight plan generation
- [API Integration Guide](../../docs/api-integration.md) - Best practices
- [Aircraft Management](../../docs/aircraft-management.md) - Aircraft configuration
