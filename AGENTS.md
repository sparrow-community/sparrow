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

Process lifecycle; UserTask; ServiceTask + Job; Manual/Receive/Send/Business Rule Task; XOR/AND/Inclusive/EventBased (catch + **instantiate exclusive start**); SubProcess; CallActivity (child instance, IO mapping, **cross-deployment**, **boundary & compensation parity**, **multi-instance**); Event Sub-Process (**incl. nested in ESP**); catch/throw/boundary (timer, message, signal, error, escalation, compensate; **multiple same-kind**); **link** throw/catch; **conditional sequence flows** (activity + gateway); **terminate end**; compensation (**into unfinished SubProcess**); **multi-instance** on User Task, Service Task, Manual/Receive/Send/Business Rule, SubProcess, Call Activity; **incident** (blocked job failure, resolve/retry); process revision coexistence.

Shipped: M1–M4c, [`001`](./specs/001-engine-completeness/)–[`014`](./specs/014-terminate-and-tasks/).

## Remaining

Official executable-completeness increment list is empty. Pick further work via `/speckit-specify` (Script Task, Complex Gateway, Transaction/cancel, process-level typed start, etc. are not on this list).

## Future

| Topic | Note |
|-------|------|
| Live instance migration | Revisions already coexist; migrating running instances is ops, not completeness |

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
