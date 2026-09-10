# Feature Specification: Conditional Flow

**Feature Branch**: `013-conditional-flow`

**Created**: 2026-09-10

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步 Remaining P3：Conditional flow — 活动多出边按 conditionExpression 选择路径（含 default）；与现有 XOR/Inclusive 条件语义一致；Recover 覆盖。"

## User Scenarios & Testing *(mandatory)*

Exclusive/Inclusive gateways already evaluate sequence-flow conditions. Authors also attach `conditionExpression` (and optional `default`) on outgoing flows from activities (e.g. User Task) without an intervening gateway. Today the engine takes the first outgoing and ignores conditions.

### User Story 1 - Activity chooses conditioned path (Priority: P1)

An author models a User Task with two outgoing flows: one conditioned (`approved == true`), one default. Completing with `approved=true` takes the conditioned path; with `approved=false` (or missing) takes the default.

**Why this priority**: Core gap vs gateway-only conditions.

**Independent Test**: Fixture without XOR; Complete with vars; assert end path.

**Acceptance Scenarios**:

1. **Given** a User Task with conditioned + default outgoing, **When** Complete sets the condition true, **Then** the conditioned path is taken.
2. **Given** the same model, **When** Complete leaves the condition false, **Then** the default path is taken.
3. **Given** conditioned outs and no matching condition and no default, **When** Complete runs, **Then** the engine rejects / fails with NO_OUTGOING_FLOW (or equivalent).

---

### User Story 2 - XOR remains unchanged (Priority: P1)

Existing exclusive gateway condition + default fixtures keep passing (same evaluation rules).

**Acceptance Scenarios**:

1. **Given** existing XOR tests, **When** the suite runs, **Then** they still pass.

---

### User Story 3 - Recover then Complete chooses path (Priority: P2)

Recover while waiting on the User Task; Complete with variables after Recover selects the same path as a continuous run.

**Acceptance Scenarios**:

1. **Given** waiting on the conditional User Task, **When** Recover then Complete with condition true, **Then** the conditioned end completes.

---

### Edge Cases

- Single outgoing with a false condition and no default → NO_OUTGOING_FLOW.
- Multiple unconditional outgoings (no conditions, no default): keep prior behavior (take first) — parallel-from-activity is out of scope.
- Service Task / Call Activity / SubProcess host completion use the same chooser when TakeOutgoing has no preferred flow.
- Conditional *events* (intermediate conditional catch) stay out of scope for this increment.

## Requirements *(mandatory)*

- **FR-001**: When leaving an activity with outgoing sequence flows that have conditions and/or an activity `default`, the engine MUST choose the outgoing like exclusive gateway (first true non-default condition, else default).
- **FR-002**: Existing XOR/Inclusive condition evaluation MUST remain behavior-compatible.
- **FR-003**: Recover MUST preserve waiting activity tokens so post-Recover Complete with variables selects the correct flow.
- **FR-004**: Unconditional single or multi-out without conditions/default MUST keep current take-first behavior when no preferred flow is set.

## Success Criteria *(mandatory)*

- **SC-001**: User Task conditioned path taken when vars match.
- **SC-002**: Default path taken when vars do not match.
- **SC-003**: Existing XOR/Inclusive suites pass.
- **SC-004**: Recover then Complete matches continuous conditioned outcome.

## Assumptions

- Selection is exclusive (one flow), not inclusive fan-out from activities.
- Expression language remains the existing `processing/expr` evaluator.
- Conditional events are a separate Remaining/future item if still needed after this.
