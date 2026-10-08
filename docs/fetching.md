---
title: Fetching a Plan
description: Fetch the latest OFP of a SimBrief user, or a specific plan by static ID, as JSON v2 or XML.
order: 2
section: guides
---

`GetFlightPlan` fetches the latest OFP of a user (or a specific one by static ID) and decodes it into `types.FlightPlanResponse`.

```go
// Latest OFP of a user, as JSON v2 (JSON: false fetches XML instead).
ofp, err := c.GetFlightPlan(&types.FetchRequest{Username: "pilot", JSON: true})
if err != nil {
	return err
}
fmt.Println(ofp.Origin.ICAO, "->", ofp.Destination.ICAO)
```

## FetchRequest

| Field | Query parameter | Meaning |
|-------|-----------------|---------|
| `UserID` | `userid` | SimBrief user ID (numeric, as a string) |
| `Username` | `username` | SimBrief username |
| `StaticID` | `static_id` | the static ID the plan was generated with |
| `JSON` | `json=v2` | `true` for JSON v2, `false` for XML |

JSON v2 and XML decode into the same struct; the only differences are the time formats (see [Times](ofp.md#times)).

## Shortcuts

All of these fetch JSON v2:

```go
_, _ = c.GetFlightPlanByUserID("123456")
_, _ = c.GetFlightPlanByUsername("pilot")
_, _ = c.GetFlightPlanByStaticID("123456", "MYAPP_42") // a plan generated with StaticID
```

## Raw XML

The raw XML, undecoded (`GetFlightPlanXML` sets `req.JSON` to `false`):

```go
raw, err := c.GetFlightPlanXML(&types.FetchRequest{UserID: "123456"})
if err != nil {
	return err
}
fmt.Println(len(raw), "bytes of <OFP> XML")
```

A fetch that SimBrief answers with an error status (unknown user, no plan generated yet) returns a `types.APIError`; see [Errors](errors.md).
