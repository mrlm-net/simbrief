# Performance

This document covers performance optimization strategies for the SimBrief SDK.

## Table of Contents

- [Client Performance](#client-performance)
- [Memory Management](#memory-management)
- [Network Optimization](#network-optimization)
- [Caching Strategies](#caching-strategies)
- [Concurrent Operations](#concurrent-operations)
- [Profiling and Monitoring](#profiling-and-monitoring)
- [Best Practices](#best-practices)

## Client Performance

### Connection Pooling

```go
import (
    "net/http"
    "time"
)

// Optimized HTTP client configuration
func createOptimizedClient(apiKey string) *client.Client {
    transport := &http.Transport{
        MaxIdleConns:        100,              // Maximum idle connections
        MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
        IdleConnTimeout:     90 * time.Second, // How long idle connections last
        TLSHandshakeTimeout: 10 * time.Second, // TLS handshake timeout
        ExpectContinueTimeout: 1 * time.Second, // Expect: 100-continue timeout
        
        // Enable HTTP/2
        ForceAttemptHTTP2: true,
    }
    
    httpClient := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
    
    return client.NewClientWithConfig(apiKey, "", httpClient)
}
```

### Request Optimization

```go
// Batch multiple operations together
type BatchProcessor struct {
    client *client.Client
    queue  chan *types.FlightPlanRequest
    batch  []*types.FlightPlanRequest
    size   int
}

func NewBatchProcessor(client *client.Client, batchSize int) *BatchProcessor {
    bp := &BatchProcessor{
        client: client,
        queue:  make(chan *types.FlightPlanRequest, batchSize*2),
        batch:  make([]*types.FlightPlanRequest, 0, batchSize),
        size:   batchSize,
    }
    
    go bp.processBatch()
    return bp
}

func (bp *BatchProcessor) processBatch() {
    ticker := time.NewTicker(5 * time.Second) // Process every 5 seconds
    defer ticker.Stop()
    
    for {
        select {
        case request := <-bp.queue:
            bp.batch = append(bp.batch, request)
            if len(bp.batch) >= bp.size {
                bp.flushBatch()
            }
            
        case <-ticker.C:
            if len(bp.batch) > 0 {
                bp.flushBatch()
            }
        }
    }
}

func (bp *BatchProcessor) flushBatch() {
    // Process all requests in the batch concurrently
    var wg sync.WaitGroup
    for _, request := range bp.batch {
        wg.Add(1)
        go func(req *types.FlightPlanRequest) {
            defer wg.Done()
            _, err := bp.client.GenerateFlightPlan(req)
            if err != nil {
                log.Printf("Error processing flight plan: %v", err)
            }
        }(request)
    }
    wg.Wait()
    
    // Clear the batch
    bp.batch = bp.batch[:0]
}
```

## Memory Management

### Object Pooling

```go
import (
    "sync"
)

// Pool for flight plan requests to reduce allocations
var requestPool = sync.Pool{
    New: func() interface{} {
        return &types.FlightPlanRequest{}
    },
}

func GetFlightPlanRequest() *types.FlightPlanRequest {
    return requestPool.Get().(*types.FlightPlanRequest)
}

func PutFlightPlanRequest(req *types.FlightPlanRequest) {
    // Reset the request before putting it back
    *req = types.FlightPlanRequest{}
    requestPool.Put(req)
}

// Usage example
func createOptimizedFlightPlan(origin, destination, aircraft string) (*types.FlightPlanResponse, error) {
    request := GetFlightPlanRequest()
    defer PutFlightPlanRequest(request)
    
    request.Origin = origin
    request.Destination = destination
    request.Aircraft = aircraft
    
    return client.GenerateFlightPlan(request)
}
```

### Memory-Efficient Response Handling

```go
// Stream large responses instead of loading everything into memory
func streamFlightPlanData(response *types.FlightPlanResponse) error {
    // Process response data in chunks instead of all at once
    if response.Navigation != nil {
        for i := 0; i < len(response.Navigation.Fixes); i += 100 {
            end := i + 100
            if end > len(response.Navigation.Fixes) {
                end = len(response.Navigation.Fixes)
            }
            
            chunk := response.Navigation.Fixes[i:end]
            if err := processFixChunk(chunk); err != nil {
                return err
            }
            
            // Allow garbage collection
            runtime.GC()
        }
    }
    
    return nil
}

func processFixChunk(fixes []types.NavFix) error {
    // Process a chunk of navigation fixes
    for _, fix := range fixes {
        // Process individual fix
        _ = fix
    }
    return nil
}
```

### Garbage Collection Optimization

```go
import (
    "runtime"
    "runtime/debug"
)

// Configure garbage collection for better performance
func init() {
    // Set initial GC target percentage
    debug.SetGCPercent(100)
    
    // Set memory limit (example: 1GB)
    debug.SetMemoryLimit(1 << 30)
}

// Force garbage collection at appropriate times
func forceGCIfNeeded() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // Force GC if heap is getting large
    if m.HeapAlloc > 500*1024*1024 { // 500MB
        runtime.GC()
    }
}
```

## Network Optimization

### Compression

```go
import (
    "compress/gzip"
    "net/http"
)

// HTTP client with compression support
func createCompressedClient(apiKey string) *client.Client {
    transport := &http.Transport{
        DisableCompression: false, // Enable compression
    }
    
    // Custom round tripper to add compression headers
    rt := &compressionRoundTripper{
        Transport: transport,
    }
    
    httpClient := &http.Client{
        Transport: rt,
        Timeout:   30 * time.Second,
    }
    
    return client.NewClientWithConfig(apiKey, "", httpClient)
}

type compressionRoundTripper struct {
    Transport http.RoundTripper
}

func (crt *compressionRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    // Add compression headers
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    
    resp, err := crt.Transport.RoundTrip(req)
    if err != nil {
        return nil, err
    }
    
    // Handle compressed response
    if resp.Header.Get("Content-Encoding") == "gzip" {
        resp.Body = &gzipReadCloser{resp.Body}
    }
    
    return resp, nil
}

type gzipReadCloser struct {
    io.ReadCloser
}

func (grc *gzipReadCloser) Read(p []byte) (n int, err error) {
    if grc.ReadCloser == nil {
        reader, err := gzip.NewReader(grc.ReadCloser)
        if err != nil {
            return 0, err
        }
        grc.ReadCloser = reader
    }
    return grc.ReadCloser.Read(p)
}
```

### Request Deduplication

```go
import (
    "crypto/md5"
    "fmt"
    "sync"
)

type RequestDeduplicator struct {
    client   *client.Client
    pending  map[string]chan result
    mutex    sync.Mutex
}

type result struct {
    response *types.FlightPlanResponse
    err      error
}

func NewRequestDeduplicator(client *client.Client) *RequestDeduplicator {
    return &RequestDeduplicator{
        client:  client,
        pending: make(map[string]chan result),
    }
}

func (rd *RequestDeduplicator) GenerateFlightPlan(request *types.FlightPlanRequest) (*types.FlightPlanResponse, error) {
    key := rd.generateKey(request)
    
    rd.mutex.Lock()
    
    // Check if request is already pending
    if ch, exists := rd.pending[key]; exists {
        rd.mutex.Unlock()
        
        // Wait for existing request to complete
        res := <-ch
        return res.response, res.err
    }
    
    // Create channel for this request
    ch := make(chan result, 1)
    rd.pending[key] = ch
    rd.mutex.Unlock()
    
    // Execute request
    go func() {
        defer func() {
            rd.mutex.Lock()
            delete(rd.pending, key)
            rd.mutex.Unlock()
        }()
        
        response, err := rd.client.GenerateFlightPlan(request)
        res := result{response: response, err: err}
        
        // Broadcast result to all waiters
        ch <- res
        close(ch)
    }()
    
    // Wait for result
    res := <-ch
    return res.response, res.err
}

func (rd *RequestDeduplicator) generateKey(request *types.FlightPlanRequest) string {
    data := fmt.Sprintf("%s-%s-%s-%s", 
        request.Origin, request.Destination, request.Aircraft, request.Route)
    hash := md5.Sum([]byte(data))
    return fmt.Sprintf("%x", hash)
}
```

## Caching Strategies

### Multi-Level Caching

```go
type MultiLevelCache struct {
    l1Cache sync.Map        // In-memory cache
    l2Cache *redis.Client   // Redis cache
    l3Cache *sql.DB         // Database cache
}

func NewMultiLevelCache(redisAddr, dbURL string) (*MultiLevelCache, error) {
    rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, err
    }
    
    return &MultiLevelCache{
        l2Cache: rdb,
        l3Cache: db,
    }, nil
}

func (mlc *MultiLevelCache) Get(key string) (*types.FlightPlanResponse, bool) {
    // Try L1 cache first (in-memory)
    if value, ok := mlc.l1Cache.Load(key); ok {
        if entry, ok := value.(cacheEntry); ok && time.Now().Before(entry.ExpiresAt) {
            return entry.Data.(*types.FlightPlanResponse), true
        }
        mlc.l1Cache.Delete(key)
    }
    
    // Try L2 cache (Redis)
    if data, err := mlc.l2Cache.Get(context.Background(), key).Result(); err == nil {
        var response types.FlightPlanResponse
        if json.Unmarshal([]byte(data), &response) == nil {
            // Store in L1 cache
            mlc.l1Cache.Store(key, cacheEntry{
                Data:      &response,
                ExpiresAt: time.Now().Add(5 * time.Minute),
            })
            return &response, true
        }
    }
    
    // Try L3 cache (Database)
    var data string
    err := mlc.l3Cache.QueryRow("SELECT data FROM flight_plans WHERE key = $1 AND expires_at > NOW()", key).Scan(&data)
    if err == nil {
        var response types.FlightPlanResponse
        if json.Unmarshal([]byte(data), &response) == nil {
            // Store in upper level caches
            mlc.Set(key, &response, time.Hour)
            return &response, true
        }
    }
    
    return nil, false
}

func (mlc *MultiLevelCache) Set(key string, data *types.FlightPlanResponse, ttl time.Duration) {
    // Store in all cache levels
    jsonData, _ := json.Marshal(data)
    
    // L1 cache
    mlc.l1Cache.Store(key, cacheEntry{
        Data:      data,
        ExpiresAt: time.Now().Add(ttl),
    })
    
    // L2 cache
    mlc.l2Cache.Set(context.Background(), key, jsonData, ttl)
    
    // L3 cache
    mlc.l3Cache.Exec("INSERT INTO flight_plans (key, data, expires_at) VALUES ($1, $2, $3) ON CONFLICT (key) DO UPDATE SET data = $2, expires_at = $3",
        key, string(jsonData), time.Now().Add(ttl))
}
```

### Smart Cache Warming

```go
type CacheWarmer struct {
    client *client.Client
    cache  *MultiLevelCache
}

func NewCacheWarmer(client *client.Client, cache *MultiLevelCache) *CacheWarmer {
    return &CacheWarmer{
        client: client,
        cache:  cache,
    }
}

// Warm cache with popular routes
func (cw *CacheWarmer) WarmPopularRoutes() error {
    popularRoutes := []struct {
        origin, destination, aircraft string
    }{
        {"KJFK", "KLAX", "B38M"},
        {"KLAX", "KJFK", "B38M"},
        {"EGLL", "KJFK", "B77W"},
        {"KJFK", "EGLL", "B77W"},
        // Add more popular routes
    }
    
    var wg sync.WaitGroup
    for _, route := range popularRoutes {
        wg.Add(1)
        go func(r struct{ origin, destination, aircraft string }) {
            defer wg.Done()
            
            request := &types.FlightPlanRequest{
                Origin:      r.origin,
                Destination: r.destination,
                Aircraft:    r.aircraft,
            }
            
            key := generateCacheKey(request)
            if _, exists := cw.cache.Get(key); !exists {
                if response, err := cw.client.GenerateFlightPlan(request); err == nil {
                    cw.cache.Set(key, response, 6*time.Hour)
                }
            }
        }(route)
    }
    
    wg.Wait()
    return nil
}
```

## Concurrent Operations

### Worker Pool Pattern

```go
type WorkerPool struct {
    client    *client.Client
    workers   int
    jobs      chan Job
    results   chan Result
    wg        sync.WaitGroup
}

type Job struct {
    ID      string
    Request *types.FlightPlanRequest
}

type Result struct {
    ID       string
    Response *types.FlightPlanResponse
    Error    error
}

func NewWorkerPool(client *client.Client, numWorkers int) *WorkerPool {
    return &WorkerPool{
        client:  client,
        workers: numWorkers,
        jobs:    make(chan Job, numWorkers*2),
        results: make(chan Result, numWorkers*2),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    
    for job := range wp.jobs {
        response, err := wp.client.GenerateFlightPlan(job.Request)
        wp.results <- Result{
            ID:       job.ID,
            Response: response,
            Error:    err,
        }
    }
}

func (wp *WorkerPool) Submit(id string, request *types.FlightPlanRequest) {
    wp.jobs <- Job{ID: id, Request: request}
}

func (wp *WorkerPool) Stop() {
    close(wp.jobs)
    wp.wg.Wait()
    close(wp.results)
}

func (wp *WorkerPool) Results() <-chan Result {
    return wp.results
}
```

### Parallel Request Processing

```go
import (
    "context"
    "golang.org/x/sync/errgroup"
)

func processMultipleFlightPlans(client *client.Client, requests []*types.FlightPlanRequest) ([]*types.FlightPlanResponse, error) {
    ctx := context.Background()
    g, ctx := errgroup.WithContext(ctx)
    
    responses := make([]*types.FlightPlanResponse, len(requests))
    
    // Process requests concurrently with limited concurrency
    sem := make(chan struct{}, 5) // Limit to 5 concurrent requests
    
    for i, request := range requests {
        i, request := i, request // Capture loop variables
        
        g.Go(func() error {
            sem <- struct{}{} // Acquire semaphore
            defer func() { <-sem }() // Release semaphore
            
            response, err := client.GenerateFlightPlan(request)
            if err != nil {
                return err
            }
            
            responses[i] = response
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    
    return responses, nil
}
```

## Profiling and Monitoring

### Performance Metrics

```go
import (
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "simbrief_request_duration_seconds",
            Help: "Duration of SimBrief API requests",
        },
        []string{"method", "status"},
    )
    
    requestTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "simbrief_requests_total",
            Help: "Total number of SimBrief API requests",
        },
        []string{"method", "status"},
    )
    
    cacheHits = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "simbrief_cache_hits_total",
            Help: "Total number of cache hits",
        },
        []string{"cache_level"},
    )
)

type MetricsClient struct {
    client *client.Client
}

func NewMetricsClient(client *client.Client) *MetricsClient {
    return &MetricsClient{client: client}
}

func (mc *MetricsClient) GenerateFlightPlan(request *types.FlightPlanRequest) (*types.FlightPlanResponse, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues("generate_flight_plan", "success").Observe(duration)
    }()
    
    response, err := mc.client.GenerateFlightPlan(request)
    
    status := "success"
    if err != nil {
        status = "error"
    }
    
    requestTotal.WithLabelValues("generate_flight_plan", status).Inc()
    
    return response, err
}
```

### CPU and Memory Profiling

```go
import (
    _ "net/http/pprof"
    "net/http"
)

// Enable profiling endpoint
func enableProfiling() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}

// Profile a specific operation
func profileOperation() {
    // CPU profiling
    cpuFile, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer cpuFile.Close()
    
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
    
    // Run your operation here
    performFlightPlanGeneration()
    
    // Memory profiling
    memFile, err := os.Create("mem.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer memFile.Close()
    
    runtime.GC()
    if err := pprof.WriteHeapProfile(memFile); err != nil {
        log.Fatal(err)
    }
}
```

## Best Practices

### 1. Use Connection Pooling
Configure HTTP transport with appropriate connection pool settings:

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
}
```

### 2. Implement Smart Caching
Use multi-level caching with appropriate TTLs:

```go
// Cache aircraft data for 24 hours
// Cache flight plans for 6 hours
// Cache routes for 1 hour
```

### 3. Batch Operations
Group similar operations together to reduce overhead:

```go
// Process multiple flight plans concurrently
// Use worker pools for high-throughput scenarios
```

### 4. Monitor Performance
Track key metrics:

```go
// Request latency
// Error rates
// Cache hit ratios
// Memory usage
// CPU utilization
```

### 5. Optimize Memory Usage
Use object pooling and efficient data structures:

```go
// Pool expensive objects
// Process large responses in chunks
// Use appropriate garbage collection settings
```

### 6. Handle Errors Gracefully
Implement proper error handling and recovery:

```go
// Use circuit breakers
// Implement exponential backoff
// Log performance metrics
```

### 7. Profile Regularly
Use Go's built-in profiling tools:

```bash
# CPU profiling
go tool pprof cpu.prof

# Memory profiling  
go tool pprof mem.prof

# Analyze allocations
go test -memprofile=mem.prof -bench=.
```
