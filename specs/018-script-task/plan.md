# Implementation Plan: Script Task

**Branch**: `018-script-task` | **Date**: 2026-09-15 | **Spec**: [spec.md](./spec.md)

## Summary

Parse `scriptTask` in bpmn element; deploy as TYPE_SCRIPT_TASK; job wait like Business Rule (job type from scriptFormat → name → id). Extend Job with Script/ScriptFormat. No in-engine script runtime.

## Technical Context

**Language/Version**: Go 1.26.5  
**Testing**: `m22_script_task.bpmn` + tests  
**Constraints**: constitution — map onto job-backed tasks  

## Constitution Check

Pass.

## Project Structure

```text
bpmn/element/script_task.go
bpmn/element/flow_element.go
processing/deploy/deploy.go
processing/handlers/script_task.go
processing/jobs.go
processing/projection/instance.go
processing/engine.go / errors.go / recover.go
processing/testdata/m22_*.bpmn
processing/script_task_test.go
```
