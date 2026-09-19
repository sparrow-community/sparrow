# processing

BPMN execution kernel. Append-only `event.v1` ledger; projections rebuild from EVENTs.

Roadmap: [`../AGENTS.md`](../AGENTS.md) · 中文: [`../AGENTS.zh.md`](../AGENTS.zh.md)

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
| Process / Start / End | Lifecycle; **process-level typed starts** (message/timer/signal/conditional) mint via PublishMessage / PublishSignal / FireDue / EvaluateConditionalStarts; CreateInstance uses none start, instantiate event-based gateway entry(ies), or instantiate receive task; process-level error start rejected; end tries scope/process complete; **terminate end** cancels remaining tokens in the enclosing scope then completes it; **message/signal end** publish then complete |
| UserTask / ManualTask / Task | Wait → Complete; multi-instance and **standardLoopCharacteristics** supported (Task = abstract `bpmn:task`) |
| ServiceTask / BusinessRuleTask / ScriptTask | Wait as Job → Complete; multi-instance and standard loop supported (Script: job type from scriptFormat → name → id; no in-engine script runtime) |
| SendTask | Publish message then continue |
| ReceiveTask | Wait for PublishMessage; **instantiate** receive may be process entry |
| Catch / throw / boundary | Timer, message, signal, error, escalation (including standalone intermediate catch), compensate, **conditional** (EvaluateConditions); multiple same-kind per activity (unique message/signal names / condition text); link throw/catch; `implicitThrowEvent` as a flow element rejected at Deploy (only valid as a multi-instance behavior event) |
| Gateways | XOR, AND, inclusive, complex (activationCondition join; inclusive-style split), event-based (catch; exclusive/parallel; instantiate; targets may be catch or receive) |
| SubProcess | Embedded scope; Event Sub-Process (including compensation and conditional start via EvaluateConditions); multi-instance |
| Transaction | Transaction SubProcess (`##Compensate`); Cancel End → in-scope compensate → interrupting Cancel Boundary; nested/MI rejected at Deploy |
| AdHocSubProcess | Flat inner activities enabled without sequence flow; `ordering` Parallel (all) / Sequential (one at a time in document order); `completionCondition` finishes the scope, `cancelRemainingInstances` (default true) cancels inner activities still enabled; exhausting all inner activities also completes the scope; non-flat bodies and a missing `completionCondition` rejected at Deploy |
| CallActivity | Child instance; IO mapping (name copy, transformation, assignment); cross-deployment; boundary & compensation (including into unfinished child); multi-instance |
| SequenceFlow | `SEQUENCE_FLOW_TAKEN`; conditions on gateway and activity outgoings (+ `default`); multiple unconditional outs fan out in parallel; endpoints and declared `incoming`/`outgoing` must resolve at Deploy (a flow into an unmodeled element such as `choreographyTask` is rejected, not discovered at runtime) |

Full snapshot + roadmap: [`../AGENTS.md`](../AGENTS.md).

## Multi-instance

User Task, Service Task, Manual/Receive/Send/Business Rule/Script Task, embedded SubProcess, Call Activity. Host `loop_instance_index = -1`; inners `0..N-1`. Parallel / sequential / collection / early completion / boundary cancel / recover. Behavior None/One/All/Complex; `noneBehaviorEventRef` / `oneBehaviorEventRef` / `complexBehaviorDefinition` publish signal or message. Fixtures: `testdata/m6_mi_*.bpmn`, `m25_mi_*.bpmn`.

```shell
go test ./processing/
```
