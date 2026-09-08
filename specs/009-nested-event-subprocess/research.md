# Research: Nested Event Sub-Process

## 1. Gap

**Decision**: Remaining item is ESP **inside ESP**, not ESP inside embedded SubProcess (already covered by `m4_nested_event_subprocess`).

**Rationale**: Deploy rejects `insideEventSubProcess` in `validateSubProcessesAt`.

## 2. Deploy

**Decision**: Remove the reject; still validate each Event Sub-Process (no sequence flows, one start, supported start kinds). Recurse with `insideEventSubProcess=true` for further nesting validation of other rules if any remain.

## 3. Arm on ESP enter

**Decision**: In `Executor.Enter` EnterChild path, call `emitEventSubProcessStartArms` for **all** SubProcess enters, including Event Sub-Process (drop `!IsEventSubProcess` guard). Mirror in `multi_instance_executor` if present for consistency (ESP is not MI).

**Rationale**: Embedded SP already arms this way; ESP enter currently skips arming, so nested ESPs never arm.

## 4. Disarm / interrupt

**Decision**: No change — `emitEventSubProcessStartDisarmInScope` and interrupting trigger already use `ParentScopeID`.

## 5. Recover

**Decision**: Arms are ledger ACTIVATED start events; replay restores. Add Open/Recover test.
