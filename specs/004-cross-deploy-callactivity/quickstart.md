# Quickstart: Cross-deployment Call Activity

Validate each user story independently. Run from the repository root (`go.work`).

## Prerequisites

- Go 1.26.5
- Buf CLI only if optional proto field is added (`cd protocol/proto && ./build.sh`)

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Story 1 — Separately deployed callee spawns child

1. Deploy callee BPMN (`PaymentFlow` or fixture `m8_cross_call_callee.bpmn`) → note `deployment_id_callee`.
2. Deploy caller BPMN (`OrderFlow` with Call Activity `calledElement=PaymentFlow`, fixture `m8_cross_call_caller.bpmn`) → `deployment_id_caller`.
3. `CreateInstance` on caller deployment.
4. `GetInstance` caller: token on Call Activity `status=waiting`, `called_process_instance_id` set.
5. `GetInstance` child: `process_id=PaymentFlow`, `deployment_id=deployment_id_callee`, `parent_process_instance_id=caller id`.
6. Complete child work → caller completes Call Activity → caller `status=completed`.

**Negative (SC-004)**: Deploy caller only (no callee). `CreateInstance` + reach Call Activity → REJECTION `NOT_FOUND`; caller not waiting on call.

## Story 2 — IO mapping across deployments

1. Deploy callee with user task or variable-producing end.
2. Deploy caller with input map (`orderId` → child `id`) and output map (`status` → `paymentStatus`).
3. Start caller with `orderId`; after call, caller has `paymentStatus` and original `orderId`; child did not receive unmapped caller fields.

Mirror assertions from existing `m4_call_activity.bpmn` IO test but with split deploy files.

## Story 3 — Recover with different deployment ids

1. Complete Story 1 steps 1–4 (parent waiting, child active).
2. `Recover` from EventLog + deploy store containing **both** BPMN deployments (mirror `processing/recover_test.go`).
3. `GetInstance` parent and child: same ids, same `deployment_id` values as before crash.
4. Complete child → caller completes.

**Redrive**: Crash simulation after parent ACTIVATED but before child PROCESS events complete — `Recover` finishes one child instance, consistent parent linkage.

## Regression (SC-005)

```shell
go test ./processing/ -run CallActivity -count=1
```

Existing tests using `m4_call_activity.bpmn`, `m5_call_*.bpmn` must pass without fixture changes.

## Expected test layout

- Fixtures: `processing/testdata/m8_cross_call_caller.bpmn`, `m8_cross_call_callee.bpmn` (names illustrative)
- Tests: extend `processing/call_activity_test.go`; add `recover_test.go` cross-deploy case
- No gateway test required unless optional proto field added

## Full suite

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```
