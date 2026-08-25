# Sparrow — 规划与状态

跨会话恢复「做到哪了、下一步」。设计见 [`processing/DESIGN.md`](processing/DESIGN.md)；夹具见 [`processing/README.md`](processing/README.md)；约定见 [`AGENTS.md`](AGENTS.md)。

**定位**：单节点、事件账本真相源；轻量且功能完备的 BPMN **执行引擎**（非 Camunda 类产品套件）。Event = 行为，Element = 主语。

---

## 路线图

| 阶段 | 内容 | 状态 |
|------|------|------|
| M1 | Start → UserTask → XOR → End；EventLog + Recover | 完成 |
| M2 | ServiceTask/Job；Timer/Message catch + boundary | 完成 |
| M3 | 网关 / SubProcess / Throw / ESP / Compensation / Signal boundary | 完成 |
| M4 | Error 事件（boundary / end / ThrowError / Recover / gateway） | 完成 |
| M4b | 嵌套 Event Sub-Process（嵌入式 SubProcess 内） | 完成 |
| M4c | CallActivity（同定义内联 called process） | 完成 |
| 其后 | CallActivity 子实例 / IO 映射 / 版本 | 规划中 |

**下一步（按序）**：① CallActivity 增强（独立子实例 / 变量映射）或 Error ESP → ② 版本管理与迁移

**暂缓**：同一活动 ≥3 个 **同类** waiting boundary；Incident；instantiate EBG；补偿进未完成 SubProcess；ESP 嵌套于 ESP；跨部署 CallActivity；集群；产品套件
（同类各一：timer + message + signal 可共存）

---

## 已实现元素

| 元素 | 范围 |
|------|------|
| Process | 启动 / 完成 |
| StartEvent | none |
| EndEvent | none；error（抛出）；**compensate** |
| SequenceFlow | 条件（XOR/OR 用） |
| UserTask | 等待；统一 Complete |
| ServiceTask | Job（Activate / Fail / Heartbeat） |
| ExclusiveGateway | 条件 + default |
| InclusiveGateway | OR-split / OR-join |
| ParallelGateway | fork / join |
| EventBasedGateway | exclusive（先到取消兄弟）；parallel（保留兄弟） |
| SubProcess | 嵌入式（可多层）；补偿 boundary 订阅 |
| CallActivity | **同定义** `calledElement`；同实例内联执行 |
| Event Sub-Process | 流程级或**嵌入式 SubProcess 内**；message / timer / signal start；打断 / 非打断 |
| IntermediateCatchEvent | timer / message / signal |
| IntermediateThrowEvent | none / message / signal / compensate |
| BoundaryEvent | timer / message / signal（打断 / 非打断）；compensate（含 SubProcess）；error（打断） |
| Association | 补偿 handler 关联 |

**未实现（常见）**：Escalation / Link / Conditional / Terminate；Send/Receive/Manual/BusinessRule Task；Multi-instance；ESP 嵌套于 ESP；Error ESP；CallActivity 独立子实例 / IO 映射 / 跨部署

---

## 运行时能力

- API：`Deploy` / `CreateInstance` / `Complete` / `ThrowError` / `FireDue` / `PublishMessage` / `PublishSignal` / Job 三件套；查询 `GetInstance` / `ListEvents`
- 持久化：EventLog + deploy.Store + runtime.Store（租约 / 消息缓冲）；`Recover` / `Open`
- 传输：`gateway`（`engine.v1` + `job.v1`）

**最新提交**

```text
（本轮）清理 DESIGN §8 已还清项；打断 scope 的 SubProcess TERMINATED 可无 token_id（审计）以免复活投影 token。
（设计还债）Unify terminateScopeTokens; EventPayload.token_wait for compensate end.
3451012 Support CallActivity for called processes in the same definitions.
```

---

## 开场 / 索引

```text
@STATUS.md @AGENTS.md @processing/DESIGN.md
本轮任务：<一件事>
```

| 问题 | 文档 |
|------|------|
| 语义与架构 | `processing/DESIGN.md` |
| 路线图 / 元素清单 / 下一步 | 本文 |
| 夹具 BPMN | `processing/README.md` |
| 协议与命令 | `AGENTS.md` |

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```
