# Sparrow

> English (canonical for specs and tooling). Chinese: [`AGENTS.zh.md`](./AGENTS.zh.md). Keep both in sync in the same change whenever Supported / Planned / Excluded, axioms, Purpose, or Completeness change.

Sparrow is a BPMN **execution and fact kernel**: it accepts a versioned process definition, runs it, and records what happened. Agents, UIs, and overlay renderers are consumers of that contract and trail—they do not rewrite them.

**Axioms** (what must stay true):

- Effective definition and instance state change only through accepted **COMMAND**s; every accepted behavior is an append-only **EVENT**. Projections are disposable and rebuildable by replay (**Recover**).
- **Event = behavior; Element = subject.** The ledger is a fact stream about process behavior, not a second copy of the diagram.
- Combinations with runtime token or side-effect semantics are either **Supported** (real BPMN semantics) or **Excluded** (explicit non-goal). Unsupported work is rejected at Deploy. Silence is not exclusion.
- The ledger records behavior only. Documentary or non-subject constructs stay on the definition; they do not receive empty Element intents.

**Scope choice** (deliberate, not derived from the axioms alone): OMG BPMN 2.0 **executable process** subset—FlowElements and event definitions that carry runtime token behavior—not Collaboration, Choreography, or the full meta-model. **Engineering choice**: single-node; one fact stream and one lock per instance.

## Purpose

Process work needs a durable, shared contract across people, time, and short-lived planners (humans or agents): a versioned definition, faithful execution, and an auditable trail. Sparrow exists so that contract cannot be silently mutated outside COMMAND → EVENT.

## Completeness

Faithful interpretation of the declared language is one axis (coverage + semantics). Kernel shape is another.

1. **Coverage** — Every executable-process FlowElement and event-definition combination is either Supported or Excluded. Unsupported combinations are rejected at deploy.
2. **Semantics** — Each Supported combination implements the BPMN runtime semantics required for that combination (token movement, waits, throws, joins, scopes, boundaries, multi-instance, compensation, and related rules). Deploy acceptance alone is insufficient.
3. **Kernel boundary** — The engine stays a ledger-backed execution kernel (COMMAND → EVENT, Recover, jobs, messages, timers, incidents). Modeler, operations console, cluster replication, and choreography/collaboration runtime are consumer or non-goals—not engine duties. Product features such as overlaying Events onto a deployment diagram belong in consumers; they must not invent ledger subjects or empty intents.

Completeness holds when (1) and (2) are true for the tables below, and every Supported combination has end-to-end and Recover tests (how we know—not a separate engine property). Completing a delivery batch does not by itself establish completeness.

| | |
|--|--|
| Governance | [`.specify/memory/constitution.md`](./.specify/memory/constitution.md) |
| Runtime | [`processing/README.md`](./processing/README.md) |

Work proceeds `/speckit-specify` → plan → tasks → implement (`.cursor/skills/`).

## Supported

Process lifecycle; process-level none start and typed starts (message, timer, signal, conditional — CreateInstance for none start, instantiate exclusive/parallel event-based gateway entries, and instantiate receive task; typed mint via PublishMessage / PublishSignal / FireDue / EvaluateConditionalStarts); User Task; Service Task + Job; Manual / abstract Task (wait → Complete); Receive / Send (including instantiate receive as process entry); Business Rule / Script Task (job-backed); exclusive / parallel / inclusive / complex / event-based gateways (catch, including instantiate entries and receive-task targets); SubProcess; Transaction SubProcess and cancel (##Compensate; Cancel End → compensate → Cancel Boundary); Ad-Hoc SubProcess (flat inner activities, `ordering` Parallel / Sequential, `completionCondition`, `cancelRemainingInstances`); Call Activity (child instance, IO mapping with name copy / transformation / assignment, cross-deployment, boundary and compensation parity, multi-instance; called process requires none start, instantiate event-based gateway, or instantiate receive); Event Sub-Process (including nested, compensation event sub-process, and conditional start via EvaluateConditions); intermediate and boundary catches/throws for timer, message, signal, error, escalation (including standalone intermediate escalation catch), compensate, conditional (multiple same-kind boundaries); link throw/catch; conditional sequence flows on activities and gateways; parallel fan-out from activities and start events with multiple unconditional outgoings; terminate end; message end; signal end; compensation into an unfinished embedded SubProcess or unfinished Call Activity child instance; multi-instance on User Task, Service Task, Manual / abstract Task / Receive / Send / Business Rule / Script, SubProcess, Call Activity (including `complexBehaviorDefinition` and none/one behavior event refs for signal/message); `standardLoopCharacteristics` on waiting tasks (testBefore / loopCondition / loopMaximum); incident open / resolve / retry; coexisting process revisions.

Shipped: M1–M4c, specs [`001`](./specs/001-engine-completeness/)–[`035`](./specs/035-get-deployment/). Every Supported combination above has an end-to-end and a Recover test ([`034`](./specs/034-recover-coverage/) closed the last coverage gaps).

## Planned

Open executable-process gaps. Priority guides sequencing; a feature plan may reorder for dependencies. Each row is one specify increment unless a plan bundles tightly related cells.

No open executable-process element gaps. The next increment is chosen when a new gap is identified; candidates are recorded here before work starts.

## Excluded

Non-goals. Moving an item out of Excluded requires an AGENTS (both language editions) and constitution amendment.

| Topic | Rationale |
|-------|-----------|
| Collaboration, message flow, and choreography execution | Outside single-process execution; a sequence flow into a choreography or otherwise unmodeled flow element is rejected at Deploy |
| Lane runtime semantics | Documentation and grouping only |
| Data Object / Data Store as ledger or token subjects | Process variables and IO mappings carry executable data |
| Modeler, operations console, product-suite UI | Separate consumers of the engine |
| Cluster and multi-node replication | Single-node design |
| In-engine DMN evaluation | Business Rule Task is job-backed |
| `implicitThrowEvent` as a flow element | No token semantics of its own; valid only as a multi-instance behavior event, so Deploy rejects it with `UNSUPPORTED_ELEMENT` instead of aliasing it to a none throw |

## Future (operations)

| Topic | Note |
|-------|------|
| Live instance migration | Revisions coexist; moving a running instance across revisions is operations, not element completeness |

## Consumers

Agents, overlay UIs, ops tools, and MCP adapters are peers: each consumes the same COMMAND surface and event log. They may draft definitions, query instances and trails, assist waits, propose COMMANDs, or derive views from definition + Events. They must not invent non-OMG core element types, skip COMMAND/EVENT, invent ledger subjects or empty intents, or mutate projections outside the log. How deeply any one consumer integrates with models or tools is outside this kernel's scope.

Static browser wiki/demo: sibling [`sparrow-playground`](../sparrow-playground/) (loads `wasm` artifacts; JS schedules timers and jobs). Cursor multi-root: [`sparrow-dev.code-workspace`](../sparrow-dev.code-workspace).

## API and persistence

`Deploy` · `CreateInstance` · `Complete` · `ThrowError` · `ResolveIncident` · `FireDue` · `PublishMessage` · `PublishSignal` · `EvaluateConditionalStarts` · `EvaluateConditions` · Job Activate / Fail / Heartbeat · `GetDeployment` · `GetInstance` · `ListEvents`

`EventLog` + `deploy.Store` + optional `runtime.Store` · `Recover` / `Open` · gRPC in `gateway` only · browser host in `wasm` (in-memory; JS schedules `FireDue` / jobs)

## Workspace

Go **1.26.5** · modules: `bpmn` · `protocol` · `processing` · `gateway` · `wasm`

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
go test ./bpmn/
cd protocol/proto && ./build.sh
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
cd wasm && ./build.sh   # → dist/sparrow.wasm (gitignored)
```

Edit `.proto`, then run `build.sh`. Do not hand-edit `*.pb.go`. Element semantics live in `handlers/` and [`processing/README.md`](./processing/README.md). Apache-2.0 headers.
