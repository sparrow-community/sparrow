# bpmn

Go module `github.com/sparrow-community/sparrow/bpmn`. BPMN 2.0 XML types and read/write helpers for the Sparrow kernel.

Go source in this module is Apache License 2.0. See [`LICENSE`](./LICENSE) and the repository [`NOTICE`](../NOTICE).

## Test fixtures

[`test/`](./test/) holds the BPMN XML and PNG files used by the MIWG load tests and the structural export/round-trip tests in this module.

Those files are unmodified copies of the **bpmn.io (Camunda Modeler) 18.6.1** results in the [BPMN Model Interchange Working Group (BPMN MIWG) test suite](https://github.com/bpmn-miwg/bpmn-miwg-test-suite/tree/master/bpmn.io%20(Camunda%20Modeler)%2018.6.1), plus `C.10.0-reference.*` copied from upstream `Reference/` (Camunda 18.6.1 has no C.10.0 vendor export).

**Upstream pin:** [`test/MIWG_UPSTREAM_SHA`](./test/MIWG_UPSTREAM_SHA) — `8416c1118ff98e9161e9e342220460be545ebf7c` (2026-09-25, bpmn-miwg/bpmn-miwg-test-suite `master`).

They are licensed under the Creative Commons Attribution 3.0 Unported License (CC BY 3.0). The notice shipped with them is [`test/LICENSE.txt`](./test/LICENSE.txt), reproduced from the upstream `LICENSE.txt`:

> This work is licensed under the Creative Commons Attribution 3.0 Unported License. To view a copy of this license, visit http://creativecommons.org/licenses/by/3.0/ or send a letter to Creative Commons, 444 Castro Street, Suite 900, Mountain View, California, 94041, USA.

Attribute them to the BPMN Model Interchange Working Group (OMG), BPMN Model Interchange Test Suite, under [CC BY 3.0](http://creativecommons.org/licenses/by/3.0/). The Apache License on this module does not apply to `test/`.

`A.2.0-rountrip.bpmn` keeps the upstream filename spelling. Tests load that path.

## Automated tests

```shell
# Parse every MIWG fixture
go test -C bpmn -run TestMIWGFixturesLoad ./...

# Kernel: Deploy + run (or intentional UNSUPPORTED / INVALID_CONDITION skip)
go test -C processing -run TestMIWGKernelFixtures ./...

# Convenience wrapper (same two packages)
./scripts/miwg-test.sh
```
