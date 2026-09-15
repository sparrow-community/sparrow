# Implementation Plan: Standard Loop Characteristics

**Branch**: `020-standard-loop` | **Date**: 2026-09-15 | **Spec**: [spec.md](./spec.md)

## Summary

Compile `standardLoopCharacteristics` at deploy; reject coexistence with MI. On waiting-task enter/complete, apply testBefore / loopCondition / loopMaximum; re-enter the same element while the loop continues. Track iteration on the token via projection replay for Recover. No new EventLog subject types.

## Technical Context

**Language/Version**: Go 1.26.5  
**Testing**: `m24_standard_loop_*.bpmn` + unit/e2e tests  
**Constraints**: constitution COMMAND→EVENT; reuse `processing/expr`

## Constitution Check

Pass.

## Project Structure

```text
specs/020-standard-loop/
bpmn/element/standard_loop_characteristics.go  # loopCondition → ExpressionUnMarshal
processing/deploy/standard_loop.go
processing/deploy/multi_instance.go            # mutual exclusion
processing/handlers/standard_loop.go
processing/handlers/waiting_task.go / user_task.go
processing/handlers/handler.go                 # Effect.ReEnter
processing/executor.go
processing/projection/instance.go              # StandardLoopIteration
processing/testdata/m24_*.bpmn
processing/standard_loop_test.go
```
