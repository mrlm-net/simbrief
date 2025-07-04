package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

// FlightParams contains parameters used to generate the flight plan
type FlightParams struct {
	RequestID string        `xml:"request_id" json:"request_id"`
	UserID    string        `xml:"user_id" json:"user_id"`
	TimeGen   string        `xml:"time_generated" json:"time_generated"`
	StaticID  StaticIDField `xml:"static_id" json:"static_id"`
	XMLFile   string        `xml:"xml_file" json:"xml_file"`
	OFPLayout string        `xml:"ofp_layout" json:"ofp_layout"`
	Units     Units         `xml:"units" json:"units"`
}

// FlightPlanResponse represents the complete response from SimBrief API
type FlightPlanResponse struct {
	XMLName xml.Name `xml:"SimBrief" json:"-"`

	// Basic flight information
	Params      FlightParams  `xml:"params" json:"params"`
	General     GeneralInfo   `xml:"general" json:"general"`
	Aircraft    AircraftInfo  `xml:"aircraft" json:"aircraft"`
	Origin      AirportInfo   `xml:"origin" json:"origin"`
	Destination AirportInfo   `xml:"destination" json:"destination"`
	Alternate   AlternateInfo `xml:"alternate" json:"alternate"`

	// Flight planning data
	Fuel    FuelInfo    `xml:"fuel" json:"fuel"`
	Weights WeightInfo  `xml:"weights" json:"weights"`
	Times   TimeInfo    `xml:"times" json:"times"`
	Weather WeatherInfo `xml:"weather" json:"weather"`
	NavLog  interface{} `xml:"navlog>fix" json:"navlog"`

	// Generated files and links
	Files FilesInfo `xml:"files" json:"files"`
	Links LinksInfo `xml:"links" json:"links"`

	// Raw response for advanced usage
	Raw map[string]interface{} `xml:"-" json:"raw,omitempty"`
}

// StaticIDField handles the static_id field which can be either a string or an empty object
type StaticIDField struct {
	Value string
}

// UnmarshalJSON implements custom JSON unmarshaling for StaticIDField
func (s *StaticIDField) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		s.Value = str
		return nil
	}

	// If that fails, try as object (empty object means no static ID)
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil {
		s.Value = "" // Empty object means no static ID
		return nil
	}

	// If both fail, return error
	return fmt.Errorf("static_id must be either string or object")
}

// MarshalJSON implements custom JSON marshaling for StaticIDField
func (s StaticIDField) MarshalJSON() ([]byte, error) {
	if s.Value == "" {
		return []byte("{}"), nil
	}
	return json.Marshal(s.Value)
}

// String returns the string value of the static ID
func (s StaticIDField) String() string {
	return s.Value
}

// GeneralInfo contains general flight information
type GeneralInfo struct {
	ICAO           string    `xml:"icao_airline" json:"icao_airline"`
	IATA           string    `xml:"iata_airline" json:"iata_airline"`
	FlightNumber   string    `xml:"flight_number" json:"flight_number"`
	CallSign       string    `xml:"callsign" json:"callsign"`
	CostIndex      string    `xml:"costindex" json:"costindex"`
	CruiseAltitude string    `xml:"initial_altitude" json:"initial_altitude"`
	StepClimbs     string    `xml:"stepclimb_string" json:"stepclimb_string"`
	Route          string    `xml:"route" json:"route"`
	RouteNAVID     string    `xml:"route_navids" json:"route_navids"`
	Distance       string    `xml:"air_distance" json:"air_distance"`
	Units          Units     `xml:"units" json:"units"`
	CreatedTime    time.Time `xml:"plan_html" json:"plan_html"`
}

// AircraftInfo contains aircraft-specific information
type AircraftInfo struct {
	ICAO         string      `xml:"icaocode" json:"icaocode"`
	Name         string      `xml:"name" json:"name"`
	Engine       string      `xml:"engine" json:"engine"`
	Registration string      `xml:"reg" json:"reg"`
	Fin          string      `xml:"fin" json:"fin"`
	SELCAL       interface{} `xml:"selcal" json:"selcal"`
	MaxPax       int         `xml:"maxpax" json:"maxpax"`
	OEW          float64     `xml:"oew" json:"oew"`         // Operating Empty Weight
	MZFW         float64     `xml:"mzfw" json:"mzfw"`       // Max Zero Fuel Weight
	MTOW         float64     `xml:"mtow" json:"mtow"`       // Max Takeoff Weight
	MLW          float64     `xml:"mlw" json:"mlw"`         // Max Landing Weight
	MaxFuel      float64     `xml:"maxfuel" json:"maxfuel"` // Max Fuel Capacity
}

// AirportInfo contains airport information
type AirportInfo struct {
	ICAO        string `xml:"icao_code" json:"icao_code"`
	IATA        string `xml:"iata_code" json:"iata_code"`
	Name        string `xml:"name" json:"name"`
	City        string `xml:"city" json:"city"`
	Country     string `xml:"country" json:"country"`
	CountryCode string `xml:"country_code" json:"country_code"`
	Elevation   string `xml:"elevation" json:"elevation"`
	Latitude    string `xml:"pos_lat" json:"pos_lat"`
	Longitude   string `xml:"pos_long" json:"pos_long"`
	Runway      string `xml:"plan_rwy" json:"plan_rwy"`
	TimeZone    string `xml:"timezone" json:"timezone"`
	UTCOffset   string `xml:"utc_offset" json:"utc_offset"`
}

// AlternateInfo contains alternate airport information
type AlternateInfo struct {
	ICAO         string `xml:"icao_code" json:"icao_code"`
	IATA         string `xml:"iata_code" json:"iata_code"`
	Name         string `xml:"name" json:"name"`
	Distance     string `xml:"distance" json:"distance"`
	Bearing      string `xml:"bearing" json:"bearing"`
	FuelRequired string `xml:"burn" json:"burn"`
}

// FuelInfo contains fuel planning information
type FuelInfo struct {
	Plan        string `xml:"plan_ramp" json:"plan_ramp"`           // Total planned fuel
	Taxi        string `xml:"taxi" json:"taxi"`                     // Taxi fuel
	Trip        string `xml:"enroute_burn" json:"enroute_burn"`     // Trip fuel
	Contingency string `xml:"contingency" json:"contingency"`       // Contingency fuel
	Alternate   string `xml:"alternate_burn" json:"alternate_burn"` // Alternate fuel
	Reserve     string `xml:"reserve" json:"reserve"`               // Reserve fuel
	Extra       string `xml:"extra" json:"extra"`                   // Extra fuel
	MinTakeoff  string `xml:"min_takeoff" json:"min_takeoff"`       // Minimum takeoff fuel
	PlanLanding string `xml:"plan_landing" json:"plan_landing"`     // Planned landing fuel
	AvgFuelFlow string `xml:"avg_fuel_flow" json:"avg_fuel_flow"`   // Average fuel flow
}

// WeightInfo contains weight and balance information
type WeightInfo struct {
	OEW       string `xml:"oew" json:"oew"`               // Operating Empty Weight
	Payload   string `xml:"payload" json:"payload"`       // Total payload
	PaxWeight string `xml:"pax_weight" json:"pax_weight"` // Passenger weight
	BagWeight string `xml:"bag_weight" json:"bag_weight"` // Baggage weight
	Cargo     string `xml:"cargo" json:"cargo"`           // Cargo weight
	ZFW       string `xml:"est_zfw" json:"est_zfw"`       // Zero Fuel Weight
	TakeoffWt string `xml:"est_tow" json:"est_tow"`       // Takeoff Weight
	LandingWt string `xml:"est_ldw" json:"est_ldw"`       // Landing Weight
	PaxCount  string `xml:"pax_count" json:"pax_count"`   // Number of passengers
}

// TimeInfo contains flight timing information
type TimeInfo struct {
	Departure  string `xml:"est_out" json:"est_out"`                       // Estimated departure
	Takeoff    string `xml:"est_off" json:"est_off"`                       // Estimated takeoff
	Landing    string `xml:"est_on" json:"est_on"`                         // Estimated landing
	Arrival    string `xml:"est_in" json:"est_in"`                         // Estimated arrival
	FlightTime string `xml:"est_time_enroute" json:"est_time_enroute"`     // Flight time
	BlockTime  string `xml:"sched_time_enroute" json:"sched_time_enroute"` // Block time
	TaxiOut    string `xml:"taxi_out" json:"taxi_out"`                     // Taxi out time (minutes)
	TaxiIn     string `xml:"taxi_in" json:"taxi_in"`                       // Taxi in time (minutes)
}

// WeatherInfo contains weather information
type WeatherInfo struct {
	OriginMETAR string `xml:"orig_metar" json:"orig_metar"`
	DestMETAR   string `xml:"dest_metar" json:"dest_metar"`
	AvgWindDir  string `xml:"avg_wind_dir" json:"avg_wind_dir"`
	AvgWindSpd  string `xml:"avg_wind_spd" json:"avg_wind_spd"`
	AvgTemp     string `xml:"avg_temp" json:"avg_temp"`
	MaxTemp     string `xml:"max_temp" json:"max_temp"`
	MinTemp     string `xml:"min_temp" json:"min_temp"`
}

// NavLogFix represents a single navigation fix in the flight plan
type NavLogFix struct {
	Ident       string  `xml:"ident" json:"ident"`
	Name        string  `xml:"name" json:"name"`
	Type        string  `xml:"type" json:"type"`
	Frequency   string  `xml:"frequency" json:"frequency"`
	Latitude    float64 `xml:"pos_lat" json:"pos_lat"`
	Longitude   float64 `xml:"pos_long" json:"pos_long"`
	Route       string  `xml:"via_airway" json:"via_airway"`
	Distance    float64 `xml:"distance_nm" json:"distance_nm"`
	Track       float64 `xml:"track_true" json:"track_true"`
	TrackMag    float64 `xml:"track_mag" json:"track_mag"`
	Altitude    int     `xml:"altitude_feet" json:"altitude_feet"`
	Wind        string  `xml:"wind" json:"wind"`
	Temperature int     `xml:"oat" json:"oat"`
	FuelFlow    float64 `xml:"fuel_flow" json:"fuel_flow"`
	FuelRemain  float64 `xml:"fuel_totalused" json:"fuel_totalused"`
	ETE         string  `xml:"time_leg" json:"time_leg"`
	ETA         string  `xml:"eta" json:"eta"`
}

// FilesInfo contains links to generated files
type FilesInfo struct {
	Directory string      `xml:"directory" json:"directory"`
	PDFLink   interface{} `xml:"pdf" json:"pdf"`
	XMLLink   interface{} `xml:"xml" json:"xml"`
	JSONLink  interface{} `xml:"json" json:"json"`
	KMLLink   interface{} `xml:"kml" json:"kml"`
	PLNLink   interface{} `xml:"pln" json:"pln"`
	FMSLink   interface{} `xml:"fms" json:"fms"`
	XPFMSLink interface{} `xml:"xpfms" json:"xpfms"`
}

// LinksInfo contains various SimBrief links
type LinksInfo struct {
	SkyVectorLink  string `xml:"skyvector" json:"skyvector"`
	EditFlightLink string `xml:"edit" json:"edit"`
	ViewFlightLink string `xml:"view" json:"view"`
}

// APIError represents an error response from the API
type APIError struct {
	XMLName xml.Name `xml:"error" json:"-"`
	Message string   `xml:",chardata" json:"message"`
	Code    int      `json:"code,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

// SupportedOptions represents the response from the inputs.list endpoint
// Based on official SimBrief API documentation at http://www.simbrief.com/api/inputs.list.json
type SupportedOptions struct {
	Aircraft    map[string]AircraftOption `json:"aircraft"`
	Layouts     map[string]LayoutOption   `json:"layouts"`
	LastUpdated string                    `json:"last_updated"`
	ProcessTime float64                   `json:"process_time"`
}

// AircraftOption represents an available aircraft type with detailed information
type AircraftOption struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Accuracy      string  `json:"accuracy"`
	ChartData     bool    `json:"chart_data"`
	CostIndexData bool    `json:"costindex_data"`
	TLRData       bool    `json:"tlr_data"`
	LastUpdated   string  `json:"last_updated"`
	PopularityPct float64 `json:"popularity_pct"`
}

// LayoutOption represents an available plan format/layout
type LayoutOption struct {
	ID            string  `json:"id"`
	NameShort     string  `json:"name_short"`
	NameLong      string  `json:"name_long"`
	PopularityPct float64 `json:"popularity_pct"`
	LastUpdated   string  `json:"last_updated"`
}
