# Implementation Plan: Terminate End and Remaining Tasks

**Branch**: `014-terminate-and-tasks` | **Date**: 2026-09-11 | **Spec**: [spec.md](./spec.md)

## Summary

Accept terminate ends and Manual/Receive/Send/Business Rule tasks at deploy. Terminate cancels remaining tokens in the enclosing scope then completes that scope successfully. Manual waits like User Task; Receive waits like message catch; Send publishes like message throw; Business Rule waits as a job like Service Task. Recover covers waiting types.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `processing/handlers`, `processing/deploy`, existing Complete / PublishMessage / Job APIs, proto `Element.Type` already reserved  
**Storage**: EventLog + optional runtime store (unchanged)  
**Testing**: `go test ./processing/` fixtures `m18_*.bpmn`  
**Target Platform**: single-node engine  
**Project Type**: Go modules (`bpmn` / `protocol` / `processing` / `gateway`)  
**Constraints**: constitution — one Type per BPMN element; no new ledger subjects; Recover tests  
**Scale/Scope**: five element kinds; no proto schema change  

## Constitution Check

Pass — Types already on `Element.Type`; handlers + registry + deploy validation; Recover stories; no product-suite UI.

Post-design: still pass. `TerminateScope` is an Effect flag, not a new Element type.

## Project Structure

```text
specs/014-terminate-and-tasks/
processing/deploy/deploy.go
processing/deploy/end.go          # terminate end spec (or deploy.go)
processing/handlers/{manual,receive,send,business_rule}_task.go
processing/handlers/end_event.go
processing/handlers/handler.go
processing/executor.go / scope_terminate.go
processing/projection/instance.go
processing/messages.go / jobs.go / engine.go / errors.go / recover.go
processing/testdata/m18_*.bpmn
processing/terminate_task_test.go
```
