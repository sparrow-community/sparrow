# Research

**Decision**: Reuse `escalationCatch` map for intermediate catch ids; `CatchKind` / `IsEscalationBoundary` distinguish by element type.

**Decision**: Scope order: event sub-process → boundary on scope → waiting intermediate catches in scope → bubble.

**Decision**: Matching waits in a scope all complete (instance-local fan-in), then stop bubbling.
