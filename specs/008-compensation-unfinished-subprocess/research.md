# Research: Compensation into Unfinished SubProcess

## 1. Semantics vs Camunda default

**Decision**: Implement parent-scope compensate **entering unfinished** embedded SubProcesses (terminate + run completed inner subscriptions). This is an intentional completeness behavior; Camunda 7 documents the opposite default for broadcast into unfinished SPs. Sparrow Remaining explicitly lists this gap.

**Rationale**: Spec FR-001–FR-003; enables concurrent “cancel in-flight book while undoing completed inners” patterns without Transaction/cancel.

**Alternatives considered**:
- Match Camunda “never into unfinished” — contradicts Remaining item.
- Transaction SubProcess + cancel boundary — larger surface; deferred (FR-007).

## 2. Where to hook

**Decision**: Extend `Executor.startCompensation` only. Before/while building the handler queue:

1. Resolve `activityRef` and `throwScope` as today.
2. If `activityRef` names an **active** SubProcess (host token still present) under throwScope → `terminateScope(that SP)` then collect subscriptions with activity under that SP.
3. If `activityRef` empty (broadcast) → collect same-scope completed subscriptions as today; also for each **active** direct-child SubProcess under throwScope, terminate and collect under-SP subscriptions.
4. If `activityRef` names a completed activity / completed SP subscription → unchanged path (no unfinished entry).
5. Merge collected items, sort by `CompensationSub.Seq` descending, emit boundary COMPLETED for consumed boundaries, set PendingCompensation, `advanceCompensation`.

**Rationale**: Single choke point for all compensate throws/ends (`TriggerCompensation`).

## 3. Finding unfinished SubProcesses

**Decision**: An unfinished embedded SubProcess P under throwScope is present when some token has `ElementID == P`, type `TYPE_SUB_PROCESS`, and `ScopeOf(P) == throwScope` (P is direct child of throw scope). Use `terminateScope` (IncludeHost) to cancel host + descendants.

**Rationale**: Matches existing SubProcess host-token model from error/ESP terminate paths.

## 4. Subscription selection under SP

**Decision**: Include subscription if `activityInScopeTree(activityID, P)` — walk `ScopeOf` upward until empty; match P. That covers direct children and nested scopes under P after terminate of the whole subtree.

**Rationale**: Spec edge case on nested unfinished; one terminate of P clears descendants.

## 5. Targeted vs broadcast

**Decision**:
- `activityRef == unfinished SP` → only that SP’s inner queue; do not consume other same-scope subs.
- `activityRef == ""` → same-scope completed + all unfinished direct-child SPs’ inners.
- `activityRef == completed activity` → existing equality on `c.ActivityID`.

**Rationale**: FR-003, FR-004.

## 6. Ordering

**Decision**: One merged list sorted by subscription `Seq` descending (global reverse completion), including inners discovered after terminate. Boundary COMPLETED emit for each consumed boundary id as today.

**Rationale**: Preserves reverse-order compensation across concurrent branches.

## 7. Recover

**Decision**: No change expected to `rebuildPendingCompensation` if throw token still waiting and handler waiting; terminate events already in log. Add Open/Recover test.

**Rationale**: FR-006; constitution II.

## 8. Protocol

**Decision**: No proto/API changes.
