# Implementation Plan

**Branch**: `023-compensation-event-subprocess` | **Spec**: [spec.md](./spec.md)

## Summary

Accept compensate start on event sub-process nested in an embedded SubProcess; subscribe that event sub-process as the SubProcess compensation handler on complete; run it from `startCompensation`; nested compensate throws inside it collect subscriptions under the enclosing SubProcess scope.

## Research

- Reuse `Compensation` map with BoundaryID = compensate start id, HandlerID = event sub-process id, ActivityID = enclosing SubProcess.
- Skip `CatchKindCompensate` in `emitEventSubProcessStartArms`.
- On event sub-process complete, `AdvanceCompensation` when handler is compensation event sub-process.
- Remap throw scope when throw lives inside a compensation event sub-process.

## Files

`processing/deploy/event_subprocess.go`, `compensation.go`, `deploy.go`, `timer.go` (CatchKind), `processing/event_subprocess.go`, `executor.go`, `handlers/sub_process.go`, fixtures `m27_*`, tests, AGENTS/README.
