# Security

Sparrow is an early public preview (alpha). There is no bug bounty and no separate security support window.

## Reporting a vulnerability

Report suspected vulnerabilities in private. Do not open a public issue or pull request with exploit details, proof-of-concept BPMN, or steps that attack a running kernel.

Use [private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability) on [github.com/sparrow-community/sparrow](https://github.com/sparrow-community/sparrow) when that channel is enabled. If it is not enabled yet, contact the maintainers directly and keep the details off public trackers.

## Trusted network

The gateway is a single-node kernel, not an internet-facing API.

`gateway/cmd/sparrow` listens with plain gRPC. The server is created without TLS and without authentication interceptors. Anyone who can reach `-listen` can deploy definitions, start and complete work, publish messages and signals, resolve incidents, and read deployments, instance state, and the event log.

Run it only on a trusted network. Do not publish the listen address to untrusted clients, and do not put it on the public internet without your own transport security and access control in front of it. The on-disk data directory (`-data-dir`) is the event log and the deployed definitions; protect it the same way you protect the process.

The WebAssembly host is in-memory inside the page that loads it. It has the same command surface and no network listener of its own. The page that embeds it is responsible for who can call that surface.

## Deployed BPMN is code the host executes

A process definition is not inert configuration.

Sparrow evaluates definition-supplied expressions inside the engine process with [github.com/expr-lang/expr](https://github.com/expr-lang/expr). That includes sequence-flow conditions, conditional catches and starts, standard-loop and multi-instance conditions, ad-hoc completion conditions, complex-gateway activation conditions, and Call Activity IO transformations and assignments. Deploying BPMN asks the Sparrow process to compile and run that expression text with the privileges of the host process, against instance variables.

Script Task bodies are not evaluated inside the engine. They are delivered to a job worker. A worker that runs the script is executing caller-supplied code.

Deploy only definitions you trust. Do not accept BPMN from untrusted authors on a shared engine process.
