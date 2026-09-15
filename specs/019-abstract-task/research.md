# Research: Abstract Task

## Manual-equivalent vs Job

**Decision**: Wait → Complete only (same as Manual Task). Not job-backed.

**Rationale**: AGENTS Planned wording; BPMN abstract task has no execution semantics beyond a generic activity; modelers use it like a placeholder human step. Mapping to Service/Script would invent job semantics the diagram does not declare.

## Ledger type

**Decision**: Use existing `TYPE_TASK`; do not emit `TYPE_MANUAL_TASK` for `<task>`.

**Rationale**: Element coverage requires distinct Supported combinations; silence/alias would blur Manual vs abstract.

## Deploy rejection removal

**Decision**: Remove `countAbstractTasks` gate in `validateM1`.

**Alternatives**: Keep reject and only accept via extension — rejected; Planned item is acceptance.
