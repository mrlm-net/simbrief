package client

import (
	"testing"
	"time"
)

func TestRouteHelper_ParseRoute(t *testing.T) {
	helper := NewRouteHelper()

	tests := []struct {
		name     string
		route    string
		expected []string
	}{
		{
			name:     "simple route",
			route:    "KJFK HAPIE J174 COATE KLAX",
			expected: []string{"KJFK", "HAPIE", "J174", "COATE", "KLAX"},
		},
		{
			name:     "route with extra spaces",
			route:    "  KJFK   HAPIE    J174   COATE   KLAX  ",
			expected: []string{"KJFK", "HAPIE", "J174", "COATE", "KLAX"},
		},
		{
			name:     "empty route",
			route:    "",
			expected: []string{},
		},
		{
			name:     "single waypoint",
			route:    "KJFK",
			expected: []string{"KJFK"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.ParseRoute(tt.route)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseRoute() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("ParseRoute()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestRouteHelper_ValidateICAOCode(t *testing.T) {
	helper := NewRouteHelper()

	tests := []struct {
		name string
		code string
		want bool
	}{
		{name: "valid ICAO", code: "KJFK", want: true},
		{name: "valid ICAO with numbers", code: "K1G3", want: true},
		{name: "too short", code: "JFK", want: false},
		{name: "too long", code: "KJFKX", want: false},
		{name: "lowercase", code: "kjfk", want: false},
		{name: "special characters", code: "K-FK", want: false},
		{name: "empty", code: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.ValidateICAOCode(tt.code); got != tt.want {
				t.Errorf("ValidateICAOCode(%s) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestRouteHelper_FormatFlightLevel(t *testing.T) {
	helper := NewRouteHelper()

	tests := []struct {
		name string
		feet int
		want string
	}{
		{name: "FL340", feet: 34000, want: "FL340"},
		{name: "FL100", feet: 10000, want: "FL100"},
		{name: "FL050", feet: 5000, want: "FL050"},
		{name: "FL001", feet: 100, want: "FL001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.FormatFlightLevel(tt.feet); got != tt.want {
				t.Errorf("FormatFlightLevel(%d) = %v, want %v", tt.feet, got, tt.want)
			}
		})
	}
}

func TestRouteHelper_ParseFlightLevel(t *testing.T) {
	helper := NewRouteHelper()

	tests := []struct {
		name    string
		flStr   string
		want    int
		wantErr bool
	}{
		{name: "FL340", flStr: "FL340", want: 34000, wantErr: false},
		{name: "fl340 lowercase", flStr: "fl340", want: 34000, wantErr: false},
		{name: "direct feet", flStr: "34000", want: 34000, wantErr: false},
		{name: "with spaces", flStr: " FL340 ", want: 34000, wantErr: false},
		{name: "invalid FL", flStr: "FLXXX", want: 0, wantErr: true},
		{name: "invalid feet", flStr: "ABC", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helper.ParseFlightLevel(tt.flStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlightLevel(%s) error = %v, wantErr %v", tt.flStr, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFlightLevel(%s) = %v, want %v", tt.flStr, got, tt.want)
			}
		})
	}
}

func TestFuelHelper_Conversions(t *testing.T) {
	helper := NewFuelHelper()

	// Test LBS to KGS
	lbs := 10000.0
	kgs := helper.ConvertLBSToKGS(lbs)
	expected := 4535.92
	if kgs < expected-0.1 || kgs > expected+0.1 {
		t.Errorf("ConvertLBSToKGS(%f) = %f, want ~%f", lbs, kgs, expected)
	}

	// Test KGS to LBS (round trip)
	backToLbs := helper.ConvertKGSToLBS(kgs)
	if backToLbs < lbs-0.1 || backToLbs > lbs+0.1 {
		t.Errorf("Round trip conversion failed: %f -> %f -> %f", lbs, kgs, backToLbs)
	}
}

func TestFuelHelper_ParseFuelValue(t *testing.T) {
	helper := NewFuelHelper()

	tests := []struct {
		name       string
		value      string
		wantWeight float64
		wantTime   string
		wantErr    bool
	}{
		{name: "simple weight", value: "5000", wantWeight: 5000, wantTime: "", wantErr: false},
		{name: "weight and time", value: "0.05/15", wantWeight: 0.05, wantTime: "15", wantErr: false},
		{name: "with spaces", value: " 5000 ", wantWeight: 5000, wantTime: "", wantErr: false},
		{name: "empty", value: "", wantWeight: 0, wantTime: "", wantErr: true},
		{name: "invalid", value: "ABC", wantWeight: 0, wantTime: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWeight, gotTime, err := helper.ParseFuelValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFuelValue(%s) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if gotWeight != tt.wantWeight {
				t.Errorf("ParseFuelValue(%s) weight = %v, want %v", tt.value, gotWeight, tt.wantWeight)
			}
			if gotTime != tt.wantTime {
				t.Errorf("ParseFuelValue(%s) time = %v, want %v", tt.value, gotTime, tt.wantTime)
			}
		})
	}
}

func TestTimeHelper_ParseTimeString(t *testing.T) {
	helper := NewTimeHelper()

	tests := []struct {
		name       string
		timeStr    string
		wantHour   int
		wantMinute int
		wantErr    bool
	}{
		{name: "valid time", timeStr: "14:30", wantHour: 14, wantMinute: 30, wantErr: false},
		{name: "midnight", timeStr: "00:00", wantHour: 0, wantMinute: 0, wantErr: false},
		{name: "late night", timeStr: "23:59", wantHour: 23, wantMinute: 59, wantErr: false},
		{name: "invalid format", timeStr: "14-30", wantHour: 0, wantMinute: 0, wantErr: true},
		{name: "invalid hour", timeStr: "25:30", wantHour: 0, wantMinute: 0, wantErr: true},
		{name: "invalid minute", timeStr: "14:60", wantHour: 0, wantMinute: 0, wantErr: true},
		{name: "non-numeric", timeStr: "AB:CD", wantHour: 0, wantMinute: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHour, gotMinute, err := helper.ParseTimeString(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeString(%s) error = %v, wantErr %v", tt.timeStr, err, tt.wantErr)
				return
			}
			if gotHour != tt.wantHour {
				t.Errorf("ParseTimeString(%s) hour = %v, want %v", tt.timeStr, gotHour, tt.wantHour)
			}
			if gotMinute != tt.wantMinute {
				t.Errorf("ParseTimeString(%s) minute = %v, want %v", tt.timeStr, gotMinute, tt.wantMinute)
			}
		})
	}
}

func TestTimeHelper_FormatTimeString(t *testing.T) {
	helper := NewTimeHelper()

	tests := []struct {
		name   string
		hour   int
		minute int
		want   string
	}{
		{name: "afternoon", hour: 14, minute: 30, want: "14:30"},
		{name: "midnight", hour: 0, minute: 0, want: "00:00"},
		{name: "single digits", hour: 9, minute: 5, want: "09:05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.FormatTimeString(tt.hour, tt.minute); got != tt.want {
				t.Errorf("FormatTimeString(%d, %d) = %v, want %v", tt.hour, tt.minute, got, tt.want)
			}
		})
	}
}

func TestTimeHelper_ParseDuration(t *testing.T) {
	helper := NewTimeHelper()

	tests := []struct {
		name        string
		durationStr string
		want        int
		wantErr     bool
	}{
		{name: "2 hours 30 minutes", durationStr: "02:30", want: 150, wantErr: false},
		{name: "1 hour", durationStr: "01:00", want: 60, wantErr: false},
		{name: "30 minutes", durationStr: "00:30", want: 30, wantErr: false},
		{name: "invalid format", durationStr: "2h30m", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helper.ParseDuration(tt.durationStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration(%s) error = %v, wantErr %v", tt.durationStr, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration(%s) = %v, want %v", tt.durationStr, got, tt.want)
			}
		})
	}
}

func TestTimeHelper_FormatDuration(t *testing.T) {
	helper := NewTimeHelper()

	tests := []struct {
		name         string
		totalMinutes int
		want         string
	}{
		{name: "2 hours 30 minutes", totalMinutes: 150, want: "02:30"},
		{name: "1 hour", totalMinutes: 60, want: "01:00"},
		{name: "30 minutes", totalMinutes: 30, want: "00:30"},
		{name: "over 24 hours", totalMinutes: 1500, want: "25:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.FormatDuration(tt.totalMinutes); got != tt.want {
				t.Errorf("FormatDuration(%d) = %v, want %v", tt.totalMinutes, got, tt.want)
			}
		})
	}
}

func TestFlightPlanBuilder_DateFromTime(t *testing.T) {
	builder := NewFlightPlan("KJFK", "KLAX", "B738")

	// Test date formatting
	testTime := time.Date(2023, 7, 15, 14, 30, 0, 0, time.UTC)
	request := builder.DateFromTime(testTime).Build()

	expected := "15Jul23"
	if request.Date != expected {
		t.Errorf("DateFromTime() = %s, want %s", request.Date, expected)
	}
}

func TestFlightPlanBuilder_AltitudeFromFeet(t *testing.T) {
	builder := NewFlightPlan("KJFK", "KLAX", "B738")

	request := builder.AltitudeFromFeet(34000).Build()

	expected := "34000"
	if request.Altitude != expected {
		t.Errorf("AltitudeFromFeet(34000) = %s, want %s", request.Altitude, expected)
	}
}

func TestFlightPlanBuilder_AltitudeFromFlightLevel(t *testing.T) {
	builder := NewFlightPlan("KJFK", "KLAX", "B738")

	request := builder.AltitudeFromFlightLevel(340).Build()

	expected := "FL340"
	if request.Altitude != expected {
		t.Errorf("AltitudeFromFlightLevel(340) = %s, want %s", request.Altitude, expected)
	}
}
