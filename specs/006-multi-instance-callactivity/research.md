# Research: Multi-Instance Call Activity

## 1. Deploy: allow and index MI on Call Activity

**Decision**: Remove the `UNSUPPORTED_ELEMENT: multi-instance callActivity … not supported` reject in `validateCallActivity`. When indexing Call Activities, call `indexMultiInstance(d, ca.ID, ca.LoopCharacteristicsElements)` (same as User/Service/SubProcess). Keep existing CallActivity IO / external-callee validation.

**Rationale**: Spec FR-001; `bpmn` already parses loop characteristics on activities including Call Activity.

**Alternatives considered**:
- Separate CallActivity-only MI compiler — duplicates `compileMultiInstance`.
- Keep reject until a different increment — contradicts Remaining P2 choice.

## 2. Handler shape: host vs inner

**Decision**: Mirror `UserTaskHandler` / `SubProcessHandler`:

- **Host enter** (`MultiInstanceSpec` present and not inner): `multiInstanceHostEnter` (arms boundaries on host payload, parks host token, `MultiInstanceStart`).
- **Inner enter** (loop index ≥ 0): existing Call Activity path — allocate child id, `attachBoundary` only if needed for non-host (inners typically no separate boundary arm; host holds timer/message/signal), `PublicationStartChild` with `HostTokenID` = **inner** token id.
- **Non-MI**: unchanged single-child path from `004`/`005`.
- **Inner complete** (resume parent after child COMPLETED): `multiInstanceInnerComplete` instead of `TakeOutgoing` on the host; host complete path (when join fires) cancels attached boundaries + `subscribeCompensation` once.

**Rationale**: Spec FR-002–FR-004, FR-013; reuses proven MI host/inner split.

**Alternatives considered**:
- One child process with internal MI — not BPMN Call Activity MI.
- Host token owns all `called_process_instance_id`s in a list — harder Recover and Complete targeting.

## 3. Publications from MI spawn must not be dropped

**Decision**: Change `runMultiInstanceStart` / `enterWithLoopIndex` usage so `Publication` values from each inner `OnEnter` are returned and appended to the parent `Enter` pubs (then `flushPublications` as today). Today `runMultiInstanceStart` discards `enterWithLoopIndex` pubs — fine for User/Service Task, **broken for Call Activity** (would never start children).

**Rationale**: Constitution VI / existing publication pattern; Call Activity child start is deferred unlock work.

**Alternatives considered**:
- Synchronously start children inside the handler — couples locks across instances; rejected.
- Special-case Call Activity inside `runMultiInstanceStart` without returning pubs — opaque; prefer generic pub collection.

## 4. Per-iteration IO / collection snapshot

**Decision**: When building `PublicationStartChild` for an MI inner, **snapshot** the mapped child input variables onto the publication (new optional field, e.g. `ChildVariables` / `InputVariables`). `startCalledInstance` prefers the snapshot over re-reading `parent.Variables` at flush time.

**Rationale**: Parallel MI currently writes `inputDataItem` into `inst.Variables` per index before each inner enter; deferred flush would otherwise see only the last index’s element. Spec FR-010 / Story 3.

**Alternatives considered**:
- Flush each `StartChild` immediately inside the MI start loop — works but differs from other publication batching; snapshot is clearer and Recover-friendly if COMMAND already recorded child id.
- Store element only on token payload — still need mapped Call Activity inputs at start.

## 5. Cancel / interrupting boundary → terminate all children

**Decision**: Extend `cancelMultiInstanceActivity` (and any early-completion cancel path) so that for each terminated inner Call Activity token with `CalledProcessInstanceID`, emit/queue `PublicationTerminateChild` (or call `terminateCalledInstance` after unlock via returned pubs). Interrupting boundary already sets `MultiInstanceCancel` on MI hosts; ensure Call Activity inners are covered.

**Rationale**: Spec FR-008, SC-004/SC-005; single-call terminate already exists for non-MI.

**Alternatives considered**:
- Rely only on parent token TERMINATED without child terminate — leaves orphan children; rejected.

## 6. Sequential mode

**Decision**: Reuse existing MI sequential spawn (index 0 first; next index after inner complete). Each sequential step is one Call Activity inner enter → one child → resume → inner complete → next.

**Rationale**: Spec US2; no Call-Activity-specific sequential design.

## 7. Completion condition and output collection

**Decision**: Reuse `MultiInstanceSpec.CompletionMet` and output-item assembly from `002`. On each successful child resume, mapped Call Activity outputs become the inner complete variables / output collection item when `loopDataOutputRef` is configured.

**Rationale**: Spec FR-006–FR-007, FR-010.

## 8. Recover

**Decision**: No new ledger subjects. Replay rebuilds `MultiInstanceLoops` and per-inner tokens with `called_process_instance_id` as today. Child instances remain separate process instances with parent linkage. Redrive unfinished child CREATE / parent resume COMMANDs unchanged. Cross-deploy: `CalledDeploymentID` on StartChild publication / child log `deployment_id` prevents rebinding to a newer callee revision mid-redrive (same as `004`).

**Rationale**: Spec FR-011, SC-006; constitution II.

**Alternatives considered**:
- Persist Pending MI state outside the log — rejected.

## 9. Protocol / API

**Decision**: No proto changes expected. Clients continue to complete **child** waiting work via existing Complete; parent Call Activity completes via resume path. `GetInstance` already exposes loop index and called instance id on tokens.

**Rationale**: Spec assumptions; token_id remains the Completes target for child work.

## 10. Compensation

**Decision**: Subscribe compensation when the **multi-instance host** completes successfully (after join), not per unfinished child. Interrupted / early-cancelled loops do not subscribe — same as User Task MI host complete path.

**Rationale**: Spec edge case + FR-013.
