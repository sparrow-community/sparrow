# Sparrow — 项目状态（跨会话速览）

> **用途**：新开 Cursor 对话时 `@STATUS.md`，快速恢复「做到哪了、设计约束、下一步」。  
> **细节以文档为准**：语义与架构看 [`processing/DESIGN.md`](processing/DESIGN.md)；模块与夹具看 [`processing/README.md`](processing/README.md)；工作区约定看 [`AGENTS.md`](AGENTS.md)；产品定位看 [`AI-Driven-BPMN.md`](AI-Driven-BPMN.md)。  
> **维护**：每合并一块功能后更新本文件（约 5 分钟），不必改长对话历史。

---

## 一句话

单节点、事件账本为真相源的轻量 BPMN 引擎；**Event = 行为，Element = 主语**；对外 API 薄（Deploy / CreateInstance / Complete / FireDue / PublishMessage / PublishSignal / Job Activate），语义在 `processing/handlers/` 扩展。

## 当前里程碑

| 阶段 | 状态 |
|------|------|
| **M1** | 完成：Start → UserTask → XOR → End；文件日志 + `Recover` |
| **M2** | 完成：ServiceTask + Job；Timer/Message catch；Timer/Message boundary（打断型 + 非打断型） |
| **M3** | 进行中：网关 / SubProcess / Throw / ESP / Parallel EBG / **Compensation** |

**`main` 最新提交**（更新时改这里）：

```text
41bc9e9 Support activity compensation (boundary + throw)
60f5fba Record parallel event-based gateway commit in STATUS.md
664fc85 Support parallel event-based gateway (keep sibling catches)
ad765c9 Support process-level event sub-process (message/timer)
```

## 已实现（运行时）

- **核心**：`Deploy` / `CreateInstance` / `Complete`；`Recover` / `Open`；实例内串行锁
- **网关**：Exclusive（条件 + default）；Inclusive（OR-split / OR-join，可达 token 防死锁）；Parallel（fork/join）；Event-Based exclusive（先到的 catch 赢，取消兄弟）；Event-Based parallel（各 catch 独立完成，不取消兄弟）
- **SubProcess**：嵌套子流程（递归支持多层），共享变量，scope completion
- **Event Sub-Process**：流程级 `triggeredByEvent`（Message / Timer；打断型与非打断型）；Signal 同路径可扩展；迟到 message 可缓冲投递
- **等待与完成**：UserTask、ServiceTask、中间 Timer / Message / Signal catch — 统一 `Complete`
- **Throw**：Intermediate Throw none（里程碑）/ Message / Signal；锁释放后投递，避免重入
- **Timer**：`timeDuration` / `timeDate` / `timeCycle`（仅首次到期）；`Engine.FireDue`；`cmd/sparrow` 轮询
- **Message**：`PublishMessage`（name + 可选 correlation_keys + 缓冲）；`Open` 时缓冲持久化到 `runtime.Store`
- **Signal**：`PublishSignal`（name；无缓冲）；throw 与外部注入共用路径
- **Boundary**：Timer / Message 可挂 UserTask、ServiceTask 或 SubProcess；支持打断型与非打断型；同一活动可同时挂 Timer + Message；**Compensation** boundary + association 处理器
- **Compensation**：活动完成后订阅；`compensate` intermediate throw 按完成逆序执行 handler（wait）；可选 `activityRef`
- **Job**：`Activate` / `Fail` / `Heartbeat`；`Open` 时租约持久化到 `runtime.Store`（非账本）
- **持久化**：`log.EventLog`（行为账本）+ `deploy.Store`（定义）+ `runtime.Store`（租约 + 消息缓冲，`dataDir/runtime/state.json`）
- **传输**：`gateway` — `engine.v1` + `job.v1` gRPC

## 明确不做（当前阶段）

- 流程定义**版本管理与迁移**（每次 `Deploy` 新 `deployment_id`）
- 同一活动上 **三个及以上 boundary**；Signal boundary
- Incident；嵌套在 embedded SubProcess 内的 Event Sub-Process；instantiate EventBasedGateway；补偿传播进未完成的 SubProcess
- 集群 / 多活
- 复制 Camunda 产品广度（建模器、Cockpit 等）

## 下一步候选（按优先级）

1. 流程定义版本管理与迁移
2. Signal boundary / 嵌套 Event Sub-Process
3. 补偿增强（activityRef 多活动、End compensate、嵌套 scope）

## 新会话开场模板

```text
@STATUS.md @AGENTS.md @processing/DESIGN.md
本轮任务：<一件事>
请先读 STATUS + DESIGN §6（元素语义）和 §11（进度表），再改代码。
```

## 文档索引

| 问题 | 读哪里 |
|------|--------|
| 元素语义、时序链、API、进度表 | `processing/DESIGN.md` §6、§7、§11 |
| 夹具 BPMN 列表 | `processing/README.md` |
| 协议字段、模块命令 | `AGENTS.md` |
| 为什么做 Sparrow | `AI-Driven-BPMN.md` |

## 常用命令

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
cd protocol/proto && ./build.sh
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```
