package client

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mrlm-net/simbrief/pkg/types"
)

const (
	// DefaultBaseURL is the default SimBrief API base URL
	DefaultBaseURL = "https://www.simbrief.com"

	// API endpoints - based on official documentation
	endpointXMLFetcher = "/api/xml.fetcher.php"  // Fetch existing flight plan data
	endpointInputsList = "/api/inputs.list.json" // Get supported aircraft/layouts (JSON)
	endpointInputsXML  = "/api/inputs.list.xml"  // Get supported aircraft/layouts (XML)
	endpointGenerate   = "/system/dispatch.php"  // Generate new flight plan

	// Default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second
)

// Client represents a SimBrief API client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new SimBrief API client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewClientWithConfig creates a new SimBrief API client with custom configuration
func NewClientWithConfig(apiKey, baseURL string, httpClient *http.Client) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultTimeout}
	}

	return &Client{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

// GetFlightPlanByUserID retrieves the latest flight plan for a specific user ID
func (c *Client) GetFlightPlanByUserID(userID string) (*types.FlightPlanResponse, error) {
	req := &types.FetchRequest{
		UserID: userID,
		JSON:   true, // Request JSON format for easier parsing
	}
	return c.fetchFlightPlan(req)
}

// GetFlightPlanByUsername retrieves the latest flight plan for a specific username
func (c *Client) GetFlightPlanByUsername(username string) (*types.FlightPlanResponse, error) {
	req := &types.FetchRequest{
		Username: username,
		JSON:     true,
	}
	return c.fetchFlightPlan(req)
}

// GetFlightPlanByStaticID retrieves a specific flight plan using static ID
func (c *Client) GetFlightPlanByStaticID(userID, staticID string) (*types.FlightPlanResponse, error) {
	req := &types.FetchRequest{
		UserID:   userID,
		StaticID: staticID,
		JSON:     true,
	}
	return c.fetchFlightPlan(req)
}

// GetFlightPlanXML retrieves flight plan data in XML format
func (c *Client) GetFlightPlanXML(req *types.FetchRequest) ([]byte, error) {
	// Force XML format
	req.JSON = false

	fullURL := c.BaseURL + endpointXMLFetcher + req.ToQueryParams()

	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse error from XML
		var apiErr types.APIError
		if err := xml.Unmarshal(body, &apiErr); err == nil {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetSupportedOptions retrieves the list of supported aircraft types and plan formats
func (c *Client) GetSupportedOptions() (*types.SupportedOptions, error) {
	fullURL := c.BaseURL + endpointInputsList

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var options types.SupportedOptions
	if err := json.NewDecoder(resp.Body).Decode(&options); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &options, nil
}

// GenerateFlightPlanURL generates a URL for flight plan generation
// Note: Actual flight plan generation requires browser popup authentication
func (c *Client) GenerateFlightPlanURL(req *types.FlightPlanRequest) string {
	values := req.ToURLValues()

	// Add API key if available
	if c.APIKey != "" {
		values.Add("api_key", c.APIKey)
	}

	return c.BaseURL + endpointGenerate + "?" + values.Encode()
}

// ValidateFlightPlanRequest validates that a flight plan request has all required fields
func (c *Client) ValidateFlightPlanRequest(req *types.FlightPlanRequest) error {
	if req.Origin == "" {
		return fmt.Errorf("origin airport (orig) is required")
	}
	if req.Destination == "" {
		return fmt.Errorf("destination airport (dest) is required")
	}
	if req.Aircraft == "" {
		return fmt.Errorf("aircraft type (type) is required")
	}

	// Validate ICAO airport codes format (basic validation)
	if len(req.Origin) != 4 {
		return fmt.Errorf("origin airport code must be 4 characters (ICAO format)")
	}
	if len(req.Destination) != 4 {
		return fmt.Errorf("destination airport code must be 4 characters (ICAO format)")
	}

	// Validate departure time if provided
	if req.DepartureHour != 0 && (req.DepartureHour < 0 || req.DepartureHour > 23) {
		return fmt.Errorf("departure hour must be between 0 and 23")
	}
	if req.DepartureMinute != 0 && (req.DepartureMinute < 0 || req.DepartureMinute > 59) {
		return fmt.Errorf("departure minute must be between 0 and 59")
	}

	return nil
}

// fetchFlightPlan is a helper method to fetch flight plan data
func (c *Client) fetchFlightPlan(req *types.FetchRequest) (*types.FlightPlanResponse, error) {
	fullURL := c.BaseURL + endpointXMLFetcher + req.ToQueryParams()

	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse error
		if req.JSON {
			var apiErr types.APIError
			if err := json.Unmarshal(body, &apiErr); err == nil {
				return nil, apiErr
			}
		} else {
			var apiErr types.APIError
			if err := xml.Unmarshal(body, &apiErr); err == nil {
				return nil, apiErr
			}
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var flightPlan types.FlightPlanResponse

	if req.JSON {
		if err := json.Unmarshal(body, &flightPlan); err != nil {
			return nil, fmt.Errorf("failed to decode JSON response: %w", err)
		}
	} else {
		if err := xml.Unmarshal(body, &flightPlan); err != nil {
			return nil, fmt.Errorf("failed to decode XML response: %w", err)
		}
	}

	return &flightPlan, nil
}

// GetDirectEditURL generates a URL to edit a specific flight plan on SimBrief website
func (c *Client) GetDirectEditURL(staticID string) string {
	return fmt.Sprintf("%s/system/dispatch.php?editflight=last&static_id=%s", c.BaseURL, url.QueryEscape(staticID))
}

// GetAircraftTypes retrieves just the aircraft types from supported options
func (c *Client) GetAircraftTypes() (map[string]types.AircraftOption, error) {
	options, err := c.GetSupportedOptions()
	if err != nil {
		return nil, err
	}
	return options.Aircraft, nil
}

// GetPlanFormats retrieves just the plan formats from supported options
func (c *Client) GetPlanFormats() (map[string]types.LayoutOption, error) {
	options, err := c.GetSupportedOptions()
	if err != nil {
		return nil, err
	}
	return options.Layouts, nil
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}

// SetUserAgent sets a custom User-Agent header for requests
func (c *Client) SetUserAgent(userAgent string) {
	// Create a custom transport that adds the User-Agent header
	originalTransport := c.HTTPClient.Transport
	if originalTransport == nil {
		originalTransport = http.DefaultTransport
	}

	c.HTTPClient.Transport = &userAgentTransport{
		Transport: originalTransport,
		UserAgent: userAgent,
	}
}

// userAgentTransport is a custom transport that adds User-Agent header
type userAgentTransport struct {
	Transport http.RoundTripper
	UserAgent string
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	return t.Transport.RoundTrip(req)
}
