# processing 程序设计

本文描述 `processing` 模块的目标形态、核心模型、包结构、处理流程与 M1 实现边界。  
实现应以本文为准；协议字段变更在 `protocol` 中进行，并保持「Event = 元素行为」的定位。

**M1 状态**：Deploy / CreateInstance / Complete、XOR 条件、内存与文件 EventLog、`Open` 重启回放已可用。  
**M2 起步**：ServiceTask 等待；Job 经 `Activate` / `Fail` / `Heartbeat`；中间捕获 Timer（`timeDuration`）；gRPC 在 `gateway`（`engine.v1` + `job.v1`），进程入口 `gateway/cmd/sparrow`（含 `FireDue` 轮询）。  
尚未实现：Message catch、Timer 的 date/cycle/boundary、更多 BPMN 元素、跨重启的 Job 租约。

---

## 1. 目标与非目标

### 1.1 目标

- 单节点可运行的 BPMN 执行内核
- 事件日志为唯一真相源；实例状态为投影
- 行为可审计、可回放、可拒绝（COMMAND → EVENT / REJECTION）
- API 表面保持简单，便于后续在不改语义的前提下加分区做水平扩展
- **少中间层**：定义层直接复用 `bpmn/element`，不为执行另造平行图模型

### 1.2 非目标（当前阶段）

- 集群、多活、跨节点分区调度
- 完整 BPMN 覆盖（DMN、CMMN、建模器、运维 UI）
- 将 Job / Timer 提升为与 Element 同级的事件主语
- 复制 Camunda 产品广度
- 为「整洁」而堆叠无必要的 adapter / Node / Flow 包装类型

---

## 2. 核心概念

### 2.1 行为账本

每条 `event.v1.Event` 描述 **一次行为**：

| 字段 | 含义 |
|------|------|
| `record_type` | COMMAND（请求）/ EVENT（已发生事实）/ REJECTION（拒绝） |
| `deployment_id` / `process_instance_id` / `process_version` | 行为发生的定义与实例上下文 |
| `element` | 行为主语：哪个元素、什么 Intent、附带 payload |
| `source_record_id` | 因果：通常指向触发本条 EVENT/REJECTION 的 COMMAND |
| `rejection` | 仅 REJECTION：机器可读 code + 说明 |

日志偏移（position/sequence）由存储层维护，可不写入 protobuf；回放按追加顺序即可。

### 2.2 Element 是行为主语

```text
谁：Type + element.id + token_id
做了什么：Intent
附带什么：payload（变量增量、选中的 sequenceFlow、后续的 job/timer 字段等）
```

`Element.Type` 与 BPMN 独立元素对齐（`PROCESS` + `FlowElements` 具体类型）。  
**不使用**「粗 Type + kind」合并不同 BPMN 元素。

Job / Timer / Message 等待等，视为 **该元素行为的载荷与阶段**，写入对应 `payload` 与 `Intent`，而不是新的顶层 value 类型。

### 2.3 串行与分区键

- 分区键：`process_instance_id`
- 同一实例上的 COMMAND 严格串行处理（每实例一把互斥锁）
- 单节点 = 一个逻辑分区的宿主；日后多节点只是多分区复制同一模型

### 2.4 定义 vs 运行时

| 层 | 来源 | 作用 |
|----|------|------|
| 定义 | `bpmn` 解析的 `element.Process` | 静态结构：节点、边、默认流等 |
| 部署 | `deploy.Deployment` | 校验后的不可变定义快照 + `deployment_id` |
| 实例 | `projection.Instance` | 一次执行的投影（状态 / 变量 / 令牌） |
| 令牌 | `token_id` → `{element_id, active\|waiting, job_type, due_unix_ms}` | 实例内控制流位置（M1 单 token） |

**不另建** `graph.Node` / `graph.Flow`。`deploy` 在 `element.Process` 上提供查询辅助（`TypeOf`、`Outgoing`、`SequenceFlow`、`ChooseExclusiveOutgoing`）。

---

## 3. 架构总览

```text
                    ┌──────────────────────────────────────────┐
                    │              Engine（薄门面）               │
                    │  Deploy / CreateInstance / Complete        │
                    │  Activate / Fail / Heartbeat（Job 拉模型）   │
                    │  实例锁 · 写 COMMAND/EVENT/REJECTION        │
                    │              │                             │
                    │              ▼                             │
                    │         Executor                           │
                    │  Enter / Complete · 应用 Effect            │
                    │  出边 → SEQUENCE_FLOW_TAKEN · 自动步进      │
                    │              │                             │
                    │              ▼                             │
                    │   handlers（按 Element.Type）               │
                    │   OnEnter / OnComplete → Effect            │
                    │                                            │
                    │  deploy.Deployment  ← element.Process      │
                    │  projection.Instance ← ApplyEvent          │
                    │  EventLog / deploy.Store（可替换实现）        │
                    └──────────────────────────────────────────┘
```

**原则**：

- 投影可丢；EventLog 不可丢。`Recover` 按 EVENT 重建投影。
- **元素语义在 handlers**；Engine 不写具体生命周期分支。
- 持久化只有两个注入点：`log.EventLog`（行为账本）与 `deploy.Store`（定义字节）。Memory / File 是开发默认实现。
- COMMAND 处理由 Engine 内联；不另留空的 Processor 接口。

---

## 4. 包结构（与仓库一致）

```text
processing/
├── README.md
├── DESIGN.md
├── engine.go                 // API、实例锁、写日志、emitter
├── jobs.go                   // Activate、Job 快照、内存租约
├── open.go                   // Recover(ctx, log, store)、Open 便捷封装
├── executor.go               // 令牌推进编排（调用 handlers）
├── id.go                     // UUIDv7
├── expr/                     // 条件表达式（expr-lang）
├── deploy/
│   ├── deploy.go             // Compile + Deployment 查询辅助
│   ├── store.go              // Store 接口 + MemoryStore
│   └── store_dir.go          // 目录实现
├── handlers/
│   ├── handler.go            // Effect / 接口 / Registry / InstantLifecycle
│   ├── process.go
│   ├── start_event.go
│   ├── end_event.go
│   ├── user_task.go
│   ├── service_task.go
│   ├── exclusive_gateway.go
│   └── sequence_flow.go
├── log/                      // EventLog（Memory / File）
├── projection/               // Instance / Token / ApplyEvent
└── testdata/                 // m1_simple.bpmn 等
```

持久化布局（`Open` 文件便捷实现，不是唯一方式）：

```text
dataDir/
  events.log                     // length-delimited protobuf Event
  deployments/<id>.bpmn          // 原始定义，供重启后 Compile
```

自定义持久化：

```go
eng, err := processing.Recover(ctx, myEventLog, myDeploymentStore)
```
扩展新 BPMN 元素时：**新增一个 handler 文件 + 注册到 `DefaultRegistry`**，并在 `deploy.validateM1`（或后续更细校验）中放开该类型。

### 4.1 关键类型

```go
// EventLog — 行为账本（Memory / File / 自实现）
type EventLog interface {
    Append(ctx context.Context, e *eventv1.Event) (position int64, err error)
    ReadByInstance(ctx context.Context, processInstanceID string) ([]*eventv1.Event, error)
    ReadAll(ctx context.Context) ([]*eventv1.Event, error)
}

// Store — BPMN 定义字节（MemoryStore / DirStore / 自实现）
type Store interface {
    Put(id string, bpmnXML []byte) error
    LoadAll() (map[string][]byte, error) // id → xml
}

// ElementHandler（handlers 包）
type ElementHandler interface {
    Type() eventv1.Element_Type
    OnEnter(in EnterInput) (*Effect, error)
    OnComplete(in CompleteInput) (*Effect, error)
}

// Recover(ctx, log, store) 加载定义并回放 EVENT；Open(ctx, dataDir) = Recover(File, DirStore)
// NewEngine(log) 仅内存定义、不回放（测试 / 无持久化会话）
```

Effect：handler 产出，由 Executor 应用（Records / Wait / OutgoingFlowID / TakeOutgoing / TryCompleteProcess）。

Engine 对外能力：Deploy / CreateInstance / Complete / Activate / Fail / Heartbeat / GetInstance / ListEvents。

ID 统一走 `NextID()`（UUIDv7 字符串）。时间戳使用 Unix millis。空 id 用空字符串表示。

### 4.2 Effect 与步进

| Effect 字段 | 含义 |
|-------------|------|
| `Records` | 要写成 EVENT 的 Element 行为（经 emitter 追加并 `ApplyEvent`） |
| `Wait` | 停止自动步进（UserTask ACTIVATED） |
| `TakeOutgoing` | 取一条出边并发 `SEQUENCE_FLOW_TAKEN`，再 Enter target |
| `OutgoingFlowID` | 指定出边（XOR 选路）；空则取第一条 outgoing |
| `TryCompleteProcess` | End 后尝试 PROCESS COMPLETING→COMPLETED |

瞬时元素（Start / XOR / End）可用 `InstantLifecycle` 写出完整 Intent 链。

---

## 5. 处理循环（详细）

### 5.1 当前路径（Engine 内联）

```text
CreateInstance / Complete:
  1. 取 deployment + instance；加实例锁
  2. 校验投影（如 UserTask 须 waiting）
  3. Append(COMMAND)
  4. Executor.Enter 或 Executor.Complete
       → handler.OnEnter / OnComplete → Effect
       → emitter(Records)  // Append EVENT + ApplyEvent
       → 按 Effect 出边 / 等待 / 尝试完成流程
  5. 校验失败路径：Append(COMMAND) + Append(REJECTION)
```

同一 COMMAND 处理中可连续写出多条 EVENT（启动链、瞬时生命周期、流转移），均共享该 COMMAND 的 `source_record_id`。

### 5.2 幂等与半截链

COMMAND 先入账，再连写多条 EVENT（共享 `source_record_id`）。崩溃可能停在「只有 COMMAND」或「EVENT 链写了一半」。

`Recover` 回放 EVENT 后按日志顺序检查每条 COMMAND：

- 已有对应 REJECTION → 跳过  
- 已有至少一条 EVENT **且** 投影稳定（waiting / completed / terminated）→ 视为该命令已做完，跳过  
- 否则 **接着执行同一条 COMMAND**（不新写 COMMAND）。Emitter 对 `(source_record_id, Type, Intent, element.id, token_id)` 已出现过的 EVENT 不再追加，因此重入 `Enter` / `Complete` 只会补上缺失的步骤。

客户端未带 `cmd.id` 的重试仍是新 COMMAND；成功后的第二次 Complete 仍是 `INVALID_STATE`。幂等保证的是 **磁盘上那条未完成的命令** 能被 `Open` 做完。

### 5.3 投影与令牌

`projection.Instance` 最少包含：

- 实例状态：`active | completed | terminated`
- `Tokens`：`token_id → {element_id, active|waiting, job_type, due_unix_ms}`
- 变量表（实例级 `name → json_value`）
- `ElementIntent`：元素级最近 Intent（辅助校验）

令牌更新 **只走 EVENT → ApplyEvent → applyToken**：

- Executor 只决定步进（`Effect.Wait` / 出边），不改 `Tokens`
- `USER_TASK` + `ACTIVATED` → `waiting`；`SERVICE_TASK` + `ACTIVATED` → `waiting` 且拷贝 `payload.job_type`；`INTERMEDIATE_CATCH_EVENT` + `ACTIVATED` → `waiting` 且拷贝 `payload.due_unix_ms`
- `SEQUENCE_FLOW_TAKEN` 把位置写到 target
- PROCESS COMPLETED/TERMINATED 清空 `Tokens`
- 在线与 `Recover` 共用同一套规则

M1 为单 token；Parallel 等多 token 时仍落在同一 map，由 gateway handler 分裂/汇合。`GetInstance` 返回投影拷贝。

### 5.4 重启恢复

```text
Recover(ctx, eventLog, deploymentStore)
  → Store.LoadAll → Compile → deployments map
  → EventLog.ReadAll
  → 按序对每条 EVENT 调用 Instance.ApplyEvent（必要时先创建投影）
  → 未完成的 COMMAND 用同一 cmd.id 接着跑（幂等 emitter 不重复写已有 EVENT）
  → 恢复等待点（如 UserTask waiting），可继续 Complete
```

`Open(ctx, dataDir)` 只是文件实现的便捷封装。M1 采用全量重放；实例量大时再引入快照。COMMAND / REJECTION 不驱动投影（仅 EVENT）。文件日志读到不完整尾包时截断，不让 Recover 失败。
---

## 6. M1 可执行语义

### 6.1 支持的元素 Type

| Type | Handler 文件 | 行为要点 |
|------|--------------|----------|
| `PROCESS` | `process.go` | 实例启动 / 正常完成 |
| `START_EVENT` | `start_event.go` | 瞬时生命周期后沿出口流出 |
| `USER_TASK` | `user_task.go` | ACTIVATING→ACTIVATED 后 `Wait`；Complete → COMPLETING→COMPLETED |
| `SERVICE_TASK` | `service_task.go` | 同上等待；ACTIVATED 带 `ActivityPayload.job_type`；经 `Complete` 完成 |
| `INTERMEDIATE_CATCH_EVENT` | `intermediate_catch_event.go` | 仅 timer + `timeDuration`（ISO-8601 `PTnHnMnS`）；ACTIVATED 后 `Wait`，`EventPayload.due_unix_ms`；经 `Complete` / `FireDue` 完成 |
| `EXCLUSIVE_GATEWAY` | `exclusive_gateway.go` | 非 default 条件按序求值，否则 default；payload 带 `taken_sequence_flow_id` |
| `SEQUENCE_FLOW` | `sequence_flow.go` | 经 transit 发 `SEQUENCE_FLOW_TAKEN`（不走 OnEnter） |
| `END_EVENT` | `end_event.go` | 完成后 `TryCompleteProcess` |

部署时 `validateM1` 拒绝尚未实现的元素（Parallel、SubProcess、throw/boundary、message catch 等）。ServiceTask 与中间捕获 Timer 已纳入可执行子集。

### Timer catch（timeDuration）时序链（M2）

以 `Start → IntermediateCatchEvent(timer PTnHnMnS) → End` 为例：

1. `Engine.CreateInstance` 写启动 `COMMAND(PROCESS)`，并进入 `StartEvent`
2. token 进入 `INTERMEDIATE_CATCH_EVENT` 后，`Executor.Enter` 调用
   `IntermediateCatchEventHandler.OnEnter`
3. `OnEnter` 从 `deploy.Deployment` 读取 `timeDuration`，计算并写出到事件载荷：
   - `EVENT(INTERMEDIATE_CATCH_EVENT, ACTIVATING)`
   - `EVENT(INTERMEDIATE_CATCH_EVENT, ACTIVATED, EventPayload{due_unix_ms, duration})`
   - `Effect.Wait=true`，因此 token 保持 `waiting`（投影里是 `Token.DueUnixMs`）
4. `gateway/cmd/sparrow` 在后台循环调用 `Engine.FireDue`
   - `FireDue` 扫描 `waiting && due_unix_ms <= now` 的 token
   - 对每个到期 token 调用统一入口 `Engine.Complete(...)`
5. `Engine.Complete` 写入 `COMMAND(COMPLETING...)` + handler 输出
   `EVENT(COMPLETING...)` + `EVENT(COMPLETED...)`，并 `TakeOutgoing=true` 推进 `SEQUENCE_FLOW_TAKEN`

重启恢复：`Recover` 通过回放 `ACTIVATED{due_unix_ms}` 恢复投影中的 `Token.DueUnixMs`，但仍需等待下一次 `FireDue`（或手动 Complete）触发推进。

### 6.2 (Type, Intent) 使用（M1）

| 场景 | Type | Intent 序列（EVENT） |
|------|------|----------------------|
| 创建实例 | PROCESS | ACTIVATING → ACTIVATED |
| 进入 Start | START_EVENT | ACTIVATING → ACTIVATED → COMPLETING → COMPLETED |
| 走过流 | SEQUENCE_FLOW | SEQUENCE_FLOW_TAKEN |
| 进入 UserTask | USER_TASK | ACTIVATING → ACTIVATED（等待） |
| 完成 UserTask | USER_TASK | COMPLETING → COMPLETED |
| 进入 ServiceTask | SERVICE_TASK | ACTIVATING → ACTIVATED（等待；payload.job_type） |
| 完成 ServiceTask | SERVICE_TASK | COMPLETING → COMPLETED |
| 进入 Timer catch | INTERMEDIATE_CATCH_EVENT | ACTIVATING → ACTIVATED（等待；payload.due_unix_ms / duration） |
| 完成 Timer catch | INTERMEDIATE_CATCH_EVENT | COMPLETING → COMPLETED |
| Job 失败 | SERVICE_TASK | FAILED（仍 waiting；payload.error_message） |
| XOR | EXCLUSIVE_GATEWAY | ACTIVATING → … → COMPLETED（payload 带 taken flow） |
| End | END_EVENT | … → COMPLETED |
| 实例结束 | PROCESS | COMPLETING → COMPLETED |

**等待点（UserTask / ServiceTask / Timer catch ACTIVATED）** 与 **SEQUENCE_FLOW_TAKEN** 必须在日志中可见。

### 6.3 变量

- `ActivityPayload.variables` / `ProcessPayload.variables`（`name` + `json_value`）
- 语义为 **delta**：合并进实例变量表
- CreateInstance 可带初始变量；Complete 可带提交变量
- XOR：按 outgoing 顺序求值非 default 的条件（`expr.Eval` / expr-lang）；都不成立则走 default
- 条件：剥掉 BPMN `${...}` 后交给 [expr-lang/expr](https://github.com/expr-lang/expr)；变量为实例 JSON 值。单引号字符串会先归一成双引号。

### 6.4 拒绝示例

| code | 场景 |
|------|------|
| `NOT_FOUND` | 实例、部署或元素不存在 |
| `INVALID_STATE` | UserTask 未处于 waiting 却 Complete |
| `UNSUPPORTED_ELEMENT` | 定义含 M1 未支持元素，或无 handler |
| `INVALID_CONDITION` | XOR 条件表达式无法解析 |
| `INVALID_ARGUMENT` | Activate 缺少 `job_type` |
| `INVALID_STATE` | Fail 作用于非 job（如 UserTask），或 Heartbeat 无锁 / worker 不匹配 |

---

## 7. API 语义（M1）

### Deploy

- 输入：BPMN XML bytes  
- 行为：`deploy.Compile` 解析、M1 校验，内存登记 `Deployment`  
- 输出：`deployment_id`（UUIDv7）  
- **不**生成平行可执行图；快照即 `element.Process`

### CreateInstance

- 输入：`deployment_id`，可选变量  
- 行为：分配 instance/token id；写启动 COMMAND + PROCESS 启动 EVENT；Enter StartEvent  
- 推进到第一个等待点（通常 UserTask）或直至结束  

### Complete

- 输入：`process_instance_id`，`element_id`，`token_id`，可选变量  
- 行为：校验 token 在该元素 waiting；`Type` 从定义读取并写入 COMMAND；handler `OnComplete` 后继续自动步进  
- 不按 BPMN 类型拆 API；UserTask / ServiceTask / Timer catch / 日后等待点都走同一入口  

### FireDue

- 输入：无（扫描投影）
- 行为：收集 **waiting 且 `token.due_unix_ms` 已到** 的中间捕获 Timer，逐个调用 `Complete`
- **不**把 Timer 写成与 Element 平级的账本主语；到期只是触发统一 Complete
- `Recover` 从 `EventPayload.due_unix_ms` 还原 due；到期后仍需 `FireDue`（或手动 Complete）
- `gateway/cmd/sparrow` 每 200ms 调用一次；无独立 Timer RPC

### Activate

- 输入：`job_type`，可选 `max_jobs` / `wait` / `worker_id` / `lock_duration`  
- 行为：从投影收集 **waiting 且 `token.job_type` 匹配** 的 ServiceTask；加上内存租约后返回 Job 快照（含 instance / element / token / variables）  
- `wait=0` 立即返回（可为空）；否则长轮询，新的 `SERVICE_TASK ACTIVATED` 会唤醒  
- **不写 EventLog**。租约不是账本事实；`Recover` 后租约为空，waiting token 仍可再次 Activate  
- Worker 完成仍调用 `Complete(instance, element, token, vars)`；租约在成功 Complete 后释放  

### Fail

- 输入：`process_instance_id`，`element_id`，`token_id`，可选 `error_message`  
- 行为：仅 waiting 且带 `job_type` 的活动（ServiceTask）；写 FAILED COMMAND + EVENT；token **保持 waiting**；释放租约并唤醒 `Activate`  
- UserTask 等非 job 等待点 → `INVALID_STATE`  
- 不推进流程；重试靠再次 Activate，完成仍走 `Complete`  

### Heartbeat

- 输入：`process_instance_id`，`token_id`，`worker_id`，可选 `lock_duration`  
- 行为：延长内存租约；`worker_id` 必须与 Activate 持有者一致  
- **不写 EventLog**  

查询：

- `GetInstance`：投影快照  
- `ListEvents(process_instance_id)`：审计时间线  

---

## 8. 与 protocol 的边界

- processing **不**手写 `.pb.go`；只依赖 `protocol/gen/go/event/v1`
- 缺字段时：先在 `protocol/proto` 增加，再 `buf generate`，再改 processing
- 演进约定：
  - 元素行为扩展 → `Intent` / `Type` / **payload 字段**
  - 非元素事实（如纯 Deployment 元数据、Job 租约）→ 另议；不默认塞进 Element
- Worker 线协议在 `protocol/proto/job/v1`（`JobService`）；进程客户端在 `engine/v1`（`EngineService`）。均由 `gateway` 适配，**不**进入 processing。进程入口：`gateway/cmd/sparrow`。

---

## 9. 后续演进（不在 M1，设计预留）

| 阶段 | 能力 | 落点 |
|------|------|------|
| M2 | Message catch；Timer date/cycle/boundary | payload + Intent 等待语义；新 handler 或扩展 catch |
| M3 | Parallel/Inclusive、SubProcess、多 token | `Tokens` 多条目；gateway fork/join |
| M4 | Boundary / 补偿 / Incident | 新 Intent + payload |

单节点串行模型保持不变；分布式仅增加分区宿主，不改变 Handler 语义。

---

## 10. 测试策略

| 层级 | 内容 |
|------|------|
| 单元 | 各 handler：给定输入 → 期望 Effect（可逐步补） |
| 日志 | Memory / File Append/Read；File 重启后可读 |
| 端到端 | `testdata/m1_simple.bpmn`：Deploy → CreateInstance → Complete → PROCESS COMPLETED |
| 回归 | 不支持元素部署失败；非法 Complete → REJECTION；XOR default / 条件选路 |
| Job | Activate 领取 / Fail 释放重领 / Heartbeat 续租 / Recover 后仍可 Activate |
| Timer | `PT0S` 后 `FireDue` 完成；`PT1H` 未到期仍 waiting；Recover 后 due 仍在 |
| 恢复 | `Open` 后仍在 UserTask waiting，Complete 可完成实例 |

---

## 11. 实现进度

| 步骤 | 内容 | 状态 |
|------|------|------|
| 1 | `log` 内存 EventLog | 已完成 |
| 2 | `projection` + `ApplyEvent` | 已完成 |
| 3 | `deploy` 持有 `element.Process` | 已完成 |
| 4 | `handlers` 按文件拆分 + Registry | 已完成 |
| 5 | `engine` + `executor` API | 已完成 |
| 6 | 文件型 EventLog 与重启回放 | 已完成 |
| 7 | XOR 条件表达式（M1 子集） | 已完成 |
| 8 | ServiceTask + job_type | 已完成 |
| 9 | Job Activate（拉模型 + 内存租约） | 已完成 |
| 10 | Job Fail / Heartbeat | 已完成 |
| 11 | Job gRPC（`protocol/job.v1` + `gateway`） | 已完成 |
| 12 | Engine gRPC + `cmd/sparrow` | 已完成 |
| 13 | COMMAND 幂等 / 半截链 Recover | 已完成 |
| 14 | 中间捕获 Timer（duration）+ `FireDue` | 已完成 |
| 15 | Message catch / 更多元素 | 未开始 |

---

## 12. 代码对照

| 路径 | 职责 |
|------|------|
| `engine.go` | 薄门面：锁、COMMAND/REJECTION、emitter |
| `timers.go` | `FireDue`：到期的 timer catch 走 `Complete` |
| `jobs.go` | `Activate` / `Fail` / `Heartbeat`：Job 快照、长轮询、内存租约 |
| `open.go` / `recover.go` | `Recover` 回放 EVENT；未完成 COMMAND 接着跑 |
| `executor.go` | Enter/Complete、出边、流程完成判定 |
| `handlers/*.go` | 每元素一类文件；语义只在此扩展 |
| `expr/` | `${...}` → expr-lang 求值 |
| `deploy/` | Compile、`Store`、XOR 选路 |
| `projection/` | Instance / Token；EVENT → 投影 |
| `log/` | `EventLog` 接口（Memory / File） |
| `id.go` | UUIDv7 |
