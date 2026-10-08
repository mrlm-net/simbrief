---
title: Errors
description: APIError, wrapped transport and decoding errors, the validation sentinels, and common problems.
order: 7
section: reference
---

- **`types.APIError`** (`Message`, `Code`): SimBrief answered but the fetch failed. `Message` is the OFP's fetch status; `Code` is the HTTP status, or 0 when the error came with a 200 reply.
- **Wrapped errors** (`fmt.Errorf` with `%w`): transport and decoding failures, e.g. a timeout or a body that is not an OFP. A non-200 reply without a fetch status is a plain error with the status and body.
- **Sentinels** from `FlightPlanRequest.Validate()`: `types.ErrMissingOrigin`, `types.ErrMissingDestination`, `types.ErrMissingAircraft`. (`ErrMissingUserID`, `ErrInvalidUserID` and `ErrInvalidAPIKey` are declared but not returned by the library.)
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

## Common problems

| Problem | Cause |
|---------|-------|
| Fetch fails with `types.APIError` | the message is SimBrief's fetch status (unknown user, no plan generated yet) |
| Invalid aircraft code | list valid codes with `GetAircraftTypes()` |
| Empty `TLR` | the plan was generated without runway analysis, or the aircraft has no `TLRData` |

## Further reading

- [SimBrief API documentation](https://developers.navigraph.com/docs/simbrief/using-the-api)
