# Research: Process-Level Typed Start

## CreateInstance vs typed triggers

**Decision**: CreateInstance mints an instance only via a **none** start (empty event definitions) or the existing instantiate exclusive EBG entry. Typed-only deployments reject CreateInstance with a stable `INVALID_ARGUMENT` / `INVALID_STATE` code. Message / timer / signal / conditional triggers mint via a shared helper that appends PROCESS ACTIVATING… then Enter at the chosen start event id.

**Rationale**: Matches BPMN alternative starts and keeps today’s none-start suite unchanged. Spec FR-003.

**Alternatives considered**: Treat CreateInstance as “force start ignoring type” — rejected (hides typed semantics). Auto-pick first start regardless of type — rejected (current bug-like behavior for typed defs).

## Waiter-first then start creation

**Decision**: On PublishMessage / PublishSignal: complete matching waiters, scope boundaries, and Event Sub-Process arms first. If `delivered > 0`, do **not** create process-level start instances for that publish. If nothing delivered, create one instance per matching process-level start subscription (across deployments), then apply existing message buffer rules only when no start subscription matched (or after starts, do not require a catch for the same publish).

**Rationale**: Spec assumptions; avoids stealing an in-flight catch. Signals remain unbuffered for catches; start creation is the unmatched path.

**Alternatives considered**: Always create starts in addition to waiter delivery — rejected (duplicate fan-out). Buffer-only without start creation — rejected (no instantiation).

## Instantiate EBG unchanged (007)

**Decision**: PublishMessage / PublishSignal / FireDue MUST NOT create instances solely because an instantiate EBG is deployed. Typed **startEvent** subscriptions are the only new auto-create registry.

**Rationale**: Spec 007 FR-010 remains binding; typed starts are a different construct.

## Timer start schedules outside the ledger

**Decision**: On Deploy (and Recover rebuild from deployments + runtime store), arm process-level timer starts in engine/runtime helper state: `{deploymentID, startEventID, dueUnixMs, timerSpec}`. FireDue creates an instance at that start when due. Cycle timers re-arm next due after create; date/duration are one-shot (remove arm after create). Duration due = arm time + duration (arm at deploy / recover restore time rules consistent with catch TimerDue).

**Rationale**: Constitution: job/timer helpers outside ledger. No fake instance waiting at start before creation.

**Alternatives considered**: Create a dormant instance at Deploy that waits on the start event — rejected (invents pre-start instance lifecycle).

## Conditional start API

**Decision**: Add Engine method `EvaluateConditionalStarts(ctx, req)` with deployment_id and/or process_id (+ version) and variables. For each matching conditional start on the resolved deployment(s), `expr.Eval` the condition; on true, mint one instance at that start with those variables. False → skip. CreateInstance still rejects typed-only. Gateway RPC optional this increment (FireDue is already engine-only).

**Rationale**: Spec needs an evaluation trigger; reuses `processing/expr`. Avoids overloading CreateInstance.

**Alternatives considered**: CreateInstance-with-vars evaluates condition — rejected by FR-003. Polling daemon — out of scope for lightweight kernel.

## Process-level error start

**Decision**: Reject at Deploy (`UNSUPPORTED_ELEMENT`). ESP error starts unchanged.

**Rationale**: OMG restricts error start to Event Sub-Process; AGENTS Planned wording updated accordingly.

## Call Activity + typed-only called process

**Decision**: When indexing Call Activity, require the resolved called process to expose a none start or instantiate exclusive EBG. Typed-only called processes → Deploy reject on the caller (or when binding the called deployment).

**Rationale**: Spec assumption; Call Activity enter path uses CreateInstance-equivalent mint at child start today.

## Mixed alternative starts

**Decision**: Index all process-level starts. `NoneStartEventID` returns the none start if any. CreateInstance uses that. Typed indexes are separate maps (message name → []start refs, etc.). Multiple starts with the same message name each get an instance on unmatched publish.

**Rationale**: BPMN alternative starts; unambiguous CreateInstance entry.

## StartEvent OnEnter

**Decision**: Keep StartEventHandler: lifecycle complete + TakeOutgoing. Typed “catch” is the mint trigger; Enter runs after the instance exists.

**Rationale**: Same as none start after creation; no wait payload on process-level start token.
