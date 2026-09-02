# Sparrow

Lightweight, **functionally complete** BPMN **execution engine**. Single-node; append-only event log is source of truth. Event = behavior; Element = subject.

**Completeness** ≈ target executable BPMN models run end-to-end, and **Recover** is covered by tests (replay + redrive unfinished COMMANDs → same outcome).

| | |
|--|--|
| Governance | [`.specify/memory/constitution.md`](./.specify/memory/constitution.md) |
| Runtime | [`processing/README.md`](./processing/README.md) |
| AI Driven (after engine) | [`AI-Driven-BPMN.md`](./AI-Driven-BPMN.md) |

New work: `/speckit-specify` → plan → tasks → implement (`.cursor/skills/`).

## Implemented

Process lifecycle; UserTask; ServiceTask + Job; XOR/AND/Inclusive/EventBased (catch) gateways; SubProcess; CallActivity (child instance, IO mapping, **cross-deployment**); Event Sub-Process; catch/throw/boundary (timer, message, signal, error, compensate); compensation; **multi-instance** on User Task, Service Task, SubProcess; **incident** (blocked job failure, resolve/retry); process revision coexistence.

Shipped: M1–M4c, [`001`](./specs/001-engine-completeness/), [`002`](./specs/002-multi-instance/), [`003`](./specs/003-incident/), [`004`](./specs/004-cross-deploy-callactivity/).

## Remaining

Pick next via `/speckit-specify`. Order flexible until a plan sets dependencies.

| P | Topic |
|---|--------|
| 2 | Live instance migration |
| 2 | Multi-instance Call Activity |
| 2 | CallActivity boundary & compensation (parity with SubProcess) |
| 2 | Instantiate event-based gateway |
| 3 | Compensation into unfinished SubProcess |
| 3 | Nested Event Sub-Process |
| 3 | Multiple same-kind boundaries on one activity |
| 3 | Escalation, Link, conditional flow, Terminate, Send/Receive/Manual/Business Rule Task (one increment each) |

## AI Driven

After engine work above: agents use the same APIs and ledger — deploy, start, query, complete via COMMAND. Details: [`AI-Driven-BPMN.md`](./AI-Driven-BPMN.md). Custom behavior maps to BPMN constructs + extensions.

## API & persist

`Deploy` · `CreateInstance` · `Complete` · `ThrowError` · `ResolveIncident` · `FireDue` · `PublishMessage` · `PublishSignal` · Job Activate/Fail/Heartbeat · `GetInstance` · `ListEvents`

`EventLog` + `deploy.Store` + optional `runtime.Store` · `Recover` / `Open` · gRPC in `gateway` only

## Workspace

Go **1.26.5** · modules: `bpmn` · `protocol` · `processing` · `gateway`

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
go test ./bpmn/
cd protocol/proto && ./build.sh
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```

Proto: edit `.proto`, run `build.sh`; never hand-edit `*.pb.go`. Semantics live in `handlers/` + [`processing/README.md`](./processing/README.md). Apache-2.0 headers.
