# processing

Sparrow’s BPMN **execution kernel**: single-node, event-driven. Append-only `event.v1` ledger is the source of truth; projections are rebuildable.

**Not here:** BPMN XML parse (`bpmn`), wire types (`protocol`), gRPC (`gateway`).

| | |
|--|--|
| Roadmap / implemented | [`../AGENTS.md`](../AGENTS.md) |
| Governance | [`../.specify/memory/constitution.md`](../.specify/memory/constitution.md) |
| Active increment | [`../specs/001-engine-completeness/`](../specs/001-engine-completeness/) |

Event = behavior; Element (Type, id, token_id, Intent, payload) = subject. Job / timer / message / signal waits are payloads—not peer ledger subjects.

## Layout

```text
processing/
├── engine.go, jobs.go, timers.go, messages.go, call_child.go, open.go, recover.go
├── executor.go          # Enter / Complete, Effects, scope completion
├── handlers/            # one handler per Element.Type (semantics live only here)
├── deploy/              # Compile + Store; queries on element.Process (no parallel graph)
├── projection/          # Instance / Token; ApplyEvent (EVENT only)
├── log/, runtime/       # EventLog; optional leases + message buffer
└── expr/                # ${...} → expr-lang
```

```text
Engine (API, instance lock, COMMAND/EVENT/REJECTION)
  → Executor (handler Effect → emit → step / wait / complete)
  → handlers.OnEnter / OnComplete
```

Extend coverage: **new handler + Registry + deploy validation**. Keep semantics out of `engine.go`.

## Execution model

- Partition key: `process_instance_id` (one lock per instance).
- COMMAND → handler → EVENT(s) with shared `source_record_id`; failure → REJECTION.
- Waiting work uses one `Complete` (type from deployment). Deferred work after unlock: message/signal publish, CallActivity child start / parent resume (`Publication`).
- Tokens update only via EVENT → `ApplyEvent`. `Recover` replays EVENTs, then redrives unfinished COMMANDs. `Open(dataDir)` = file EventLog + stores.

```go
eng, err := processing.Recover(ctx, eventLog, deploymentStore, runtimeStore)
```

## Element contract (summary)

| Type | Behavior |
|------|----------|
| `PROCESS` | Start / complete (or terminate) |
| `START_EVENT` | Instant; Event Sub-Process starts arm on scope open (`ACTIVATED`; message/timer/signal/error) |
| `END_EVENT` | Then try complete scope/process |
| `USER_TASK` / `SERVICE_TASK` | Wait → `Complete`; ServiceTask `job_type` |
| Catch / timer·message·signal | Wait → `FireDue` / `Publish*` / `Complete` |
| Throw | Instant; message/signal via Publication; compensate handlers |
| `BOUNDARY_EVENT` | Interrupt terminates host; non-interrupt spawns token; compensate after COMPLETED |
| Gateways | Exclusive / parallel / inclusive / event-based (no instantiate) |
| `SUB_PROCESS` | Embedded; may host Event Sub-Process |
| `CALL_ACTIVITY` | Same-definition `calledElement` → **child process instance**; host waits with `called_process_instance_id`; IO name mappings optional |
| `SEQUENCE_FLOW` | `SEQUENCE_FLOW_TAKEN` |

APIs: `Deploy`, `CreateInstance` (by `deployment_id` or `process_id` + optional version), `Complete`, `ThrowError`, `FireDue`, `PublishMessage`, `PublishSignal`, Job `Activate`/`Fail`/`Heartbeat`, `GetInstance`, `ListEvents`. Transport is `gateway` only.

**Open debt:** CallActivity boundary/compensation thinner than SubProcess; cross-deployment call; live version migration — see deferred list in `AGENTS.md`.

## Test

```shell
go test ./processing/
```

Fixtures under `testdata/`.
