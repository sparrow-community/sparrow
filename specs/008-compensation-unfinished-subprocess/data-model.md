# Data Model: Compensation into Unfinished SubProcess

## Entities (unchanged types)

### CompensationSub / PendingCompensation

Existing projection fields. Queue may include handlers whose subscriptions were discovered via unfinished-SubProcess entry (activity scope under terminated SP).

### Unfinished SubProcess (runtime)

| Signal | Meaning |
|--------|---------|
| Host token | Token with ElementID = SubProcess id, still in instance Tokens |
| Scope | `ScopeOf(SubProcessID) == throwScope` for direct-child entry |

### Handler queue item

| Field | Notes |
|-------|--------|
| boundaryID | Compensation boundary to COMPLETE |
| handlerID | isForCompensation activity |
| seq | from CompensationSub |

## State transitions

```text
Broadcast compensate while SP unfinished:
  → terminateScope(SP)  // host + inner tokens TERMINATED
  → collect same-scope completed subs + under-SP completed subs
  → sort Seq desc → PendingCompensation
  → Enter first handler → … → Complete throw

Targeted activityRef=SP (unfinished):
  → terminateScope(SP)
  → collect under-SP subs only
  → PendingCompensation → handlers → Complete throw
```

No new Event Element types.
