# Research

**Decision**: Compensation event sub-process is nested only in an embedded SubProcess for this increment (matches C.6.0; process-root deferred).

**Decision**: Register as `deploy.Compensation` with start event id as subscription key so `subscribeCompensation` / `startCompensation` / projection stay unchanged.

**Decision**: Nested compensate throws remap `throwScope` from the compensation event sub-process id to its `ParentScopeID`.

**Decision**: Do not arm compensate starts as live event sub-process waits; they activate only via compensation subscription after the enclosing SubProcess completes.

**Decision**: Spell out “event sub-process” in AGENTS/README (no ESB abbreviation).
