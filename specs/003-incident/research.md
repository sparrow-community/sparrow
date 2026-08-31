# Research: Incident (blocked execution state)

## 1. Ledger subject for incidents

**Decision**: Incident lifecycle is recorded on the existing `SERVICE_TASK` element with two new intents: `INTENT_INCIDENT_OPENED` and `INTENT_INCIDENT_RESOLVED`. No new `Element.Type` and no separate `INCIDENT` ledger subject.

**Rationale**: Constitution III — job failure and operator resolve are lifecycle stages of the same activity token, analogous to FAILED and ACTIVATED. Audit shows element id, token id, and intent without a second ontology.

**Alternatives considered**:

- Reuse unused `INTENT_SUSPENDED` / `INTENT_RESUMED` — rejected because suspend/resume reads as process-wide BPMN suspend, not job-failure incident.
- Projection-only blocked flag without EVENT — rejected; FR-001 and Recover require append-only open/close facts.
- Separate incident UUID as ledger subject — second ontology; rejected.

## 2. Distinction from FAILED

**Decision**: Keep `INTENT_FAILED` for every worker-reported failure (audit trail unchanged). When policy says the token is blocked, append one additional `INTENT_INCIDENT_OPENED` EVENT in the same Fail command chain (after FAILED). Resolve appends `INTENT_INCIDENT_RESOLVED` only (no duplicate FAILED).

**Rationale**: Operators can see retry history (FAILED × N) plus a single clear blocked marker (INCIDENT_OPENED). Matches spec acceptance “FAILED/resolve events appear in the audit trail.”

**Alternatives considered**:

- Replace FAILED with INCIDENT_OPENED on threshold — loses per-attempt audit.
- Only increment a counter on FAILED payload without INCIDENT_OPENED — FR-002 blocked state not distinct in intent stream.

## 3. Retry threshold and non-retriable fail

**Decision**:

- **Default threshold**: `3` retriable fails per job activation cycle (count FAILED events since last ACTIVATED or INCIDENT_RESOLVED). The **3rd** FAILED opens an incident immediately after that FAILED EVENT.
- **Non-retriable**: `FailJobRequest.no_retry = true` opens an incident after that fail regardless of count (equivalent to threshold 1 for that command).
- **Per-activity override (optional MVP)**: BPMN extension `sparrow:failedJobIncidentThreshold` on `serviceTask` (positive integer). Deploy compiles into `deploy.ServiceTaskSpec.IncidentThreshold`; unset uses engine default 3.

**Rationale**: Small default matches spec assumption; explicit `no_retry` covers “human must intervene now”; extension hook avoids global-only policy without Camunda product surface.

**Alternatives considered**:

- Zeebe-style `retries` on each activated job — requires runtime store or recomputing from lease metadata; fail-count from ledger is Recover-safe without new store fields.
- Infinite retries until operator — current behavior; incident increment fixes visibility.

## 4. Projection and query surface

**Decision**:

- New token status string `blocked` (projection `TokenBlocked`; wire `Token.status = "blocked"`).
- While blocked: `job_type` retained for resolve; `ActivateJobs` skips blocked tokens; `FailJob` rejects.
- Token fields: `incident_error_message` (last error at open), `job_fail_count` (count at open, for operator context).
- `Complete` on blocked token → REJECTION code `INCIDENT_OPEN`.

**Rationale**: FR-002 and FR-007 — blocked is visibly distinct from `waiting` job retry; minimal new fields on existing `Token` message.

**Alternatives considered**:

- Keep status `waiting` + boolean — operators and tests must inspect nested flag; easier to miss in clients.
- Top-level `Instance.incidents` repeated message only — duplicates token identity; token-level fields sufficient for MVP (single incident per token max).

## 5. Resolve / retry command

**Decision**: New engine command `ResolveIncident(process_instance_id, element_id, token_id)` → COMMAND/EVENT pair with `INTENT_INCIDENT_RESOLVED` on `SERVICE_TASK`. Effect: token status `waiting`, `job_fail_count` reset to 0, incident fields cleared; job can be activated again.

**Rationale**: FR-004 explicit auditable COMMAND separate from anonymous Fail/Complete; same serial lock as other instance commands.

**Alternatives considered**:

- Overload `Complete` with empty vars as resolve — conflates success path with operator recovery; rejected.
- JobService-only RPC — engine API should exist in `processing`; gateway exposes both Engine and Job surfaces as today.

## 6. Interaction with boundaries and termination

**Decision**:

- **Interrupting boundary** on host with open incident: existing boundary path TERMINATES host token → projection clears blocked state on TERMINATED (no separate INCIDENT_RESOLVED required; optional INCIDENT_RESOLVED not emitted to avoid duplicate close facts).
- **Resolve on completed/terminated instance**: REJECTION `INVALID_STATE`.
- **Fail on User Task / non-job**: unchanged; no incident.

**Rationale**: Spec edge case “boundary semantics win”; one open incident per token; termination is authoritative close of host work.

**Alternatives considered**:

- Emit INCIDENT_RESOLVED before boundary terminate — extra EVENT with little audit value; host TERMINATED suffices.

## 7. Recover and idempotency

**Decision**: Rebuild `job_fail_count` and blocked state by replaying FAILED / INCIDENT_OPENED / INCIDENT_RESOLVED / ACTIVATED on each token. Redrive uses existing `alreadySeen(sourceCmdID, element)` — no duplicate INCIDENT_OPENED for the same Fail command.

**Rationale**: FR-006; same pattern as Complete and Fail today (`processing/recover.go` redrive paths).

**Alternatives considered**:

- Persist fail count only in runtime store — not ledger-backed; fails Recover test.

## 8. Activate during incident

**Decision**: `ActivateJobs` and in-memory lease claim skip tokens with status `blocked`. Direct activate path (if any) rejects with `INCIDENT_OPEN`.

**Rationale**: Blocked means operator must resolve first; prevents silent worker retry loops spec Story 1 rejects.

**Alternatives considered**:

- Allow activate but fail on complete — worker wastes cycles; worse UX.
