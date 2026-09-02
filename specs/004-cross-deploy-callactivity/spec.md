# Feature Specification: Cross-deployment Call Activity

**Feature Branch**: `004-cross-deploy-callactivity`

**Created**: 2026-09-02

**Status**: Draft

**Input**: User description: "Engine increment 004 — Cross-deployment CallActivity: Call Activity resolves called process from a separately deployed definition; child instance uses callee deployment; IO mapping and Recover preserved. Next P1 per AGENTS.md Remaining list."

## User Scenarios & Testing *(mandatory)*

Callers can already deploy a definitions file that contains both caller and called processes, start a caller instance, and have a Call Activity spawn a **separate child instance** with input/output variable mapping, boundary interruption, and Recover-safe parent/child linkage. Today the called process **must** live in the same BPMN definitions as the caller: deploy compiles both together, and the child instance shares the caller's deployment id. Authors who split processes into independent deployable units (shared subprocess libraries, versioned services) cannot call a process that was deployed on its own.

### User Story 1 - Call a separately deployed process (Priority: P1)

An author deploys process `OrderFlow` (caller) and process `PaymentFlow` (callee) as **two deployments**. `OrderFlow` contains a Call Activity whose `calledElement` is `PaymentFlow`. When an instance reaches that Call Activity, the engine starts a new instance of `PaymentFlow` using the **callee's deployment** (not the caller's). The caller waits at the Call Activity until the child completes; parent/child instance linkage is queryable as today.

**Why this priority**: Without cross-deployment resolution, every reusable subprocess must be bundled into the caller's BPMN file — blocking independent versioning and deployment of shared flows.

**Independent Test**: Deploy callee BPMN, then deploy caller-only BPMN referencing `calledElement=PaymentFlow`. Start caller. `GetInstance` on caller shows waiting Call Activity with `called_process_instance_id`; `GetInstance` on child shows `PaymentFlow`, callee deployment id, and parent linkage.

**Acceptance Scenarios**:

1. **Given** callee `PaymentFlow` deployed and caller `OrderFlow` deployed separately, **When** a caller instance reaches the Call Activity, **Then** a child instance of `PaymentFlow` is created with the callee deployment id and the caller waits at the Call Activity.
2. **Given** a running cross-deployment call, **When** the child instance completes, **Then** the caller completes the Call Activity and continues along its outgoing flow.
3. **Given** caller deploy with no embedded callee process, **When** deploy is requested, **Then** deploy succeeds as long as the Call Activity declares a non-empty `calledElement` (resolution is runtime, not compile-time embed).
4. **Given** callee process is not deployed at call time, **When** the Call Activity activates, **Then** the command fails with a stable rejection code and the caller does not silently skip the call.

---

### User Story 2 - Variable mapping across deployment boundary (Priority: P1)

Input and output mappings declared on the Call Activity continue to work when caller and callee are in different deployments. Unmapped caller variables do not appear on the child; unmapped child variables do not leak back to the caller.

**Why this priority**: Cross-deployment calls without IO mapping would produce isolated but useless child instances.

**Independent Test**: Caller with `{orderId: "A1"}`, map `orderId` → child `id`, child sets `status`, map `status` → caller `paymentStatus`. After call, caller has `paymentStatus` and `orderId`; child never had unrelated caller fields.

**Acceptance Scenarios**:

1. **Given** input mapping on the Call Activity, **When** the child instance starts, **Then** only mapped variables are copied from caller to child.
2. **Given** output mapping on the Call Activity, **When** the child completes, **Then** mapped variables are copied onto the caller before the Call Activity completes.
3. **Given** no mappings, **When** the child starts and completes, **Then** the child starts with an empty variable set and the caller's variables are unchanged by the call.

---

### User Story 3 - Recover with cross-deployment parent and child (Priority: P2)

After a crash, Recover rebuilds a waiting Call Activity host and its child instance even when they reference **different** deployment ids. Completing the child after Recover still resumes the caller.

**Why this priority**: Completeness definition requires Recover tests; cross-deployment must not be a special-case memory-only linkage.

**Independent Test**: Start cross-deployment call, stop before child completes, Recover from EventLog only, complete child, caller reaches end.

**Acceptance Scenarios**:

1. **Given** a persisted parent waiting on Call Activity and an active child in another deployment, **When** Recover runs, **Then** both instances and deployment ids match pre-crash projection.
2. **Given** Recover after child create was commanded but child PROCESS events were partial, **When** redrive finishes, **Then** exactly one child instance exists and parent linkage is consistent (same idempotency rules as same-deployment calls).

---

### Edge Cases

- What happens when multiple revisions of the callee process are deployed? Default uses **latest** revision for `calledElement` (same rule as starting by `process_id` without version); exact pinning deferred unless BPMN extension added in plan.
- What happens when an interrupting boundary fires on the Call Activity while the cross-deployment child runs? Child is terminated; caller follows boundary path (parity with same-deployment Call Activity).
- What happens when the callee deployment is removed from the engine while a child instance is still active? Existing child instances continue; new calls fail until callee is deployed again.
- What happens when `calledElement` names a process that exists only inside another caller's bundled definitions but was never deployed standalone? Call activation fails at runtime with NOT_FOUND semantics.
- Nested Call Activity: child may call another cross-deployment or same-deployment process; each invocation creates a new instance with correct deployment resolution per level.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Deploy MUST accept a caller definitions file whose Call Activity `calledElement` references a process id **not** embedded in that file, when the reference is a non-empty process id string.
- **FR-002**: On Call Activity activation, the engine MUST resolve the called process to a deployed definition (default: latest revision of `calledElement`) and start a **new** child instance bound to that callee deployment id.
- **FR-003**: Child instance EVENTs MUST record the callee `deployment_id` and `process_version`; parent linkage (`parent_process_instance_id`, `parent_element_id`, `parent_token_id`) MUST remain on the child PROCESS activation as today.
- **FR-004**: Caller instance MUST remain waiting at the Call Activity with `called_process_instance_id` set until the child completes or is terminated (boundary/error path).
- **FR-005**: Input/output variable mappings on the Call Activity MUST behave identically to same-deployment calls.
- **FR-006**: If no deployed definition matches `calledElement` at activation time, the engine MUST reject with a stable code and MUST NOT mutate caller projection as if the call succeeded.
- **FR-007**: Recover MUST rebuild cross-deployment parent/child state from the log; redrive MUST not duplicate child instance creation for the same command.
- **FR-008**: `GetInstance` on parent and child MUST expose enough linkage and deployment information to trace a cross-deployment call without reading raw logs only.

### Key Entities

- **Caller deployment**: Deploy artifact for the parent process; may omit callee process body.
- **Callee deployment**: Independently deployed definition for the called process id.
- **Call resolution**: Runtime binding from `calledElement` process id → deployment id + process version (default latest).
- **Cross-deployment call link**: Parent token ↔ child instance with distinct deployment ids on each instance record.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Two separately deployed BPMN files: caller reaches Call Activity and spawns a child whose deployment id equals the callee deploy, visible in instance query.
- **SC-002**: Mapped IO round-trip across deployments matches same-deployment fixture outcomes in one test pass.
- **SC-003**: Recover after mid-call crash restores parent wait + child active; completing child finishes caller without manual repair.
- **SC-004**: Activating Call Activity when callee is not deployed produces REJECTION in 100% of tested negative cases.
- **SC-005**: Existing same-definitions Call Activity fixtures continue to pass unchanged (no regression).

## Assumptions

- MVP resolves callee by **process id** with **latest revision** when multiple callee deployments exist; explicit version pin on Call Activity is optional follow-up in plan (BPMN extension or attribute).
- Caller deploy does not require callee BPMN XML to be present; optional deploy-time warning if callee missing is out of scope (runtime failure is sufficient).
- Callee must be deployed to the same engine node before the call activates (no remote/cluster registry).
- Multi-instance Call Activity, live migration, and Call Activity boundary/compensation parity gaps beyond current same-deployment behavior stay separate increments unless they block MVP paths.
- Child instance command serialization remains per `process_instance_id`; parent and child locks are independent as today.
