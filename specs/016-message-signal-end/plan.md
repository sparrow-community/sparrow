# Implementation Plan: Message End and Signal End

**Branch**: `016-message-signal-end` | **Date**: 2026-09-15 | **Spec**: [spec.md](./spec.md)

## Summary

Accept process (and SubProcess) end events with exactly one message or signal event definition. On enter, publish via existing PublicationMessage / PublicationSignal (same as intermediate throw / Send Task), complete the end lifecycle, and try-complete the enclosing scope. Reject mixed end definitions. Recover covered via existing publication + waiter paths.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `processing/deploy`, `processing/handlers` EndEventHandler, existing flushPublications  
**Storage**: EventLog (unchanged)  
**Testing**: `go test ./processing/` fixtures `m20_*.bpmn`  
**Target Platform**: single-node engine  
**Project Type**: Go modules  
**Constraints**: constitution — no new Element.Type; reuse Publication kinds; Recover tests  
**Scale/Scope**: two end kinds; no proto change  

## Constitution Check

Pass — EndEvent Type already exists; deploy index + handler Effect.Publish; no new ledger subjects.

Post-design: still pass.

## Project Structure

```text
specs/016-message-signal-end/
processing/deploy/end.go            # message/signal end specs + indexes
processing/deploy/deploy.go         # index messageEnds / signalEnds
processing/handlers/end_event.go    # publish on enter
processing/testdata/m20_*.bpmn
processing/message_signal_end_test.go
```

## Complexity Tracking

> None
