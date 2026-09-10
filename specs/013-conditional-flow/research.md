# Research: Conditional Flow

## Gap

XOR/Inclusive already choose by condition. Activity multi-out ignores conditions (`outs[0]`).

## Decision

Shared exclusive chooser: first matching non-default condition, else default. Activities use `Activity.default`; exclusive gateway uses `ExclusiveGateway.default`. Inclusive unchanged.

## takeOutgoing

Pass instance variables; if `preferredFlowID == ""`, call chooser.
