# Implementation Plan: Abstract Task

**Branch**: `019-abstract-task` | **Date**: 2026-09-15 | **Spec**: [spec.md](./spec.md)

## Summary

Accept abstract `bpmn:task` at Deploy as `TYPE_TASK`. Runtime wait → Complete mirrors Manual Task (`waitingTaskEnter` / `waitingTaskComplete`). Support MI and boundary hosts like Manual. Remove the Deploy rejection that treated any `<task>` as outside the subset.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: existing `processing` handlers, `eventv1.Element_TYPE_TASK`  
**Storage**: EventLog unchanged  
**Testing**: `m23_abstract_task.bpmn` + wait/complete/Recover tests; flip prior reject test  
**Target Platform**: single-node engine  
**Project Type**: library (`processing`)  
**Constraints**: constitution — no new Element invent; Manual-equivalent wait only  
**Scale/Scope**: one FlowElement type

## Constitution Check

Pass: executable FlowElement coverage; ledger uses distinct `TYPE_TASK`; no UI/cluster; reuse Complete command.

## Project Structure

### Documentation (this feature)

```text
specs/019-abstract-task/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/engine.md
└── tasks.md
```

### Source Code

```text
processing/deploy/deploy.go          # index Tasks; drop abstract reject; boundary/default
processing/handlers/task.go          # TYPE_TASK handler
processing/handlers/handler.go       # registry
processing/projection/instance.go    # waitingActivation
processing/errors.go                 # ThrowError hosts if Manual-listed
processing/testdata/m23_*.bpmn
processing/abstract_task_test.go
AGENTS.md / processing/README.md
```
