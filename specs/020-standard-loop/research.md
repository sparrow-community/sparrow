# Research: Standard Loop

## Semantics

**Decision**: BPMN `testBefore` (default false), `loopCondition`, optional `loopMaximum`. `loopCounter` is 1-based and incremented before each condition evaluation.

**Empty condition**: once if no max; exactly `loopMaximum` times if max set.

## vs multi-instance

**Decision**: Mutual exclusion at deploy. Standard loop reuses one token and re-enters; MI keeps host/inner tokens.

## Ledger / Recover

**Decision**: No new proto fields. Projection increments `Token.StandardLoopIteration` on each non-MI activity ACTIVATED while the token remains on that element; Recover rebuilds by replay.

## Re-enter

**Decision**: `Effect.ReEnter` after COMPLETED cancels iteration boundaries then calls OnEnter again without TakeOutgoing; compensation subscribe only when leaving the loop.
