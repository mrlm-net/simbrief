# Examples

Two runnable programs using the library. Both need an internet connection.

| Example | What it does | Needs `SIMBRIEF_USER_ID` |
|---------|--------------|--------------------------|
| [basic](basic/) | Lists aircraft types and layouts, builds and validates a plan request, prints its dispatch URL, then fetches the user's latest OFP | Optional (stops before the fetch without it) |
| [advanced](advanced/) | Builds a plan URL with custom aircraft data, then fetches the user's latest OFP and prints ATC, aircraft, route, times, weights, fuel, weather, alternates, navlog and file links | Required (exits at the start without it) |

## Running

```bash
export SIMBRIEF_USER_ID=your_user_id   # numeric ID from your SimBrief account settings

cd examples/basic && go run .
cd examples/advanced && go run .
```

`SIMBRIEF_USER_ID` is the only environment variable the examples read. The
fetch returns the latest plan generated on SimBrief for that user, so
generate one first.

## Troubleshooting

| Message | Cause |
|---------|-------|
| `Set SIMBRIEF_USER_ID environment variable ...` | the variable is not set |
| `Failed to fetch flight plan: ...` | SimBrief's fetch status (unknown user, no plan yet) or a network error |
| `Failed to get supported options: ...` | the aircraft list request failed (basic) |

## Further reading

- [Usage guide](../docs/usage.md) - the whole API
- [Main README](../README.md)
- [pkg/client](../pkg/client/) and [pkg/types](../pkg/types/) - source
