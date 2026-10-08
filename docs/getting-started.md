---
title: Getting Started
description: Install the simbrief Go library, create a client and fetch your first OFP.
order: 1
section: getting-started
---

`github.com/mrlm-net/simbrief` is a Go library for [SimBrief](https://www.simbrief.com): it fetches and decodes OFPs (JSON v2 and XML), reads the take-off and landing speeds of the runway analysis, and builds the dispatch URLs that create new plans. It is standard library only and works on any platform Go runs on.

No API key is needed: fetching a plan uses SimBrief's public `/api/xml.fetcher.php` endpoint with a SimBrief user ID or username.

## Requirements

- Go 1.21 or later
- An internet connection (for the SimBrief API)
- A web browser logged in to SimBrief, to generate new plans

## Installation

```bash
go get github.com/mrlm-net/simbrief
```

The library has two packages:

| Package | What it holds |
|---------|---------------|
| `pkg/client` | `Client` (fetching, supported options, URLs), `FlightPlanBuilder`, the route, fuel and time helpers |
| `pkg/types` | `FlightPlanRequest` and `FetchRequest`, `FlightPlanResponse` (the OFP, navlog, runway analysis), `Number`, `APIError`, sentinel errors, enums |

## Imports

Code blocks in these docs are function bodies. Unless a block creates them itself, `c` is a `*client.Client`, `ofp` is a `*types.FlightPlanResponse`, and the surrounding function returns `error`. Imports:

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

## Client

```go
c := client.NewClient() // https://www.simbrief.com, 30 s timeout
c.SetTimeout(10 * time.Second)
c.SetUserAgent("my-app/1.0")

// Custom base URL and http.Client, e.g. for an httptest server:
c = client.NewClientWithConfig("http://127.0.0.1:8080", &http.Client{Timeout: 5 * time.Second})
fmt.Println(c.BaseURL)
```

`NewClient()` takes no arguments. `NewClientWithConfig` falls back to the defaults (`client.DefaultBaseURL`, `client.DefaultTimeout`) for an empty URL or a nil `*http.Client`. `SetUserAgent` wraps the client's transport, so call it after any custom transport is in place.

## Your first OFP

Fetch the latest plan a user generated on SimBrief:

```go
c := client.NewClient()
ofp, err := c.GetFlightPlan(&types.FetchRequest{UserID: "123456", JSON: true})
if err != nil {
	return err
}
fmt.Println(ofp.ATC.Callsign, ofp.Origin.ICAO, "->", ofp.Destination.ICAO)
```

Your user ID is in your SimBrief account settings. The fetch returns the latest plan generated for that user, so generate one first.

## Next

- [Fetching a plan](fetching.md): the fetch options and shortcuts
- [The OFP](ofp.md): what the decoded plan holds
- [Runway analysis](runway-analysis.md): V1, VR, V2 and Vref
- [Generating a plan](generating-plans.md): dispatch URLs for new plans
