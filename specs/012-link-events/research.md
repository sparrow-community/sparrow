# Research: Link Events

## Semantics

**Decision**: Camunda/Zeebe-style — throw transfers to all catches with the same name in the process; catch is instantaneous (no wait).

## Pairing

**Decision**: Match `linkEventDefinition/@name`, else element `name`, else element id. Scope = entire process tree (nested SubProcesses included).

## Topology

**Decision**: Reject link throw with outgoing sequence flows and link catch with incoming sequence flows (classic off-page connector shape).

## Executor

**Decision**: New `Effect.LinkContinue []string` — Enter each catch; first reuses throw token, others mint tokens.
