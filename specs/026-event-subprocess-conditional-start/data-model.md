# Data model

| Concept | Representation |
|---------|----------------|
| Conditional event sub-process | `EventSubProcess.Kind = conditional`, `Condition` text |
| Arm | Existing `EventSubProcesses[id]` after START_EVENT ACTIVATED |
| Trigger | EvaluateConditions → `triggerEventSubProcess` |

No proto changes.
