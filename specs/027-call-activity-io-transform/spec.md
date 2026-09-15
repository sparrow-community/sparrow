# Feature Specification: Call Activity IO Mapping with Transformation / Assignment

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned：Call Activity IO mapping with transformation / assignment。保持在主分支。"

## User Scenarios & Testing *(mandatory)*

Call Activity already copies variables by `sourceRef` → `targetRef`. Associations with `transformation` or `assignment` are rejected at Deploy. This increment evaluates those expressions with the same language as conditions (`expr`, `${...}` wrappers) and writes JSON variable values on input (caller → child) and output (child → caller).

### User Story 1 - Input transformation (Priority: P1)

An author maps caller `orderId` into child `id` via a transformation expression (for example concatenate a suffix). Deploy succeeds. CreateInstance passes the evaluated value into the child.

### User Story 2 - Output assignment (Priority: P1)

An author maps child `total` into caller `amount` via an assignment `from`/`to` (or transformation). Completing the child applies the evaluated value onto the caller.

### User Story 3 - Existing name copy and Recover (Priority: P2)

Simple sourceRef/targetRef mappings and Recover mid-call remain green.

### Edge Cases

- Transformation without sourceRef is allowed when the expression does not need a source name.
- Assignment `to` must resolve to a non-empty variable name (expression body or association targetRef).
- Transformation and assignment on the same association: assignments win (run assignments only).
- Expression evaluation errors surface as command failures with a stable prefix (`INVALID_MAPPING:`).
- In-engine DMN/FEEL remains out of scope; Sparrow `expr` language only.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept Call Activity dataInputAssociation / dataOutputAssociation with transformation and/or assignment.
- **FR-002**: Input mappings MUST evaluate transformation/assignment against caller variables before starting the child.
- **FR-003**: Output mappings MUST evaluate transformation/assignment against child variables when resuming the caller.
- **FR-004**: Simple sourceRef→targetRef copy MUST remain unchanged when no transformation/assignment is present.
- **FR-005**: Existing Call Activity IO fixtures (`m5_call_io`, cross-deploy) MUST remain green.

## Success Criteria *(mandatory)*

- **SC-001**: Input transformation fixture: child receives evaluated variable.
- **SC-002**: Output assignment/transformation fixture: caller receives evaluated variable.
- **SC-003**: Simple IO and regression suites remain green.
