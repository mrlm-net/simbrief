package types

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// FlightPlanRequest represents all possible parameters for generating a flight plan
// Based on official SimBrief API documentation at https://developers.navigraph.com/docs/simbrief/using-the-api
type FlightPlanRequest struct {
	// Required parameters
	Origin      string `form:"orig"` // ICAO origin airport (required)
	Destination string `form:"dest"` // ICAO destination airport (required)
	Aircraft    string `form:"type"` // Aircraft type (required)

	// Basic flight information
	Airline         string `form:"airline"` // Airline code (e.g., "ABC")
	FlightNumber    string `form:"fltnum"`  // Flight number (e.g., "1234")
	Date            string `form:"date"`    // Date format: 11JUL13
	DepartureHour   int    `form:"deph"`    // Departure hour (0-23)
	DepartureMinute int    `form:"depm"`    // Departure minute (0-59)
	Route           string `form:"route"`   // Flight route (e.g., "PLL GAROT OAL MOD4")
	ScheduledHour   int    `form:"steh"`    // Scheduled time hour
	ScheduledMinute int    `form:"stem"`    // Scheduled time minute

	// Aircraft details
	Registration string `form:"reg"`      // Aircraft registration (e.g., "N123XX")
	FinNumber    string `form:"fin"`      // Aircraft fin number (e.g., "123")
	SELCAL       string `form:"selcal"`   // Aircraft SELCAL (e.g., "ABCD")
	ATCCallsign  string `form:"callsign"` // ATC callsign (e.g., "ABC1234")

	// Crew and passenger info
	Passengers     int    `form:"pax"`    // Number of passengers (e.g., 100)
	CaptainName    string `form:"cpt"`    // Captain's name (e.g., "JOHN DOE")
	DispatcherName string `form:"dxname"` // Dispatcher's name (e.g., "JANE DOE")
	PilotID        string `form:"pid"`    // Pilot ID number (e.g., "12345")

	// Alternates and routing
	Alternate string `form:"altn"`       // Primary alternate airport (e.g., "KLAX")
	Altitude  string `form:"fl"`         // Cruise altitude (e.g., "34000", "FL340")
	AltnCount int    `form:"altn_count"` // Number of alternates (e.g., 4)
	AltnAvoid string `form:"altn_avoid"` // Avoid alternate airports (e.g., "KJFK KPHL KBWI")

	// Alternate 1-4 specific fields
	Altn1ID     string `form:"altn_1_id"`    // Alternate 1 identifier
	Altn1Runway string `form:"altn_1_rwy"`   // Alternate 1 runway
	Altn1Route  string `form:"altn_1_route"` // Alternate 1 routing
	Altn2ID     string `form:"altn_2_id"`    // Alternate 2 identifier
	Altn2Runway string `form:"altn_2_rwy"`   // Alternate 2 runway
	Altn2Route  string `form:"altn_2_route"` // Alternate 2 routing
	Altn3ID     string `form:"altn_3_id"`    // Alternate 3 identifier
	Altn3Runway string `form:"altn_3_rwy"`   // Alternate 3 runway
	Altn3Route  string `form:"altn_3_route"` // Alternate 3 routing
	Altn4ID     string `form:"altn_4_id"`    // Alternate 4 identifier
	Altn4Runway string `form:"altn_4_rwy"`   // Alternate 4 runway
	Altn4Route  string `form:"altn_4_route"` // Alternate 4 routing

	// Fuel and weight
	FuelFactor     string  `form:"fuelfactor"`      // Fuel factor (e.g., "P00")
	ManualZFW      float64 `form:"manualzfw"`       // Manual zero fuel weight (e.g., 40.1)
	AddedFuel      string  `form:"addedfuel"`       // Extra fuel (e.g., "0.5", "20")
	AddedFuelUnits string  `form:"addedfuel_units"` // Extra fuel units ("wgt" or "min")
	ContFuelPct    string  `form:"contpct"`         // Contingency fuel (e.g., "0.05", "0.05/15")
	ReserveFuel    int     `form:"resvrule"`        // Reserve fuel minutes (e.g., 45)
	Cargo          float64 `form:"cargo"`           // Cargo weight (e.g., 5.0)

	// Taxi and runway
	TaxiOut      int    `form:"taxiout"` // Taxi out time minutes (e.g., 10)
	TaxiIn       int    `form:"taxiin"`  // Taxi in time minutes (e.g., 4)
	OriginRunway string `form:"origrwy"` // Departure runway (e.g., "06L")
	DestRunway   string `form:"destrwy"` // Arrival runway (e.g., "36R")

	// Performance profiles
	ClimbProfile   string `form:"climb"`   // Climb profile (e.g., "250/300/78")
	DescentProfile string `form:"descent"` // Descent profile (e.g., "84/280/250")
	CruiseProfile  string `form:"cruise"`  // Cruise profile ("LRC", "CI")
	CostIndex      string `form:"civalue"` // Cost index (e.g., "25", "AUTO")

	// Aircraft data (JSON string) - see official docs for structure
	AircraftData *AircraftData `form:"acdata,omitempty"` // Custom aircraft data

	// ETOPS
	ETOPSRule string `form:"etopsrule"` // ETOPS rule (e.g., "180", "207")

	// Custom remarks and static ID
	ManualRemarks string `form:"manualrmk"` // Custom remarks (can include \n for line breaks)
	StaticID      string `form:"static_id"` // Static reference ID (e.g., "ABC_12345")

	// OFP Options
	PlanFormat     string `form:"planformat"`   // Plan format (e.g., "LIDO")
	Units          Units  `form:"units"`        // Units ("LBS" or "KGS")
	NavLog         *bool  `form:"navlog"`       // Detailed navlog (1 or 0)
	ETOPS          *bool  `form:"etops"`        // ETOPS planning (1 or 0)
	StepClimbs     *bool  `form:"stepclimbs"`   // Plan stepclimbs (1 or 0)
	RunwayAnalysis *bool  `form:"tlr"`          // Runway analysis (1 or 0)
	NOTAMs         *bool  `form:"notams"`       // Include NOTAMs (1 or 0)
	FIRNOTAMs      *bool  `form:"firnot"`       // FIR NOTAMs (1 or 0)
	Maps           string `form:"maps"`         // Flight maps ("detail", "simple", "none")
	OmitSIDs       *bool  `form:"omit_sids"`    // Disable SIDs (1 or 0)
	OmitSTARs      *bool  `form:"omit_stars"`   // Disable STARs (1 or 0)
	FindSIDSTAR    string `form:"find_sidstar"` // Auto-insert SID/STARs ("R" or "C")
}

// AircraftData represents custom aircraft data as JSON
// Based on official SimBrief API documentation
type AircraftData struct {
	// ICAO flight plan fields - all 3 must be specified together
	Category    string `json:"cat,omitempty"`         // Aircraft weight category (L, M, H, J)
	Equipment   string `json:"equip,omitempty"`       // Onboard equipment string
	Transponder string `json:"transponder,omitempty"` // Transponder type

	// Performance based navigation (must start with "PBN/")
	PBN string `json:"pbn,omitempty"` // Performance based navigation

	// Additional Section 18 information
	ExtraRemark string `json:"extrarmk,omitempty"` // Additional Section 18 info

	// Weight and capacity (all weights in thousands of pounds, up to 3 decimal places)
	MaxPax  string  `json:"maxpax,omitempty"`  // Maximum passenger count
	OEW     float64 `json:"oew,omitempty"`     // Operating empty weight (thousands of pounds)
	MZFW    float64 `json:"mzfw,omitempty"`    // Max zero fuel weight (thousands of pounds)
	MTOW    float64 `json:"mtow,omitempty"`    // Max takeoff weight (thousands of pounds)
	MLW     float64 `json:"mlw,omitempty"`     // Max landing weight (thousands of pounds)
	MaxFuel float64 `json:"maxfuel,omitempty"` // Max fuel capacity (thousands of pounds)

	// Additional aircraft identification
	HexCode string `json:"hexcode,omitempty"` // ICAO Mode-S Code
	Per     string `json:"per,omitempty"`     // ICAO performance category
	PaxWgt  int    `json:"paxwgt,omitempty"`  // Average passenger weight (pounds)

	// For approximating unsupported aircraft types
	ICAO    string `json:"icao,omitempty"`    // ICAO aircraft identifier (max 4 chars)
	Name    string `json:"name,omitempty"`    // Aircraft name (max 12 chars)
	Engines string `json:"engines,omitempty"` // Engine type (max 12 chars)
}

// String returns the JSON string representation of AircraftData
func (ad *AircraftData) String() string {
	if ad == nil {
		return ""
	}
	data, _ := json.Marshal(ad)
	return string(data)
}

// FetchRequest represents parameters for fetching existing flight plan data
type FetchRequest struct {
	UserID   string `form:"userid,omitempty"`    // SimBrief user ID
	Username string `form:"username,omitempty"`  // SimBrief username
	StaticID string `form:"static_id,omitempty"` // Static reference ID
	JSON     bool   `form:"json,omitempty"`      // Request JSON format (default: XML)
}

// ToQueryParams converts FetchRequest to URL query parameters
func (fr *FetchRequest) ToQueryParams() string {
	values := url.Values{}

	if fr.UserID != "" {
		values.Add("userid", fr.UserID)
	}
	if fr.Username != "" {
		values.Add("username", fr.Username)
	}
	if fr.StaticID != "" {
		values.Add("static_id", fr.StaticID)
	}
	if fr.JSON {
		values.Add("json", "1")
	}

	if len(values) == 0 {
		return ""
	}
	return "?" + values.Encode()
}

// NewFlightPlanRequest creates a new flight plan request with required fields
func NewFlightPlanRequest(origin, destination, aircraft string) *FlightPlanRequest {
	return &FlightPlanRequest{
		Origin:      origin,
		Destination: destination,
		Aircraft:    aircraft,
	}
}

// ToURLValues converts the FlightPlanRequest to url.Values for form submission
func (fpr *FlightPlanRequest) ToURLValues() url.Values {
	values := url.Values{}

	// Helper function to add non-empty string values
	addString := func(key, value string) {
		if value != "" {
			values.Add(key, value)
		}
	}

	// Helper function to add int values (skip if 0)
	addInt := func(key string, value int) {
		if value != 0 {
			values.Add(key, strconv.Itoa(value))
		}
	}

	// Helper function to add float values (skip if 0)
	addFloat := func(key string, value float64) {
		if value != 0 {
			values.Add(key, strconv.FormatFloat(value, 'f', -1, 64))
		}
	}

	// Helper function to add bool pointer values
	addBool := func(key string, value *bool) {
		if value != nil {
			if *value {
				values.Add(key, "1")
			} else {
				values.Add(key, "0")
			}
		}
	}

	// Required fields
	addString("orig", fpr.Origin)
	addString("dest", fpr.Destination)
	addString("type", fpr.Aircraft)

	// Basic flight information
	addString("airline", fpr.Airline)
	addString("fltnum", fpr.FlightNumber)
	addString("date", fpr.Date)
	addInt("deph", fpr.DepartureHour)
	addInt("depm", fpr.DepartureMinute)
	addString("route", fpr.Route)
	addInt("steh", fpr.ScheduledHour)
	addInt("stem", fpr.ScheduledMinute)

	// Aircraft details
	addString("reg", fpr.Registration)
	addString("fin", fpr.FinNumber)
	addString("selcal", fpr.SELCAL)
	addString("callsign", fpr.ATCCallsign)

	// Crew and passengers
	addInt("pax", fpr.Passengers)
	addString("cpt", fpr.CaptainName)
	addString("dxname", fpr.DispatcherName)
	addString("pid", fpr.PilotID)

	// Alternates and routing
	addString("altn", fpr.Alternate)
	addString("fl", fpr.Altitude)
	addInt("altn_count", fpr.AltnCount)
	addString("altn_avoid", fpr.AltnAvoid)

	// Alternate airports 1-4
	addString("altn_1_id", fpr.Altn1ID)
	addString("altn_1_rwy", fpr.Altn1Runway)
	addString("altn_1_route", fpr.Altn1Route)
	addString("altn_2_id", fpr.Altn2ID)
	addString("altn_2_rwy", fpr.Altn2Runway)
	addString("altn_2_route", fpr.Altn2Route)
	addString("altn_3_id", fpr.Altn3ID)
	addString("altn_3_rwy", fpr.Altn3Runway)
	addString("altn_3_route", fpr.Altn3Route)
	addString("altn_4_id", fpr.Altn4ID)
	addString("altn_4_rwy", fpr.Altn4Runway)
	addString("altn_4_route", fpr.Altn4Route)

	// Fuel and weight
	addString("fuelfactor", fpr.FuelFactor)
	addFloat("manualzfw", fpr.ManualZFW)
	addString("addedfuel", fpr.AddedFuel)
	addString("addedfuel_units", fpr.AddedFuelUnits)
	addString("contpct", fpr.ContFuelPct)
	addInt("resvrule", fpr.ReserveFuel)
	addFloat("cargo", fpr.Cargo)

	// Taxi and runways
	addInt("taxiout", fpr.TaxiOut)
	addInt("taxiin", fpr.TaxiIn)
	addString("origrwy", fpr.OriginRunway)
	addString("destrwy", fpr.DestRunway)

	// Performance profiles
	addString("climb", fpr.ClimbProfile)
	addString("descent", fpr.DescentProfile)
	addString("cruise", fpr.CruiseProfile)
	addString("civalue", fpr.CostIndex)

	// Aircraft data
	if fpr.AircraftData != nil {
		addString("acdata", fpr.AircraftData.String())
	}

	// ETOPS
	addString("etopsrule", fpr.ETOPSRule)

	// Custom remarks and static ID
	addString("manualrmk", fpr.ManualRemarks)
	addString("static_id", fpr.StaticID)

	// OFP Options
	addString("planformat", fpr.PlanFormat)
	addString("units", string(fpr.Units))
	addBool("navlog", fpr.NavLog)
	addBool("etops", fpr.ETOPS)
	addBool("stepclimbs", fpr.StepClimbs)
	addBool("tlr", fpr.RunwayAnalysis)
	addBool("notams", fpr.NOTAMs)
	addBool("firnot", fpr.FIRNOTAMs)
	addString("maps", fpr.Maps)
	addBool("omit_sids", fpr.OmitSIDs)
	addBool("omit_stars", fpr.OmitSTARs)
	addString("find_sidstar", fpr.FindSIDSTAR)

	return values
}

// Validate checks if the flight plan request has all required fields
func (fpr *FlightPlanRequest) Validate() error {
	if fpr.Origin == "" {
		return ErrMissingOrigin
	}
	if fpr.Destination == "" {
		return ErrMissingDestination
	}
	if fpr.Aircraft == "" {
		return ErrMissingAircraft
	}
	return nil
}
