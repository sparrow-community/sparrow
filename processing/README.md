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
| Process / Start / End | Lifecycle; **process-level typed starts** (message/timer/signal/conditional) mint via PublishMessage / PublishSignal / FireDue / EvaluateConditionalStarts; CreateInstance uses none start or instantiate EBG only; process-level error start rejected; end tries scope/process complete; **terminate end** cancels remaining tokens in the enclosing scope then completes it; **message/signal end** publish then complete |
| UserTask / ManualTask / Task | Wait → Complete; MI supported (Task = abstract `bpmn:task`) |
| ServiceTask / BusinessRuleTask / ScriptTask | Wait as Job → Complete; MI supported (Script: job type from scriptFormat → name → id; no in-engine script runtime) |
| SendTask | Publish message then continue |
| ReceiveTask | Wait for PublishMessage; instantiate receive rejected |
| Catch / throw / boundary | Timer, message, signal, error, escalation, compensate, **conditional** (EvaluateConditions); multiple same-kind per activity (unique message/signal names / condition text); link throw/catch |
| Gateways | XOR, AND, inclusive, event-based (catch) |
| SubProcess | Embedded scope; Event Sub-Process; MI |
| CallActivity | Child instance; IO mapping; cross-deployment; boundary & compensation; MI |
| SequenceFlow | `SEQUENCE_FLOW_TAKEN`; conditions on gateway and activity outgoings (+ `default`) |

Full snapshot + roadmap: [`../AGENTS.md`](../AGENTS.md).

## Multi-instance

User Task, Service Task, Manual/Receive/Send/Business Rule/Script Task, embedded SubProcess, Call Activity. Host `loop_instance_index = -1`; inners `0..N-1`. Parallel / sequential / collection / early completion / boundary cancel / recover. Fixtures: `testdata/m6_mi_*.bpmn`.

```shell
go test ./processing/
```
