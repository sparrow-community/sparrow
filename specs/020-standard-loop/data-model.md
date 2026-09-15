# Data model: Standard Loop

## Definition

`StandardLoopSpec`: ElementID, TestBefore, LoopMaximum (0 = unset), ConditionExpr.

Indexed on Deployment beside MultiInstanceSpec.

## Runtime

`Token.StandardLoopIteration`: 1-based count of ACTIVATED on this element for the current visit (0 before first activation).

## No new EventLog record types
