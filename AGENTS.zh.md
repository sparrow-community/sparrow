# Sparrow

> 中文版。英文版见 [`AGENTS.md`](./AGENTS.md)（规范引用与工具默认以英文版为准）。两版须在同一次变更中同步更新 Supported / Planned / Excluded、公理、Purpose 与 Completeness。

Sparrow 是 BPMN **执行与事实内核**：接受版本化的流程定义、执行它，并记录发生了什么。Agent、UI、叠加渲染器都是该契约与轨迹的**消费者**——它们不得改写契约与轨迹。

**公理**（必须始终成立）：

- 有效定义与实例状态只能通过被接受的 **COMMAND** 改变；每一条被接受的行为都是追加写的 **EVENT**。投影可丢弃，必须能通过回放（**Recover**）重建。
- **Event = 行为；Element = 主体。** 账本是关于过程行为的事实流，不是图的第二份拷贝。
- 具有运行时 token 或副作用语义的组合要么是 **Supported**（真实 BPMN 语义），要么是 **Excluded**（显式非目标）。不支持的工作在 Deploy 时拒绝。沉默不等于排除。
- 账本只记录行为。文档性或非主体构造留在定义上；它们不接收空的 Element intent。

**范围选择**（有意为之，并非单由公理推出）：OMG BPMN 2.0 **可执行过程**子集——带有运行时 token 行为的 FlowElement 与事件定义——不含 Collaboration、Choreography 或完整元模型。**工程选择**：单机；每实例一条事实流、一把锁。

## Purpose（目的）

过程工作需要跨人、跨时间、跨短命规划者（人或 agent）共享的持久契约：版本化定义、忠实执行、可审计轨迹。Sparrow 存在的意义，是使该契约无法在 COMMAND → EVENT 之外被静默改写。

## Completeness（完备性）

对已声明语言的忠实解释是一条轴（覆盖 + 语义）。内核形态是另一条轴。

1. **Coverage（覆盖）** — 每一个可执行过程的 FlowElement 与事件定义组合要么是 Supported，要么是 Excluded。不支持的组合在 deploy 时拒绝。
2. **Semantics（语义）** — 每一个 Supported 组合实现该组合所需的 BPMN 运行时语义（token 移动、等待、抛出、汇合、作用域、边界、多实例、补偿及相关规则）。仅 Deploy 接受不足以称为完备。
3. **Kernel boundary（内核边界）** — 引擎保持为账本支撑的执行内核（COMMAND → EVENT、Recover、jobs、messages、timers、incidents）。建模器、运维控制台、集群复制、编排/协作运行时属于消费者或非目标——不是引擎职责。将 Event 叠加到部署图上渲染等产品能力属于消费者；它们不得发明账本主体或空 intent。

当下列表中 (1) 与 (2) 成立，且每一个 Supported 组合都有端到端与 Recover 测试时，完备性成立（测试是**如何知道**完备——不是引擎的另一属性）。完成一批交付本身并不等于完备。

| | |
|--|--|
| 治理 | [`.specify/memory/constitution.md`](./.specify/memory/constitution.md) |
| 运行时 | [`processing/README.md`](./processing/README.md) |

工作流：`/speckit-specify` → plan → tasks → implement（`.cursor/skills/`）。

## Supported（已支持）

过程生命周期；过程级 none 开始与类型化开始（message、timer、signal、conditional — none 开始用 CreateInstance；instantiate 的排他/并行基于事件的网关入口与 instantiate receive task；类型化铸造经 PublishMessage / PublishSignal / FireDue / EvaluateConditionalStarts）；User Task；Service Task + Job；Manual / 抽象 Task（等待 → Complete）；Receive / Send（含作为过程入口的 instantiate receive）；Business Rule / Script Task（job 支撑）；排他 / 并行 / 包容 / 复杂 / 基于事件的网关（catch，含 instantiate 入口与 receive-task 目标）；SubProcess；Transaction SubProcess 与 cancel（##Compensate；Cancel End → 补偿 → Cancel Boundary）；Ad-Hoc SubProcess（扁平内部活动，`ordering` Parallel / Sequential，`completionCondition`，`cancelRemainingInstances`）；Call Activity（子实例，IO 映射含按名拷贝 / 变换 / 赋值，跨部署，边界与补偿对等，多实例；被调过程需要 none 开始、instantiate 基于事件的网关或 instantiate receive）；Event Sub-Process（含嵌套、补偿事件子过程，以及经 EvaluateConditions 的条件开始）；timer、message、signal、error、escalation（含独立中间 escalation catch）、compensate、conditional 的中间与边界 catch/throw（同种多边界）；link throw/catch；活动与网关上的条件顺序流；活动与开始事件多条无条件出向的并行扇出；terminate end；message end；signal end；对未完成嵌入式 SubProcess 或未完成 Call Activity 子实例的补偿；User Task、Service Task、Manual / 抽象 Task / Receive / Send / Business Rule / Script、SubProcess、Call Activity 上的多实例（含 `complexBehaviorDefinition` 与 signal/message 的 none/one behavior event refs）；等待类任务上的 `standardLoopCharacteristics`（testBefore / loopCondition / loopMaximum）；incident 打开 / 解决 / 重试；共存的过程修订版。

已交付：M1–M4c，规格 [`001`](./specs/001-engine-completeness/)–[`035`](./specs/035-get-deployment/)。上列每一个 Supported 组合都有端到端与 Recover 测试（[`034`](./specs/034-recover-coverage/) 闭合了最后的覆盖缺口）。

## Planned（计划中）

尚未覆盖的可执行过程缺口。优先级指导排序；特性计划可因依赖重排。除非计划捆绑紧密相关的格子，否则每一行是一次 specify 增量。

当前无未覆盖的可执行过程元素缺口。下一步在识别出新缺口时选定；候选在开工前记入此处。

## Excluded（排除）

非目标。将某项移出 Excluded 需要修订 AGENTS 与 constitution。

| 主题 | 理由 |
|------|------|
| Collaboration、message flow 与 choreography 执行 | 超出单过程执行；进入 choreography 或其他未建模 flow element 的顺序流在 Deploy 时拒绝 |
| Lane 运行时语义 | 仅文档与分组 |
| Data Object / Data Store 作为账本或 token 主体 | 可执行数据由过程变量与 IO 映射承载 |
| 建模器、运维控制台、产品套件 UI | 引擎的独立消费者 |
| 集群与多节点复制 | 单机设计 |
| 引擎内 DMN 求值 | Business Rule Task 由 job 支撑 |
| 作为 flow element 的 `implicitThrowEvent` | 自身无 token 语义；仅在多实例行为事件中合法，故 Deploy 以 `UNSUPPORTED_ELEMENT` 拒绝，而非别名为 none throw |

## Future（运维）

| 主题 | 说明 |
|------|------|
| 运行中实例迁移 | 修订版可共存；将运行中实例迁到另一修订版属于运维，不是元素完备性 |

## Consumers（消费者）

Agent、叠加 UI、运维工具与 MCP 适配器是对等消费者：各自消费同一套 COMMAND 面与事件日志。它们可以起草定义、查询实例与轨迹、协助等待、提议 COMMAND，或从定义 + Event 派生视图。它们不得发明非 OMG 核心元素类型、跳过 COMMAND/EVENT、发明账本主体或空 intent，或在日志之外改写投影。某一消费者与模型或工具集成到多深，不在本内核范围之内。

## API and persistence（API 与持久化）

`Deploy` · `CreateInstance` · `Complete` · `ThrowError` · `ResolveIncident` · `FireDue` · `PublishMessage` · `PublishSignal` · `EvaluateConditionalStarts` · `EvaluateConditions` · Job Activate / Fail / Heartbeat · `GetDeployment` · `GetInstance` · `ListEvents`

`EventLog` + `deploy.Store` + 可选 `runtime.Store` · `Recover` / `Open` · gRPC 仅在 `gateway`

## Workspace（工作区）

Go **1.26.5** · 模块：`bpmn` · `protocol` · `processing` · `gateway`

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
go test ./bpmn/
cd protocol/proto && ./build.sh
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```

修改 `.proto` 后运行 `build.sh`。不要手改 `*.pb.go`。元素语义在 `handlers/` 与 [`processing/README.md`](./processing/README.md)。Apache-2.0 文件头。
