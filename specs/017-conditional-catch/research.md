# Research: Conditional Catch

## Trigger surface

**Decision**: Add `EvaluateConditions` (optional ProcessInstanceID + Variables overlay). Separate from `EvaluateConditionalStarts`. On success, Complete the catch token or fire the boundary with Variables.

**Rationale**: Matches starts pattern; no continuous poller in the kernel.

## True at enter

**Decision**: Intermediate catch OnEnter evaluates against instance variables; if true, InstantLifecycle + TakeOutgoing; else Wait. Boundaries always arm; caller EvaluateConditions (even immediately after start) fires them.

**Rationale**: Spec SC-001; boundaries share arm-then-trigger with timer (FireDue).

## Condition storage

**Decision**: `conditionalCatch[id]=expression` on Deployment. Projection does not need ConditionText; evaluate looks up by element/boundary id. WaitingBoundary.Kind = `"conditional"`.

**Rationale**: Avoids proto change; Recover rebuilds waits from ACTIVATED payload kinds.
