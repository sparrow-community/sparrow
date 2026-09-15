# Implementation Plan: MI Complex Behavior

**Branch**: `021-mi-complex-behavior` | **Date**: 2026-09-15 | **Spec**: [spec.md](./spec.md)

## Summary

Extend MI deploy compilation to accept None/Complex behavior, none/one event refs, and complexBehaviorDefinition (signal/message only). After each inner Complete, publish milestone events; track complex one-fire in MultiInstanceLoop for Recover.

## Technical Context

**Language/Version**: Go 1.26.5  
**Testing**: `m25_mi_behavior_*.bpmn`  
**Constraints**: reuse Publication message/signal; no new Element types

## Constitution Check

Pass.

## Project Structure

```text
specs/021-mi-complex-behavior/
processing/deploy/multi_instance.go
processing/multi_instance_executor.go
processing/projection/instance.go
processing/testdata/m25_*.bpmn
processing/mi_complex_behavior_test.go
```
