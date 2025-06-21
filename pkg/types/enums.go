package types

// Units represents the weight/fuel units
type Units string

const (
	UnitsLBS Units = "LBS"
	UnitsKGS Units = "KGS"
)

// PlanFormat represents the OFP layout format
type PlanFormat string

const (
	PlanFormatLIDO    PlanFormat = "LIDO"
	PlanFormatDefault PlanFormat = ""
)

// FlightRules represents IFR/VFR rules
type FlightRules string

const (
	FlightRulesIFR FlightRules = "I"
	FlightRulesVFR FlightRules = "V"
)

// FlightType represents scheduled/charter flight type
type FlightType string

const (
	FlightTypeScheduled FlightType = "S"
	FlightTypeCharter   FlightType = "X"
)

// MapDetail represents flight map detail level
type MapDetail string

const (
	MapDetailDetailed MapDetail = "detail"
	MapDetailSimple   MapDetail = "simple"
	MapDetailNone     MapDetail = "none"
)

// FuelUnits represents fuel quantity units
type FuelUnits string

const (
	FuelUnitsWeight FuelUnits = "wgt"
	FuelUnitsTime   FuelUnits = "min"
)

// SIDSTARPreference represents SID/STAR routing preference
type SIDSTARPreference string

const (
	SIDSTARPreferRNAV SIDSTARPreference = "R"
	SIDSTARPreferConv SIDSTARPreference = "C"
)

// AircraftCategory represents ICAO aircraft weight category
type AircraftCategory string

const (
	AircraftCategoryLight  AircraftCategory = "L"
	AircraftCategoryMedium AircraftCategory = "M"
	AircraftCategoryHeavy  AircraftCategory = "H"
	AircraftCategorySuper  AircraftCategory = "J"
)

// PerformanceCategory represents ICAO performance category
type PerformanceCategory string

const (
	PerformanceCategoryA PerformanceCategory = "A"
	PerformanceCategoryB PerformanceCategory = "B"
	PerformanceCategoryC PerformanceCategory = "C"
	PerformanceCategoryD PerformanceCategory = "D"
	PerformanceCategoryE PerformanceCategory = "E"
)
