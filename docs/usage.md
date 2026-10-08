# Usage

This guide covers the whole public API of `github.com/mrlm-net/simbrief`:
the client in `pkg/client` and the request/response types in `pkg/types`.

Code blocks are function bodies. Unless a block creates them itself, `c` is
a `*client.Client`, `ofp` is a `*types.FlightPlanResponse`, and the
surrounding function returns `error`. Imports:

```go
import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mrlm-net/simbrief/pkg/client"
	"github.com/mrlm-net/simbrief/pkg/types"
)
```

No API key is needed: fetching a plan uses the public
`/api/xml.fetcher.php` endpoint with a SimBrief user ID or username.

## Client

```go
c := client.NewClient() // https://www.simbrief.com, 30 s timeout
c.SetTimeout(10 * time.Second)
c.SetUserAgent("my-app/1.0")

// Custom base URL and http.Client, e.g. for an httptest server:
c = client.NewClientWithConfig("http://127.0.0.1:8080", &http.Client{Timeout: 5 * time.Second})
fmt.Println(c.BaseURL)
```

`NewClient()` takes no arguments. `NewClientWithConfig` falls back to the
defaults (`client.DefaultBaseURL`, `client.DefaultTimeout`) for an empty URL
or a nil `*http.Client`. `SetUserAgent` wraps the client's transport, so call
it after any custom transport is in place.

## Fetching a plan

`GetFlightPlan` fetches the latest OFP of a user (or a specific one by static
ID) and decodes it into `types.FlightPlanResponse`.

```go
// Latest OFP of a user, as JSON v2 (JSON: false fetches XML instead).
ofp, err := c.GetFlightPlan(&types.FetchRequest{Username: "pilot", JSON: true})
if err != nil {
	return err
}
fmt.Println(ofp.Origin.ICAO, "->", ofp.Destination.ICAO)
```

`types.FetchRequest` fields:

| Field | Query parameter | Meaning |
|-------|-----------------|---------|
| `UserID` | `userid` | SimBrief user ID (numeric, as a string) |
| `Username` | `username` | SimBrief username |
| `StaticID` | `static_id` | the static ID the plan was generated with |
| `JSON` | `json=v2` | `true` for JSON v2, `false` for XML |

JSON v2 and XML decode into the same struct; the only differences are the
time formats (see [Times](#times)).

Shortcuts, all fetching JSON v2:

```go
_, _ = c.GetFlightPlanByUserID("123456")
_, _ = c.GetFlightPlanByUsername("pilot")
_, _ = c.GetFlightPlanByStaticID("123456", "MYAPP_42") // a plan generated with StaticID
```

The raw XML, undecoded (`GetFlightPlanXML` sets `req.JSON` to `false`):

```go
raw, err := c.GetFlightPlanXML(&types.FetchRequest{UserID: "123456"})
if err != nil {
	return err
}
fmt.Println(len(raw), "bytes of <OFP> XML")
```

## The OFP

`types.FlightPlanResponse` mirrors the OFP blocks:

| Field | Type | Content |
|-------|------|---------|
| `Fetch` | `FetchInfo` | fetch outcome: `Status` ("Success" or an error), user ID, static ID |
| `Params` | `FlightParams` | request/sequence IDs, `TimeGen`, `AIRAC`, `OFPLayout`, `Units` ("kgs"/"lbs") |
| `General` | `GeneralInfo` | flight number, route, SID/STAR, cost index, initial altitude, distances, burn, average wind |
| `Origin`, `Destination` | `AirportInfo` | ICAO/IATA, position, elevation, planned `Runway`, `TransAlt`/`TransLevel`, METAR/TAF |
| `Alternates` | `[]AlternateInfo` | each alternate airport plus its diversion leg (distance, track, burn, route) |
| `NavLog` | `NavLog` | the route fixes |
| `AlternateNavLogs` | `[]NavLog` | one navlog per alternate, in `Alternates` order |
| `ATC` | `ATCInfo` | filed ICAO flight plan: callsign, rules, type, initial speed/level and units, route, FPL text |
| `Aircraft` | `AircraftInfo` | ICAO/IATA type, name, engines, registration, SELCAL, equipment |
| `Fuel` | `FuelInfo` | taxi, trip (`EnrouteBurn`), contingency, alternate, reserve, extra, `PlanRamp` (block) ... in `Params.Units` |
| `Weights` | `WeightInfo` | OEW, pax/bags, payload, ZFW/TOW/LDW estimates and limits, in `Params.Units` |
| `Times` | `TimeInfo` | scheduled/estimated out/off/on/in, block and enroute times |
| `Weather` | `WeatherInfo` | origin/destination METAR and TAF, one per alternate |
| `TLR` | `TLR` | runway analysis (take-off and landing); zero when the plan has none |
| `Files`, `Links` | `FilesInfo`, `LinksInfo` | PDF and exported plan files (relative to `Files.Directory`), SkyVector link |

```go
fmt.Println(ofp.ATC.Callsign)                                 // "CEF007"
fmt.Println(ofp.ATC.InitialAltUnit, ofp.ATC.InitialAlt)       // "F" "150"
fmt.Println(ofp.General.Route, ofp.General.STARIdent)         // route, "BEKV1Q"
fmt.Println(ofp.General.InitialAltitude.Int())                // 15000
fmt.Println(ofp.Origin.ICAO, ofp.Origin.Runway)               // "LKPR" "06"
fmt.Println(ofp.Origin.TransAlt.Int())                        // 5000
fmt.Println(ofp.Aircraft.ICAOCode, ofp.Aircraft.Registration) // "A319" ...
fmt.Println(ofp.Fuel.PlanRamp.Int(), ofp.Params.Units)        // block fuel, "kgs"
fmt.Println(ofp.Weights.EstTOW.Int(), ofp.Times.EstBlock)     // weight, duration text
```

### Number

SimBrief sends numbers as text, often zero-padded ("0365", "05000"), without
a leading zero (".57") or empty. Numeric fields are `types.Number`: the raw
text, decoded from a JSON string, JSON number, boolean, null or XML text.

```go
n := types.Number("0365")
fmt.Println(n.Int(), n.Float(), n.IsSet(), n.String()) // 365 365 true "0365"
fmt.Println(types.Number(".57").Float())               // 0.57
fmt.Println(types.Number("").Int())                    // 0
fmt.Println(types.Number("1").Bool())                  // true
```

`Int` truncates fractions; `Int`, `Float` and `Bool` return zero/false for an
empty or non-numeric value, so use `IsSet` when "missing" and "0" differ.

### NavLog

`NavLog` is a `[]NavLogFix`. JSON v2 sends a plain array, XML sends
`<navlog><fix>...</fix></navlog>`; both decode the same way. A fix carries
its ident, name, `Type` (wpt, vor, ndb, apt, ltlg for TOC/TOD ...), `Stage`
(CLB, CRZ, DSC), `ViaAirway`, position, altitude, speeds, wind, leg and total
time and fuel, FIR and MORA.

```go
for _, fix := range ofp.NavLog {
	if fix.IsSIDSTAR.Bool() {
		continue
	}
	fmt.Println(fix.Ident, fix.Type, fix.ViaAirway, fix.Stage,
		fix.AltitudeFeet.Int(), fix.Latitude.Float(), fix.Longitude.Float())
}
```

### Alternates

`Alternate()` returns the first alternate or nil. `AlternateNavLogs[i]`
belongs to `Alternates[i]`.

```go
if altn := ofp.Alternate(); altn != nil {
	fmt.Println(altn.ICAO, altn.Runway, altn.Distance.Int(), altn.Burn.Int())
}
for i, altn := range ofp.Alternates {
	if i < len(ofp.AlternateNavLogs) {
		fmt.Println(altn.ICAO, len(ofp.AlternateNavLogs[i]), "fixes")
	}
}
```

### Times

Time values are strings because the formats differ: JSON v2 sends ISO 8601
timestamps and `HH:MM:SS` durations, XML sends Unix seconds and durations
in seconds. This applies to `Times`, `Params.TimeGen`, `AlternateInfo.ETE`
and the navlog's `TimeLeg`/`TimeTotal`.

### Runway analysis (TLR)

`ofp.TLR` holds SimBrief's take-off and landing report. It is filled only
when the plan was generated with runway analysis on
(`FlightPlanRequest.RunwayAnalysis`) for an aircraft that supports it
(`AircraftOption.TLRData`); otherwise it is the zero value and the helpers
return `false`. Weights are in `Params.Units`, lengths in feet, speeds in
knots, temperatures in degrees Celsius.

```go
tlr := ofp.TLR
if rwy, ok := tlr.TakeoffRunway(ofp.Origin.Runway); ok { // "06"; "6" works too
	fmt.Println(rwy.SpeedsV1.Int(), rwy.SpeedsVR.Int(), rwy.SpeedsV2.Int())   // 129 129 134
	fmt.Println(rwy.FlapSetting, rwy.ThrustSetting, rwy.FlexTemperature.Int()) // "1" "FLEX" 65
}
if vref, ok := tlr.LandingVref(false); ok { // false = dry runway, true = wet
	fmt.Println(vref) // 125
}
fmt.Println(tlr.Landing.DistanceDry.FactoredDistance.Int()) // 3546 (feet)
```

`TakeoffRunway(ident)` tries the exact identifier first, then ignores case,
a "RW"/"RWY" prefix and leading zeros ("6l" finds "06L"); a number without a
side letter finds a lettered runway only when exactly one matches.
`LandingVref(wet)` returns the Vref of the dry or wet landing distance.

The rest of the report is plain data: `TLR.Takeoff.Conditions` and
`TLR.Landing.Conditions` (`TLRConditions`: planned runway and weight, wind,
temperature, altimeter, surface), `TLR.Takeoff.Runways` (`[]TakeoffRunway`:
flaps, thrust, flex temperature, V-speeds, limits and distances),
`TLR.Landing.DistanceDry`/`DistanceWet` (`TLRLandingDistance`) and
`TLR.Landing.Runways` (`[]LandingRunway`). Both runway types embed
`TLRRunway` (lengths, course, wind components, ILS frequency).

## Generating a plan

The library does not create plans itself. It builds the URL of SimBrief's
dispatch page with all parameters filled in; open it in a browser that is
logged in to SimBrief. Give the plan a `StaticID` to fetch exactly that plan
afterwards with `GetFlightPlanByStaticID`.

```go
req := client.NewFlightPlan("LKPR", "EGLL", "A320").
	Route("VOZ7B VOZ DCT ...").
	Airline("CSA").
	FlightNumber("123").
	AltitudeFromFlightLevel(360).
	DepartureTime(14, 30).
	Passengers(150).
	Units(types.UnitsKGS).
	StaticID("MYAPP_42").
	EnableNavLog().
	Build()

runwayAnalysis := true
req.RunwayAnalysis = &runwayAnalysis // TLR in the OFP (aircraft with TLRData)

if err := c.ValidateFlightPlanRequest(req); err != nil {
	return err
}
fmt.Println(c.GenerateFlightPlanURL(req)) // open in a browser logged in to SimBrief
```

`FlightPlanBuilder` (from `client.NewFlightPlan(origin, destination,
aircraft)`) also has `Alternate`, `Date`/`DateFromTime`, `Altitude`,
`AltitudeFromFeet`, `Cargo`, `PlanFormat`, `Registration`, `CallSign`,
`Captain`, `Dispatcher`, `DisableNavLog`, `EnableETOPS`, `EnableStepClimbs`,
`CustomAircraftData` (`*types.AircraftData`), `TaxiTimes` and `Runways`.
Every other SimBrief input (fuel, profiles, alternates 1-4, NOTAMs, maps ...)
is a field on `types.FlightPlanRequest`; set it on the built request or build
the struct directly:

```go
req := types.NewFlightPlanRequest("KJFK", "KLAX", "B38M")
req.Route = "DCT"
req.CostIndex = "25"
fmt.Println(c.GenerateFlightPlanURL(req))
```

`ValidateFlightPlanRequest` checks that origin, destination and aircraft are
set, that both airports are 4-letter codes and that a departure time is in
range. `req.Validate()` only checks the three required fields and returns the
sentinel errors below.

To send the user to edit a plan they already generated with a static ID:

```go
fmt.Println(c.GetDirectEditURL("MYAPP_42"))
// https://www.simbrief.com/system/dispatch.php?editflight=last&static_id=MYAPP_42
```

## Aircraft types and plan formats

```go
aircraft, err := c.GetAircraftTypes() // map[string]types.AircraftOption
if err != nil {
	return err
}
for id, a := range aircraft {
	fmt.Println(id, a.Name, a.Accuracy, a.TLRData)
}

layouts, err := c.GetPlanFormats() // map[string]types.LayoutOption
if err != nil {
	return err
}
for id, l := range layouts {
	fmt.Println(id, l.NameShort, l.NameLong)
}
```

Both read `/api/inputs.list.json`; `GetSupportedOptions()` returns the whole
`types.SupportedOptions` (aircraft, layouts, last update) in one request.

## Helpers

Small stateless helpers in `pkg/client`:

```go
rh := client.NewRouteHelper()
fmt.Println(rh.ParseRoute("VOZ7B  VOZ  DCT")) // [VOZ7B VOZ DCT]
fmt.Println(rh.ValidateICAOCode("LKPR"))      // true
fmt.Println(rh.FormatFlightLevel(36000))      // FL360
feet, err := rh.ParseFlightLevel("FL360")     // 36000 (plain "5000" works too)
if err != nil {
	return err
}
fmt.Println(feet)

fh := client.NewFuelHelper()
fmt.Println(fh.ConvertLBSToKGS(1000), fh.ConvertKGSToLBS(453.592)) // 453.592 1000
pct, minutes, err := fh.ParseFuelValue("0.05/15")                  // 0.05 "15"
if err != nil {
	return err
}
fmt.Println(pct, minutes)

th := client.NewTimeHelper()
h, m, err := th.ParseTimeString("14:30") // 14 30
if err != nil {
	return err
}
fmt.Println(h, m, th.FormatTimeString(h, m)) // 14 30 14:30
mins, err := th.ParseDuration("01:27")       // 87
if err != nil {
	return err
}
fmt.Println(mins, th.FormatDuration(mins)) // 87 01:27
```

`ValidateICAOCode` accepts upper-case letters and digits only.

## Errors

- **`types.APIError`** (`Message`, `Code`): SimBrief answered but the fetch
  failed. `Message` is the OFP's fetch status; `Code` is the HTTP status, or
  0 when the error came with a 200 reply.
- **Wrapped errors** (`fmt.Errorf` with `%w`): transport and decoding
  failures, e.g. a timeout or a body that is not an OFP. A non-200 reply
  without a fetch status is a plain error with the status and body.
- **Sentinels** from `FlightPlanRequest.Validate()`:
  `types.ErrMissingOrigin`, `types.ErrMissingDestination`,
  `types.ErrMissingAircraft`. (`ErrMissingUserID`, `ErrInvalidUserID` and
  `ErrInvalidAPIKey` are declared but not returned by the library.)
- `ValidateFlightPlanRequest` returns plain `fmt.Errorf` messages.

```go
_, err := c.GetFlightPlan(&types.FetchRequest{UserID: "123456", JSON: true})
var apiErr types.APIError
if errors.As(err, &apiErr) {
	fmt.Println("SimBrief:", apiErr.Message, apiErr.Code) // fetch status text, HTTP status (0 on a 200)
}
```

```go
err := types.NewFlightPlanRequest("LKPR", "", "A320").Validate()
if errors.Is(err, types.ErrMissingDestination) {
	fmt.Println("destination missing")
}
```

## Further reading

- [SimBrief API documentation](https://developers.navigraph.com/docs/simbrief/using-the-api)
- [examples/](../examples/) - runnable programs
