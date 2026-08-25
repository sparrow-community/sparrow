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
| M4（进行中） | Error 事件 | 进行中 |
| 其后 | 补偿收尾 → CallActivity / 版本 | 规划中 |

**下一步（按序）**：① Error 事件（进行中）→ ② 补偿收尾（End / 嵌套 scope）→ ③ 嵌套 ESP / CallActivity → ④ 版本管理与迁移

**暂缓**：同一活动 ≥3 个 **同类** waiting boundary；Incident；嵌套 ESP；instantiate EBG；补偿进未完成 SubProcess；集群；产品套件
（同类各一：timer + message + signal 可共存）

---

## 已实现元素

| 元素 | 范围 |
|------|------|
| Process | 启动 / 完成 |
| StartEvent | none |
| EndEvent | none；**error（抛出）** |
| SequenceFlow | 条件（XOR/OR 用） |
| UserTask | 等待；统一 Complete |
| ServiceTask | Job（Activate / Fail / Heartbeat） |
| ExclusiveGateway | 条件 + default |
| InclusiveGateway | OR-split / OR-join |
| ParallelGateway | fork / join |
| EventBasedGateway | exclusive（先到取消兄弟）；parallel（保留兄弟） |
| SubProcess | 嵌入式（可多层） |
| Event Sub-Process | 流程级；message / timer start；打断 / 非打断 |
| IntermediateCatchEvent | timer / message / signal |
| IntermediateThrowEvent | none / message / signal / compensate |
| BoundaryEvent | timer / message / signal（打断 / 非打断）；compensate；**error（打断）** |
| Association | 补偿 handler 关联 |

**未实现（常见）**：Escalation / Link / Conditional / Terminate；CallActivity；Send/Receive/Manual/BusinessRule Task；Multi-instance；嵌套 ESP；Error ESP

---

## 运行时能力

- API：`Deploy` / `CreateInstance` / `Complete` / `ThrowError` / `FireDue` / `PublishMessage` / `PublishSignal` / Job 三件套；查询 `GetInstance` / `ListEvents`
- 持久化：EventLog + deploy.Store + runtime.Store（租约 / 消息缓冲）；`Recover` / `Open`
- 传输：`gateway`（`engine.v1` + `job.v1`）

**最新提交**

```text
（待提交）Support BPMN error events (boundary, end, ThrowError).
aea6050 Support signal boundary on waiting activities and scopes.
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
