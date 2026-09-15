# Engine contract: Conditional catch

No proto change.

## Deploy

Accepts intermediate/boundary `conditionalEventDefinition` alone with non-empty condition. Rejects duplicate identical conditions on one activity.

## Commands

| Situation | Command |
|-----------|---------|
| Waiting intermediate conditional | `EvaluateConditions` |
| Armed conditional boundary | `EvaluateConditions` |
| True at intermediate enter | none (completes in Enter) |

## Recover

Waiting tokens / BoundaryWaits with kind conditional rebuild from events; EvaluateConditions then works.
