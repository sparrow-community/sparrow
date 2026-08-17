# processing 程序设计

本文描述 `processing` 模块的目标形态、核心模型、包结构、处理流程与 M1 实现边界。  
实现应以本文为准；协议字段变更在 `protocol` 中进行，并保持「Event = 元素行为」的定位。

**M1 状态**：Deploy / CreateInstance / CompleteUserTask、XOR 条件、内存与文件 EventLog、`Open` 重启回放已可用。  
**M2 起步**：CompleteServiceTask；Job 细节在 `ActivityPayload.job_type`。  
尚未实现：Timer / Message、COMMAND 幂等、更多 BPMN 元素。

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
| 令牌 | `token_id` → `{element_id, active\|waiting}` | 实例内控制流位置（M1 单 token） |

**不另建** `graph.Node` / `graph.Flow`。`deploy` 在 `element.Process` 上提供查询辅助（`TypeOf`、`Outgoing`、`SequenceFlow`、`ChooseExclusiveOutgoing`）。

---

## 3. 架构总览

```text
                    ┌──────────────────────────────────────────┐
                    │              Engine（薄门面）               │
                    │  Deploy / CreateInstance / Complete…       │
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

Engine 对外能力：Deploy / CreateInstance / CompleteUserTask / GetInstance / ListEvents。

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
CreateInstance / CompleteUserTask:
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

### 5.2 幂等（未做）

同一 `cmd.id` 已成功处理则直接返回。崩溃导致半截 EVENT 链时，回放得到部分投影；文件日志会丢掉不完整的最后一条记录。

### 5.3 投影与令牌

`projection.Instance` 最少包含：

- 实例状态：`active | completed | terminated`
- `Tokens`：`token_id → {element_id, active|waiting}`
- 变量表（实例级 `name → json_value`）
- `ElementIntent`：元素级最近 Intent（辅助校验）

令牌更新 **只走 EVENT → ApplyEvent → applyToken**：

- Executor 只决定步进（`Effect.Wait` / 出边），不改 `Tokens`
- `USER_TASK` + `ACTIVATED` → `waiting`；`SEQUENCE_FLOW_TAKEN` 把位置写到 target
- PROCESS COMPLETED/TERMINATED 清空 `Tokens`
- 在线与 `Recover` 共用同一套规则

M1 为单 token；Parallel 等多 token 时仍落在同一 map，由 gateway handler 分裂/汇合。`GetInstance` 返回投影拷贝。

### 5.4 重启恢复

```text
Recover(ctx, eventLog, deploymentStore)
  → Store.LoadAll → Compile → deployments map
  → EventLog.ReadAll
  → 按序对每条 EVENT 调用 Instance.ApplyEvent（必要时先创建投影）
  → 恢复等待点（如 UserTask waiting），可继续 CompleteUserTask
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
| `SERVICE_TASK` | `service_task.go` | 同上等待；ACTIVATED 带 `ActivityPayload.job_type`；`CompleteServiceTask` |
| `EXCLUSIVE_GATEWAY` | `exclusive_gateway.go` | 非 default 条件按序求值，否则 default；payload 带 `taken_sequence_flow_id` |
| `SEQUENCE_FLOW` | `sequence_flow.go` | 经 transit 发 `SEQUENCE_FLOW_TAKEN`（不走 OnEnter） |
| `END_EVENT` | `end_event.go` | 完成后 `TryCompleteProcess` |

部署时 `validateM1` 拒绝尚未实现的元素（Parallel、SubProcess 等）。ServiceTask 已纳入可执行子集。

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
| XOR | EXCLUSIVE_GATEWAY | ACTIVATING → … → COMPLETED（payload 带 taken flow） |
| End | END_EVENT | … → COMPLETED |
| 实例结束 | PROCESS | COMPLETING → COMPLETED |

**等待点（UserTask / ServiceTask ACTIVATED）** 与 **SEQUENCE_FLOW_TAKEN** 必须在日志中可见。

### 6.3 变量

- `ActivityPayload.variables` / `ProcessPayload.variables`（`name` + `json_value`）
- 语义为 **delta**：合并进实例变量表
- CreateInstance 可带初始变量；CompleteUserTask 可带提交变量
- XOR：按 outgoing 顺序求值非 default 的条件（`expr.Eval` / expr-lang）；都不成立则走 default
- 条件：剥掉 BPMN `${...}` 后交给 [expr-lang/expr](https://github.com/expr-lang/expr)；变量为实例 JSON 值。单引号字符串会先归一成双引号。

### 6.4 拒绝示例

| code | 场景 |
|------|------|
| `NOT_FOUND` | 实例、部署或元素不存在 |
| `INVALID_STATE` | UserTask 未处于 waiting 却 Complete |
| `UNSUPPORTED_ELEMENT` | 定义含 M1 未支持元素，或无 handler |
| `INVALID_CONDITION` | XOR 条件表达式无法解析 |
| `NO_OUTGOING_FLOW` | 无法选出边 |

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

### CompleteUserTask / CompleteServiceTask

- 输入：`process_instance_id`，`element_id`，`token_id`，可选变量  
- 行为：写 COMPLETING COMMAND；`Executor.Complete` 后继续自动步进  
- ServiceTask 须处于 waiting；ACTIVATED 时 `job_type` 来自 BPMN `implementation`（否则 name / id）  

查询：

- `GetInstance`：投影快照  
- `ListEvents(process_instance_id)`：审计时间线  

---

## 8. 与 protocol 的边界

- processing **不**手写 `.pb.go`；只依赖 `protocol/gen/go/event/v1`
- 缺字段时：先在 `protocol/proto` 增加，再 `buf generate`，再改 processing
- 演进约定：
  - 元素行为扩展 → `Intent` / `Type` / **payload 字段**
  - 非元素事实（如纯 Deployment 元数据）→ 另议；不默认塞进 Element

---

## 9. 后续演进（不在 M1，设计预留）

| 阶段 | 能力 | 落点 |
|------|------|------|
| M2 | Timer / Message | payload + Intent 等待语义；新 handler 文件 |
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
| 9 | Timer / Message / 更多元素 | 未开始 |

---

## 12. 代码对照

| 路径 | 职责 |
|------|------|
| `engine.go` | 薄门面：锁、COMMAND/REJECTION、emitter |
| `open.go` | `Recover(log, store)`；`Open` 为文件便捷封装 |
| `executor.go` | Enter/Complete、出边、流程完成判定 |
| `handlers/*.go` | 每元素一类文件；语义只在此扩展 |
| `expr/` | `${...}` → expr-lang 求值 |
| `deploy/` | Compile、`Store`、XOR 选路 |
| `projection/` | Instance / Token；EVENT → 投影 |
| `log/` | `EventLog` 接口（Memory / File） |
| `id.go` | UUIDv7 |
