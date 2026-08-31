# AI Driven BPMN

Engine first ([`AGENTS.md`](./AGENTS.md)); AI Driven is how agents **consume** the engine — same BPMN definitions, same event log, explicit COMMANDs only.

## Idea

AI drafts and drives; BPMN is the versioned contract; Sparrow executes and records facts.

```text
Agent / 对话
    ↓  tools · MCP
Sparrow（COMMAND → EVENT）
    ↓
BPMN 定义 + 事件日志
```

**生成可变，契约不可默改。** 定义与实例的生效变更必须落 COMMAND/EVENT，不能只在对话里发生。

## Why BPMN

流程需要跨人、跨时间的一致符号与版本；Agent 计划短命，契约层提供共享上下文、责任边界与回放事实。Sparrow 提供：**版本化定义 + 确定性执行 + 审计回放**。

## Sparrow 的角色

- **执行与审计** — Deploy、启停实例、Complete 等待点、REJECTION
- **Agent 接口** — gRPC（`engine.v1` / `job.v1`）；MCP 等作为外部 consumer
- **行为映射** — Agent 逻辑落在 Service Task、User Task、BPMN 扩展上

**Drive 模式：** 对话起草定义、查询实例、辅助待办；引擎只收 COMMAND。

**前置：** [`AGENTS.md`](./AGENTS.md) § Remaining 中引擎项优先完成。

## 后续（引擎就绪后）

| 项 | 目标 |
|----|------|
| Task extensions | Agent 元数据约定（扩展属性，非新 core 元素） |
| Correlation 查询 | 按业务键查实例 |
| Agent 操作指南 | deploy → create → complete → list events |
| 失败恢复 | 对齐引擎 Incident increment |

AI 相关能力不得绕过 COMMAND/EVENT 或无声改投影。
