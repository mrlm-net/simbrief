# simbrief

Go library `github.com/mrlm-net/simbrief`: fetches and decodes SimBrief OFPs
(JSON v2 and XML) and builds SimBrief dispatch URLs. No API key; fetching uses
the public `/api/xml.fetcher.php`.

## Layout

- `pkg/client/` - `Client` (fetch, supported options, URLs), `FlightPlanBuilder`, Route/Fuel/Time helpers
- `pkg/types/` - `FlightPlanRequest`/`FetchRequest`, `FlightPlanResponse` (OFP, NavLog, TLR), `Number`, `APIError`, sentinel errors, enums
- `pkg/types/testdata/` - real OFP (CEF007 LKPR-LKPD) as JSON v2 and XML; golden tests decode both
- `examples/basic`, `examples/advanced` - runnable programs (`SIMBRIEF_USER_ID`)
- `docs/*.md` - the API guide, one page per topic, also the website (`website/`, SvelteKit; `npm run build` there); `CHANGELOG.md`

## Consumers

MyCrew app (`mycrew-online/app`, `internal/simbrief/simbrief.go`), pinned to
v0.2.0. It uses `NewClient().GetFlightPlan(&types.FetchRequest{...JSON: true})`
and reads OFP fields plus `TLR.TakeoffRunway` / `TLR.LandingVref`. Keep those
stable; breaking them needs a coordinated app bump.

## Build and test

```bash
go build ./... && go vet ./... && go test ./...
```

## Conventions

- Stdlib only. Never add dependencies (testify is already there for tests only).
- Docs describe the real API only. Every snippet in `docs/*.md` must
  compile: check it in a throwaway module with a `replace` to this repo.
- OFP numbers are `types.Number`, times stay strings (formats differ by JSON/XML).
- Both formats must decode into the same struct; add golden test values for new fields.

## Release

Features = minor, fixes = patch. Add a CHANGELOG entry (move Unreleased under
the new version with the date). Squash-merge the PR, verify the squash commit
on `main`, then tag that commit `vX.Y.Z` and push the tag.
