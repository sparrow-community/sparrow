# Implementation Plan: Conditional Flow

**Branch**: `013-conditional-flow` | **Date**: 2026-09-10 | **Spec**: [spec.md](./spec.md)

## Summary

When `takeOutgoing` has no preferred flow, choose among outgoings using conditionExpression + element `default` (activity or exclusive gateway). Reuse XOR selection rules; Inclusive unchanged.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `deploy.ChooseExclusiveOutgoing` / new shared chooser, `executor.takeOutgoing`, `expr.Eval`  
**Testing**: `go test ./processing/ -run 'Conditional|Exclusive|Inclusive'`; fixtures `m17_conditional_*.bpmn`  

## Constitution Check

Pass — no new Element types; Recover covered.

## Project Structure

```text
processing/deploy/deploy.go   # ChooseConditionalOutgoing / DefaultFlow
processing/executor.go        # pass vars into takeOutgoing
processing/testdata/m17_*.bpmn
processing/conditional_flow_test.go
```
