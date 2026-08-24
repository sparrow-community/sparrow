# processing

Sparrow 运行时模块：单节点、事件驱动的 BPMN 执行引擎。

## 职责

- 接收并处理 `protocol` 中的行为记录（`event.v1.Event`）
- 将 BPMN 定义（来自 `bpmn`）变为可执行语义
- 以 **append-only 事件日志** 作为执行与审计的真相源
- 在内存中维护可由日志重建的实例投影（状态）

本模块 **不负责** BPMN XML 解析细节（见 `bpmn`），**不负责** 事件序列化契约（见 `protocol`），也 **不负责** gRPC（见 `gateway`）。

## 文档

- **设计**（模型、架构、语义契约）：[DESIGN.md](./DESIGN.md)
- **规划与现状**（路线图、已实现、下一步）：[../STATUS.md](../STATUS.md)

## 端到端夹具

`testdata/m1_simple.bpmn`（Start → UserTask → XOR → End）；`testdata/m2_service_task.bpmn`；`testdata/m2_timer_catch.bpmn`（Start → Timer catch `PT0S` → End）；`testdata/m2_timer_catch_date.bpmn`（`timeDate` 已过期）；`testdata/m2_timer_catch_cycle.bpmn`（`R/PT0S`）；`testdata/m2_timer_boundary.bpmn`（UserTask + interrupting `PT0S` boundary）；`testdata/m2_timer_boundary_non_interrupt.bpmn`（UserTask + non-interrupting `PT0S` boundary）；`testdata/m2_timer_boundary_non_interrupt_cycle.bpmn`（UserTask + non-interrupting `R2/PT0S` timer boundary）；`testdata/m2_message_catch.bpmn`（Start → Message catch → End）；`testdata/m2_message_after_task.bpmn`（UserTask 后再 Message catch）；`testdata/m2_message_boundary.bpmn`（UserTask + interrupting message boundary）；`testdata/m2_message_boundary_non_interrupt.bpmn`（UserTask + non-interrupting message boundary）；`testdata/m3_signal_boundary.bpmn` / `m3_signal_boundary_non_interrupt.bpmn` / `m3_timer_signal_boundary.bpmn` / `m3_subprocess_signal_boundary.bpmn`；`testdata/m3_parallel_fork_join.bpmn`（Parallel fork → 两路 UserTask → join → End）；`testdata/m3_parallel_timer_boundary.bpmn` 与 `testdata/m3_parallel_message_boundary.bpmn`（Parallel fork + boundary 分支回归）；`testdata/m2_dual_boundary.bpmn`（UserTask + Timer + Message dual boundary）；`testdata/m2_subprocess.bpmn`（Start → SubProcess(Start → UserTask → End) → End）；`testdata/m3_none_throw.bpmn` / `m3_message_throw.bpmn` / `m3_signal_throw.bpmn`；`testdata/m3_signal_catch.bpmn`；`testdata/m3_parallel_message_throw_catch.bpmn`（并行分支 throw 唤醒 sibling catch）；`testdata/m3_event_subprocess_message.bpmn` / `m3_event_subprocess_message_ni.bpmn` / `m3_event_subprocess_timer.bpmn`；`testdata/m2_event_based_gateway.bpmn`（exclusive）；`testdata/m3_parallel_event_based_gateway.bpmn`（parallel，不取消兄弟）；`testdata/m3_compensation.bpmn`（compensation boundary + throw）。

## 包布局（摘要）

```text
processing/
├── engine.go / jobs.go / timers.go / messages.go / executor.go / open.go   # API、Activate、FireDue、PublishMessage、Recover
├── expr/            # 条件表达式（expr-lang）
├── deploy/          # Compile + Store；查询走 Process
├── handlers/        # 一元素一文件 + Registry
├── projection/      # Instance / Token / ApplyEvent
└── log/             # EventLog 接口（Memory / File）
```

无平行 `graph` 模型；元素语义只加在 `handlers/`。

## 依赖

```text
processing
  ├── bpmn          // 流程定义模型
  └── protocol      // event.v1 Event / Element / Intent / Type
```

Go module: `github.com/sparrow-community/sparrow/processing`

## 设计要点（摘要）

1. **Event = 行为描述**；`Element` 是行为主语（哪个 BPMN 元素、哪个 Intent）。
2. Job / Timer / Message 等机制细节放在 **Element.payload**，不与 Element 平级另起主语。
3. 同一 `process_instance_id` **严格串行** 处理；单节点即单个（逻辑）分区。
4. 状态可丢，只要 `EventLog` + `deploy.Store` 在，即可 `Recover` 重建。
5. **少中间层**：部署直接复用 `bpmn/element.Process`；语义只加在 `handlers/`。

持久化注入：

```go
eng, err := processing.Recover(ctx, myEventLog, myDeploymentStore)
```

Memory / File 只是内置实现。`Open(dataDir)` 等于 `Recover(File, DirStore)`。

## 开发命令

在仓库根目录（`go.work`）：

```shell
go test ./processing/
```

实现推进以元素 Handler 扩展为主（见 `handlers/`），避免把语义堆回 `engine.go`。
