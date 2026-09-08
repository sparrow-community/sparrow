# Data Model: Nested Event Sub-Process

## Entities

### EventSubProcess (compile-time)

| Field | Notes |
|-------|--------|
| ID | Nested ESP id |
| ParentScopeID | Outer Event Sub-Process id (not only process / embedded SP) |

### Arm (runtime)

Existing `Instance.EventSubProcesses[nestedID]` while outer ESP scope is active.

## Transitions

```text
outer ESP triggered → Enter(outer) → arm nested ESP(s)
nested interrupting trigger → terminate tokens in outer ESP scope → Enter(nested)
outer ESP completes → disarm nested ESP(s) in that scope
```
