# WASM host adapter

Thin browser host for `processing.Engine`. Same COMMAND surface as the gRPC
gateway; timers and job/script workers are scheduled by the JS page.

## JS contract

`sparrow.d.ts` is the TypeScript surface for `globalThis.sparrow`. Playground syncs it via `../sparrow-playground/scripts/sync-wasm.sh`.

## Build

```bash
./build.sh
# → dist/sparrow.wasm (+ .gz) and dist/wasm_exec.js
```

Requires Go 1.26+ with `js/wasm` support. Artifacts under `dist/` are gitignored.

## Load (browser)

```html
<script src="wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("sparrow.wasm"), go.importObject).then((r) => {
    go.run(r.instance);
    // globalThis.sparrow is ready
  });
</script>
```

## Timer host (JS)

```js
function armTimers() {
  const due = sparrow.nextDueUnixMs();
  if (!due) return;
  const delay = Math.max(0, due - Date.now());
  setTimeout(() => {
    sparrow.fireDue();
    armTimers();
  }, delay);
}
// call armTimers after createInstance / complete / fireDue
```

Optional: `sparrow.setNowUnixMs(ms)` injects the engine clock (pass `0`/`null` to clear).

## Script / job host (JS only for now)

Activate does **not** long-poll (`wait` is forced to `0`) so the JS thread stays free.

```js
async function drainJobs(jobType) {
  const { jobs } = sparrow.activate({ jobType, maxJobs: 8, workerId: "browser" });
  for (const job of jobs) {
    try {
      if (job.scriptFormat && /javascript/i.test(job.scriptFormat)) {
        // demo only — sandbox in real products
        const fn = new Function("variables", job.script + "\nreturn variables;");
        const variables = fn({ ...job.variables });
        sparrow.complete({
          instanceId: job.instanceId,
          elementId: job.elementId,
          tokenId: job.tokenId,
          variables,
        });
      } else if (!job.script) {
        // Service Task: host implements work, then complete
        sparrow.complete({
          instanceId: job.instanceId,
          elementId: job.elementId,
          tokenId: job.tokenId,
        });
      } else {
        sparrow.fail({
          instanceId: job.instanceId,
          elementId: job.elementId,
          tokenId: job.tokenId,
          message: "unsupported scriptFormat in browser",
          noRetry: true,
        });
      }
    } catch (e) {
      sparrow.fail({
        instanceId: job.instanceId,
        elementId: job.elementId,
        tokenId: job.tokenId,
        message: String(e),
        noRetry: false,
      });
    }
  }
}
```

## Errors

Host methods return a normal result on success. On failure they return
`{ "$error": "<message>" }` instead of panicking (a panic would exit the Go
WASM runtime). JS should treat `$error` as a thrown Error. Empty BPMN is
rejected with `INVALID_ARGUMENT: bpmnXml is empty`.

## API surface

| Method | Notes |
|--------|--------|
| `deploy(bpmnXml)` | → `{ deploymentId, processId }` |
| `createInstance({ deploymentId, processId?, processVersion?, variables? })` | → `{ instanceId }` |
| `complete` / `throwError` / `resolveIncident` | request object |
| `fireDue()` / `nextDueUnixMs()` / `setNowUnixMs(ms?)` | timer host |
| `publishMessage` / `publishSignal` | |
| `evaluateConditions` / `evaluateConditionalStarts` | |
| `activate` / `fail` / `heartbeat` | jobs; activate wait=0 |
| `getDeployment(id)` | includes `bpmnXml` |
| `getInstance(id)` / `listInstanceIds()` / `listEvents(id)` | |

In-memory only (no filesystem Recover). Playground / MCP consumers live outside this module.
