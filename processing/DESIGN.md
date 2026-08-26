# processing 程序设计

`processing` 是 Sparrow 的 BPMN 执行内核。本文只描述**核心设计**：模型、架构、语义契约与模块边界。  
**Milestone / implemented snapshot**: [`AGENTS.md`](../AGENTS.md). **Active increment** (child CallActivity instance / IO mapping / Error Event Sub-Process / revision coexistence): [`specs/001-engine-completeness/`](../specs/001-engine-completeness/). **Governance**: [`.specify/memory/constitution.md`](../.specify/memory/constitution.md).

实现以本文的语义为准；协议字段变更在 `protocol` 中进行。约定：**Event = 元素行为**，**Element = 行为主语**。

---

## 1. 目标与约束

**目标**

- 轻量级且功能完备的 BPMN **执行引擎**（可执行语义完备，而非长期停留在最小子集）
- 单节点可运行；少中间层、薄 API；append-only 事件账本为真相源
- 行为可审计、可回放、可拒绝（COMMAND → EVENT / REJECTION）
- 定义层直接复用 `bpmn/element`，不为执行另造平行图模型
- API 保持简单，日后可在不改语义的前提下按 `process_instance_id` 分区扩展

**不做（设计层）**

- 把 Job / Timer / Message 提升为与 Element 同级的账本主语
- 为「整洁」堆叠无必要的 adapter / Node / Flow 包装类型
- Camunda 类产品套件、DMN / CMMN（非本引擎本体）

Deferred capabilities (version migration, cluster, etc.) are listed in `AGENTS.md` and the active spec Assumptions; not expanded here.

---

## 2. 核心概念

### 2.1 行为账本

每条 `event.v1.Event` 描述一次行为：

| 字段 | 含义 |
|------|------|
| `record_type` | COMMAND / EVENT / REJECTION |
| `deployment_id` / `process_instance_id` / `process_version` | 定义与实例上下文 |
| `element` | 主语：Type、id、token_id、Intent、payload |
| `source_record_id` | 因果（通常指向触发本条的 COMMAND） |
| `rejection` | 仅 REJECTION：code + message |

日志偏移由存储层维护；回放按追加顺序。投影可丢，EventLog 不可丢。

### 2.2 Element 是行为主语

```text
谁：Type + element.id + token_id
做了什么：Intent
附带什么：payload（变量、选中的 sequenceFlow、job/timer/message 等）
```

`Element.Type` 与 BPMN 独立元素对齐（`PROCESS` + 具体 FlowElements）。  
不使用「粗 Type + kind」合并不同 BPMN 元素。  
Job / Timer / Message 等待是该元素行为的**载荷与阶段**，不是新的顶层主语。

### 2.3 串行与分区键

- 分区键：`process_instance_id`
- 同一实例上的 COMMAND 严格串行（每实例一把锁）
- 单节点 = 一个逻辑分区；多节点只是多分区宿主，Handler 语义不变

### 2.4 定义 vs 运行时

| 层 | 来源 | 作用 |
|----|------|------|
| 定义 | `bpmn` → `element.Process` | 静态结构 |
| 部署 | `deploy.Deployment` | 校验后的不可变快照 + `deployment_id` |
| 实例 | `projection.Instance` | 状态 / 变量 / 令牌（投影） |
| 令牌 | `token_id` → 元素位置与等待载荷 | 控制流位置 |

不另建 `graph.Node` / `graph.Flow`。`deploy` 在 `element.Process` 上查询，并缓存已解析的 timer / message / signal 等。

---

## 3. 架构

```text
Engine（薄门面）
  Deploy / CreateInstance / Complete
  Activate / Fail / Heartbeat
  FireDue / PublishMessage / PublishSignal
  实例锁 · COMMAND / EVENT / REJECTION
        │
        ▼
   Executor — Enter / Complete · 应用 Effect · 出边步进
        │
        ▼
   handlers（按 Element.Type）— OnEnter / OnComplete → Effect

  deploy.Deployment  ← element.Process
  projection.Instance ← ApplyEvent（仅 EVENT）
  EventLog / deploy.Store / runtime.Store（runtime 可选，非账本）
```

**原则**

- 元素语义只在 `handlers/`；Engine 不写具体生命周期分支
- 持久化三点注入：`log.EventLog`（账本）、`deploy.Store`（定义）、`runtime.Store`（Job 租约 + 消息缓冲）
- `runtime.Store == nil` 时租约与缓冲仅进程内有效；丢失后流程仍可 replay，辅助语义退化
- COMMAND 处理内联在 Engine；不另留空 Processor 接口

### 包职责

| 路径 | 职责 |
|------|------|
| `engine.go` 等 | API、锁、写日志、FireDue / Publish* / Job |
| `executor.go` | Enter / Complete、出边、流程完成判定 |
| `handlers/` | 每元素一类；语义只在此扩展 |
| `deploy/` | Compile、校验、结构查询与缓存 |
| `projection/` | Instance / Token；EVENT → 投影 |
| `log/` / `runtime/` | EventLog；非账本辅助状态 |
| `open.go` / `recover.go` | 回放 EVENT；未完成 COMMAND 接着跑 |
| `expr/` | `${...}` → expr-lang |

扩展新元素：**新增 handler + 注册 Registry + 在 deploy 校验中放开**。

---

## 4. Effect 与处理循环

### 4.1 Effect

Handler 产出，由 Executor 应用：

| 字段 | 含义 |
|------|------|
| `Records` | 写成 EVENT 的元素行为 |
| `Wait` | 停止自动步进（等待点） |
| `TakeOutgoing` / `OutgoingFlowID` | 出边推进（XOR 可指定 flow） |
| `TryCompleteProcess` | End 后尝试 PROCESS 完成 |
| `Publication` | 锁释放后投递 message/signal（避免重入） |

瞬时元素可用完整 Intent 链一次写出。

### 4.2 COMMAND 路径

```text
CreateInstance / Complete / …
  取 deployment + instance → 加锁 → 校验投影
  Append(COMMAND)
  Executor.Enter | Complete → handler → Effect
    → emitter(Records)  // Append EVENT + ApplyEvent
    → 出边 / 等待 / 完成流程
  失败：COMMAND + REJECTION
```

同一 COMMAND 可连写多条 EVENT，共享 `source_record_id`。

### 4.3 幂等与半截链

COMMAND 先入账再写 EVENT 链；崩溃可能只剩 COMMAND 或写了一半。

`Recover` 回放 EVENT 后检查每条 COMMAND：

- 已有 REJECTION → 跳过
- 已有 EVENT 且投影稳定（waiting / completed / terminated）→ 跳过
- 否则**接着执行同一 COMMAND**（不新写 COMMAND）；emitter 对已出现的 `(source, Type, Intent, id, token)` 不再追加

客户端未带同一 `cmd.id` 的重试仍是新 COMMAND。

### 4.4 投影与令牌

令牌更新**只走** EVENT → `ApplyEvent`：

- Executor 只决定步进，不直接改 `Tokens`
- 等待点：活动 / catch 的 `ACTIVATED` 写入 `waiting` 及 due / message_name / signal_name / job_type / boundary_id 等
- `SEQUENCE_FLOW_TAKEN` 更新位置；PROCESS 结束清空 Tokens
- 在线与 `Recover` 共用同一套规则

多 token（Parallel 等）落在同一 `Tokens` map，由 gateway handler 分裂/汇合。

### 4.5 恢复

```text
Recover → 加载定义 → 回放 EVENT → 加载 runtime（可选）
       → 重驱未完成 COMMAND → 等待点可继续 Complete / Activate / FireDue …
```

runtime 加载时丢弃过期/无效 lease 与过期/死实例消息缓冲；实例结束时 sweep。  
`Open(dataDir)` 是文件实现的便捷封装。COMMAND / REJECTION 不驱动投影。

---

## 5. 可执行语义（契约）

Deploy validation rejects unsupported elements. The following is the current semantic contract (implementation progress: `AGENTS.md`).

### 5.1 元素与行为

| Type | 要点 |
|------|------|
| `PROCESS` | 实例启动 / 完成 |
| `START_EVENT` | 瞬时生命周期后出边；**Event Sub-Process start** 在 parent scope 打开时 `ACTIVATED`（`EventPayload.event_sub_process_element_id` + message/timer/signal） |
| `END_EVENT` | 完成后 `TryCompleteProcess` |
| `USER_TASK` / `SERVICE_TASK` | ACTIVATED 后 Wait；统一 `Complete` 完成；ServiceTask 带 `job_type` |
| `INTERMEDIATE_CATCH_EVENT` | Timer（due）/ Message / Signal → Wait；由 FireDue / Publish* / Complete 完成 |
| `INTERMEDIATE_THROW_EVENT` | None 瞬时；Message/Signal 经 Publication 锁外投递；Compensate 逆序执行 handler |
| `BOUNDARY_EVENT` | Timer / Message / Signal / Compensate；打断型 TERMINATE 活动；非打断型 mint 新 token；补偿在活动 COMPLETED 后订阅 |
| `EXCLUSIVE_GATEWAY` | 条件顺序求值，否则 default |
| `PARALLEL_GATEWAY` | fork / join |
| `INCLUSIVE_GATEWAY` | OR-split / OR-join（可达 token 防死锁） |
| `EVENT_BASED_GATEWAY` | 先到 cancel 兄弟，或 Parallel 保留兄弟；instantiate 不支持 |
| `SUB_PROCESS` | 嵌入式子流程；Event Sub-Process（`triggeredByEvent`）可挂在流程或嵌入式子流程上 |
| `CALL_ACTIVITY` | 同一定义内 `calledElement` → 另一 process；同一实例 host token 停在 CallActivity、child 进被调流程；被调流程内 Event Sub-Process 在进入时武装 |
| `SEQUENCE_FLOW` | `SEQUENCE_FLOW_TAKEN`（经 transit，不走 OnEnter） |

**统一完成入口**：等待点（UserTask / ServiceTask / catch / 部分 throw·compensate）都走 `Complete`；类型来自部署，不按类型拆 API。

### 5.2 Intent 模式

| 场景 | Intent（EVENT） |
|------|-----------------|
| 创建实例 | PROCESS ACTIVATING → ACTIVATED |
| 瞬时元素 | ACTIVATING → ACTIVATED → COMPLETING → COMPLETED |
| 等待进入 | … → ACTIVATED（Wait） |
| 等待完成 | COMPLETING → COMPLETED |
| 走过流 | SEQUENCE_FLOW_TAKEN |
| 打断 boundary | 活动 TERMINATING → TERMINATED；boundary COMPLETING → COMPLETED |
| Job 失败 | SERVICE_TASK FAILED（仍 waiting） |
| 实例结束 | PROCESS COMPLETING → COMPLETED |

等待点的 `ACTIVATED` 与 `SEQUENCE_FLOW_TAKEN` 必须在账本中可见。

### 5.3 变量与条件

- payload 中的 `variables` 为 **delta**，合并进实例变量表
- XOR / Inclusive 条件：剥掉 `${...}` 后用 expr-lang；失败可走 default（若有）

### 5.4 辅助触发（非账本主语）

| API | 行为 |
|-----|------|
| `FireDue` | 到期 timer catch / timer boundary → `Complete` |
| `PublishMessage` | 按 name（+ 可选 correlation_keys / instance）匹配 → `Complete`；无 waiter 则 FIFO 缓冲（可持久化到 runtime） |
| `PublishSignal` | 按 name 匹配 catch / signal boundary → `Complete`；**无缓冲** |
| `Activate` / `Fail` / `Heartbeat` | Job 拉模型与租约；**不写** EventLog（除 Fail 的 FAILED） |

---

## 6. API 契约

| API | 语义 |
|-----|------|
| `Deploy` | Compile + 校验 → `deployment_id`；快照即 `element.Process` |
| `CreateInstance` | 启动 COMMAND + PROCESS；Enter Start；步进至等待点或结束 |
| `Complete` | 校验 waiting；handler OnComplete 后自动步进 |
| `FireDue` / `PublishMessage` / `PublishSignal` | 见上；gateway 暴露对应 RPC |
| `ThrowError` | 等待中的 UserTask/ServiceTask 抛 BPMN error；gateway `engine.v1` 暴露 |
| `Activate` / `Fail` / `Heartbeat` | Job 租约生命周期 |
| `GetInstance` / `ListEvents` | 投影快照 / 审计时间线 |

拒绝码示例：`NOT_FOUND`、`INVALID_STATE`、`UNSUPPORTED_ELEMENT`、`INVALID_CONDITION`、`INVALID_ARGUMENT`。

---

## 7. 与 protocol / gateway 的边界

- processing **不**手写 `.pb.go`；只依赖 `protocol/gen/go`
- 元素行为扩展 → `Intent` / `Type` / **payload 字段**
- 非元素事实（租约、消息缓冲）→ `runtime.Store`，不默认进 EventLog
- `job.v1` / `engine.v1` 由 `gateway` 适配；**不**进入 processing
- 进程入口：`gateway/cmd/sparrow`（含 `FireDue` 轮询）

---

## 8. 已知设计张力（实现中标记）

下列不是阻塞缺陷，但是当前模型下需要注意的不合理点 / 技术债；后续可重构时优先处理。已还清的条目直接删除，不保留「已还清」坟场。

1. **CallActivity v1 = 同实例内联，不是子 process instance**  
   当前 CallActivity 在**同一** `process_instance_id` 上把 token 送入同一定义文件中的被调 process；host token 停在 CallActivity，内部用 child token（与嵌入式 SubProcess 同构）。被调流程 scope 的 host 解析与打断/终止、以及进入时按 called process id 武装 Event Sub-Process，已与嵌入式 SubProcess 对齐。这与 BPMN「独立 called process instance」不完全一致。限制：同一被调 process 仅允许一个 CallActivity、禁止递归 CallActivity、元素 id 必须在 Definitions 内全局唯一、尚无 IO 映射 / 跨部署 calledElement / 版本选择；CallActivity 上的 boundary / 补偿订阅仍弱于嵌入式 SubProcess。独立子实例与版本管理见 [`specs/001-engine-completeness/`](../specs/001-engine-completeness/)，不在本文展开实现步骤。
