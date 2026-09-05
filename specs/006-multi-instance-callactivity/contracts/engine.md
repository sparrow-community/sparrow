# Contracts: Multi-Instance Call Activity

No new `EngineService` or `JobService` RPCs. Behavior is deploy validation + existing Call Activity / multi-instance COMMAND/EVENT semantics.

## Deploy contract

- `callActivity` MAY include `multiInstanceLoopCharacteristics` (`isSequential`, `loopCardinality`, collection refs, `completionCondition`, behavior All/One).
- Deploy MUST accept MI Call Activity when characteristics compile; MUST reject complex behavior / unsupported behavior event refs (same as other MI hosts).
- Deploy MUST NOT reject solely because Call Activity has multi-instance loop characteristics.
- Bundled vs external `calledElement` rules unchanged (`004`).

## Runtime contract

### Parallel / sequential spawn

1. Parent token enters Call Activity → if MI, host ACTIVATED (parked) then inner ACTIVATED records (one per started index).
2. Each inner ACTIVATED carries `called_process_instance_id` and `loop_instance_index`.
3. Engine starts one child process instance per inner (same-file or cross-deploy).
4. Sequential: at most one active child at a time.

### Child completion → parent join

1. Child process COMPLETED → resume parent on the **inner** token.
2. Inner CALL_ACTIVITY COMPLETED increments loop counters / evaluates completion condition.
3. When join satisfied: host Call Activity COMPLETED once; single outgoing sequence flow.
4. Straggler inners and their children terminated when early completion applies.

### IO mapping

1. Input mappings (+ MI collection element for that index) applied when each child starts.
2. Output mappings applied when each child completes; optional MI output collection assembled on host complete.

### Boundaries / compensation

1. Timer/message/signal/error/compensation boundaries attach to the Call Activity **host** (parity with `005` + MI host arming).
2. Interrupting boundary cancels all inners and terminates all active children; boundary path taken once.
3. Compensation subscription on successful **host** complete only.

### Recover

Mid-loop Recover restores host + remaining inners + child instances; finishing remaining children yields the same parent outcome as a continuous run.

## event.v1

No schema changes expected. Relevant existing fields:

- `ActivityPayload.loop_instance_index`, `loop_total_instances`, `loop_sequential`
- `ActivityPayload.called_process_instance_id`
- Boundary fields on host ACTIVATED

## engine.v1

No changes. Operators complete **child** waiting elements with existing `Complete`; inspect loop/child linkage via `GetInstance` tokens.

## processing Go API (non-wire)

- `Publication` may carry snapped child input variables for `StartChild`.
- `runMultiInstanceStart` must propagate publications from inner enters.
- `cancelMultiInstanceActivity` must terminate called children for Call Activity inners.
