# Quickstart: Compensation into Unfinished SubProcess

Validate from repository root (`go.work`).

## Prerequisites

```shell
go test ./processing/ -run 'Compensation|Compensate' -count=1
```

## Story 1 — Broadcast into unfinished SP

1. Deploy `m12_compensate_unfinished_sp_parallel.bpmn` (AND split: SP with compensatable A then waiting B; sibling path to compensate throw).
2. CreateInstance → complete A inside SP → B waiting; sibling ready or complete sibling to throw.
3. Expect: B (and SP host) terminated; Undo_A waiting/completable; process completes after Undo_A.

## Story 2 — Targeted activityRef

1. Deploy `m12_compensate_unfinished_sp_targeted.bpmn` (unfinished SP + completed sibling Task_X; throw activityRef=SP).
2. Expect: SP cancelled + inner undo; Task_X subscription not consumed.

## Story 3 — Recover

1. Drive Story 1 to Undo_A waiting after throw.
2. Close/Open Recover; Complete Undo_A → completed instance.

## Regression

```shell
go test ./processing/ -run 'Compensation|Compensate|CallActivityCompensation' -count=1
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```
