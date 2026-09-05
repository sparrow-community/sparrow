# Quickstart: Multi-Instance Call Activity

Validate each user story independently. Run from the repository root (`go.work`).

## Prerequisites

- Go 1.26.5
- No proto changes expected

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Story 1 — Parallel MI Call Activity

1. Deploy fixture with Call Activity `multiInstanceLoopCharacteristics` parallel, `loopCardinality` 3, plus called process (e.g. `m10_mi_call_parallel.bpmn`).
2. `CreateInstance` on caller.
3. Expect three distinct child instance ids waiting on the called user task; parent waiting on Call Activity host/inners.
4. Complete all three children → parent completed once; one Call Activity host COMPLETED.

## Story 2 — Sequential MI Call Activity

1. Deploy sequential MI Call Activity cardinality 3 (`m10_mi_call_sequential.bpmn`).
2. Start parent → exactly one child active.
3. Complete child → next child starts; repeat thrice → parent ends.

## Story 3 — Collection + IO

1. Deploy parallel MI Call Activity with `loopDataInputRef` / `inputDataItem` and Call Activity IO mappings (`m10_mi_call_collection.bpmn`).
2. Start with collection variable of size N (and empty-collection case).
3. Verify N children (or immediate host complete when empty) and mapped inputs on children; optional output collection on parent after join.

## Story 4 — Completion condition + interrupting boundary

1. Parallel MI Call Activity cardinality 5 with completion condition `nrOfCompletedInstances >= 2`.
2. Complete any two children → remaining children terminated; parent continues.
3. Separate fixture: interrupting timer on MI Call Activity → `FireDue` → all children terminated; timeout path taken.

## Story 5 — Recover

1. Parallel MI Call Activity cardinality 3; complete one child.
2. `Recover` / `Open` from EventLog + deployments.
3. Complete remaining children → parent completed; child deployment ids unchanged (include a cross-deploy variant if split).

## Regression

```shell
go test ./processing/ -run 'CallActivity|CrossDeploy|MultiInstance' -count=1
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Expected test layout

- `processing/testdata/m10_mi_call_*.bpmn`
- `processing/call_activity_test.go` and/or `multi_instance_test.go`
- `processing/recover_test.go` — mid-loop MI call
- Boundary coverage in `boundary_test.go` or call activity tests
