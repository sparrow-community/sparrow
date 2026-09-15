# Data model

No proto or ledger schema changes. Runtime: multiple SEQUENCE_FLOW_TAKEN + tokens after leave.

| Concept | Representation |
|---------|----------------|
| ChooseOutgoingFlows | []flowID — all outs or single exclusive |
| Fork leave | existing Effect.Fork path / shared helper |
