package client

import (
	"testing"

	"github.com/mrlm-net/simbrief/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	assert.Equal(t, apiKey, client.APIKey)
	assert.Equal(t, DefaultBaseURL, client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, DefaultTimeout, client.HTTPClient.Timeout)
}

func TestNewClientWithConfig(t *testing.T) {
	apiKey := "test-api-key"
	baseURL := "https://custom.simbrief.com"

	client := NewClientWithConfig(apiKey, baseURL, nil)

	assert.Equal(t, apiKey, client.APIKey)
	assert.Equal(t, baseURL, client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
}

func TestValidateFlightPlanRequest(t *testing.T) {
	client := NewClient("test-key")

	tests := []struct {
		name    string
		request *types.FlightPlanRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &types.FlightPlanRequest{
				Origin:      "KJFK",
				Destination: "KLAX",
				Aircraft:    "B738",
			},
			wantErr: false,
		},
		{
			name: "missing origin",
			request: &types.FlightPlanRequest{
				Destination: "KLAX",
				Aircraft:    "B738",
			},
			wantErr: true,
			errMsg:  "origin airport (orig) is required",
		},
		{
			name: "missing destination",
			request: &types.FlightPlanRequest{
				Origin:   "KJFK",
				Aircraft: "B738",
			},
			wantErr: true,
			errMsg:  "destination airport (dest) is required",
		},
		{
			name: "missing aircraft",
			request: &types.FlightPlanRequest{
				Origin:      "KJFK",
				Destination: "KLAX",
			},
			wantErr: true,
			errMsg:  "aircraft type (type) is required",
		},
		{
			name: "invalid origin format",
			request: &types.FlightPlanRequest{
				Origin:      "JFK",
				Destination: "KLAX",
				Aircraft:    "B738",
			},
			wantErr: true,
			errMsg:  "origin airport code must be 4 characters (ICAO format)",
		},
		{
			name: "invalid departure hour",
			request: &types.FlightPlanRequest{
				Origin:        "KJFK",
				Destination:   "KLAX",
				Aircraft:      "B738",
				DepartureHour: 25,
			},
			wantErr: true,
			errMsg:  "departure hour must be between 0 and 23",
		},
		{
			name: "invalid departure minute",
			request: &types.FlightPlanRequest{
				Origin:          "KJFK",
				Destination:     "KLAX",
				Aircraft:        "B738",
				DepartureMinute: 60,
			},
			wantErr: true,
			errMsg:  "departure minute must be between 0 and 59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ValidateFlightPlanRequest(tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateFlightPlanURL(t *testing.T) {
	client := NewClient("test-api-key")

	request := &types.FlightPlanRequest{
		Origin:       "KJFK",
		Destination:  "KLAX",
		Aircraft:     "B738",
		Route:        "HAPIE6 HAPIE J174 COATE",
		Airline:      "UAL",
		FlightNumber: "1234",
	}

	url := client.GenerateFlightPlanURL(request)

	assert.Contains(t, url, DefaultBaseURL)
	assert.Contains(t, url, "orig=KJFK")
	assert.Contains(t, url, "dest=KLAX")
	assert.Contains(t, url, "type=B738")
	assert.Contains(t, url, "route=HAPIE6+HAPIE+J174+COATE")
	assert.Contains(t, url, "airline=UAL")
	assert.Contains(t, url, "fltnum=1234")
	assert.Contains(t, url, "api_key=test-api-key")
}

func TestGetDirectEditURL(t *testing.T) {
	client := NewClient("test-key")
	staticID := "UAL_1234_TEST"

	url := client.GetDirectEditURL(staticID)

	expected := DefaultBaseURL + "/system/dispatch.php?editflight=last&static_id=UAL_1234_TEST"
	assert.Equal(t, expected, url)
}

func TestFlightPlanBuilder(t *testing.T) {
	builder := NewFlightPlan("KJFK", "KLAX", "B738")

	request := builder.
		Route("HAPIE6 HAPIE J174 COATE").
		Airline("UAL").
		FlightNumber("1234").
		Alternate("KLAS").
		DepartureTime(14, 30).
		AltitudeFromFlightLevel(340).
		Passengers(150).
		Units(types.UnitsLBS).
		Registration("N12345").
		StaticID("TEST_FLIGHT").
		Build()

	assert.Equal(t, "KJFK", request.Origin)
	assert.Equal(t, "KLAX", request.Destination)
	assert.Equal(t, "B738", request.Aircraft)
	assert.Equal(t, "HAPIE6 HAPIE J174 COATE", request.Route)
	assert.Equal(t, "UAL", request.Airline)
	assert.Equal(t, "1234", request.FlightNumber)
	assert.Equal(t, "KLAS", request.Alternate)
	assert.Equal(t, 14, request.DepartureHour)
	assert.Equal(t, 30, request.DepartureMinute)
	assert.Equal(t, "FL340", request.Altitude)
	assert.Equal(t, 150, request.Passengers)
	assert.Equal(t, types.UnitsLBS, request.Units)
	assert.Equal(t, "N12345", request.Registration)
	assert.Equal(t, "TEST_FLIGHT", request.StaticID)
}

func TestFlightPlanBuilderBooleans(t *testing.T) {
	builder := NewFlightPlan("KJFK", "KLAX", "B738")

	request := builder.
		EnableNavLog().
		EnableETOPS().
		EnableStepClimbs().
		Build()

	require.NotNil(t, request.NavLog)
	assert.True(t, *request.NavLog)

	require.NotNil(t, request.ETOPS)
	assert.True(t, *request.ETOPS)

	require.NotNil(t, request.StepClimbs)
	assert.True(t, *request.StepClimbs)
}

func TestFlightPlanBuilderCustomAircraft(t *testing.T) {
	customAircraft := &types.AircraftData{
		ICAO:     "B38M",
		Name:     "737 MAX 8",
		Category: "M",
		OEW:      99.3,
		MTOW:     181.2,
		MaxFuel:  46.0,
	}

	builder := NewFlightPlan("KJFK", "KLAX", "B38M")
	request := builder.CustomAircraftData(customAircraft).Build()

	require.NotNil(t, request.AircraftData)
	assert.Equal(t, "B38M", request.AircraftData.ICAO)
	assert.Equal(t, "737 MAX 8", request.AircraftData.Name)
	assert.Equal(t, "M", request.AircraftData.Category)
	assert.Equal(t, 99.3, request.AircraftData.OEW)
}

func TestAircraftDataString(t *testing.T) {
	data := &types.AircraftData{
		ICAO:     "B738",
		Name:     "737-800",
		Category: "M",
		OEW:      90.7,
		MTOW:     174.2,
	}

	jsonStr := data.String()

	assert.Contains(t, jsonStr, `"icao":"B738"`)
	assert.Contains(t, jsonStr, `"name":"737-800"`)
	assert.Contains(t, jsonStr, `"cat":"M"`)
	assert.Contains(t, jsonStr, `"oew":90.7`)
	assert.Contains(t, jsonStr, `"mtow":174.2`)
}

func TestFetchRequestToQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		request  *types.FetchRequest
		contains []string
		empty    bool
	}{
		{
			name: "user ID only",
			request: &types.FetchRequest{
				UserID: "123456",
			},
			contains: []string{"userid=123456"},
		},
		{
			name: "username only",
			request: &types.FetchRequest{
				Username: "testuser",
			},
			contains: []string{"username=testuser"},
		},
		{
			name: "user ID with static ID",
			request: &types.FetchRequest{
				UserID:   "123456",
				StaticID: "TEST_FLIGHT",
			},
			contains: []string{"userid=123456", "static_id=TEST_FLIGHT"},
		},
		{
			name: "username with JSON",
			request: &types.FetchRequest{
				Username: "testuser",
				JSON:     true,
			},
			contains: []string{"username=testuser", "json=1"},
		},
		{
			name: "all parameters",
			request: &types.FetchRequest{
				UserID:   "123456",
				StaticID: "TEST_FLIGHT",
				JSON:     true,
			},
			contains: []string{"userid=123456", "static_id=TEST_FLIGHT", "json=1"},
		},
		{
			name:    "empty request",
			request: &types.FetchRequest{},
			empty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.ToQueryParams()
			if tt.empty {
				assert.Equal(t, "", result)
			} else {
				assert.True(t, len(result) > 0)
				assert.True(t, result[0] == '?')
				for _, param := range tt.contains {
					assert.Contains(t, result, param)
				}
			}
		})
	}
}
