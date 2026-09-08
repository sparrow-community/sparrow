# processing

BPMN execution kernel. Append-only `event.v1` ledger; projections rebuild from EVENTs.

Roadmap: [`../AGENTS.md`](../AGENTS.md) · AI Driven: [`../AI-Driven-BPMN.md`](../AI-Driven-BPMN.md)

Event = behavior; Element (Type, id, token_id, Intent, payload) = subject.

## Layout

```text
engine.go, executor.go, handlers/, deploy/, projection/, log/, runtime/, expr/
```

```text
Engine → Executor (Effect → emit) → handlers.OnEnter / OnComplete
```

New element: handler + registry + deploy validation. Semantics stay out of `engine.go`.

## Execution

- One lock per `process_instance_id`
- COMMAND → EVENT(s) (`source_record_id`); invalid → REJECTION
- Waiting: one `Complete`. After unlock: `Publication` (message/signal, CallActivity child)
- Tokens: EVENT → `ApplyEvent` only. `Recover`: replay + redrive unfinished COMMANDs

## Elements (summary)

| Type | Behavior |
|------|----------|
| Process / Start / End | Lifecycle; end tries scope/process complete |
| UserTask / ServiceTask | Wait → Complete; MI supported |
| Catch / throw / boundary | Timer, message, signal, error, escalation, compensate; multiple same-kind per activity (unique message/signal names) |
| Gateways | XOR, AND, inclusive, event-based (catch) |
| SubProcess | Embedded scope; Event Sub-Process; MI |
| CallActivity | Child instance; IO mapping; same-definition today |
| SequenceFlow | `SEQUENCE_FLOW_TAKEN` |

Full snapshot + roadmap: [`../AGENTS.md`](../AGENTS.md).

## Multi-instance

User Task, Service Task, embedded SubProcess. Host `loop_instance_index = -1`; inners `0..N-1`. Parallel / sequential / collection / early completion / boundary cancel / recover. Fixtures: `testdata/m6_mi_*.bpmn`.

```shell
go test ./processing/
```
