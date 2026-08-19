# Sparrow — 项目状态（跨会话速览）

> **用途**：新开 Cursor 对话时 `@STATUS.md`，快速恢复「做到哪了、设计约束、下一步」。  
> **细节以文档为准**：语义与架构看 [`processing/DESIGN.md`](processing/DESIGN.md)；模块与夹具看 [`processing/README.md`](processing/README.md)；工作区约定看 [`AGENTS.md`](AGENTS.md)；产品定位看 [`AI-Driven-BPMN.md`](AI-Driven-BPMN.md)。  
> **维护**：每合并一块功能后更新本文件（约 5 分钟），不必改长对话历史。

---

## 一句话

单节点、事件账本为真相源的轻量 BPMN 引擎；**Event = 行为，Element = 主语**；对外 API 薄（Deploy / CreateInstance / Complete / FireDue / PublishMessage / Job Activate），语义在 `processing/handlers/` 扩展。

## 当前里程碑

| 阶段 | 状态 |
|------|------|
| **M1** | 完成：Start → UserTask → XOR → End；文件日志 + `Recover` |
| **M2** | 完成：ServiceTask + Job；Timer/Message catch；Timer/Message boundary（打断型 + 非打断型） |
| **M3** | 起步：Parallel gateway fork/join（多 token）+ boundary 组合回归 |

**`main` 最新提交**（更新时改这里）：

```text
1a93436 Add timeCycle re-arm for non-interrupting timer boundaries
xxxxxxx Support dual boundary (timer+message) on same activity
6e0d45e Support non-interrupting timer and message boundaries
cf47acf Add STATUS.md for cross-session project continuity.
3eab8ef Document parallel gateway in processing design.
```

## 已实现（运行时）

- **核心**：`Deploy` / `CreateInstance` / `Complete`；`Recover` / `Open`；实例内串行锁
- **网关**：Exclusive（条件 + default）；Parallel（fork mint token / join 同步）
- **等待与完成**：UserTask、ServiceTask、中间 Timer catch、中间 Message catch — 统一 `Complete`
- **Timer**：`timeDuration` / `timeDate` / `timeCycle`（仅首次到期）；`Engine.FireDue`；`cmd/sparrow` 轮询
- **Message**：`PublishMessage`（name + 可选 correlation_keys + 内存缓冲，Recover 后缓冲丢失）
- **Boundary**：Timer / Message 可挂 UserTask 或 ServiceTask；支持打断型与非打断型；同一活动可同时挂 Timer + Message
- **Job**：`Activate` / `Fail` / `Heartbeat`（内存租约，非账本）
- **传输**：`gateway` — `engine.v1` + `job.v1` gRPC

## 明确不做（当前阶段）

- 流程定义**版本管理与迁移**（每次 `Deploy` 新 `deployment_id`）
- 同一活动上 **三个及以上 boundary**
- Inclusive / EventBased gateway、SubProcess、Throw、补偿、Incident
- 集群 / 多活；跨重启的 **Job 租约** 与 **消息缓冲** 持久化
- 复制 Camunda 产品广度（建模器、Cockpit 等）

## 下一步候选（按优先级）

1. **SubProcess / Inclusive gateway** — 范围更大，后置
2. 跨重启的 Job 租约 / 消息缓冲持久化
3. EventBased gateway / Signal / Compensation

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
