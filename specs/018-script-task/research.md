# Research: Script Task

## In-engine vs job

**Decision**: Job-backed only; workers execute scripts outside the engine.

**Rationale**: Matches Business Rule / constitution; avoids embedding language VMs.

## Job type

**Decision**: scriptFormat → name → id.

**Alternatives**: Always element id — weaker for worker subscription by language.
