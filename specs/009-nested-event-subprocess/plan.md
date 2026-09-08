# Implementation Plan: Nested Event Sub-Process

**Branch**: `009-nested-event-subprocess` | **Date**: 2026-09-08 | **Spec**: [spec.md](./spec.md)

## Summary

Allow Event Sub-Process nested inside another Event Sub-Process. Remove deploy reject; when entering any SubProcess scope (including Event Sub-Process), arm child Event Sub-Processes for that scope. Existing interrupt/disarm/complete paths already key off ParentScopeID.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `processing/deploy/deploy.go` (`validateSubProcessesAt`), `processing/executor.go` (arm on EnterChild), `processing/event_subprocess.go`  
**Storage**: EventLog; existing EventSubProcesses arms  
**Testing**: `go test ./processing/ -run 'EventSubProcess|Nested'`; fixtures `m13_esp_in_esp_*.bpmn`  
**Constraints**: no new Element.Type; no proto changes  

## Constitution Check

All principles Pass — deploy + arm wiring only; ledger subjects unchanged.

## Project Structure

```text
processing/
├── deploy/deploy.go
├── executor.go
├── event_subprocess_test.go
└── testdata/m13_esp_in_esp_*.bpmn
```

## Complexity Tracking

None.
