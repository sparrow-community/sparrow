# Protocol breaking baseline (bootstrap)

Frozen copy of `protocol/proto` `.proto` sources at the 036 compatibility
baseline. CI runs:

```shell
cd protocol/proto && buf breaking --against ../breaking-baseline
```

This exists because no published protocol release tag is available yet.
When a protocol tag is cut, switch CI to
`buf breaking --against '.git#tag=<tag>,subdir=protocol/proto'` and update
[`COMPATIBILITY.md`](../../COMPATIBILITY.md) in the same change.
