package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// FlightPlanResponse is a SimBrief OFP as returned by
// /api/xml.fetcher.php, either as XML (root element <OFP>) or as JSON v2
// (json=v2). The same struct decodes both with encoding/xml and
// encoding/json.
//
// Value conventions:
//   - Text values (identifiers, routes, runways, METARs) are string.
//   - Numeric values are Number: the raw SimBrief text ("0365", "05000",
//     ".57") with Int, Float and Bool helpers.
//   - Times are string because the formats differ: JSON v2 sends ISO 8601
//     timestamps and HH:MM:SS durations, XML sends Unix seconds and
//     durations in seconds.
//   - Lists (alternates, navlog, remarks) are slices in both formats.
type FlightPlanResponse struct {
	XMLName xml.Name `xml:"OFP" json:"-"`

	Fetch       FetchInfo       `xml:"fetch" json:"fetch"`
	Params      FlightParams    `xml:"params" json:"params"`
	General     GeneralInfo     `xml:"general" json:"general"`
	Origin      AirportInfo     `xml:"origin" json:"origin"`
	Destination AirportInfo     `xml:"destination" json:"destination"`
	Alternates  []AlternateInfo `xml:"alternate" json:"alternate"`
	NavLog      NavLog          `xml:"navlog" json:"navlog"`
	// AlternateNavLogs holds one navlog per alternate, in Alternates order.
	AlternateNavLogs []NavLog     `xml:"alternate_navlog" json:"alternate_navlog"`
	ATC              ATCInfo      `xml:"atc" json:"atc"`
	Aircraft         AircraftInfo `xml:"aircraft" json:"aircraft"`
	Fuel             FuelInfo     `xml:"fuel" json:"fuel"`
	Times            TimeInfo     `xml:"times" json:"times"`
	Weights          WeightInfo   `xml:"weights" json:"weights"`
	Weather          WeatherInfo  `xml:"weather" json:"weather"`
	Files            FilesInfo    `xml:"files" json:"files"`
	Links            LinksInfo    `xml:"links" json:"links"`
}

// Alternate returns the first alternate, or nil when the plan has none.
func (f *FlightPlanResponse) Alternate() *AlternateInfo {
	if len(f.Alternates) == 0 {
		return nil
	}
	return &f.Alternates[0]
}

// FetchInfo is the <fetch> block: the outcome of the fetch request.
// Status is "Success" or an error message.
type FetchInfo struct {
	UserID   string `xml:"userid" json:"userid"`
	StaticID string `xml:"static_id" json:"static_id"`
	Status   string `xml:"status" json:"status"`
	Time     Number `xml:"time" json:"time"`
}

// FlightParams contains parameters used to generate the flight plan.
type FlightParams struct {
	RequestID  string        `xml:"request_id" json:"request_id"`
	SequenceID string        `xml:"sequence_id" json:"sequence_id"`
	StaticID   StaticIDField `xml:"static_id" json:"static_id"`
	UserID     string        `xml:"user_id" json:"user_id"`
	// TimeGen is ISO 8601 in JSON, Unix seconds in XML.
	TimeGen   string `xml:"time_generated" json:"time_generated"`
	XMLFile   string `xml:"xml_file" json:"xml_file"`
	OFPLayout string `xml:"ofp_layout" json:"ofp_layout"`
	AIRAC     string `xml:"airac" json:"airac"`
	// Units is the weight unit of the plan as sent ("kgs" or "lbs").
	Units string `xml:"units" json:"units"`
}

// StaticIDField handles the static_id field which can be either a string or
// an empty object (older JSON replies).
type StaticIDField struct {
	Value string
}

// UnmarshalJSON implements custom JSON unmarshaling for StaticIDField
func (s *StaticIDField) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		s.Value = str
		return nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil {
		s.Value = "" // Empty object means no static ID
		return nil
	}
	return fmt.Errorf("static_id must be either string or object")
}

// MarshalJSON implements custom JSON marshaling for StaticIDField
func (s StaticIDField) MarshalJSON() ([]byte, error) {
	if s.Value == "" {
		return []byte("{}"), nil
	}
	return json.Marshal(s.Value)
}

// UnmarshalText lets encoding/xml read the element text.
func (s *StaticIDField) UnmarshalText(text []byte) error {
	s.Value = string(text)
	return nil
}

// MarshalText lets encoding/xml write the element text.
func (s StaticIDField) MarshalText() ([]byte, error) {
	return []byte(s.Value), nil
}

// String returns the string value of the static ID
func (s StaticIDField) String() string {
	return s.Value
}

// GeneralInfo is the <general> block.
type GeneralInfo struct {
	Release          string   `xml:"release" json:"release"`
	ICAOAirline      string   `xml:"icao_airline" json:"icao_airline"`
	FlightNumber     string   `xml:"flight_number" json:"flight_number"`
	DispatchRemarks  []string `xml:"dx_rmk" json:"dx_rmk"`
	SystemRemarks    []string `xml:"sys_rmk" json:"sys_rmk"`
	CruiseProfile    string   `xml:"cruise_profile" json:"cruise_profile"`
	ClimbProfile     string   `xml:"climb_profile" json:"climb_profile"`
	DescentProfile   string   `xml:"descent_profile" json:"descent_profile"`
	AlternateProfile string   `xml:"alternate_profile" json:"alternate_profile"`
	ReserveProfile   string   `xml:"reserve_profile" json:"reserve_profile"`
	CostIndex        Number   `xml:"costindex" json:"costindex"`
	ContRule         string   `xml:"cont_rule" json:"cont_rule"`
	// InitialAltitude is the initial cruise altitude in feet ("15000").
	InitialAltitude Number `xml:"initial_altitude" json:"initial_altitude"`
	StepClimbs      string `xml:"stepclimb_string" json:"stepclimb_string"`
	AvgTempDev      Number `xml:"avg_temp_dev" json:"avg_temp_dev"`
	AvgTropopause   Number `xml:"avg_tropopause" json:"avg_tropopause"`
	AvgWindComp     Number `xml:"avg_wind_comp" json:"avg_wind_comp"`
	AvgWindDir      Number `xml:"avg_wind_dir" json:"avg_wind_dir"`
	AvgWindSpd      Number `xml:"avg_wind_spd" json:"avg_wind_spd"`
	GCDistance      Number `xml:"gc_distance" json:"gc_distance"`
	RouteDistance   Number `xml:"route_distance" json:"route_distance"`
	AirDistance     Number `xml:"air_distance" json:"air_distance"`
	TotalBurn       Number `xml:"total_burn" json:"total_burn"`
	CruiseTAS       Number `xml:"cruise_tas" json:"cruise_tas"`
	CruiseMach      Number `xml:"cruise_mach" json:"cruise_mach"`
	Passengers      Number `xml:"passengers" json:"passengers"`
	Route           string `xml:"route" json:"route"`
	RouteIFPS       string `xml:"route_ifps" json:"route_ifps"`
	RouteNavigraph  string `xml:"route_navigraph" json:"route_navigraph"`
	SIDIdent        string `xml:"sid_ident" json:"sid_ident"`
	SIDTrans        string `xml:"sid_trans" json:"sid_trans"`
	STARIdent       string `xml:"star_ident" json:"star_ident"`
	STARTrans       string `xml:"star_trans" json:"star_trans"`
}

// ATCInfo is the <atc> block: the ICAO flight plan as filed.
type ATCInfo struct {
	// FlightPlanText is the full ICAO FPL message, "(FPL-...)".
	FlightPlanText string `xml:"flightplan_text" json:"flightplan_text"`
	Route          string `xml:"route" json:"route"`
	RouteIFPS      string `xml:"route_ifps" json:"route_ifps"`
	Callsign       string `xml:"callsign" json:"callsign"`
	FlightType     string `xml:"flight_type" json:"flight_type"`   // ICAO item 8b: S, N, G, M, X
	FlightRules    string `xml:"flight_rules" json:"flight_rules"` // ICAO item 8a: I, V, Y, Z
	InitialSpeed   Number `xml:"initial_spd" json:"initial_spd"`
	// InitialSpeedUnit is the ICAO speed unit: N (knots), K (km/h), M (Mach).
	InitialSpeedUnit string `xml:"initial_spd_unit" json:"initial_spd_unit"`
	// InitialAlt is in the unit of InitialAltUnit: "150" with "F" is FL150.
	InitialAlt Number `xml:"initial_alt" json:"initial_alt"`
	// InitialAltUnit is the ICAO level unit: F (flight level), A (altitude
	// in hundreds of feet), S or M (metric).
	InitialAltUnit string   `xml:"initial_alt_unit" json:"initial_alt_unit"`
	Section18      string   `xml:"section18" json:"section18"`
	FIROrig        string   `xml:"fir_orig" json:"fir_orig"`
	FIRDest        string   `xml:"fir_dest" json:"fir_dest"`
	FIRAltn        []string `xml:"fir_altn" json:"fir_altn"`
}

// AircraftInfo is the <aircraft> block.
type AircraftInfo struct {
	// ICAOCode is the icaocode element, e.g. "A319".
	ICAOCode string `xml:"icaocode" json:"icaocode"`
	// IATACode is the iatacode element, e.g. "319".
	IATACode string `xml:"iatacode" json:"iatacode"`
	BaseType string `xml:"base_type" json:"base_type"`
	ListType string `xml:"list_type" json:"list_type"`
	// ICAO is the icao_code element; it equals ICAOCode in observed plans.
	ICAO          string `xml:"icao_code" json:"icao_code"`
	IATA          string `xml:"iata_code" json:"iata_code"`
	Name          string `xml:"name" json:"name"`
	Engines       string `xml:"engines" json:"engines"`
	Registration  string `xml:"reg" json:"reg"`
	Fin           string `xml:"fin" json:"fin"`
	SELCAL        string `xml:"selcal" json:"selcal"`
	Equipment     string `xml:"equip" json:"equip"`
	EquipCategory string `xml:"equip_category" json:"equip_category"`
	EquipNav      string `xml:"equip_navigation" json:"equip_navigation"`
	EquipXpdr     string `xml:"equip_transponder" json:"equip_transponder"`
	FuelFactor    Number `xml:"fuelfactor" json:"fuelfactor"`
	MaxPassengers Number `xml:"max_passengers" json:"max_passengers"`
	InternalID    string `xml:"internal_id" json:"internal_id"`
	IsCustom      Number `xml:"is_custom" json:"is_custom"`
}

// AirportInfo is the <origin> or <destination> block (and the common part
// of an <alternate>).
type AirportInfo struct {
	ICAO       string `xml:"icao_code" json:"icao_code"`
	IATA       string `xml:"iata_code" json:"iata_code"`
	FAA        string `xml:"faa_code" json:"faa_code"`
	ICAORegion string `xml:"icao_region" json:"icao_region"`
	Elevation  Number `xml:"elevation" json:"elevation"` // feet
	Latitude   Number `xml:"pos_lat" json:"pos_lat"`
	Longitude  Number `xml:"pos_long" json:"pos_long"`
	Name       string `xml:"name" json:"name"`
	TimeZone   Number `xml:"timezone" json:"timezone"` // UTC offset, hours
	// Runway is the planned runway ("06").
	Runway          string `xml:"plan_rwy" json:"plan_rwy"`
	TransAlt        Number `xml:"trans_alt" json:"trans_alt"`     // feet
	TransLevel      Number `xml:"trans_level" json:"trans_level"` // feet
	METAR           string `xml:"metar" json:"metar"`
	METARTime       string `xml:"metar_time" json:"metar_time"`
	METARCategory   string `xml:"metar_category" json:"metar_category"`
	METARVisibility Number `xml:"metar_visibility" json:"metar_visibility"`
	METARCeiling    Number `xml:"metar_ceiling" json:"metar_ceiling"`
	TAF             string `xml:"taf" json:"taf"`
	TAFTime         string `xml:"taf_time" json:"taf_time"`
}

// AlternateInfo is one <alternate> block: the airport plus the diversion
// leg from the destination.
type AlternateInfo struct {
	AirportInfo
	CruiseAltitude Number `xml:"cruise_altitude" json:"cruise_altitude"`
	Distance       Number `xml:"distance" json:"distance"`
	GCDistance     Number `xml:"gc_distance" json:"gc_distance"`
	AirDistance    Number `xml:"air_distance" json:"air_distance"`
	TrackTrue      Number `xml:"track_true" json:"track_true"`
	TrackMag       Number `xml:"track_mag" json:"track_mag"`
	TAS            Number `xml:"tas" json:"tas"`
	GS             Number `xml:"gs" json:"gs"`
	// AvgWindComp is signed text such as "M010" (10 kt headwind).
	AvgWindComp string `xml:"avg_wind_comp" json:"avg_wind_comp"`
	AvgWindDir  Number `xml:"avg_wind_dir" json:"avg_wind_dir"`
	AvgWindSpd  Number `xml:"avg_wind_spd" json:"avg_wind_spd"`
	// ETE is HH:MM:SS in JSON, seconds in XML.
	ETE       string `xml:"ete" json:"ete"`
	Burn      Number `xml:"burn" json:"burn"`
	Route     string `xml:"route" json:"route"`
	RouteIFPS string `xml:"route_ifps" json:"route_ifps"`
}

// NavLog is the list of navlog fixes. JSON v2 sends it as a plain array,
// XML as <navlog><fix>...</fix></navlog>; both decode into NavLog.
type NavLog []NavLogFix

// UnmarshalXML reads the <fix> children of a navlog element.
func (n *NavLog) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var wrapper struct {
		Fix []NavLogFix `xml:"fix"`
	}
	if err := d.DecodeElement(&wrapper, &start); err != nil {
		return err
	}
	*n = wrapper.Fix
	return nil
}

// MarshalXML writes the fixes as <fix> children.
func (n NavLog) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	wrapper := struct {
		Fix []NavLogFix `xml:"fix"`
	}{Fix: n}
	return e.EncodeElement(wrapper, start)
}

// NavLogFix is one fix of the navlog.
type NavLogFix struct {
	Ident      string `xml:"ident" json:"ident"`
	Name       string `xml:"name" json:"name"`
	Type       string `xml:"type" json:"type"` // wpt, vor, ndb, apt, ltlg (TOC/TOD), ...
	ICAORegion string `xml:"icao_region" json:"icao_region"`
	RegionCode string `xml:"region_code" json:"region_code"`
	Frequency  Number `xml:"frequency" json:"frequency"`
	Latitude   Number `xml:"pos_lat" json:"pos_lat"`
	Longitude  Number `xml:"pos_long" json:"pos_long"`
	Stage      string `xml:"stage" json:"stage"` // CLB, CRZ, DSC
	ViaAirway  string `xml:"via_airway" json:"via_airway"`
	// IsSIDSTAR is "1" when the fix belongs to a SID or STAR; use Bool().
	IsSIDSTAR     Number `xml:"is_sid_star" json:"is_sid_star"`
	Distance      Number `xml:"distance" json:"distance"`
	TrackTrue     Number `xml:"track_true" json:"track_true"`
	TrackMag      Number `xml:"track_mag" json:"track_mag"`
	HeadingTrue   Number `xml:"heading_true" json:"heading_true"`
	HeadingMag    Number `xml:"heading_mag" json:"heading_mag"`
	AltitudeFeet  Number `xml:"altitude_feet" json:"altitude_feet"`
	IndAirspeed   Number `xml:"ind_airspeed" json:"ind_airspeed"`
	TrueAirspeed  Number `xml:"true_airspeed" json:"true_airspeed"`
	Mach          Number `xml:"mach" json:"mach"`
	WindComponent Number `xml:"wind_component" json:"wind_component"`
	Groundspeed   Number `xml:"groundspeed" json:"groundspeed"`
	// TimeLeg and TimeTotal are HH:MM:SS in JSON, seconds in XML.
	TimeLeg         string `xml:"time_leg" json:"time_leg"`
	TimeTotal       string `xml:"time_total" json:"time_total"`
	FuelFlow        Number `xml:"fuel_flow" json:"fuel_flow"`
	FuelLeg         Number `xml:"fuel_leg" json:"fuel_leg"`
	FuelTotalUsed   Number `xml:"fuel_totalused" json:"fuel_totalused"`
	FuelMinOnboard  Number `xml:"fuel_min_onboard" json:"fuel_min_onboard"`
	FuelPlanOnboard Number `xml:"fuel_plan_onboard" json:"fuel_plan_onboard"`
	OAT             Number `xml:"oat" json:"oat"`
	OATISADev       Number `xml:"oat_isa_dev" json:"oat_isa_dev"`
	WindDir         Number `xml:"wind_dir" json:"wind_dir"`
	WindSpd         Number `xml:"wind_spd" json:"wind_spd"`
	Shear           Number `xml:"shear" json:"shear"`
	TropopauseFeet  Number `xml:"tropopause_feet" json:"tropopause_feet"`
	GroundHeight    Number `xml:"ground_height" json:"ground_height"`
	FIR             string `xml:"fir" json:"fir"`
	MORA            Number `xml:"mora" json:"mora"`
}

// FuelInfo is the <fuel> block, in Params.Units.
type FuelInfo struct {
	Taxi          Number `xml:"taxi" json:"taxi"`
	EnrouteBurn   Number `xml:"enroute_burn" json:"enroute_burn"` // trip fuel
	Contingency   Number `xml:"contingency" json:"contingency"`
	AlternateBurn Number `xml:"alternate_burn" json:"alternate_burn"`
	Reserve       Number `xml:"reserve" json:"reserve"`
	ETOPS         Number `xml:"etops" json:"etops"`
	Extra         Number `xml:"extra" json:"extra"`
	ExtraRequired Number `xml:"extra_required" json:"extra_required"`
	ExtraOptional Number `xml:"extra_optional" json:"extra_optional"`
	MinTakeoff    Number `xml:"min_takeoff" json:"min_takeoff"`
	PlanTakeoff   Number `xml:"plan_takeoff" json:"plan_takeoff"`
	PlanRamp      Number `xml:"plan_ramp" json:"plan_ramp"` // block fuel
	PlanLanding   Number `xml:"plan_landing" json:"plan_landing"`
	AvgFuelFlow   Number `xml:"avg_fuel_flow" json:"avg_fuel_flow"`
	MaxTanks      Number `xml:"max_tanks" json:"max_tanks"`
}

// WeightInfo is the <weights> block, in Params.Units.
type WeightInfo struct {
	OEW            Number `xml:"oew" json:"oew"`
	PaxCount       Number `xml:"pax_count" json:"pax_count"`
	BagCount       Number `xml:"bag_count" json:"bag_count"`
	PaxCountActual Number `xml:"pax_count_actual" json:"pax_count_actual"`
	BagCountActual Number `xml:"bag_count_actual" json:"bag_count_actual"`
	PaxWeight      Number `xml:"pax_weight" json:"pax_weight"`
	BagWeight      Number `xml:"bag_weight" json:"bag_weight"`
	FreightAdded   Number `xml:"freight_added" json:"freight_added"`
	Cargo          Number `xml:"cargo" json:"cargo"`
	Payload        Number `xml:"payload" json:"payload"`
	EstZFW         Number `xml:"est_zfw" json:"est_zfw"`
	MaxZFW         Number `xml:"max_zfw" json:"max_zfw"`
	EstTOW         Number `xml:"est_tow" json:"est_tow"`
	MaxTOW         Number `xml:"max_tow" json:"max_tow"`
	MaxTOWStruct   Number `xml:"max_tow_struct" json:"max_tow_struct"`
	TOWLimitCode   string `xml:"tow_limit_code" json:"tow_limit_code"`
	EstLDW         Number `xml:"est_ldw" json:"est_ldw"`
	MaxLDW         Number `xml:"max_ldw" json:"max_ldw"`
	EstRamp        Number `xml:"est_ramp" json:"est_ramp"`
}

// TimeInfo is the <times> block. Points in time are ISO 8601 in JSON and
// Unix seconds in XML; durations are HH:MM:SS in JSON and seconds in XML.
type TimeInfo struct {
	EstTimeEnroute   string `xml:"est_time_enroute" json:"est_time_enroute"`
	SchedTimeEnroute string `xml:"sched_time_enroute" json:"sched_time_enroute"`
	SchedOut         string `xml:"sched_out" json:"sched_out"`
	SchedOff         string `xml:"sched_off" json:"sched_off"`
	SchedOn          string `xml:"sched_on" json:"sched_on"`
	SchedIn          string `xml:"sched_in" json:"sched_in"`
	SchedBlock       string `xml:"sched_block" json:"sched_block"`
	EstOut           string `xml:"est_out" json:"est_out"`
	EstOff           string `xml:"est_off" json:"est_off"`
	EstOn            string `xml:"est_on" json:"est_on"`
	EstIn            string `xml:"est_in" json:"est_in"`
	EstBlock         string `xml:"est_block" json:"est_block"`
	OrigTimezone     string `xml:"orig_timezone" json:"orig_timezone"`
	DestTimezone     string `xml:"dest_timezone" json:"dest_timezone"`
	TaxiOut          string `xml:"taxi_out" json:"taxi_out"`
	TaxiIn           string `xml:"taxi_in" json:"taxi_in"`
	ReserveTime      string `xml:"reserve_time" json:"reserve_time"`
	Endurance        string `xml:"endurance" json:"endurance"`
	ContFuelTime     string `xml:"contfuel_time" json:"contfuel_time"`
	ETOPSFuelTime    string `xml:"etopsfuel_time" json:"etopsfuel_time"`
	ExtraFuelTime    string `xml:"extrafuel_time" json:"extrafuel_time"`
}

// WeatherInfo is the <weather> block. Alternate reports are one per
// alternate, in Alternates order.
type WeatherInfo struct {
	OrigMETAR string   `xml:"orig_metar" json:"orig_metar"`
	OrigTAF   string   `xml:"orig_taf" json:"orig_taf"`
	DestMETAR string   `xml:"dest_metar" json:"dest_metar"`
	DestTAF   string   `xml:"dest_taf" json:"dest_taf"`
	AltnMETAR []string `xml:"altn_metar" json:"altn_metar"`
	AltnTAF   []string `xml:"altn_taf" json:"altn_taf"`
}

// FileLink is one generated file; Link is relative to FilesInfo.Directory.
type FileLink struct {
	Name string `xml:"name" json:"name"`
	Link string `xml:"link" json:"link"`
}

// FilesInfo is the <files> block: the PDF OFP and the exported plan files.
type FilesInfo struct {
	Directory string     `xml:"directory" json:"directory"`
	PDF       FileLink   `xml:"pdf" json:"pdf"`
	Files     []FileLink `xml:"file" json:"file"`
}

// LinksInfo is the <links> block.
type LinksInfo struct {
	SkyVector string `xml:"skyvector" json:"skyvector"`
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
