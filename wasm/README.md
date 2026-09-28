# WASM host adapter

Thin browser host for `processing.Engine`. Same COMMAND surface as the gRPC
gateway; timers and job/script workers are scheduled by the JS page.

## npm package

Published as **`@sparrow-community/wasm`** on [npmjs.com](https://www.npmjs.com/)
under dist-tag **`alpha`**.

```bash
npm install @sparrow-community/wasm@alpha
```

### Version scheme (alpha)

`0.0.0-alpha.YYYYMMDD.N`

- `YYYYMMDD` — UTC calendar day of the publish
- `N` — starts at **1** each UTC day; increments for each publish that day
- Committed `package.json` keeps placeholder `0.0.0-alpha.0`; CI runs
  `scripts/bump-alpha-version.sh` before publish

Releases run from GitHub Actions (`.github/workflows/publish-wasm.yml`) via npm
**Trusted Publishing** (OIDC). No long-lived npm tokens in the repo or Actions
secrets. The first-ever package version must be bootstrapped interactively once
(npm cannot attach a Trusted Publisher until the package exists); after that,
Actions publishes without tokens.

Package contents (after `npm install`):

| Path | Role |
|------|------|
| `dist/sparrow.wasm.gz` | Engine binary (gzip) |
| `dist/wasm_exec.js` | Go JS/WASM support script (BSD) |
| `dist/wasm_exec.LICENSE` | License for `wasm_exec.js` |
| `sparrow.d.ts` | TypeScript surface for `globalThis.sparrow` |

The package ships artifacts only. Load `wasm_exec.js` as a classic script (it
defines `globalThis.Go`), inflate `sparrow.wasm.gz`, instantiate, then use
`globalThis.sparrow`. See **Load (browser)** below. Playground still vendors a
copy via `scripts/sync-wasm.sh`; switching it to the npm package is optional.

```bash
cd wasm
./build.sh          # → dist/sparrow.wasm(.gz), dist/wasm_exec.js(+LICENSE)
npm run pack:dry    # verify tarball contents (no registry auth)
npm run bump:alpha  # optional local preview of next 0.0.0-alpha.YYYYMMDD.N
```

`prepack` runs `./build.sh`, so `npm pack` / `npm publish` always rebuild.

## JS contract

`sparrow.d.ts` is the TypeScript surface for `globalThis.sparrow`. Playground
syncs it via `../sparrow-playground/scripts/sync-wasm.sh`.

## Build

```bash
./build.sh
# → dist/sparrow.wasm (+ .gz), dist/wasm_exec.js, dist/wasm_exec.LICENSE
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

With the npm package and a gzipped binary (modern browsers):

```js
import wasmGzUrl from "@sparrow-community/wasm/sparrow.wasm.gz?url"; // bundler-dependent
// or: new URL("@sparrow-community/wasm/dist/sparrow.wasm.gz", import.meta.url)
```

Prefer resolving package file URLs with your bundler; then inflate with
`DecompressionStream("gzip")` before `WebAssembly.instantiate`, and load
`wasm_exec.js` via a classic `<script>` or equivalent.

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

## License

Sparrow sources: Apache-2.0 (`LICENSE`). `wasm_exec.js`: BSD (`dist/wasm_exec.LICENSE`). See `NOTICE`.
