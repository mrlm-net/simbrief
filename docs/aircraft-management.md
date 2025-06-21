# Aircraft Management

This document covers aircraft configuration and management features of the SimBrief SDK.

## Table of Contents

- [Getting Supported Aircraft](#getting-supported-aircraft)
- [Custom Aircraft Configuration](#custom-aircraft-configuration)
- [Aircraft Performance Data](#aircraft-performance-data)
- [Aircraft Categories](#aircraft-categories)
- [Layout Management](#layout-management)
- [Validation](#validation)

## Getting Supported Aircraft

### Retrieving All Supported Aircraft

```go
client := client.NewClient("")  // API key not required for this operation

options, err := client.GetSupportedOptions()
if err != nil {
    log.Fatalf("Failed to get supported options: %v", err)
}

// List all aircraft types
for icao, aircraft := range options.Aircraft {
    fmt.Printf("ICAO: %s, Name: %s, Engine: %s\n", 
        icao, aircraft.Name, aircraft.Engine)
}
```

### Finding Specific Aircraft

```go
options, err := client.GetSupportedOptions()
if err != nil {
    return err
}

// Find Boeing 737 MAX aircraft
for icao, aircraft := range options.Aircraft {
    if strings.Contains(strings.ToLower(aircraft.Name), "737 max") {
        fmt.Printf("Found: %s - %s\n", icao, aircraft.Name)
    }
}
```

### Aircraft Search by Category

```go
// Find all wide-body aircraft (Heavy category)
heavyAircraft := make(map[string]types.Aircraft)
for icao, aircraft := range options.Aircraft {
    if aircraft.Category == "H" {  // Heavy category
        heavyAircraft[icao] = aircraft
    }
}
```

## Custom Aircraft Configuration

### Basic Custom Aircraft

```go
customAircraft := &types.AircraftData{
    ICAO:        "B39M",        // Custom ICAO code
    Name:        "737 MAX 9",   // Aircraft name
    Engines:     "CFM LEAP-1B", // Engine type
    Category:    "M",           // Medium category
    WakeClass:   "M",           // Medium wake turbulence
}
```

### Complete Aircraft Configuration

```go
customAircraft := &types.AircraftData{
    // Basic identification
    ICAO:            "B39M",
    Name:            "Boeing 737 MAX 9",
    Manufacturer:    "Boeing",
    Series:          "737 MAX",
    Variant:         "MAX 9",
    
    // Engine configuration
    Engines:         "CFM LEAP-1B",
    EngineCount:     2,
    ThrustPerEngine: 28000,      // lbs thrust per engine
    
    // Aircraft categories
    Category:        "M",         // ICAO aircraft category (L/M/H/J)
    WakeClass:       "M",         // Wake turbulence category
    EquipmentCode:   "SDFGHIRWY", // ICAO equipment codes
    
    // Performance limits
    MaxAltitude:     41000,       // Service ceiling (feet)
    CruiseSpeed:     0.785,       // Cruise Mach number
    MaxSpeed:        0.82,        // Maximum Mach number
    StallSpeed:      108,         // Stall speed (knots)
    
    // Weight limits (pounds)
    MaxTakeoffWeight:  194700,    // MTOW
    MaxLandingWeight:  155500,    // MLW  
    MaxZeroFuelWeight: 138500,    // MZFW
    EmptyWeight:       98000,     // OEW
    
    // Fuel system
    MaxFuelCapacity: 25820,       // Maximum fuel (pounds)
    FuelDensity:     6.7,         // Fuel density (lbs/gallon)
    
    // Performance data
    FuelBurnRate:    5500,        // Cruise fuel burn (lbs/hour)
    ClimbRate:       2500,        // Initial climb rate (ft/min)
    
    // Operating costs
    OperatingCost:   12000,       // Operating cost per hour ($)
    FuelCostPerGallon: 3.50,      // Fuel cost per gallon ($)
}
```

## Aircraft Performance Data

### Fuel Burn Configuration

```go
aircraftData := &types.AircraftData{
    // ... basic configuration
    
    // Detailed fuel burn by flight phase
    TaxiFuelBurn:    200,         // Taxi fuel burn (lbs)
    ClimbFuelBurn:   8000,        // Climb fuel burn (lbs)
    CruiseFuelBurn:  5500,        // Cruise fuel burn (lbs/hour)
    DescentFuelBurn: 3000,        // Descent fuel burn (lbs)
    
    // Fuel planning factors
    ContingencyPercent: 5,        // Contingency fuel percentage
    AlternateFuelTime:  45,       // Alternate fuel time (minutes)
    HoldingFuelTime:    30,       // Holding fuel time (minutes)
}
```

### Performance by Altitude

```go
// Configure performance data for different flight levels
aircraftData := &types.AircraftData{
    // ... basic configuration
    
    // Altitude-specific performance
    ServiceCeiling:    41000,     // Maximum service ceiling
    OptimalAltitude:   37000,     // Optimal cruise altitude
    SecondaryAltitude: 39000,     // Secondary cruise altitude
    
    // Speed by altitude
    ClimbSpeed:        250,       // Climb speed below 10,000ft
    CruiseSpeedKnots:  470,       // Cruise speed in knots
    DescentSpeed:      290,       // Descent speed
    
    // Fuel efficiency by altitude
    FuelEfficiencyFL350: 5.2,     // NM per gallon at FL350
    FuelEfficiencyFL370: 5.5,     // NM per gallon at FL370
    FuelEfficiencyFL390: 5.8,     // NM per gallon at FL390
}
```

## Aircraft Categories

### ICAO Aircraft Categories

```go
const (
    CategoryLight     = "L"  // Light aircraft (< 7,000 kg)
    CategoryMedium    = "M"  // Medium aircraft (7,000 - 136,000 kg)
    CategoryHeavy     = "H"  // Heavy aircraft (> 136,000 kg)
    CategorySuper     = "J"  // Super heavy (A380, An-225)
)

// Set aircraft category
aircraftData.Category = CategoryMedium
```

### Wake Turbulence Categories

```go
const (
    WakeLight  = "L"  // Light wake turbulence
    WakeMedium = "M"  // Medium wake turbulence  
    WakeHeavy  = "H"  // Heavy wake turbulence
    WakeSuper  = "J"  // Super wake turbulence
)

// Set wake turbulence category
aircraftData.WakeClass = WakeMedium
```

### Equipment Codes

```go
// Common ICAO equipment codes
const (
    EquipGPS     = "G"  // GPS capability
    EquipRNAV    = "R"  // RNAV capability
    EquipTCAS    = "C"  // TCAS equipped
    EquipMode_S  = "S"  // Mode S transponder
    EquipADS_B   = "B"  // ADS-B capability
)

// Combine equipment codes
aircraftData.EquipmentCode = "SDFGHIRWY"  // Full RNAV/GPS equipped
```

## Layout Management

### Getting Available Layouts

```go
options, err := client.GetSupportedOptions()
if err != nil {
    return err
}

// List all layout options
for id, layout := range options.Layouts {
    fmt.Printf("Layout ID: %s, Name: %s\n", id, layout.Name)
}
```

### Using Custom Layouts

```go
request := &types.FlightPlanRequest{
    Origin:      "KJFK",
    Destination: "KLAX", 
    Aircraft:    "B38M",
    Layout:      "lido",     // Use LIDO layout
}
```

### Layout-Specific Options

```go
// Available layout types
const (
    LayoutDefault  = "default"   // Standard SimBrief layout
    LayoutLido     = "lido"      // Lido mPilot layout
    LayoutJeppesen = "jeppesen"  // Jeppesen layout
    LayoutCustom   = "custom"    // Custom layout
)
```

## Validation

### Aircraft Data Validation

```go
func validateAircraftData(aircraft *types.AircraftData) error {
    if aircraft.ICAO == "" {
        return errors.New("aircraft ICAO code is required")
    }
    
    if aircraft.MaxTakeoffWeight <= 0 {
        return errors.New("invalid maximum takeoff weight")
    }
    
    if aircraft.FuelBurnRate <= 0 {
        return errors.New("invalid fuel burn rate")
    }
    
    // Validate category
    validCategories := []string{"L", "M", "H", "J"}
    if !contains(validCategories, aircraft.Category) {
        return errors.New("invalid aircraft category")
    }
    
    return nil
}
```

### Pre-Flight Validation

```go
func validateForFlight(aircraft *types.AircraftData, request *types.FlightPlanRequest) error {
    // Check if aircraft is suitable for route distance
    estimatedDistance := calculateDistance(request.Origin, request.Destination)
    maxRange := aircraft.MaxFuelCapacity / aircraft.FuelBurnRate * aircraft.CruiseSpeedKnots
    
    if estimatedDistance > maxRange {
        return fmt.Errorf("route distance (%0.f nm) exceeds aircraft range (%0.f nm)", 
            estimatedDistance, maxRange)
    }
    
    return nil
}
```

## Advanced Aircraft Management

### Aircraft Performance Profiles

```go
type PerformanceProfile struct {
    Aircraft     *types.AircraftData
    ClimbProfile []AltitudePerformance
    CruiseProfile []AltitudePerformance
    DescentProfile []AltitudePerformance
}

type AltitudePerformance struct {
    Altitude    int     // Flight level
    Speed       int     // Speed in knots
    FuelBurn    float64 // Fuel burn rate
    ClimbRate   int     // Climb/descent rate
}

// Create detailed performance profile
profile := &PerformanceProfile{
    Aircraft: customAircraft,
    ClimbProfile: []AltitudePerformance{
        {Altitude: 10000, Speed: 250, FuelBurn: 8000, ClimbRate: 2500},
        {Altitude: 20000, Speed: 300, FuelBurn: 7000, ClimbRate: 2000},
        {Altitude: 30000, Speed: 350, FuelBurn: 6000, ClimbRate: 1500},
    },
    // ... cruise and descent profiles
}
```

### Fleet Management

```go
type FleetManager struct {
    Aircraft map[string]*types.AircraftData
}

func NewFleetManager() *FleetManager {
    return &FleetManager{
        Aircraft: make(map[string]*types.AircraftData),
    }
}

func (fm *FleetManager) AddAircraft(icao string, aircraft *types.AircraftData) {
    fm.Aircraft[icao] = aircraft
}

func (fm *FleetManager) GetSuitableAircraft(distance float64, passengers int) []*types.AircraftData {
    var suitable []*types.AircraftData
    
    for _, aircraft := range fm.Aircraft {
        // Check range and passenger capacity
        maxRange := aircraft.MaxFuelCapacity / aircraft.FuelBurnRate * aircraft.CruiseSpeedKnots
        if maxRange >= distance && aircraft.MaxPassengers >= passengers {
            suitable = append(suitable, aircraft)
        }
    }
    
    return suitable
}
```

## Best Practices

1. **Use Standard ICAO Codes**: Always use official ICAO aircraft type codes when possible
2. **Validate Performance Data**: Ensure fuel burn rates and performance data are realistic
3. **Cache Aircraft Data**: Store frequently used aircraft configurations locally
4. **Regular Updates**: Keep aircraft data updated with latest performance specifications
5. **Error Handling**: Always validate aircraft data before using in flight plans
6. **Documentation**: Document custom aircraft configurations for team use
