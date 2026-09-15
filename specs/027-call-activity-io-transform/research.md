# Research

**Decision**: Use Sparrow `expr` for transformation/assignment values (same as conditions); marshal results as JSON variable values.

**Decision**: Assignments on an association take precedence over transformation; otherwise transformation into targetRef; else source→target copy.

**Decision**: Fix Assignment `from`/`to` to ExpressionUnMarshal nested elements.
