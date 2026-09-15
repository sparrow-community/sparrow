# Data Model: Conditional Catch

| Construct | Index | Runtime |
|-----------|-------|---------|
| Intermediate conditional catch | `conditionalCatch[id]=expr` | Wait or immediate complete |
| Conditional boundary | same map + activity attach | BoundaryWait kind `conditional` |

Discriminator for multiple conditional boundaries on one activity: condition text.
