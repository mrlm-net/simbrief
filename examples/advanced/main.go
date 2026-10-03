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
	fmt.Printf("Origin: %s (%s) runway %s, elevation %d ft\n",
		flightPlan.Origin.Name, flightPlan.Origin.ICAO,
		flightPlan.Origin.Runway, flightPlan.Origin.Elevation.Int())
	fmt.Printf("Destination: %s (%s) runway %s, elevation %d ft\n",
		flightPlan.Destination.Name, flightPlan.Destination.ICAO,
		flightPlan.Destination.Runway, flightPlan.Destination.Elevation.Int())

	// ATC flight plan
	fmt.Printf("\nATC:\n")
	fmt.Printf("Callsign: %s (rules %s, type %s)\n",
		flightPlan.ATC.Callsign, flightPlan.ATC.FlightRules, flightPlan.ATC.FlightType)
	fmt.Printf("Initial level: %s%s\n", flightPlan.ATC.InitialAltUnit, flightPlan.ATC.InitialAlt)
	fmt.Printf("Route: %s\n", flightPlan.ATC.Route)

	// Aircraft details
	fmt.Printf("\nAircraft: %s (%s)\n", flightPlan.Aircraft.Name, flightPlan.Aircraft.ICAOCode)
	fmt.Printf("Registration: %s\n", flightPlan.Aircraft.Registration)
	fmt.Printf("Engines: %s\n", flightPlan.Aircraft.Engines)

	// Flight planning details
	fmt.Printf("\nFlight Planning:\n")
	fmt.Printf("Distance: %d nm\n", flightPlan.General.RouteDistance.Int())
	fmt.Printf("Route: %s\n", flightPlan.General.Route)
	fmt.Printf("SID: %s  STAR: %s\n", flightPlan.General.SIDIdent, flightPlan.General.STARIdent)
	fmt.Printf("Cruise Altitude: %d ft\n", flightPlan.General.InitialAltitude.Int())
	fmt.Printf("Cost Index: %d\n", flightPlan.General.CostIndex.Int())
	fmt.Printf("Average Wind: %s°/%s kts\n",
		flightPlan.General.AvgWindDir, flightPlan.General.AvgWindSpd)

	// Timing information (JSON v2: ISO 8601 times, HH:MM:SS durations)
	fmt.Printf("\nTiming:\n")
	fmt.Printf("Flight Time: %s\n", flightPlan.Times.EstTimeEnroute)
	fmt.Printf("Block Time: %s\n", flightPlan.Times.EstBlock)
	fmt.Printf("Taxi Out: %s\n", flightPlan.Times.TaxiOut)
	fmt.Printf("Taxi In: %s\n", flightPlan.Times.TaxiIn)

	// Weight and balance
	units := flightPlan.Params.Units
	fmt.Printf("\nWeight & Balance (%s):\n", units)
	fmt.Printf("Operating Empty Weight: %s\n", flightPlan.Weights.OEW)
	fmt.Printf("Zero Fuel Weight: %s\n", flightPlan.Weights.EstZFW)
	fmt.Printf("Takeoff Weight: %s\n", flightPlan.Weights.EstTOW)
	fmt.Printf("Landing Weight: %s\n", flightPlan.Weights.EstLDW)
	fmt.Printf("Payload: %s (Pax: %s @ %s)\n",
		flightPlan.Weights.Payload, flightPlan.Weights.PaxCount, flightPlan.Weights.PaxWeight)

	// Fuel planning
	fmt.Printf("\nFuel Planning (%s):\n", units)
	fmt.Printf("Block Fuel: %s\n", flightPlan.Fuel.PlanRamp)
	fmt.Printf("Trip Fuel: %s\n", flightPlan.Fuel.EnrouteBurn)
	fmt.Printf("Taxi Fuel: %s\n", flightPlan.Fuel.Taxi)
	fmt.Printf("Alternate Fuel: %s\n", flightPlan.Fuel.AlternateBurn)
	fmt.Printf("Contingency: %s\n", flightPlan.Fuel.Contingency)
	fmt.Printf("Reserve: %s\n", flightPlan.Fuel.Reserve)
	fmt.Printf("Extra: %s\n", flightPlan.Fuel.Extra)
	fmt.Printf("Average Flow: %s\n", flightPlan.Fuel.AvgFuelFlow)

	// Weather information
	if flightPlan.Weather.OrigMETAR != "" || flightPlan.Weather.DestMETAR != "" {
		fmt.Printf("\nWeather:\n")
		if flightPlan.Weather.OrigMETAR != "" {
			fmt.Printf("Origin METAR: %s\n", flightPlan.Weather.OrigMETAR)
		}
		if flightPlan.Weather.DestMETAR != "" {
			fmt.Printf("Destination METAR: %s\n", flightPlan.Weather.DestMETAR)
		}
	}

	// Alternates
	for _, altn := range flightPlan.Alternates {
		fmt.Printf("\nAlternate: %s (%s) runway %s\n", altn.Name, altn.ICAO, altn.Runway)
		fmt.Printf("Distance: %s nm, Track: %s°\n", altn.Distance, altn.TrackMag)
		fmt.Printf("Fuel Required: %s %s\n", altn.Burn, units)
	}

	// Navigation log
	fmt.Printf("\nNavigation Log (%d fixes):\n", len(flightPlan.NavLog))
	for _, fix := range flightPlan.NavLog {
		fmt.Printf("  %-6s %-4s %-8s %6d ft  %s\n",
			fix.Ident, fix.Type, fix.ViaAirway, fix.AltitudeFeet.Int(), fix.Stage)
	}

	// File links
	if flightPlan.Files.PDF.Link != "" {
		fmt.Printf("\nGenerated Files:\n")
		fmt.Printf("PDF: %s%s\n", flightPlan.Files.Directory, flightPlan.Files.PDF.Link)
		fmt.Printf("Other formats: %d\n", len(flightPlan.Files.Files))
	}

	fmt.Println("\n✅ Advanced analysis completed!")
}
