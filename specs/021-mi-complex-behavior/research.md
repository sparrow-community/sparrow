# Research: MI Complex Behavior

## Behavior vs completion

**Decision**: Keep existing One/All completion mapping. None completes like All (or completionCondition). Complex does not change completion by itself—`completionCondition` / All-style join still apply unless One.

**oneBehaviorEventRef** is additive publish on first complete; One completion already cancels siblings.

## Throw kinds

**Decision**: Signal and message only via existing `Publication`.

## Complex one-fire / Recover

**Decision**: `MultiInstanceLoop.ComplexFired []bool`; on Recover, mark definitions whose condition is already true as fired without publishing.
