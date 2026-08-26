# Sparrow

Apache-2.0 Go workspace for a lightweight, functionally complete BPMN **execution engine** (not a Camunda-style product suite). Single-node; append-only Protobuf event log as source of truth. Event = behavior; Element = subject. End state is engine completeness, not a permanent minimal subset.

Why it exists: [`AI-Driven-BPMN.md`](./AI-Driven-BPMN.md).

| | |
|--|--|
| Governance | [`.specify/memory/constitution.md`](./.specify/memory/constitution.md) |
| Active increment | [`specs/001-engine-completeness/`](./specs/001-engine-completeness/) (spec → plan → tasks) |
| Runtime design | [`processing/README.md`](./processing/README.md) |

New work: Spec Kit skills in `.cursor/skills/` (`/speckit-specify` → `/speckit-plan` → `/speckit-tasks` → `/speckit-implement`).

```text
@AGENTS.md @.specify/memory/constitution.md @specs/001-engine-completeness/
This session: <one thing; prefer next unchecked US in tasks.md>
```

## Roadmap

**Done (M1–M4c + CallActivity child instance + IO mapping + error Event Sub-Process + process revisions):** Start/UserTask/XOR/End; ServiceTask/Job; timer/message/signal catch & boundary; gateways; SubProcess; Event Sub-Process (message/timer/signal/error); throw/compensation; error boundary/end/ThrowError; CallActivity as **independent child instance** with optional IO name mappings; same `process_id` revision coexistence (`CreateInstance` by deployment or process id/version).

**Next:** New feature increment via Spec Kit (`/speckit-specify`) — deferred items below stay out of scope until specified.

**Deferred:** cross-deployment CallActivity; live migration; ≥3 same-kind waiting boundaries on one activity; Incident; instantiate EBG; compensation into unfinished SubProcess; Event Sub-Process nested in Event Sub-Process; cluster; product suite. (One timer + one message + one signal boundary may coexist.)

## Implemented elements (snapshot)

Process; Start (none); End (none/error/compensate); SequenceFlow; UserTask; ServiceTask (+ Job Activate/Fail/Heartbeat); Exclusive/Inclusive/Parallel/EventBased gateways; embedded SubProcess; CallActivity | Same-definition `calledElement`; **child process instance**; host parked with `called_process_instance_id`; optional IO name mappings; Event Sub-Process (message/timer/signal/error); intermediate catch/throw (timer/message/signal/compensate); Boundary (timer/message/signal/compensate/error); Association (compensation); process revision index (latest or explicit version at start).

**Gaps covered by active spec:** CallActivity child instance + IO mapping; Error Event Sub-Process; process revision coexistence — **implemented**. **Still out of scope:** Escalation/Link/Conditional/Terminate; Send/Receive/Manual/BusinessRule Task; Multi-instance; cross-deployment CallActivity.

## Runtime surface

- API: `Deploy` / `CreateInstance` / `Complete` / `ThrowError` / `FireDue` / `PublishMessage` / `PublishSignal` / Job trio; `GetInstance` / `ListEvents`
- Persist: `EventLog` + `deploy.Store` + optional `runtime.Store`; `Recover` / `Open`
- Transport: `gateway` only (`engine.v1` + `job.v1`) — not in `processing`

## Workspace

Go **1.26.5** (`go.work`; commit it, ignore `go.work.sum`). Module prefix `github.com/sparrow-community/sparrow/...`. Local `replace` for sibling modules.

| Module | Role |
|--------|------|
| `bpmn` | BPMN 2.0 XML → `element` types; MIWG tests |
| `protocol` | Protobuf only: `event.v1` ledger, `job.v1` / `engine.v1` RPC (not Event types). Generated Go committed under `protocol/gen/go` |
| `processing` | Engine: `deploy` / `handlers` / `executor` / `projection` / `log`; UUIDv7 ids; no parallel graph |
| `gateway` | gRPC adapter + `cmd/sparrow` |

**protocol:** change `.proto`, run `cd protocol/proto && ./build.sh` (`buf lint` + generate). Never hand-edit `*.pb.go`. Keep `buf` FILE wire compatibility.

**processing:** semantics live in handlers + `processing/README.md`; waiting work uses one `Complete`; Job leases / late message buffers are runtime store, not ledger subjects.

## Commands

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
go test ./bpmn/   # MIWG; slower
cd protocol/proto && ./build.sh
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```

## Conventions

- Apache-2.0 (Sparrow community) license headers.
- Do not add top-level modules without need; do not put gRPC in `processing`.
- AI/custom behavior maps onto existing BPMN constructs (e.g. Service Task + extensions), not new OMG-core element types.
