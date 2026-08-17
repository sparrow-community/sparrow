# processing

Sparrow 运行时模块：单节点、事件驱动的 BPMN 执行引擎。

## 职责

- 接收并处理 `protocol` 中的行为记录（`event.v1.Event`）
- 将 BPMN 定义（来自 `bpmn`）变为可执行语义
- 以 **append-only 事件日志** 作为执行与审计的真相源
- 在内存中维护可由日志重建的实例投影（状态）

本模块 **不负责** BPMN XML 解析细节（见 `bpmn`），也 **不负责** 事件序列化契约定义（见 `protocol`）。

## 现状

| 能力 | 状态 |
|------|------|
| UUIDv7 ID（`id.go`） | 已有 |
| 内存 EventLog（`log`） | 已有 |
| 实例投影与令牌（`projection`） | 已有 |
| Deploy（`deploy`，直接持有 `element.Process`） | 已有 |
| 元素 Handler（`handlers/`，按类型分文件） | 已有 |
| Executor 令牌推进 | 已有 |
| Engine：Deploy / CreateInstance / CompleteUserTask | 已有 |
| 文件 / 内存 EventLog + Store（可替换） | 已有 |
| `Recover` 回放 / `Open` 文件便捷入口 | 已有 |
| XOR 条件选路（`expr`，default 回退） | 已有 |

端到端夹具：`testdata/m1_simple.bpmn`（Start → UserTask → XOR → End）。

详细程序设计见 [DESIGN.md](./DESIGN.md)。

## 包布局（摘要）

```text
processing/
├── engine.go / executor.go / open.go   # API、Recover、Open
├── expr/            # M1 条件表达式
├── deploy/          # Compile + Store 接口
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
5. M1 可执行子集：`Start → UserTask → ExclusiveGateway → End`（含 SequenceFlow）。
6. **少中间层**：部署直接复用 `bpmn/element.Process`。

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
