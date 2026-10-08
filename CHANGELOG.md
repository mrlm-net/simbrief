# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Docs rewritten to the real API: `docs/usage.md` replaces
  `flight-planning.md`, `api-integration.md`, `performance.md` and
  `aircraft-management.md`, which described functions and types the library
  does not have. Every snippet in `docs/usage.md` compiles.
- README and example READMEs match the code; the unused `SIMBRIEF_DEBUG`
  variable is no longer mentioned.

### Added

- `CHANGELOG.md` and `CLAUDE.md`.

## [0.2.0] - 2026-10-07

### Added

- Runway analysis: `FlightPlanResponse.TLR` decodes the `<tlr>` block from
  JSON v2 and XML (`TLRTakeoff`, `TLRLanding`, `TLRConditions`, `TLRRunway`,
  `TakeoffRunway`, `LandingRunway`, `TLRLandingDistance`).
- `TLR.TakeoffRunway(ident)` finds a runway's take-off data (V1, VR, V2,
  flaps, thrust, flex temperature), tolerant of case and leading zeros.
- `TLR.LandingVref(wet)` returns the landing Vref for a dry or wet runway.

## [0.1.0] - 2026-10-03

### Added

- `client.NewClient()` / `NewClientWithConfig`, `SetTimeout`,
  `SetUserAgent`.
- Fetching: `GetFlightPlan(*types.FetchRequest)` (JSON v2 or XML),
  `GetFlightPlanByUserID`, `GetFlightPlanByUsername`,
  `GetFlightPlanByStaticID`, `GetFlightPlanXML`.
- `types.FlightPlanResponse` decoding real OFPs in both formats (root
  `<OFP>`): `FetchInfo`, `ATCInfo` (callsign, rules, type, initial
  speed/level with units, routes, FPL text), `Alternates` with
  `Alternate()`, `NavLog`/`NavLogFix`, `AlternateNavLogs`, fuel, weights,
  times, weather, files.
- `types.Number`: numeric values kept as SimBrief's raw text with `Int`,
  `Float`, `Bool`, `IsSet`.
- Fetch status errors returned as `types.APIError`.
- Plan generation URLs: `FlightPlanBuilder` (`client.NewFlightPlan`),
  `types.FlightPlanRequest`, `ValidateFlightPlanRequest`,
  `GenerateFlightPlanURL`, `GetDirectEditURL`.
- `GetSupportedOptions`, `GetAircraftTypes`, `GetPlanFormats`.
- `RouteHelper`, `FuelHelper`, `TimeHelper`.
- Golden tests on a real OFP (CEF007 LKPR-LKPD) in JSON v2 and XML.
- Basic and advanced examples.

[Unreleased]: https://github.com/mrlm-net/simbrief/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/mrlm-net/simbrief/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/mrlm-net/simbrief/releases/tag/v0.1.0
