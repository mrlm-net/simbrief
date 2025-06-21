package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mrlm-net/simbrief/pkg/types"
)

// FlightPlanBuilder provides a fluent interface for building flight plan requests
type FlightPlanBuilder struct {
	request *types.FlightPlanRequest
}

// NewFlightPlan creates a new flight plan builder with required fields
func NewFlightPlan(origin, destination, aircraft string) *FlightPlanBuilder {
	return &FlightPlanBuilder{
		request: types.NewFlightPlanRequest(origin, destination, aircraft),
	}
}

// Route sets the flight route
func (b *FlightPlanBuilder) Route(route string) *FlightPlanBuilder {
	b.request.Route = route
	return b
}

// Airline sets the airline code
func (b *FlightPlanBuilder) Airline(airline string) *FlightPlanBuilder {
	b.request.Airline = airline
	return b
}

// FlightNumber sets the flight number
func (b *FlightPlanBuilder) FlightNumber(flightNumber string) *FlightPlanBuilder {
	b.request.FlightNumber = flightNumber
	return b
}

// Alternate sets the alternate airport
func (b *FlightPlanBuilder) Alternate(alternate string) *FlightPlanBuilder {
	b.request.Alternate = alternate
	return b
}

// DepartureTime sets the departure time
func (b *FlightPlanBuilder) DepartureTime(hour, minute int) *FlightPlanBuilder {
	b.request.DepartureHour = hour
	b.request.DepartureMinute = minute
	return b
}

// Date sets the departure date
func (b *FlightPlanBuilder) Date(date string) *FlightPlanBuilder {
	b.request.Date = date
	return b
}

// DateFromTime sets the departure date from a time.Time
func (b *FlightPlanBuilder) DateFromTime(t time.Time) *FlightPlanBuilder {
	b.request.Date = t.Format("02Jan06")
	return b
}

// Altitude sets the cruise altitude
func (b *FlightPlanBuilder) Altitude(altitude string) *FlightPlanBuilder {
	b.request.Altitude = altitude
	return b
}

// AltitudeFromFeet sets the cruise altitude from feet
func (b *FlightPlanBuilder) AltitudeFromFeet(feet int) *FlightPlanBuilder {
	b.request.Altitude = fmt.Sprintf("%d", feet)
	return b
}

// AltitudeFromFlightLevel sets the cruise altitude from flight level
func (b *FlightPlanBuilder) AltitudeFromFlightLevel(fl int) *FlightPlanBuilder {
	b.request.Altitude = fmt.Sprintf("FL%03d", fl)
	return b
}

// Passengers sets the number of passengers
func (b *FlightPlanBuilder) Passengers(pax int) *FlightPlanBuilder {
	b.request.Passengers = pax
	return b
}

// Cargo sets the cargo weight
func (b *FlightPlanBuilder) Cargo(cargo float64) *FlightPlanBuilder {
	b.request.Cargo = cargo
	return b
}

// Units sets the weight/fuel units
func (b *FlightPlanBuilder) Units(units types.Units) *FlightPlanBuilder {
	b.request.Units = units
	return b
}

// PlanFormat sets the OFP format
func (b *FlightPlanBuilder) PlanFormat(format types.PlanFormat) *FlightPlanBuilder {
	b.request.PlanFormat = string(format)
	return b
}

// Registration sets the aircraft registration
func (b *FlightPlanBuilder) Registration(reg string) *FlightPlanBuilder {
	b.request.Registration = reg
	return b
}

// CallSign sets the ATC callsign
func (b *FlightPlanBuilder) CallSign(callsign string) *FlightPlanBuilder {
	b.request.ATCCallsign = callsign
	return b
}

// Captain sets the captain's name
func (b *FlightPlanBuilder) Captain(name string) *FlightPlanBuilder {
	b.request.CaptainName = name
	return b
}

// Dispatcher sets the dispatcher's name
func (b *FlightPlanBuilder) Dispatcher(name string) *FlightPlanBuilder {
	b.request.DispatcherName = name
	return b
}

// StaticID sets a static reference ID
func (b *FlightPlanBuilder) StaticID(id string) *FlightPlanBuilder {
	b.request.StaticID = id
	return b
}

// EnableNavLog enables detailed navigation log
func (b *FlightPlanBuilder) EnableNavLog() *FlightPlanBuilder {
	enable := true
	b.request.NavLog = &enable
	return b
}

// DisableNavLog disables detailed navigation log
func (b *FlightPlanBuilder) DisableNavLog() *FlightPlanBuilder {
	disable := false
	b.request.NavLog = &disable
	return b
}

// EnableETOPS enables ETOPS planning
func (b *FlightPlanBuilder) EnableETOPS() *FlightPlanBuilder {
	enable := true
	b.request.ETOPS = &enable
	return b
}

// EnableStepClimbs enables step climb planning
func (b *FlightPlanBuilder) EnableStepClimbs() *FlightPlanBuilder {
	enable := true
	b.request.StepClimbs = &enable
	return b
}

// CustomAircraftData sets custom aircraft data
func (b *FlightPlanBuilder) CustomAircraftData(data *types.AircraftData) *FlightPlanBuilder {
	b.request.AircraftData = data
	return b
}

// TaxiTimes sets taxi out and taxi in times in minutes
func (b *FlightPlanBuilder) TaxiTimes(taxiOut, taxiIn int) *FlightPlanBuilder {
	b.request.TaxiOut = taxiOut
	b.request.TaxiIn = taxiIn
	return b
}

// Runways sets departure and arrival runways
func (b *FlightPlanBuilder) Runways(departure, arrival string) *FlightPlanBuilder {
	b.request.OriginRunway = departure
	b.request.DestRunway = arrival
	return b
}

// Build returns the completed flight plan request
func (b *FlightPlanBuilder) Build() *types.FlightPlanRequest {
	return b.request
}

// RouteHelper provides utilities for working with flight routes
type RouteHelper struct{}

// NewRouteHelper creates a new route helper
func NewRouteHelper() *RouteHelper {
	return &RouteHelper{}
}

// ParseRoute parses a route string and returns individual waypoints
func (rh *RouteHelper) ParseRoute(route string) []string {
	// Remove extra spaces and split by spaces
	route = strings.TrimSpace(route)
	if route == "" {
		return []string{}
	}

	parts := strings.Fields(route)
	waypoints := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			waypoints = append(waypoints, part)
		}
	}

	return waypoints
}

// ValidateICAOCode validates an ICAO airport code format
func (rh *RouteHelper) ValidateICAOCode(code string) bool {
	if len(code) != 4 {
		return false
	}

	// ICAO codes should be alphanumeric
	for _, r := range code {
		if !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}

	return true
}

// FormatFlightLevel formats a flight level from feet
func (rh *RouteHelper) FormatFlightLevel(feet int) string {
	fl := feet / 100
	return fmt.Sprintf("FL%03d", fl)
}

// ParseFlightLevel parses a flight level string and returns feet
func (rh *RouteHelper) ParseFlightLevel(flStr string) (int, error) {
	flStr = strings.ToUpper(strings.TrimSpace(flStr))

	if strings.HasPrefix(flStr, "FL") {
		flStr = strings.TrimPrefix(flStr, "FL")
		fl, err := strconv.Atoi(flStr)
		if err != nil {
			return 0, fmt.Errorf("invalid flight level format: %s", flStr)
		}
		return fl * 100, nil
	}

	// Try to parse as direct feet
	feet, err := strconv.Atoi(flStr)
	if err != nil {
		return 0, fmt.Errorf("invalid altitude format: %s", flStr)
	}

	return feet, nil
}

// FuelHelper provides utilities for fuel calculations
type FuelHelper struct{}

// NewFuelHelper creates a new fuel helper
func NewFuelHelper() *FuelHelper {
	return &FuelHelper{}
}

// ConvertLBSToKGS converts pounds to kilograms
func (fh *FuelHelper) ConvertLBSToKGS(lbs float64) float64 {
	return lbs * 0.453592
}

// ConvertKGSToLBS converts kilograms to pounds
func (fh *FuelHelper) ConvertKGSToLBS(kgs float64) float64 {
	return kgs / 0.453592
}

// ParseFuelValue parses a fuel value string that might contain weight or time
func (fh *FuelHelper) ParseFuelValue(value string) (float64, string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, "", fmt.Errorf("empty fuel value")
	}

	// Check if it contains a time component (e.g., "0.05/15" for percentage and minutes)
	if strings.Contains(value, "/") {
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			weight, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return 0, "", fmt.Errorf("invalid weight component: %s", parts[0])
			}
			return weight, parts[1], nil
		}
	}

	// Try to parse as a simple float
	weight, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid fuel value: %s", value)
	}

	return weight, "", nil
}

// TimeHelper provides utilities for time calculations
type TimeHelper struct{}

// NewTimeHelper creates a new time helper
func NewTimeHelper() *TimeHelper {
	return &TimeHelper{}
}

// ParseTimeString parses a time string in HH:MM format
func (th *TimeHelper) ParseTimeString(timeStr string) (int, int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format, expected HH:MM")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid hour: %s", parts[0])
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minute: %s", parts[1])
	}

	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("hour must be between 0 and 23")
	}

	if minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("minute must be between 0 and 59")
	}

	return hour, minute, nil
}

// FormatTimeString formats hour and minute into HH:MM string
func (th *TimeHelper) FormatTimeString(hour, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// ParseDuration parses a duration string in HH:MM format and returns total minutes
func (th *TimeHelper) ParseDuration(durationStr string) (int, error) {
	hour, minute, err := th.ParseTimeString(durationStr)
	if err != nil {
		return 0, err
	}

	return hour*60 + minute, nil
}

// FormatDuration formats total minutes into HH:MM string
func (th *TimeHelper) FormatDuration(totalMinutes int) string {
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
