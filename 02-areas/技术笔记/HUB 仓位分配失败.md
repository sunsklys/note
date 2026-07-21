# HUB 仓位分配失败

> Source: WPS 知识库 / 领域/技术笔记/HUB 仓位分配失败
> file_id: h3399arbJxMghuvmBdm2rxuZdDYkhc28u

---

- 事故类型：单订单永久卡死（仓位分配持续失败，无限重试无法自愈）
- 影响范围：单台机器出杯完全停摆，直到人工清理 Redis 锁
- 涉及模块：internal/server/DetailSpace.go、internal/service/ros/robot_system.go
- 根因类型：双重缺陷叠加（锁清理缺陷 + 算法容错缺陷）
- 验证状态：已通过独立复现（runtime）、Oracle 对抗审查、路径核查三路交叉验证

## 一、TL;DR

一次取餐/换杯流程的异常退出，在 Redis 里留下了一把永不过期的"取餐目标位锁"（位置 2）。而仓位选择算法遇到"目标位被锁"时没有回退到次优位的能力——它用粗暴截断候选列表的方式处理锁定，既跳过了已扫描的空闲位置（位置 1），又破坏了自身的重试条件。  
锁残留让目标位不可用，算法不容错让次优位也选不到，两者叠加，一把残留锁就锁死了整台机器的出杯流程。订单 971759（门店 215 / 机器 168）连续 3 次重试（12:43:37 / 38 / 39）完全相同失败，即此机制所致。  

## 二、问题现象

### 2.1 故障订单


| 字段 | 值 |
| --- | --- |
| 门店 entry | 215 |
| 机器 ID | 168 |
| 订单详情 ID | 971759 |
| 机器机型 | G1 桌面型（MachineModelIDDeskB = 6） |
| 订单杯数 | 单杯（orderDetails 长度 = 1） |

### 2.2 关键日志链（已按时间排序）


```javascript
INFO  DetailSpace.go:481  仓位信息读取成功: {"free":4,"order":{},"ros":{"1":0,"2":0,"3":0,"4":0}}
INFO  DetailSpace.go:785  次数数据3(getCupStandCountInfo 返回): [4, 3, 1, 2]
ERROR DetailSpace.go:660  获取的位置和最少使用次数位置不同: i=1, cupCount=[1,2]
ERROR DetailSpace.go:748  to cup space is locked: space=2          ← 位置 2 被锁
ERROR DetailSpace.go:669  位置已锁定: i=2, cupCount=[1]            ← 截断发生
ERROR DetailSpace.go:660  获取的位置和最少使用次数位置不同: i=3, cupCount=[1]
ERROR DetailSpace.go:660  获取的位置和最少使用次数位置不同: i=4, cupCount=[1]
ERROR DetailSpace.go:720  创建仓位数据不匹配: detailSpace=[]
ERROR DetailSpace.go:597  创建仓位信息失败: empty slice found
ERROR saas.go:136         仓位分配失败: 创建仓位数据为0
（以上链路在 37s/38s/39s 连续重试 3 次，每次完全相同）
```

### 2.3 状态不一致（核心矛盾）


| 状态来源 | 位置 2 的状态 |
| --- | --- |
| SHOP_MACHINE_SPACE_INFO_KEY（仓位主状态） | ✅ 空闲（ros[2]=0，free=4，order={}） |
| SHOP_MACHINE_TO_SPACE_LOCK_KEY（取餐目标位锁） | 🔒 锁定（SIsMember(2)=true） |

仓位物理空闲，但分配锁判定它被占用 —— 两个状态层面不一致。  

## 三、根因分析

这是两个独立缺陷的组合。任一单独存在都不会造成本次故障；两者同时命中才导致"一把残留锁 = 整单永久卡死"。  

### 3.1 缺陷一：取餐目标位锁的清理路径有漏洞（锁残留）

锁的写入（internal/service/ros/robot_system.go 的 checkBackupCup）：  

| 行号 | 触发条件 | 写入 |
| --- | --- | --- |
| 1741 | 备用杯托且 changeCup==0 && toCup!=0 | SAdd(sKey, toCup) |
| 1749 | 备用杯托且 changeCup!=0 | SAdd(sKey, toCup, changeCup) |
| 1760 | 常规杯托（Space<=10 && Space>0） | SAdd(sKey, oDetailSpace.Space) |

锁的清理（SRem，共 7 处）：panel/tool.go:351、ros/recover.go:365、robot_system.go:2101/2150/2151/2170/2171。  
漏洞：PickOrderDetail（robot_system.go:2067）在以下三处提前 return，而所有 SRem 都位于这三处之后：  

| 行号 | 返回条件 | 后果 |
| --- | --- | --- |
| 2083 | len(req.Space)==0（无位置信息） | 锁残留 |
| 2088 | FindListByMap DB 错误 | 锁残留 |
| 2143 | syncOrderPickStatus 失败（代码注释明写"不释放位置"，设计如此等待人工干预） | 锁残留 |

致命点：所有 SAdd均无 TTL，仅靠显式 SRem 释放。一旦流程在 SRem 之前崩溃/提前返回/异常分支，锁永久残留，唯一回收途径是人工面板 ClearLocation 或 recover 流程——而这两者又依赖订单状态正确写入，若中间环节崩溃同样覆盖不到。  
机器 168 某次取餐流程异常退出（最可能是 PickOrderDetail:2143 同步订单状态失败），位置 2 的锁就此永久残留。  

### 3.2 缺陷二：仓位选择算法对"目标位被锁"零容错（无法回退）

算法目标：选使用次数最少的空闲位置（负载均衡）。  
执行过程还原（getDetailSpaceByCupCount，DetailSpace.go:607-734）：  
G1 机型非奶油 SKU，候选位置经 dealOrderSpace:582-588 过滤，从 [4,3,1,2] 变为 [1,2]（仅保留位置编号 <3 的）。此时 cupCount=[1,2]，detailNum=1（单杯），dLen=1。  

| 内层 for i=1..4 | loc = cupCount[len-detailNum] | 判定 | cupCount 变化 |
| --- | --- | --- | --- |
| i=1 | cupCount[1]=2 | loc(2)≠i(1) → continue | [1,2] |
| i=2 | cupCount[1]=2 | loc==i ✅ → 查锁 → locked → 行668 截断 | → [1] |
| i=3 | cupCount[0]=1 | loc(1)≠i(3) → continue | [1] |
| i=4 | cupCount[0]=1 | loc(1)≠i(4) → continue | [1] |

循环结束 detailSpace=[]。两个 goto loop 重试点都不触发：  
- 行 703：continueFlag==1 —— 单杯 continueFlag=0，不触发；
- 行 712：len(cupCount)=1 > dLen=1 → false（被行 668 的截断破坏），不触发。
最终行 719 len(detailSpace)!=dLen → 返回 nil,nil → "创建仓位数据不匹配" → 失败。  
缺陷本质：  
1. 行 668cupCount = cupCount[0 : len(cupCount)-detailNum]（截断末尾）而非"移除被锁的那个位置元素"。截断发生在遍历中途（i=2），把目标从位置 2 改成了位置 1，但 for 循环已越过 i=1，再也回不去选空闲的位置 1。
2. 截断同时把 cupCount 长度从 2 砍到 1，破坏了行 712 的重试条件，自愈失败。
3. 潜伏 bug：startIndex（行 620）在 loop: 标签（行 622）之前初始化，goto loop 不会重置它。本故障因 goto 从未触发而无影响，但多杯 + 重试场景会状态污染。
位置 1 物理空闲（ros[1]=0、order 空、未锁），却因算法回退路径断裂而永远选不到。  

### 3.3 因果链


```javascript
取餐流程异常退出（PickOrderDetail 提前 return）
        │
        ▼
位置 2 的 SAdd 无配对 SRem + 无 TTL  →  锁永久残留（缺陷一）
        │
        ▼
下次仓位分配：算法选中使用次数最少的位置 2（正确）
        │
        ▼
位置 2 被锁 → 行 668 截断 cupCount（错误处理）
        │
        ▼
候选列表中途变性 + 循环已过 i=1 + 重试条件被破坏（缺陷二）
        │
        ▼
明明位置 1/3/4 空闲，却返回空 → 整单失败
        │
        ▼
worker 每秒重试，但锁不过期、算法逻辑确定 → 永久卡死
```

## 四、验证证据

本根因已经三路正交验证，结论一致：  

| 验证路 | 方法 | 结论 |
| --- | --- | --- |
| 独立复现 | 独立 Go 程序逐行复制 getDetailSpaceByCupCount，用日志真实输入 go run | 输出与生产日志逐条吻合，缺陷复现成功（runtime truth） |
| Oracle 对抗审查 | 独立 Read 源码，逐假设证伪 | 核心结论成立，无根本性错误；算法缺陷证据确凿；幽灵锁存在确认 |
| 路径核查 | 全量审计 SAdd/SRem 配对、常量定义、函数签名 | 找到确凿的锁残留路径（PickOrderDetail:2083/2088/2143） |

复现脚本：/var/folders/4w/xcv9wk5s29b7l3fztxfxpbh40000gn/T/opencode/repro_detailspace/main.go（独立 module repro，零依赖，go run main.go 即可复现）。  

## 五、修复方案

按"风险 / 收益 / 紧迫度"分三层。三层都应做——只修一层，另一层仍是定时炸弹。  

### 方案 A：紧急止血（运维，立即执行）

目标：解锁当前卡死的订单 971759，恢复机器 168 出杯。  

```javascript
SREM brain:shop_machine_space_to_lock:215_168 2
```

执行后，订单的下一次重试即可正常分配位置 2（或位置 1）成功。  
配套排查：检索机器 168 最近 PickOrderDetail 的日志，定位是 2083 / 2088 / 2143 哪一处触发了锁残留（重点看 2143 的 syncOrderPickStatus 失败记录），找到上游异常根因。  
⚠️ 止血不解决根本问题——下次任何取餐异常都会再次残留锁。  

### 方案 B：算法容错根治（核心，必修）

目标：即使目标位被锁，算法也能回退到次优空闲位，不再因单把锁卡死整单。  
文件：internal/server/DetailSpace.go  
⚠️ 历史回归风险（必须论证）：git blame 确认，行 668 的截断、行 666 的 checkCupSpaceLock 调用、行 630-634 的提前返回，全部来自 commit 3eb93bb08f "备用取餐和新杯并发bug"（2025-08-01, PR !143）——该 commit 之前 getDetailSpaceByCupCount 完全不检查仓位锁。本方案的改动必须论证不会回退该 commit 修复的并发场景，且修复 PR 必须补该场景的回归测试。  
⚠️ 关键认知（review 实测）：仅做"截断→移除被锁元素"无法修复——review 用独立 Go 程序实测证明，改动后算法仍返回 nil。根因：行 668 移除位置 2 后 cupCount=[1]，但行 713 cupCount = cupCount[:len-1] 的重试剥离会把仅剩的 [1] 剥成 []，重试命中行 630 返回空。正确的修复必须同时绕过行 713 的破坏性剥离（见改动 1）。  
改动 1 — 命中锁后：移除被锁元素 + 直接 goto loop（绕过行 713 剥离）  

```go
// 修改前（行 666-672）：
check := checkCupSpaceLock(requestID, shopID, machineID, i)
if check != nil {
    cupCount = cupCount[0 : len(cupCount)-detailNum]   // ❌ 截断末尾
    logger.Errorw(requestID, serverName, fn, "位置已锁定", i, cupCount, ...)
    continue
}

// 修改后：
check := checkCupSpaceLock(requestID, shopID, machineID, i)
if check != nil {
    for idx, v := range cupCount {                       // ✅ 精确移除被锁的位置 i
        if v == i {
            cupCount = append(cupCount[:idx], cupCount[idx+1:]...)
            break
        }
    }
    logger.Errorw(requestID, serverName, fn, "位置已锁定", i, cupCount, ...)
    goto loop                                            // ✅ 带剩余 cupCount 直接重试，绕过行 712-713
}
```

goto loop 会重置 startIndex（改动 3）、detailNum、nRosSpace；重试时 cupCount=[1] 使 loc=cupCount[0]=1==i=1 命中位置 1。review 用此变体实测得到 detailSpace=[1]，证明有效。  
改动 2 — 补越界保护（防 panic）  
行 658 loc := cupCount[len(cupCount)-detailNum]，当 len(cupCount) < detailNum 时负下标 panic。行 630 只保护 len==0，挡不住 0 < len < detailNum。放宽保护：  

```go
// 修改前（行 630）：
if len(cupCount) == 0 { return nil, nil }

// 修改后：
if len(cupCount) < dLen { return nil, nil }
```

改动 3 — 修复 startIndex 不随 goto 重置的潜伏 bug（编译合法）  
行 620 startIndex := 0 在 loop: 标签（行 622）之前，goto 不重置它（多杯重试场景状态污染）。移到 loop 标签之后：  

```go
// 修改前：
startIndex := 0          // 行 620（loop 标签之前）
dLen := len(orderDetails)
loop:
    detailSpace := make([]models.OrderDetailSpace, 0)

// 修改后：
dLen := len(orderDetails)
loop:
    startIndex := 0      // 移到 loop 之后，每次 goto 重置
    detailSpace := make([]models.OrderDetailSpace, 0)
```

✅ Go 编译合法性：goto loop 是向后跳，Go 的 "goto jumps over declaration" 只针对向前跳越过声明。现有源码本身就在 loop: 后有 5 处 := 配合向后 goto 且正常编译，review 已用 go build 验证本改动通过。  
改动 4（推荐）— 内层循环改为遍历 cupCount  
内层 for i:=1; i<=lLen; i++ 按物理位置遍历再用 loc 反查是缺陷根源。改为直接遍历已排序的 cupCount，锁定一个跳一个，天然支持回退。但涉及 continueFlag 连续放置逻辑，必须配 commit 3eb93bb08f 场景回归测试。  
撤销原"改动 2（行 712 >→>=）"：review 实测证明它不仅无效反而有害（加速行 713 错误剥离），应删除。本方案改为改动 1 的"命中锁直接 goto loop 绕过行 713"。  

### 方案 C：锁机制根治（防止残留复发）

目标：从源头杜绝"取餐流程异常退出 → 锁永久残留"。  
改动 1（推荐，低风险）— 给取餐目标位锁加 TTL 兜底  
在 checkBackupCup 的三处 SAdd 之后，对 sKey 补一个 EXPIRE，确保即使 SRem 漏执行，锁也会自动过期（建议 TTL 略大于一次取餐流程的最长耗时，如 10 分钟）：  

```go
// robot_system.go，在每处 SAdd 之后追加（行 1741/1749/1760 各加一次）：
drives.Rds.Expire(context.Background(), sKey, 10*time.Minute)
```

⚠️ TTL 刷新局限（review 发现）：Expire 作用于整个 Set（每台机器一个 key），Redis Set 不支持成员级 TTL。每次 SAdd 后 Expire 会刷新整 set TTL——只要 10 分钟内有任意新取餐，残留锁就被反复续命。因此"TTL 兜底"只在机器完全静默 10min时生效，正常营业机器几乎不会静默。成员级正确性仍依赖改动 2（补 SRem）的完备性。如需真正的成员级过期，应为每个被锁位置用独立 key（...:lock:<space>）各自设 TTL。  
改动 2（治本，中等风险）— 补全 PickOrderDetail 提前 return 的 SRem  
在 PickOrderDetail 的行 2083 / 2088 / 2143 三处提前 return 之前，根据已锁定的位置补 SRem，确保任何退出路径都清理锁。需结合当时已 SAdd 的位置集合逐分支处理（注意 2143 当前是"设计如此等待人工"，加 SRem 会改变语义，需产品确认是否改为自动释放）。  
方案 C-改动 2 涉及行为语义变更（尤其行 2143），属产品级决策，建议先落地方案 C-改动 1（TTL 兜底）作为安全网，再评估改动 2。  

## 六、验证方法

### 6.1 算法修复（方案 B）的验证

单元测试（推荐新增 internal/server/DetailSpace_test.go）：  
将 getDetailSpaceByCupCount 的锁检查抽为可注入依赖（或对 checkCupSpaceLock 做 redis mock），覆盖以下场景：  

| 场景 | 输入 | 期望结果 |
| --- | --- | --- |
| 正常单杯 | cupCount=[1,2]，无锁 | 分配位置 2（最少使用）成功 |
| 目标位被锁（本次故障） | cupCount=[1,2]，位置 2 锁定 | 回退到位置 1 成功（修复前返回空） |
| 全部候选被锁 | cupCount=[1,2]，位置 1、2 均锁 | 返回空（合理失败，不应卡死算法） |
| 多杯连续放置 | cupCount=[1,2,3,4]，2 杯 | 分配连续位置成功 |
| 多杯 + 部分锁定 | 2 杯，其中一目标位锁 | 回退到次优连续位置 |

回归验证：复现脚本（repro_detailspace/main.go）在修复后应输出"分配位置 1 成功"而非"创建仓位数据不匹配"。  
⚠️ 实施前必须实测：review 发现原"改动 1+2+3"组合实测仍返回 nil（详见方案 B 关键认知）。落地前必须先用脚本验证修正后的算法（改动 1 的 goto loop 变体）确实输出 detailSpace=[1]，再写单测。  

### 6.2 锁 TTL（方案 C-改动 1）的验证

- 构造一次取餐流程，SAdd 后检查 TTL(sKey) 约为 10 分钟；
- 模拟 PickOrderDetail 提前 return，确认 10 分钟后锁自动消失，SIsMember 返回 false。

### 6.3 端到端 QA

- 在测试环境模拟"位置 2 残留锁 + 下单"，确认订单能回退到位置 1 并完成出杯；
- 模拟"取餐同步失败（2143）"，确认 TTL 到期后锁释放，机器恢复可用。

## 七、历史先例与已知问题

### 7.1 锁残留是本代码库反复出现的模式

- 2025-02-05 commit a081ad8f6 / d75f2ad50（PR !2）"修复：locked_value 在清理订单时未释放的BUG"：之前已修过同类问题（清理订单时 MaterialLockedValue 原料锁未释放）。本次是仓位锁（SHOP_MACHINE_TO_SPACE_LOCK_KEY）的同类问题。建议系统性审视所有 *LOCK_KEY* 的释放完整性，避免打地鼠式修复。
- internal/service/ros/robot_system.go:265 有直接相关的锁难题 TODO：// TODO 加制作锁 怎么解决锁的问题 ？？？ 锁可能已经超时释放 也可能别的订单已占据锁——开发者早就意识到锁超时/锁被占是未解难题，佐证本故障紧迫性。

### 7.2 平行仓位分配路径（本次未涉及，但存在同类风险）

DealFallSpaceNew（行 487）判断 (HL100 && Free<=1) || Free==0 走 getDetailBackupSpace（备用杯托，DetailSpace.go:383-439），否则走 getDetailSpaceByCupCount（常规）。本次故障机器是 G1 桌面型、Free=4，走常规路径（本文档分析对象）。但 getDetailBackupSpace 不调用 checkCupSpaceLock、不检查任何锁——若未来备用杯托路径出现锁残留（如 robot_system.go:1749 的 toCup+changeCup 双锁），同样会选到被锁位置。建议一并审视。  

### 7.3 影响范围补充

worker/cmd/order_make.go 的 DealOrder 用 goto retry 在同一协程内死循环重试（行 337-347），期间不再从 channel 取下一条消息。因此本故障不仅是"单订单卡死"，而是该机器的订单处理协程被卡死，后续订单全部排队等待——单台机器出杯完全停摆的机制在此。  

## 八、运维落地（回滚 / 监控 / 灰度）

### 8.1 回滚

- 止血回滚：方案 A 的 SREM 是单命令，无需回滚（若误清，重新走取餐流程会重新加锁）。
- 算法改动回滚：方案 B 改动 1/2/3 均为局部代码改动，回滚 = git revert + 重新部署。
- 锁 TTL 回滚：方案 C 改动 1 的 Expire 是纯增量，删除即可，无数据影响。

### 8.2 监控告警

- 锁残留监控：定时扫描 brain:shop_machine_space_to_lock:* 集合 size，超过机器仓位总数（如 >4）告警——正常集合 size ≤ 仓位数。
- 仓位分配失败率：监控 saas.go:136 "仓位分配失败" 日志频率，超阈值告警。
- worker 卡死监控：监控 DealOrder 单消息处理耗时，超过 N 秒告警（goto retry 死循环征兆）。

### 8.3 灰度

- 方案 B（DetailSpace.go）影响全门店仓位分配，必须灰度：先 1 台 G1 机器试点 24h，观察分配成功率 + 无 panic，再逐步放量。
- 方案 C 改动 1（TTL）影响全门店取餐锁，同上灰度。
- 灰度期重点观察：commit 3eb93bb08f 的备用取餐并发场景是否回归。

## 九、附录

### 9.1 关键代码位置


| 文件 | 行号 | 说明 |
| --- | --- | --- |
| internal/server/DetailSpace.go | 452-541 | DealFallSpaceNew（仓位分配入口） |
| internal/server/DetailSpace.go | 543-605 | dealOrderSpace（含 G1 过滤分支 562-590） |
| internal/server/DetailSpace.go | 607-734 | getDetailSpaceByCupCount（缺陷二核心） |
| internal/server/DetailSpace.go | 668 | ❌ 锁定截断（缺陷二） |
| internal/server/DetailSpace.go | 712 | ❌ 重试条件过早失效 |
| internal/server/DetailSpace.go | 620 vs 622 | ⚠️ startIndex 不随 goto 重置（潜伏 bug） |
| internal/server/DetailSpace.go | 736-753 | checkCupSpaceLock（server 版，日志来源） |
| internal/server/DetailSpace.go | 755-787 | getCupStandCountInfo（杯次数排序） |
| internal/service/ros/robot_system.go | 1703-1764 | checkBackupCup（锁写入 SAdd） |
| internal/service/ros/robot_system.go | 2067+ | PickOrderDetail（缺陷一，提前 return 漏 SRem） |
| internal/service/ros/robot_system.go | 2083/2088/2143 | ❌ 三处提前 return 致锁残留 |
| internal/service/ros/robot_system.go | 2101/2150/2151/2170/2171 | SRem 清理点 |
| internal/service/panel/tool.go | 351 | 人工 ClearLocation 的 SRem |
| internal/service/ros/recover.go | 365 | recover 流程的 SRem |

### 9.2 关键常量


| 常量 | 值 | 位置 |
| --- | --- | --- |
| enum.MachineModelIDDeskB | 6 | gitee.com/ihaoin_enterprise/hooloo-protocol@v1.1.9/enum/machine.go:265 |
| NoRobotForG1Flag | 默认 false（brain 服务保持 false） | internal/server/OpRobotic.go:84 |
| SHOP_MACHINE_TO_SPACE_LOCK_KEY | brain:shop_machine_space_to_lock:%d_%d | enum 包 |

### 9.3 复现脚本

路径：/var/folders/4w/xcv9wk5s29b7l3fztxfxpbh40000gn/T/opencode/repro_detailspace/main.go  
独立 Go module（repro），零外部依赖（不连 redis/mysql），逐行复制 getDetailSpaceByCupCount 算法语义，checkLock 注入为函数变量（位置 2 返回 locked）。go run main.go 输出与生产日志逐条对应，证明算法缺陷。  

### 9.4 修复优先级建议


| 优先级 | 方案 | 理由 |
| --- | --- | --- |
| P0（立即） | 方案 A 止血 | 恢复卡死订单 |
| P1（本周） | 方案 B 算法根治 | 消除"单锁卡死整机"的灾难性后果 |
| P1（本周） | 方案 C-改动 1 TTL 兜底 | 防止锁残留复发，低风险 |
| P2（后续） | 方案 C-改动 2 补 SRem | 治本，但涉及语义变更需产品确认 |
| P3（重构） | 方案 B-改动 4 遍历 cupCount | 彻底消除索引映射脆弱性 |

文档基于生产日志（请求 ID c3504b08-d6b7-4e00-992b-312442ee6b86）与源码（版本 v1.1.9-90-g648e81440）分析，经独立复现 / Oracle 对抗审查 / 路径核查三路交叉验证。  
📌 来源：[Notion 原文](https://app.notion.com/p/38d7618faf3a8170a64ce75247d7332e)  
