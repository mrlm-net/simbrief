package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrlm-net/simbrief/pkg/client"
	"github.com/mrlm-net/simbrief/pkg/types"
)

func main() {
	// This example demonstrates advanced SimBrief SDK features
	simbrief := client.NewClient()

	userID := os.Getenv("SIMBRIEF_USER_ID")
	if userID == "" {
		fmt.Println("Set SIMBRIEF_USER_ID environment variable to run this example")
		fmt.Println("Example: export SIMBRIEF_USER_ID=857341")
		return
	}

	// Example 1: Custom Aircraft Data
	fmt.Println("=== Creating Flight Plan with Custom Aircraft Data ===")

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

	// Build a comprehensive flight plan
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

	// Print generated URL
	url := simbrief.GenerateFlightPlanURL(request)
	fmt.Printf("Advanced flight plan URL: %s\n", url)

	// Example 2: Fetch and analyze detailed flight data
	fmt.Println("\n=== Detailed Flight Plan Analysis ===")

	flightPlan, err := simbrief.GetFlightPlanByUserID(userID)
	if err != nil {
		log.Fatalf("Failed to fetch flight plan: %v", err)
	}

	// Basic flight information
	fmt.Printf("Flight Plan Analysis for %s\n", userID)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Route: %s → %s\n", flightPlan.Origin.ICAO, flightPlan.Destination.ICAO)
	fmt.Printf("Origin: %s (%s) - %s, %s\n",
		flightPlan.Origin.Name, flightPlan.Origin.ICAO,
		flightPlan.Origin.City, flightPlan.Origin.Country)
	fmt.Printf("Destination: %s (%s) - %s, %s\n",
		flightPlan.Destination.Name, flightPlan.Destination.ICAO,
		flightPlan.Destination.City, flightPlan.Destination.Country)

	// Aircraft details
	fmt.Printf("\nAircraft: %s (%s)\n", flightPlan.Aircraft.Name, flightPlan.Aircraft.ICAO)
	fmt.Printf("Registration: %s\n", flightPlan.Aircraft.Registration)
	fmt.Printf("Engine: %s\n", flightPlan.Aircraft.Engine)

	// Flight planning details
	fmt.Printf("\nFlight Planning:\n")
	fmt.Printf("Distance: %s nm\n", flightPlan.General.Distance)
	fmt.Printf("Route: %s\n", flightPlan.General.Route)
	fmt.Printf("Cruise Altitude: %s\n", flightPlan.General.CruiseAltitude)
	fmt.Printf("Cost Index: %s\n", flightPlan.General.CostIndex)

	// Timing information
	fmt.Printf("\nTiming:\n")
	fmt.Printf("Flight Time: %s\n", flightPlan.Times.FlightTime)
	fmt.Printf("Block Time: %s\n", flightPlan.Times.BlockTime)
	fmt.Printf("Taxi Out: %s min\n", flightPlan.Times.TaxiOut)
	fmt.Printf("Taxi In: %s min\n", flightPlan.Times.TaxiIn)

	// Weight and balance
	fmt.Printf("\nWeight & Balance (%s):\n", flightPlan.General.Units)
	fmt.Printf("Operating Empty Weight: %s\n", flightPlan.Weights.OEW)
	fmt.Printf("Zero Fuel Weight: %s\n", flightPlan.Weights.ZFW)
	fmt.Printf("Takeoff Weight: %s\n", flightPlan.Weights.TakeoffWt)
	fmt.Printf("Landing Weight: %s\n", flightPlan.Weights.LandingWt)
	fmt.Printf("Payload: %s (Pax: %s @ %s)\n",
		flightPlan.Weights.Payload, flightPlan.Weights.PaxCount, flightPlan.Weights.PaxWeight)

	// Fuel planning
	fmt.Printf("\nFuel Planning (%s):\n", flightPlan.General.Units)
	fmt.Printf("Total Planned: %s\n", flightPlan.Fuel.Plan)
	fmt.Printf("Trip Fuel: %s\n", flightPlan.Fuel.Trip)
	fmt.Printf("Taxi Fuel: %s\n", flightPlan.Fuel.Taxi)
	fmt.Printf("Alternate Fuel: %s\n", flightPlan.Fuel.Alternate)
	fmt.Printf("Contingency: %s\n", flightPlan.Fuel.Contingency)
	fmt.Printf("Reserve: %s\n", flightPlan.Fuel.Reserve)
	fmt.Printf("Extra: %s\n", flightPlan.Fuel.Extra)
	fmt.Printf("Average Flow: %s\n", flightPlan.Fuel.AvgFuelFlow)

	// Weather information
	if flightPlan.Weather.OriginMETAR != "" || flightPlan.Weather.DestMETAR != "" {
		fmt.Printf("\nWeather:\n")
		if flightPlan.Weather.OriginMETAR != "" {
			fmt.Printf("Origin METAR: %s\n", flightPlan.Weather.OriginMETAR)
		}
		if flightPlan.Weather.DestMETAR != "" {
			fmt.Printf("Destination METAR: %s\n", flightPlan.Weather.DestMETAR)
		}
		fmt.Printf("Average Wind: %s°/%s kts\n",
			flightPlan.Weather.AvgWindDir, flightPlan.Weather.AvgWindSpd)
		fmt.Printf("Temperature Range: %s°C to %s°C (avg: %s°C)\n",
			flightPlan.Weather.MinTemp, flightPlan.Weather.MaxTemp, flightPlan.Weather.AvgTemp)
	}

	// Alternate information
	if flightPlan.Alternate.ICAO != "" {
		fmt.Printf("\nAlternate: %s (%s)\n", flightPlan.Alternate.Name, flightPlan.Alternate.ICAO)
		fmt.Printf("Distance: %s nm, Bearing: %s°\n",
			flightPlan.Alternate.Distance, flightPlan.Alternate.Bearing)
		fmt.Printf("Fuel Required: %s %s\n",
			flightPlan.Alternate.FuelRequired, flightPlan.General.Units)
	}

	// Navigation log summary
	if flightPlan.NavLog != nil {
		fmt.Printf("\nNavigation Log: Available (structure may vary)\n")
	}

	// File links
	if flightPlan.Files.PDFLink != nil {
		fmt.Printf("\nGenerated Files:\n")
		if flightPlan.Files.PDFLink != nil {
			fmt.Printf("PDF: Available\n")
		}
		if flightPlan.Files.XMLLink != nil {
			fmt.Printf("XML: Available\n")
		}
		if flightPlan.Files.PLNLink != nil {
			fmt.Printf("FSX/P3D: Available\n")
		}
	}

	fmt.Println("\n✅ Advanced analysis completed!")
}
