# Sparrow

Sparrow is a single-node BPMN **execution engine**. It executes OMG BPMN 2.0 **process** models, records every accepted behavior on an append-only event log, and rebuilds state by replay. Event = behavior; Element = subject.

## Purpose

Sparrow provides three properties together:

1. **Element coverage** — Every executable process FlowElement and event-definition combination is either Supported or Excluded. Unsupported combinations are rejected at deploy. Silence is not exclusion.
2. **Semantic completeness** — Each Supported combination implements the BPMN runtime semantics required for that combination (token movement, waits, throws, joins, scopes, boundaries, multi-instance, compensation, and related rules). Deploy acceptance alone is insufficient.
3. **Executable lightness** — The engine runs processes through COMMAND → EVENT with Recover, jobs, messages, timers, and incidents. It stays a ledger-backed execution kernel: no modeler, ops console, cluster, or choreography/collaboration runtime.

Completeness holds when (1) and (2) are true for the tables below, and every Supported combination has end-to-end and Recover tests. Completing a delivery batch does not by itself establish completeness.

The ledger records process behavior only. Documentary or non-subject constructs are Excluded; they do not receive empty Element intents.

| | |
|--|--|
| Governance | [`.specify/memory/constitution.md`](./.specify/memory/constitution.md) |
| Runtime | [`processing/README.md`](./processing/README.md) |
| AI Driven (consumer) | [`AI-Driven-BPMN.md`](./AI-Driven-BPMN.md) |

Work proceeds `/speckit-specify` → plan → tasks → implement (`.cursor/skills/`).

## Supported

Process lifecycle; User Task; Service Task + Job; Manual / Receive / Send / Business Rule Task; exclusive / parallel / inclusive / event-based gateways (catch, including exclusive instantiate entry); SubProcess; Call Activity (child instance, IO mapping, cross-deployment, boundary and compensation parity, multi-instance); Event Sub-Process (including nested); intermediate and boundary catches/throws for timer, message, signal, error, escalation, compensate (multiple same-kind boundaries); link throw/catch; conditional sequence flows on activities and gateways; terminate end; compensation into an unfinished embedded SubProcess; multi-instance on User Task, Service Task, Manual / Receive / Send / Business Rule, SubProcess, Call Activity; incident open / resolve / retry; coexisting process revisions.

Shipped: M1–M4c, specs [`001`](./specs/001-engine-completeness/)–[`014`](./specs/014-terminate-and-tasks/).

## Planned

Open executable-process gaps. Priority guides sequencing; a feature plan may reorder for dependencies. Each row is one specify increment unless a plan bundles tightly related cells.

| P | Topic |
|---|--------|
| 1 | Process-level typed start: message, timer, signal, conditional, error |
| 1 | Message end; signal end |
| 1 | Conditional catch: intermediate and boundary |
| 1 | Script Task |
| 1 | Abstract `bpmn:task` as Manual-equivalent wait/complete |
| 2 | `standardLoopCharacteristics` |
| 2 | Multi-instance `complexBehaviorDefinition` and none/one behavior event refs |
| 2 | Instantiate Receive Task; parallel instantiate event-based gateway; event-based gateway targeting Receive Task; multiple instantiate event-based gateway entries |
| 2 | Compensation Event Sub-Process |
| 2 | Compensation into an unfinished Call Activity child instance |
| 2 | Standalone intermediate escalation catch |
| 2 | Event Sub-Process conditional start |
| 2 | Call Activity IO mapping with transformation / assignment |
| 2 | Parallel fan-out from an activity with multiple unconditional outgoings |
| 3 | Transaction SubProcess and cancel |
| 3 | Complex Gateway |
| 3 | Ad-Hoc SubProcess |
| 3 | Implicit throw event — reject at deploy with a stable code (no silent ignore, no none-throw alias) |

## Excluded

Non-goals. Moving an item out of Excluded requires an AGENTS and constitution amendment.

| Topic | Rationale |
|-------|-----------|
| Collaboration, message flow, and choreography execution | Outside single-process execution |
| Lane runtime semantics | Documentation and grouping only |
| Data Object / Data Store as ledger or token subjects | Process variables and IO mappings carry executable data |
| Modeler, operations console, product-suite UI | Separate consumers of the engine |
| Cluster and multi-node replication | Single-node design |
| In-engine DMN evaluation | Business Rule Task is job-backed |

## Future (operations)

| Topic | Note |
|-------|------|
| Live instance migration | Revisions coexist; moving a running instance across revisions is operations, not element completeness |

## AI Driven

[`AI-Driven-BPMN.md`](./AI-Driven-BPMN.md) states why Sparrow exists and how agents consume it. BPMN element coverage and semantic completeness in this file are the prerequisite. Agents use the same COMMAND surface and event log. Agent behavior maps onto BPMN constructs and extensions. Planned engine work remains the completeness mainline. AI Driven does not invent non-OMG core element types and does not skip COMMAND/EVENT.

## API and persistence

`Deploy` · `CreateInstance` · `Complete` · `ThrowError` · `ResolveIncident` · `FireDue` · `PublishMessage` · `PublishSignal` · Job Activate / Fail / Heartbeat · `GetInstance` · `ListEvents`

`EventLog` + `deploy.Store` + optional `runtime.Store` · `Recover` / `Open` · gRPC in `gateway` only

## Workspace

Go **1.26.5** · modules: `bpmn` · `protocol` · `processing` · `gateway`

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
go test ./bpmn/
cd protocol/proto && ./build.sh
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```

Edit `.proto`, then run `build.sh`. Do not hand-edit `*.pb.go`. Element semantics live in `handlers/` and [`processing/README.md`](./processing/README.md). Apache-2.0 headers.
