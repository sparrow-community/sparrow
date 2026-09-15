# Data model

| Concept | Representation |
|---------|----------------|
| VariableMapping | Source, Target, Transformation, Assignments[] |
| MappingAssignment | From expr, To variable name |
| Runtime | applyMappings evaluates into map[string]json |

No proto changes.
