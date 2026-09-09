# Implementation Plan: Link Events

**Branch**: `012-link-events` | **Date**: 2026-09-09 | **Spec**: [spec.md](./spec.md)

## Summary

Add intermediate link throw/catch. Deploy indexes by link name and rejects orphan throws. Runtime: throw InstantLifecycle then `LinkContinue` into each matching catch; catch InstantLifecycle + TakeOutgoing.

## Technical Context

**Language/Version**: Go 1.26.5  
**Primary Dependencies**: `processing/deploy`, handlers, executor  
**Testing**: `go test ./processing/ -run Link`; fixtures `m16_link_*.bpmn`  
**Constraints**: Apache-2.0 headers; optional `link_name` on EventPayload  

## Constitution Check

Pass — ledger subjects unchanged; Recover covered.

## Project Structure

```text
processing/deploy/link.go
processing/handlers/intermediate_*.go
processing/executor.go  # LinkContinue
processing/testdata/m16_link_*.bpmn
processing/link_test.go
```
