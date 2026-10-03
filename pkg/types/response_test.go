package types

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The golden files are a real OFP (CEF007 LKPR-LKPD, A319) fetched from
// /api/xml.fetcher.php as JSON v2 (json=v2) and as XML.

func loadJSON(t *testing.T) *FlightPlanResponse {
	t.Helper()
	data, err := os.ReadFile("testdata/ofp_v2.json")
	require.NoError(t, err)
	var ofp FlightPlanResponse
	require.NoError(t, json.Unmarshal(data, &ofp))
	return &ofp
}

func loadXML(t *testing.T) *FlightPlanResponse {
	t.Helper()
	data, err := os.ReadFile("testdata/ofp.xml")
	require.NoError(t, err)
	var ofp FlightPlanResponse
	require.NoError(t, xml.Unmarshal(data, &ofp))
	return &ofp
}

func TestGoldenOFP(t *testing.T) {
	for name, load := range map[string]func(*testing.T) *FlightPlanResponse{
		"json-v2": loadJSON,
		"xml":     loadXML,
	} {
		t.Run(name, func(t *testing.T) {
			ofp := load(t)

			assert.Equal(t, "Success", ofp.Fetch.Status)
			assert.Equal(t, "857341", ofp.Params.UserID)
			assert.Equal(t, "2610", ofp.Params.AIRAC)
			assert.Equal(t, "kgs", ofp.Params.Units)
			assert.Equal(t, "", ofp.Params.StaticID.String())

			// ATC
			assert.Equal(t, "CEF007", ofp.ATC.Callsign)
			assert.Equal(t, "I", ofp.ATC.FlightRules)
			assert.Equal(t, "S", ofp.ATC.FlightType)
			assert.Equal(t, Number("150"), ofp.ATC.InitialAlt)
			assert.Equal(t, 150, ofp.ATC.InitialAlt.Int())
			assert.Equal(t, "F", ofp.ATC.InitialAltUnit)
			assert.Equal(t, 365, ofp.ATC.InitialSpeed.Int())
			assert.Equal(t, "N", ofp.ATC.InitialSpeedUnit)
			assert.Equal(t, "N0365F150 DCT BEKVI BEKV1Q", ofp.ATC.Route)
			assert.Equal(t, "N0365F150 DCT BEKVI BEKVI1Q", ofp.ATC.RouteIFPS)
			assert.Contains(t, ofp.ATC.FlightPlanText, "(FPL-CEF007-IS")
			assert.Equal(t, []string{"LKAA"}, ofp.ATC.FIRAltn)

			// General
			g := ofp.General
			assert.Equal(t, "CEF", g.ICAOAirline)
			assert.Equal(t, "007", g.FlightNumber)
			assert.Equal(t, "DCT BEKVI BEKV1Q", g.Route)
			assert.Equal(t, "DCT BEKVI BEKVI1Q", g.RouteIFPS)
			assert.Equal(t, "", g.SIDIdent)
			assert.Equal(t, "", g.SIDTrans)
			assert.Equal(t, "BEKV1Q", g.STARIdent)
			assert.Equal(t, "", g.STARTrans)
			assert.Equal(t, 15000, g.InitialAltitude.Int())
			assert.Equal(t, 26, g.CostIndex.Int())
			assert.Equal(t, 142, g.Passengers.Int())
			assert.InDelta(t, 0.57, g.CruiseMach.Float(), 1e-9)
			assert.Equal(t, []string{"PLANNED OPTIMUM FLIGHT LEVEL"}, g.DispatchRemarks)

			// Origin / destination
			assert.Equal(t, "LKPR", ofp.Origin.ICAO)
			assert.Equal(t, "PRG", ofp.Origin.IATA)
			assert.Equal(t, "RUZYNE", ofp.Origin.Name)
			assert.Equal(t, "06", ofp.Origin.Runway)
			assert.Equal(t, 1234, ofp.Origin.Elevation.Int())
			assert.Equal(t, 5000, ofp.Origin.TransAlt.Int())
			assert.Equal(t, 6000, ofp.Origin.TransLevel.Int())
			assert.InDelta(t, 50.100833, ofp.Origin.Latitude.Float(), 1e-9)
			assert.InDelta(t, 14.26, ofp.Origin.Longitude.Float(), 1e-9)
			assert.Equal(t, "LKPR 032100Z VRB01KT CAVOK 12/08 Q1029 NOSIG", ofp.Origin.METAR)

			assert.Equal(t, "LKPD", ofp.Destination.ICAO)
			assert.Equal(t, "PED", ofp.Destination.IATA)
			assert.Equal(t, "PARDUBICE", ofp.Destination.Name)
			assert.Equal(t, "09", ofp.Destination.Runway)
			assert.Equal(t, 741, ofp.Destination.Elevation.Int())

			// Alternates
			require.Len(t, ofp.Alternates, 1)
			altn := ofp.Alternate()
			require.NotNil(t, altn)
			assert.Equal(t, "LKPR", altn.ICAO)
			assert.Equal(t, "06", altn.Runway)
			assert.Equal(t, 5000, altn.TransAlt.Int())
			assert.Equal(t, 1025, altn.Burn.Int())
			assert.Equal(t, "BEKV4V BEKVI DCT", altn.Route)
			assert.Equal(t, "M010", altn.AvgWindComp)
			require.Len(t, ofp.AlternateNavLogs, 1)
			require.NotEmpty(t, ofp.AlternateNavLogs[0])
			assert.Equal(t, "PD903", ofp.AlternateNavLogs[0][0].Ident)

			// Navlog
			require.Len(t, ofp.NavLog, 10)
			first, last := ofp.NavLog[0], ofp.NavLog[len(ofp.NavLog)-1]
			assert.Equal(t, "BEKVI", first.Ident)
			assert.Equal(t, "wpt", first.Type)
			assert.Equal(t, "DCT", first.ViaAirway)
			assert.Equal(t, "CLB", first.Stage)
			assert.False(t, first.IsSIDSTAR.Bool())
			assert.Equal(t, 13500, first.AltitudeFeet.Int())
			assert.InDelta(t, 50.073356, first.Latitude.Float(), 1e-9)
			assert.InDelta(t, 14.722358, first.Longitude.Float(), 1e-9)
			assert.False(t, first.Frequency.IsSet())
			assert.Equal(t, "TOC", ofp.NavLog[1].Ident)
			assert.True(t, ofp.NavLog[1].IsSIDSTAR.Bool())
			assert.Equal(t, "LKPD", last.Ident)
			assert.Equal(t, "apt", last.Type)
			assert.Equal(t, "DSC", last.Stage)

			// Aircraft
			assert.Equal(t, "A319", ofp.Aircraft.ICAOCode)
			assert.Equal(t, "A319", ofp.Aircraft.ICAO)
			assert.Equal(t, "G-SMOL", ofp.Aircraft.Registration)
			assert.Equal(t, "FENIX A319", ofp.Aircraft.Name)
			assert.Equal(t, "CFM56-5B6", ofp.Aircraft.Engines)
			assert.Equal(t, 150, ofp.Aircraft.MaxPassengers.Int())

			// Fuel, weights, weather, files
			assert.Equal(t, 3979, ofp.Fuel.PlanRamp.Int())
			assert.Equal(t, 1057, ofp.Fuel.EnrouteBurn.Int())
			assert.Equal(t, 55400, ofp.Weights.EstZFW.Int())
			assert.Equal(t, 142, ofp.Weights.PaxCount.Int())
			assert.Equal(t, ofp.Origin.METAR, ofp.Weather.OrigMETAR)
			assert.Len(t, ofp.Weather.AltnMETAR, 1)
			assert.Equal(t, "LKPRLKPD_PDF_1791062318.pdf", ofp.Files.PDF.Link)
			assert.NotEmpty(t, ofp.Files.Files)
			assert.Contains(t, ofp.Links.SkyVector, "skyvector.com")
		})
	}
}

// Times differ by format: JSON v2 uses ISO 8601 and HH:MM:SS, XML uses Unix
// seconds and seconds.
func TestGoldenOFPTimeFormats(t *testing.T) {
	j, x := loadJSON(t), loadXML(t)
	assert.Equal(t, "2026-10-03T21:50:00Z", j.Times.EstOut)
	assert.Equal(t, "1791064200", x.Times.EstOut)
	assert.Equal(t, "00:23:03", j.Times.EstTimeEnroute)
	assert.Equal(t, "1383", x.Times.EstTimeEnroute)
	assert.Equal(t, "00:04:55", j.NavLog[0].TimeLeg)
	assert.Equal(t, "295", x.NavLog[0].TimeLeg)
}

func TestNumber(t *testing.T) {
	var v struct {
		A, B, C, D, E Number
	}
	require.NoError(t, json.Unmarshal([]byte(`{"A":"0365","B":15000,"C":".57","D":null,"E":""}`), &v))
	assert.Equal(t, 365, v.A.Int())
	assert.Equal(t, "0365", v.A.String())
	assert.Equal(t, 15000, v.B.Int())
	assert.Equal(t, Number("15000"), v.B)
	assert.InDelta(t, 0.57, v.C.Float(), 1e-9)
	assert.Equal(t, 0, v.C.Int())
	assert.False(t, v.D.IsSet())
	assert.Equal(t, 0, v.E.Int())
	assert.False(t, v.E.Bool())
	assert.True(t, Number("1").Bool())
	assert.Equal(t, 50, Number("50.9").Int())

	var x struct {
		N Number `xml:"n"`
	}
	require.NoError(t, xml.Unmarshal([]byte(`<r><n>05000</n></r>`), &x))
	assert.Equal(t, 5000, x.N.Int())

	assert.Error(t, json.Unmarshal([]byte(`{"A":{}}`), &v))
}

func TestNavLogXMLRoundTrip(t *testing.T) {
	in := struct {
		XMLName xml.Name `xml:"r"`
		NavLog  NavLog   `xml:"navlog"`
	}{NavLog: NavLog{{Ident: "A"}, {Ident: "B"}}}
	data, err := xml.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(data), "<navlog><fix><ident>A</ident>")

	out := in
	out.NavLog = nil
	require.NoError(t, xml.Unmarshal(data, &out))
	require.Len(t, out.NavLog, 2)
	assert.Equal(t, "B", out.NavLog[1].Ident)
}
