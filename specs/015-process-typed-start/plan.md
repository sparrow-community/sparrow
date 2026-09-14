# Implementation Plan: Process-Level Typed Start

**Branch**: `015-process-typed-start` | **Date**: 2026-09-14 | **Spec**: [spec.md](./spec.md)

## Summary

Index process-level message, timer, signal, and conditional start events at Deploy. CreateInstance enters only via none start (or existing instantiate exclusive EBG). PublishMessage / PublishSignal / FireDue / EvaluateConditionalStarts create instances at the matching typed start when waiters do not consume the trigger. Process-level error start is rejected; Event Sub-Process error starts unchanged. Recover covers armed timer-start schedules and post-creation waits.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `processing/deploy`, `processing/expr`, existing CreateInstance / PublishMessage / PublishSignal / FireDue / Recover, optional `runtime.Store`  
**Storage**: EventLog + deploy.Store + optional runtime.Store (timer-start schedules / subscription helpers outside the ledger)  
**Testing**: `go test ./processing/` fixtures `m19_*.bpmn`  
**Target Platform**: single-node engine  
**Project Type**: Go modules (`bpmn` / `protocol` / `processing` / `gateway`)  
**Constraints**: constitution — no new ledger subjects; typed creation reuses PROCESS/START_EVENT intents; Recover tests; instantiate EBG must not auto-create (007 unchanged)  
**Scale/Scope**: four typed start kinds + error reject; shared instance-mint path; optional gateway RPC only if we expose EvaluateConditionalStarts (FireDue remains engine-only today)

## Constitution Check

Pass — StartEvent Type already exists; deploy indexing + Engine trigger paths; timer-start arms live in runtime helpers (like message buffers), not as ledger subjects; no product-suite UI; Call Activity still requires none/instantiate entry on called process.

Post-design: still pass. Shared `createInstanceAt` is an Engine helper, not a parallel graph. EvaluateConditionalStarts mirrors FireDue (processing API; gateway optional).

## Project Structure

```text
specs/015-process-typed-start/
processing/deploy/start.go          # process-level start specs / NoneStartEventID
processing/deploy/deploy.go         # validate + index typed starts; reject process-level error start
processing/deploy/call_activity.go  # called process must have none or instantiate entry
processing/engine.go               # CreateInstance uses none start only; shared mint helper
processing/messages.go             # after waiters: create message-start instances
processing/signals.go              # after waiters: create signal-start instances
processing/timers.go               # FireDue: due process-level timer starts
processing/conditional_start.go    # EvaluateConditionalStarts
processing/runtime/                # persist timer-start schedules if store present
processing/testdata/m19_*.bpmn
processing/typed_start_test.go
gateway/                            # optional EvaluateConditionalStarts RPC (defer if FireDue pattern)
```

**Structure Decision**: Extend `processing` deploy + Engine trigger surfaces; no new module; no proto required for MVP (align with FireDue).

## Complexity Tracking

> None — no constitution violations.
