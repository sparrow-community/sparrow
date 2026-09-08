# Implementation Plan: Multiple Same-Kind Boundaries

**Branch**: `010-multiple-same-kind-boundaries` | **Date**: 2026-09-08 | **Spec**: [spec.md](./spec.md)

## Summary

Lift the one-per-kind deploy reject for timer/message/signal boundaries. Arm all attached boundaries on activity enter via a repeated `waiting_boundaries` ActivityPayload (legacy singular fields keep the first of each kind). Collectors and interrupting Complete terminate all siblings. Unique message/signal names enforced at deploy; compensation remains one-per-activity.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `deploy/deploy.go`, `handlers/boundary_event.go`, `projection/instance.go`, `timers.go`/`messages.go`/`signals.go`, `protocol/proto/event/v1/payloads.proto`  
**Testing**: `m14_multi_*_boundary*.bpmn`; regression Boundary/Signal suites  
**Constraints**: proto via `build.sh`; no new Element.Type  

## Constitution Check

Pass — additive payload field; Element subjects unchanged.

## Project Structure

```text
protocol/proto/event/v1/payloads.proto
processing/deploy/{deploy,timer}.go
processing/handlers/boundary_event.go
processing/projection/instance.go
processing/{timers,messages,signals,engine}.go
processing/testdata/m14_multi_*.bpmn
```

## Complexity Tracking

None beyond justified repeated waiting-boundary payload.
