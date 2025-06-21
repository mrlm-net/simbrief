# API Integration

This document covers best practices for integrating with the SimBrief API using the Go SDK.

## Table of Contents

- [Client Configuration](#client-configuration)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Caching Strategies](#caching-strategies)
- [Timeout Management](#timeout-management)
- [Retry Logic](#retry-logic)
- [Testing](#testing)

## Client Configuration

### Basic Client Setup

```go
import (
    "github.com/mrlm-net/simbrief/pkg/client"
    "github.com/mrlm-net/simbrief/pkg/types"
)

// Simple client creation
simbrief := client.NewClient("your-api-key")
```

### Advanced Client Configuration

```go
import (
    "net/http"
    "time"
)

// Custom HTTP client with specific timeouts
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        10,
        MaxIdleConnsPerHost: 2,
        IdleConnTimeout:     30 * time.Second,
    },
}

// Create client with custom configuration
simbrief := client.NewClientWithConfig(
    "your-api-key",
    "https://www.simbrief.com",  // Custom base URL if needed
    httpClient,
)
```

### Environment-Based Configuration

```go
import (
    "os"
    "strconv"
    "time"
)

func createClient() *client.Client {
    apiKey := os.Getenv("SIMBRIEF_API_KEY")
    if apiKey == "" {
        log.Fatal("SIMBRIEF_API_KEY environment variable required")
    }
    
    // Parse timeout from environment
    timeoutStr := os.Getenv("SIMBRIEF_TIMEOUT")
    timeout := 30 * time.Second
    if timeoutStr != "" {
        if t, err := time.ParseDuration(timeoutStr); err == nil {
            timeout = t
        }
    }
    
    httpClient := &http.Client{Timeout: timeout}
    return client.NewClientWithConfig(apiKey, "", httpClient)
}
```

## Authentication

### API Key Management

```go
type Config struct {
    APIKey  string
    UserID  string
    PilotID string
}

func LoadConfig() (*Config, error) {
    return &Config{
        APIKey:  os.Getenv("SIMBRIEF_API_KEY"),
        UserID:  os.Getenv("SIMBRIEF_USER_ID"),
        PilotID: os.Getenv("SIMBRIEF_PILOT_ID"),
    }, nil
}

// Validate configuration
func (c *Config) Validate() error {
    if c.APIKey == "" {
        return errors.New("API key is required")
    }
    return nil
}
```

### Secure Key Storage

```go
// Example using a secrets manager
func getAPIKeyFromVault() (string, error) {
    // Implementation depends on your secrets management system
    // This is just an example structure
    return secretsManager.GetSecret("simbrief-api-key")
}

func createSecureClient() (*client.Client, error) {
    apiKey, err := getAPIKeyFromVault()
    if err != nil {
        return nil, fmt.Errorf("failed to get API key: %w", err)
    }
    
    return client.NewClient(apiKey), nil
}
```

## Error Handling

### Structured Error Handling

```go
import (
    "errors"
    "fmt"
    "net/http"
)

func handleSimBriefError(err error) {
    if err == nil {
        return
    }
    
    // Check for specific SimBrief errors
    var sbErr *types.SimBriefError
    if errors.As(err, &sbErr) {
        switch sbErr.Code {
        case types.ErrCodeInvalidAPIKey:
            log.Printf("Authentication failed: %s", sbErr.Message)
            // Handle authentication error
            
        case types.ErrCodeRateLimit:
            log.Printf("Rate limit exceeded: %s", sbErr.Message)
            // Implement backoff strategy
            
        case types.ErrCodeInvalidRoute:
            log.Printf("Invalid route: %s", sbErr.Message)
            // Handle route validation error
            
        default:
            log.Printf("SimBrief API error: %s", sbErr.Message)
        }
        return
    }
    
    // Check for HTTP errors
    var httpErr *types.HTTPError
    if errors.As(err, &httpErr) {
        switch httpErr.StatusCode {
        case http.StatusTooManyRequests:
            log.Printf("Rate limited, retry after: %s", httpErr.RetryAfter)
            
        case http.StatusBadRequest:
            log.Printf("Bad request: %s", httpErr.Message)
            
        case http.StatusUnauthorized:
            log.Printf("Unauthorized: check API key")
            
        case http.StatusInternalServerError:
            log.Printf("Server error: %s", httpErr.Message)
            
        default:
            log.Printf("HTTP error %d: %s", httpErr.StatusCode, httpErr.Message)
        }
        return
    }
    
    // Handle other errors
    log.Printf("Unexpected error: %v", err)
}
```

### Error Recovery

```go
func generateFlightPlanWithRecovery(client *client.Client, request *types.FlightPlanRequest) (*types.FlightPlanResponse, error) {
    flightPlan, err := client.GenerateFlightPlan(request)
    if err != nil {
        // Try to recover from certain errors
        var sbErr *types.SimBriefError
        if errors.As(err, &sbErr) {
            switch sbErr.Code {
            case types.ErrCodeInvalidRoute:
                // Try with direct routing
                log.Println("Invalid route, trying direct routing")
                request.Route = "DCT"
                return client.GenerateFlightPlan(request)
                
            case types.ErrCodeInvalidAircraft:
                // Try with a default aircraft
                log.Println("Invalid aircraft, trying with B38M")
                request.Aircraft = "B38M"
                return client.GenerateFlightPlan(request)
            }
        }
        return nil, err
    }
    return flightPlan, nil
}
```

## Rate Limiting

### Rate Limiter Implementation

```go
import (
    "golang.org/x/time/rate"
    "context"
    "time"
)

type RateLimitedClient struct {
    client  *client.Client
    limiter *rate.Limiter
}

func NewRateLimitedClient(apiKey string, requestsPerSecond float64) *RateLimitedClient {
    return &RateLimitedClient{
        client:  client.NewClient(apiKey),
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
    }
}

func (rlc *RateLimitedClient) GenerateFlightPlan(ctx context.Context, request *types.FlightPlanRequest) (*types.FlightPlanResponse, error) {
    // Wait for rate limiter
    if err := rlc.limiter.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limiter error: %w", err)
    }
    
    return rlc.client.GenerateFlightPlan(request)
}
```

### Adaptive Rate Limiting

```go
type AdaptiveRateLimiter struct {
    client       *client.Client
    limiter      *rate.Limiter
    lastError    time.Time
    errorCount   int
    successCount int
}

func (arl *AdaptiveRateLimiter) GenerateFlightPlan(ctx context.Context, request *types.FlightPlanRequest) (*types.FlightPlanResponse, error) {
    // Adjust rate based on recent errors
    if time.Since(arl.lastError) < 5*time.Minute && arl.errorCount > 3 {
        // Slow down if we've had recent errors
        arl.limiter.SetLimit(rate.Limit(0.5)) // 1 request per 2 seconds
    } else if arl.successCount > 10 {
        // Speed up if we've had recent successes
        arl.limiter.SetLimit(rate.Limit(2.0)) // 2 requests per second
    }
    
    if err := arl.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    
    flightPlan, err := arl.client.GenerateFlightPlan(request)
    if err != nil {
        arl.errorCount++
        arl.lastError = time.Now()
        arl.successCount = 0
    } else {
        arl.successCount++
        arl.errorCount = 0
    }
    
    return flightPlan, err
}
```

## Caching Strategies

### Simple In-Memory Cache

```go
import (
    "sync"
    "time"
)

type CachedClient struct {
    client *client.Client
    cache  sync.Map
}

type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

func NewCachedClient(apiKey string) *CachedClient {
    return &CachedClient{
        client: client.NewClient(apiKey),
    }
}

func (cc *CachedClient) GetSupportedOptions() (*types.SupportedOptions, error) {
    cacheKey := "supported_options"
    
    // Check cache first
    if entry, ok := cc.cache.Load(cacheKey); ok {
        cacheEntry := entry.(CacheEntry)
        if time.Now().Before(cacheEntry.ExpiresAt) {
            return cacheEntry.Data.(*types.SupportedOptions), nil
        }
        // Cache expired, remove it
        cc.cache.Delete(cacheKey)
    }
    
    // Fetch from API
    options, err := cc.client.GetSupportedOptions()
    if err != nil {
        return nil, err
    }
    
    // Store in cache (valid for 1 hour)
    cc.cache.Store(cacheKey, CacheEntry{
        Data:      options,
        ExpiresAt: time.Now().Add(1 * time.Hour),
    })
    
    return options, nil
}
```

### Redis Cache Implementation

```go
import (
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "context"
)

type RedisCachedClient struct {
    client     *client.Client
    redis      *redis.Client
    expiration time.Duration
}

func NewRedisCachedClient(apiKey, redisAddr string) *RedisCachedClient {
    rdb := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    
    return &RedisCachedClient{
        client:     client.NewClient(apiKey),
        redis:      rdb,
        expiration: 1 * time.Hour,
    }
}

func (rcc *RedisCachedClient) GetSupportedOptions(ctx context.Context) (*types.SupportedOptions, error) {
    cacheKey := "simbrief:supported_options"
    
    // Try to get from cache
    cached, err := rcc.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var options types.SupportedOptions
        if err := json.Unmarshal([]byte(cached), &options); err == nil {
            return &options, nil
        }
    }
    
    // Fetch from API
    options, err := rcc.client.GetSupportedOptions()
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    data, _ := json.Marshal(options)
    rcc.redis.Set(ctx, cacheKey, data, rcc.expiration)
    
    return options, nil
}
```

## Timeout Management

### Context-Based Timeouts

```go
func generateFlightPlanWithTimeout(client *client.Client, request *types.FlightPlanRequest, timeout time.Duration) (*types.FlightPlanResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    // Use a channel to handle the API call
    type result struct {
        flightPlan *types.FlightPlanResponse
        err        error
    }
    
    resultChan := make(chan result, 1)
    
    go func() {
        flightPlan, err := client.GenerateFlightPlan(request)
        resultChan <- result{flightPlan: flightPlan, err: err}
    }()
    
    select {
    case res := <-resultChan:
        return res.flightPlan, res.err
    case <-ctx.Done():
        return nil, fmt.Errorf("operation timed out after %v: %w", timeout, ctx.Err())
    }
}
```

### Progressive Timeouts

```go
type TimeoutStrategy struct {
    Initial   time.Duration
    Maximum   time.Duration
    Increment time.Duration
    attempts  int
}

func (ts *TimeoutStrategy) NextTimeout() time.Duration {
    timeout := ts.Initial + time.Duration(ts.attempts)*ts.Increment
    if timeout > ts.Maximum {
        timeout = ts.Maximum
    }
    ts.attempts++
    return timeout
}

func (ts *TimeoutStrategy) Reset() {
    ts.attempts = 0
}
```

## Retry Logic

### Exponential Backoff

```go
import (
    "math"
    "math/rand"
    "time"
)

type RetryConfig struct {
    MaxAttempts int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
}

func retryWithBackoff(config RetryConfig, operation func() error) error {
    var lastErr error
    
    for attempt := 0; attempt < config.MaxAttempts; attempt++ {
        if attempt > 0 {
            delay := calculateBackoffDelay(config, attempt)
            time.Sleep(delay)
        }
        
        if err := operation(); err != nil {
            lastErr = err
            
            // Don't retry certain errors
            var sbErr *types.SimBriefError
            if errors.As(err, &sbErr) {
                switch sbErr.Code {
                case types.ErrCodeInvalidAPIKey:
                    // Don't retry authentication errors
                    return err
                case types.ErrCodeInvalidRoute:
                    // Don't retry validation errors
                    return err
                }
            }
            
            continue
        }
        
        return nil
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

func calculateBackoffDelay(config RetryConfig, attempt int) time.Duration {
    delay := float64(config.BaseDelay) * math.Pow(config.Multiplier, float64(attempt-1))
    
    // Add jitter to prevent thundering herd
    jitter := rand.Float64() * 0.1 * delay
    delay += jitter
    
    if delay > float64(config.MaxDelay) {
        delay = float64(config.MaxDelay)
    }
    
    return time.Duration(delay)
}
```

### Retry with Circuit Breaker

```go
type CircuitBreaker struct {
    maxFailures  int
    timeout      time.Duration
    failures     int
    lastFailTime time.Time
    state        string // "closed", "open", "half-open"
    mutex        sync.Mutex
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       "closed",
    }
}

func (cb *CircuitBreaker) Call(operation func() error) error {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    switch cb.state {
    case "open":
        if time.Since(cb.lastFailTime) > cb.timeout {
            cb.state = "half-open"
            cb.failures = 0
        } else {
            return errors.New("circuit breaker is open")
        }
    }
    
    err := operation()
    
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }
    
    // Success
    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

## Testing

### Mock Client for Testing

```go
type MockClient struct {
    SupportedOptions *types.SupportedOptions
    FlightPlan       *types.FlightPlanResponse
    Errors           map[string]error
}

func (mc *MockClient) GetSupportedOptions() (*types.SupportedOptions, error) {
    if err := mc.Errors["GetSupportedOptions"]; err != nil {
        return nil, err
    }
    return mc.SupportedOptions, nil
}

func (mc *MockClient) GenerateFlightPlan(request *types.FlightPlanRequest) (*types.FlightPlanResponse, error) {
    if err := mc.Errors["GenerateFlightPlan"]; err != nil {
        return nil, err
    }
    return mc.FlightPlan, nil
}

// Test example
func TestFlightPlanGeneration(t *testing.T) {
    mockClient := &MockClient{
        FlightPlan: &types.FlightPlanResponse{
            General: struct {
                DistanceNM float64 `xml:"distance"`
            }{
                DistanceNM: 2475.0,
            },
        },
    }
    
    request := &types.FlightPlanRequest{
        Origin:      "KJFK",
        Destination: "KLAX",
        Aircraft:    "B38M",
    }
    
    flightPlan, err := mockClient.GenerateFlightPlan(request)
    assert.NoError(t, err)
    assert.Equal(t, 2475.0, flightPlan.General.DistanceNM)
}
```

### Integration Testing

```go
func TestRealAPIIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    apiKey := os.Getenv("SIMBRIEF_API_KEY")
    if apiKey == "" {
        t.Skip("SIMBRIEF_API_KEY not set")
    }
    
    client := client.NewClient(apiKey)
    
    // Test supported options
    options, err := client.GetSupportedOptions()
    require.NoError(t, err)
    assert.NotEmpty(t, options.Aircraft)
    
    // Test flight plan generation
    request := &types.FlightPlanRequest{
        Origin:      "KJFK",
        Destination: "KLAX",
        Aircraft:    "B38M",
    }
    
    flightPlan, err := client.GenerateFlightPlan(request)
    require.NoError(t, err)
    assert.NotEmpty(t, flightPlan.General.Route)
}
```

## Best Practices

1. **Always use timeouts**: Set appropriate timeouts for all API calls
2. **Implement retry logic**: Use exponential backoff for transient failures
3. **Cache frequently accessed data**: Cache aircraft types and other static data
4. **Handle rate limits gracefully**: Implement proper rate limiting and backoff
5. **Use structured error handling**: Handle specific error types appropriately
6. **Monitor API usage**: Track API calls and response times
7. **Test with mocks**: Use mock clients for unit testing
8. **Secure API keys**: Never hardcode API keys, use environment variables or secure storage
