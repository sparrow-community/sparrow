---
description: "Task list for Nested Event Sub-Process"
---

# Tasks: Nested Event Sub-Process

## Phase 1: Setup

- [x] T001 Baseline: `go test ./processing/ -run 'EventSubProcess' -count=1`

## Phase 2: Foundational

- [x] T002 Allow ESP nested in ESP in `validateSubProcessesAt` (`processing/deploy/deploy.go`)
- [x] T003 Arm Event Sub-Processes when entering Event Sub-Process in `processing/executor.go` (and MI executor if needed)

## Phase 3: US1 — Run ESP-in-ESP

- [x] T004 Add `processing/testdata/m13_esp_in_esp_message.bpmn`
- [x] T005 Add deploy + outer trigger + nested interrupting tests in `processing/event_subprocess_test.go`

## Phase 4: US2 — Disarm

- [x] T006 Add outer-complete-disarms-nested test in `processing/event_subprocess_test.go`

## Phase 5: US3 — Recover

- [x] T007 Add Recover mid nested-armed test in `processing/event_subprocess_test.go`

## Phase 6: Polish

- [x] T008 Update `AGENTS.md` for 009 shipped
- [x] T009 Full suite `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`
