package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrlm-net/simbrief/pkg/client"
	"github.com/mrlm-net/simbrief/pkg/types"
)

func main() {
	// Initialize the SimBrief client
	// Note: API key is not required for fetching existing flight plans
	simbrief := client.NewClient("")

	// Example 1: Get supported aircraft types and layouts
	fmt.Println("=== Getting Supported Options ===")
	options, err := simbrief.GetSupportedOptions()
	if err != nil {
		log.Fatalf("Failed to get supported options: %v", err)
	}

	fmt.Printf("Found %d aircraft types and %d layouts\n", len(options.Aircraft), len(options.Layouts))

	// Show a few aircraft examples
	count := 0
	for id, aircraft := range options.Aircraft {
		if count >= 5 {
			break
		}
		fmt.Printf("  Aircraft: %s - %s (Accuracy: %s)\n", id, aircraft.Name, aircraft.Accuracy)
		count++
	}

	// Show layout examples
	count = 0
	for id, layout := range options.Layouts {
		if count >= 3 {
			break
		}
		fmt.Printf("  Layout: %s - %s\n", id, layout.NameLong)
		count++
	}

	// Example 2: Build a flight plan request
	fmt.Println("\n=== Building Flight Plan Request ===")

	// Create a flight plan using the builder pattern
	builder := client.NewFlightPlan("KJFK", "KLAX", "B738")
	request := builder.
		Route("HUSKY6 PENNS CRP LAS1").
		Altitude("FL380").
		Registration("N123AB").
		Captain("JOHN DOE").
		Passengers(148).
		StaticID("EXAMPLE_FLIGHT_123").
		EnableNavLog().
		Units(types.UnitsLBS).
		Build()

	// Validate the request
	if err := simbrief.ValidateFlightPlanRequest(request); err != nil {
		log.Fatalf("Flight plan validation failed: %v", err)
	}

	// Generate URL for flight plan creation (requires browser authentication)
	url := simbrief.GenerateFlightPlanURL(request)
	fmt.Printf("Flight plan generation URL: %s\n", url)

	// Example 3: Fetch existing flight plan data
	fmt.Println("\n=== Fetching Flight Plan Data ===")

	// Use a real user ID if provided via environment variable
	userID := os.Getenv("SIMBRIEF_USER_ID")
	if userID == "" {
		fmt.Println("Set SIMBRIEF_USER_ID environment variable to test data fetching")
		fmt.Println("Example: export SIMBRIEF_USER_ID=857341")
		return
	}

	// Fetch the latest flight plan for the user
	flightPlan, err := simbrief.GetFlightPlanByUserID(userID)
	if err != nil {
		log.Printf("Failed to fetch flight plan: %v", err)
		return
	}

	// Display flight plan information
	fmt.Printf("Flight: %s → %s\n", flightPlan.Origin.ICAO, flightPlan.Destination.ICAO)
	fmt.Printf("Aircraft: %s (%s)\n", flightPlan.Aircraft.ICAO, flightPlan.Aircraft.Name)
	fmt.Printf("Route: %s\n", flightPlan.General.Route)
	fmt.Printf("Distance: %s nm\n", flightPlan.General.Distance)
	fmt.Printf("Planned Fuel: %s %s\n", flightPlan.Fuel.Plan, flightPlan.General.Units)
	fmt.Printf("Flight Time: %s\n", flightPlan.Times.FlightTime)

	// Show first few navigation fixes
	if flightPlan.NavLog != nil {
		fmt.Println("\nNavigation log available (structure may vary)")
	}

	fmt.Println("\n✅ Example completed successfully!")
}
