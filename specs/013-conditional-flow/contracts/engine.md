# Contract: Conditional outgoing selection

1. If preferred flow id is set → take it.
2. Else evaluate non-default outgoings in model order: empty condition = match; else `expr.Eval`.
3. Else take `default` if present.
4. Else if no conditions exist on any outgoing → take first (legacy).
5. Else → `NO_OUTGOING_FLOW`.
