# Implementation Plan: Escalation Events

**Branch**: `011-escalation-events` | **Date**: 2026-09-08 | **Spec**: [spec.md](./spec.md)

## Summary

Add BPMN escalation throw/catch: intermediate throw, escalation end, escalation boundary (interrupting + non-interrupting), and escalation Event Sub-Process. Propagation mirrors error scope walk; uncaught escalation does not terminate the instance. Ledger: `INTENT_ESCALATION_THROWN` + `escalation_code`.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `bpmn/element` (Escalation root), `processing/deploy`, `handlers`, `executor`, `errors.go` pattern, proto event/payloads  
**Storage**: EventLog; no new runtime indexes beyond deploy maps  
**Testing**: `go test ./processing/ -run Escalation`; fixtures `m15_escalation_*.bpmn`  
**Constraints**: proto via `build.sh`; Apache-2.0 headers on new Go files  

## Constitution Check

Pass — Event subjects unchanged; new intent/payload fields only; Recover covered by tests.

## Project Structure

```text
bpmn/element/escalation.go          # <escalation> root
protocol/proto/event/v1/…           # INTENT_ESCALATION_THROWN, escalation_code
processing/deploy/escalation.go     # compile + match helpers
processing/escalations.go           # propagateEscalation
processing/handlers/…               # ThrowEscalation effect
processing/testdata/m15_*.bpmn
processing/escalation_test.go
```

## Complexity Tracking

None beyond mirroring error with non-failing uncaught path and non-interrupting boundary spawn.
