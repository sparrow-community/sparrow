# Implementation Plan: Conditional Catch

**Branch**: `017-conditional-catch` | **Date**: 2026-09-15 | **Spec**: [spec.md](./spec.md)

## Summary

Index intermediate and boundary conditional catches. Intermediate: wait if condition false at enter, else take outgoing; EvaluateConditions completes waiting catches when true. Boundaries: arm with kind `conditional`; EvaluateConditions fires like message/signal boundaries. Reuse `processing/expr`. No proto change.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: deploy indexes, handlers catch/boundary, expr.Eval, Complete / boundary fire paths  
**Testing**: `go test ./processing/` fixtures `m21_*.bpmn`  
**Constraints**: constitution; ESP conditional start out of scope  

## Constitution Check

Pass — existing Element types; deploy + handler + Engine evaluate API; no new ledger subjects.

## Project Structure

```text
specs/017-conditional-catch/
processing/deploy/conditional.go   # catch/boundary specs
processing/deploy/deploy.go        # index + validate
processing/deploy/timer.go         # CatchKindConditional
processing/handlers/intermediate_catch_event.go
processing/handlers/boundary_event.go
processing/conditional_eval.go     # EvaluateConditions
processing/testdata/m21_*.bpmn
processing/conditional_catch_test.go
```

## Complexity Tracking

> None
