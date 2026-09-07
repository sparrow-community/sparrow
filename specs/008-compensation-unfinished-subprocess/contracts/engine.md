# Contracts: Compensation into Unfinished SubProcess

No new RPCs or proto fields.

## Compensate throw / end (runtime)

### Broadcast (empty activityRef)

1. Run compensation for completed activities in the throw’s scope (existing).
2. For each unfinished direct-child embedded SubProcess of that scope: terminate its active work; enqueue handlers for completed activities under that SubProcess (reverse completion order).
3. Throw waits until the combined handler queue is empty, then continues.

### Targeted activityRef

1. If activityRef names an **unfinished** embedded SubProcess in scope: terminate that SubProcess; enqueue only under-SP completed handlers; do not consume other same-scope subscriptions.
2. If activityRef names a completed subscribed activity (including completed SubProcess host subscription): existing behavior.
3. If activityRef matches nothing applicable: empty queue; throw completes (existing empty-queue behavior).

## Recover

After terminate + handler Enter recorded, Recover restores waiting throw and active compensation handler; Complete on handler finishes the chain.

## Out of scope

Call Activity unfinished children; Transaction/cancel; compensation Event Sub-Process handlers.
