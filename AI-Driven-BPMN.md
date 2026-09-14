# AI Driven BPMN

This note explains **why Sparrow exists** and **how agents may use it**. It is not an element-completeness roadmap. Element coverage and semantic completeness for OMG BPMN process execution are defined in [`AGENTS.md`](./AGENTS.md) and remain the engine prerequisite.

## Why Sparrow

Process work needs a durable, shared contract across people and time: a versioned definition, deterministic execution, and an auditable trail. Agent plans are short-lived. Sparrow is the execution and fact layer—BPMN in, COMMAND → EVENT out—so generated or conversational intent cannot silently rewrite a running contract.

```text
Agent / tools / MCP
        ↓
   Sparrow (COMMAND → EVENT)
        ↓
BPMN definitions + event log
```

Generation may change. The effective definition and instance state change only through COMMAND/EVENT.

## How AI combines with Sparrow

| Role | Responsibility |
|------|----------------|
| Agent | Draft definitions, query instances, assist waits, propose next COMMANDs |
| Sparrow | Deploy, create instances, complete waits, reject invalid work, recover from the log |
| BPMN | Versioned contract; agent behavior maps onto existing constructs and extensions |

Drive mode: conversation and tools prepare work; the engine accepts only COMMANDs. Custom agent behavior maps onto Service Task, User Task, and BPMN extensions—not new non-OMG core element types.

## Prerequisite

BPMN executable-process completeness ([`AGENTS.md`](./AGENTS.md) Supported / Planned / Excluded) is the foundation. AI Driven consumes that engine; it does not replace Planned element work, invent ledger subjects, or mutate projections outside COMMAND/EVENT.

## Later consumer work (after / beside Planned)

| Item | Goal |
|------|------|
| Task extensions | Agent metadata via extension attributes (not new core types) |
| Correlation query | Look up instances by business key |
| Agent playbook | Deploy → create → complete → list events |
| Failure handling | Align with engine incident semantics |

AI features must not bypass COMMAND/EVENT or silently rewrite projections.
