# Quickstart: Call Activity boundary and compensation parity

Validate each user story independently. Run from the repository root (`go.work`).

## Prerequisites

- Go 1.26.5
- No proto changes expected for this increment

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Story 1 — Interrupting boundaries on Call Activity

1. Deploy fixture with caller + called process and Call Activity with **interrupting timer** boundary (e.g. `m9_call_timer_boundary.bpmn`).
2. `CreateInstance` on caller; confirm child active and caller waiting at Call Activity.
3. `FireDue` (or wait for timer) on caller instance.
4. `GetInstance` child: `status=terminated` (or completed via terminate path).
5. `GetInstance` caller: token left Call Activity; process at boundary end or following path.
6. Repeat subset for **message** and **signal** interrupting fixtures if split.

**Error boundary variant**: Deploy Call Activity with error boundary; `ThrowError` on Call Activity host token → child terminated, boundary path taken.

**Cross-deploy variant**: Interrupting timer on `m8_cross_call_caller` + callee deployed → same outcome.

## Story 2 — Non-interrupting boundaries while call waits

1. Deploy fixture with non-interrupting message boundary on Call Activity (`m9_call_message_non_interrupt.bpmn`).
2. Start caller; child active, caller waiting.
3. `PublishMessage` matching boundary.
4. `GetInstance` caller: extra token on boundary path; host still on Call Activity; child still active.
5. Complete child user task → caller completes Call Activity; verify boundary cleanup per existing non-interrupt rules.

## Story 3 — Compensation on completed call

1. Deploy fixture: Call Activity + compensation boundary → handler user task (`m9_call_compensate.bpmn`).
2. Run call to completion (child done).
3. `GetInstance` caller: `CompensationSubs` contains boundary entry (or boundary ACTIVATED in event log).
4. Throw compensate in scope → handler user task becomes active; complete handler.

**Negative**: Interrupt call via timer before complete → no compensation subscription; compensate throw does not run that handler.

## Story 4 — Recover

1. Reach waiting Call Activity with armed timer + active child (Story 1 setup, before FireDue).
2. `Recover` from EventLog + deployments (mirror `recover_test.go`).
3. `GetInstance`: same `due_unix_ms`, `called_process_instance_id`, boundary ids.
4. `FireDue` after recover → same as non-recover interrupt outcome.

**Compensation recover**: Complete call with compensation armed; Recover; compensate throw still runs handler.

## Regression (SC-006)

```shell
go test ./processing/ -run 'CallActivity|CrossDeploy' -count=1
```

## Full suite

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Expected test layout

- `processing/testdata/m9_call_*.bpmn`
- `processing/call_activity_test.go` or `boundary_test.go` — interrupting cases
- `processing/compensation_test.go` — call compensation case
- `processing/recover_test.go` — boundary + compensation recover
