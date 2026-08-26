# Quickstart: Engine Completeness (next increment)

Validate each story independently. Run from the repository root (`go.work`).

## Prerequisites

- Go 1.26.5
- Buf CLI only if proto files changed (`cd protocol/proto && ./build.sh`)

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Story 1 — Called process as its own instance

1. Deploy a definitions file with `Caller` (start → callActivity `Child` → end) and `Child` (start → userTask → end). Two Call Activities naming `Child` MUST deploy.
2. `CreateInstance` on `Caller`.
3. `GetInstance(caller)`: one waiting host token on the Call Activity with a non-empty called instance id.
4. `GetInstance(child)`: status active, `parent_process_instance_id` = caller, waiting on the user task.
5. `Complete` the child user task.
6. Both instances complete; caller's audit shows CALL_ACTIVITY COMPLETED after the child's PROCESS COMPLETED.

Recover: after step 4, rebuild from the EventLog and repeat 5–6.

## Story 2 — IO mapping

1. Same shape as story 1, with input `orderId`→`id` and output `total`→`amount`.
2. Start caller with `{orderId: "o1", extra: 1}`.
3. Child variables include `id`, exclude `extra`.
4. Complete child with `{total: 9}`.
5. Caller variables include `amount=9` and `extra=1`.

No-mapping fixture: child starts empty; caller variables unchanged after the call.

## Story 3 — Error Event Sub-Process

1. Process with an error end (`E1`) on the default path and an interrupting Event Sub-Process whose error start names `E1`, then a user task.
2. Start: default path interrupted; wait in the Event Sub-Process.
3. Complete the handler user task; process completes.
4. Add a fixture where `E2` is thrown and only `E1` is caught: instance terminates (or bubbles) without starting the ESP.
5. Non-interrupting fixture: waiting user task on the default path remains while the ESP runs.

## Story 4 — Revisions

1. Deploy `P` (user task `TaskA`); start instance A.
2. Deploy `P` again (user task `TaskB`).
3. Start instance B without naming a version (latest).
4. A waits on `TaskA`; B waits on `TaskB`.
5. Start instance C with `process_id=P` and the first version; C waits on `TaskA`.
6. Unknown version → rejected, no instance.

## Expected test layout

Fixtures under `processing/testdata/` (`m5_call_instance*.bpmn`, `m5_call_io*.bpmn`, `m5_error_esp*.bpmn`, `m5_version*.bpmn`). Tests in `processing/call_activity_test.go` (extend), new `processing/error_event_subprocess_test.go`, `processing/version_test.go`. Gateway mapping covered in `gateway/engine_server_test.go`.
