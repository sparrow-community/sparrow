# Research: Message End and Signal End

## Publish on end vs throw

**Decision**: Message/signal end OnEnter emits InstantLifecycle (or ACTIVATING/ACTIVATED/COMPLETING/COMPLETED with payload name), sets `Publish` PublicationMessage/Signal, and `TryCompleteProcess` — same pattern as intermediate throw plus end scope completion.

**Rationale**: BPMN end with message/signal definition throws then ends the token. Existing flushPublications already delivers to waiters and typed starts.

**Alternatives considered**: Separate end-only publication kind — rejected (duplicate). Job-backed send — rejected (parity with throw/Send Task).

## Indexing

**Decision**: `messageEnds[id]=name`, `signalEnds[id]=name` on Deployment; lookup helpers `MessageEndName` / `SignalEndName`. Parse with same resolve rules as throw (`messageRef`/`signalRef`).

**Rationale**: Matches errorEnds / escalationEnds indexing style.

## Mixed definitions

**Decision**: Exactly one message or one signal definition; any other defs (terminate, error, timer, etc.) → UNSUPPORTED_ELEMENT at Deploy (extend existing `endHasEventDefinitions` reject path by accepting message/signal first).

**Rationale**: Spec FR-005; silence is not exclusion.
