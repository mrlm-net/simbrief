package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mrlm-net/simbrief/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ofpServer serves the golden OFP files the way xml.fetcher.php does:
// JSON v2 for json=v2, XML otherwise.
func ofpServer(t *testing.T) *httptest.Server {
	t.Helper()
	jsonBody, err := os.ReadFile("../types/testdata/ofp_v2.json")
	require.NoError(t, err)
	xmlBody, err := os.ReadFile("../types/testdata/ofp.xml")
	require.NoError(t, err)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.URL.Path != endpointXMLFetcher || q.Get("userid") != "857341" {
			w.WriteHeader(http.StatusBadRequest)
			if q.Get("json") == "v2" {
				_, _ = w.Write([]byte(`{"fetch":{"status":"Error: Unknown UserID"}}`))
			} else {
				_, _ = w.Write([]byte(`<?xml version="1.0"?><OFP><fetch><status>Error: Unknown UserID</status></fetch></OFP>`))
			}
			return
		}
		if q.Get("json") == "v2" {
			_, _ = w.Write(jsonBody)
			return
		}
		_, _ = w.Write(xmlBody)
	}))
}

func TestFetchFlightPlanGolden(t *testing.T) {
	srv := ofpServer(t)
	defer srv.Close()
	c := NewClientWithConfig(srv.URL, nil)

	viaJSON, err := c.GetFlightPlanByUserID("857341")
	require.NoError(t, err)
	viaXML, err := c.GetFlightPlan(&types.FetchRequest{UserID: "857341"})
	require.NoError(t, err)

	for name, ofp := range map[string]*types.FlightPlanResponse{"json": viaJSON, "xml": viaXML} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, "CEF007", ofp.ATC.Callsign)
			assert.Equal(t, "06", ofp.Origin.Runway)
			assert.Equal(t, "09", ofp.Destination.Runway)
			assert.Equal(t, "BEKV1Q", ofp.General.STARIdent)
			assert.Equal(t, 15000, ofp.General.InitialAltitude.Int())
			assert.Equal(t, "A319", ofp.Aircraft.ICAOCode)
			require.NotEmpty(t, ofp.NavLog)
			assert.Equal(t, "BEKVI", ofp.NavLog[0].Ident)
			assert.Equal(t, "LKPD", ofp.NavLog[len(ofp.NavLog)-1].Ident)
		})
	}
}

func TestFetchFlightPlanError(t *testing.T) {
	srv := ofpServer(t)
	defer srv.Close()
	c := NewClientWithConfig(srv.URL, nil)

	_, err := c.GetFlightPlanByUserID("1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown UserID")

	_, err = c.GetFlightPlan(&types.FetchRequest{UserID: "1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown UserID")
}
