# IoT系统设计方案：云端+工控机分层架构

> Source: WPS 知识库 / 领域/技术笔记/IoT系统设计方案：云端+工控机分层架构
> file_id: zgLeafdCurM9aPssWxhAxxop8FZgqbvMx

---

IoT系统设计方案：云端+工控机分层架构  
## 一、核心设计原则

一句话原则：云端管"交易"，工控机管"执行"，两者解耦，断网可运行。  
具体含义：  
- 云端：商品/订单/支付/营销/数据，不依赖任何单台机器在线
- 工控机：制作/取餐/硬件控制/本地状态，断网时独立运行
- 通信：事件驱动，最终一致，消息队列兜底
## 二、职责划分总览


```plaintext
┌─────────────────────────────────────────────────────────────────────┐
│                           云端 (Cloud)                              │
│                                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ 商品中心  │ │ 订单中心  │ │ 支付中心  │ │ 营销中心  │ │ 数据中心  │ │
│  │ Catalog  │ │ Order    │ │ Payment  │ │ Marketing│ │ Analytics│ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ 门店中心  │ │ 用户中心  │ │ 设备中心  │ │ 消息中心  │ │ 对账中心  │ │
│  │ Store    │ │ Member   │ │ Device   │ │ Notify   │ │ Finance  │ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
│                                                                     │
│  职责：所有跨门店共享的数据 + 所有涉及金钱的操作                       │
│  特点：水平扩展、高可用、不依赖单台工控机                              │
└────────────────────────────────┬────────────────────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │   消息总线 / API 网关     │
                    │   (MQTT + HTTP + gRPC)  │
                    └────────────┬────────────┘
                                 │
┌────────────────────────────────┴────────────────────────────────────┐
│                      工控机 (Edge / IPC)                            │
│                                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ 本地订单  │ │ 制作调度  │ │ 硬件控制  │ │ 取餐管理  │ │ 本地缓存  │ │
│  │ 缓存     │ │ Producer  │ │ Hardware │ │ Pickup   │ │ SQLite   │ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                          │
│  │ 状态上报  │ │ 离线降级  │ │ 日志收集  │                          │
│  │ Reporter  │ │ Offline  │ │ Logger   │                          │
│  └──────────┘ └──────────┘ └──────────┘                          │
│                                                                     │
│  职责：硬件实时控制 + 制作编排 + 本地状态管理                         │
│  特点：实时性保障、断网自治、故障自恢复                               │
└─────────────────────────────────────────────────────────────────────┘
```

### 云端职责（跨门店、跨机器共享）


| 模块 | 职责 | 为什么在云端 |
| --- | --- | --- |
| 商品中心 | 商品定义、SKU、图片、分类、上下架 | 所有门店共享同一套商品 |
| 价格中心 | 门店定价、促销价格、会员价、生效时间 | 多渠道需统一价格 |
| 订单中心 | 订单创建、状态查询、历史记录 | 用户多终端查看、客服查询 |
| 支付中心 | 微信/支付宝/POS支付、退款、对账 | 涉及金钱，必须集中管控 |
| 营销中心 | 优惠券、满减、会员积分、活动 | 跨门店通用 |
| 用户中心 | 注册、登录、会员等级、积分 | 用户跨门店使用 |
| 门店中心 | 门店信息、营业时间、机器绑定 | 运营管理需要全局视图 |
| 设备中心 | 设备注册、远程配置、OTA升级、监控大盘 | 远程运维 |
| 消息中心 | 微信模板消息、短信、APP推送 | 依赖外部服务 |
| 数据中心 | 销售报表、用户分析、商品分析 | 跨门店分析 |
| 对账中心 | 收款对账、退款对账、分账 | 财务合规 |

### 工控机职责（单店、单机、实时）


| 模块 | 职责 | 为什么在工控机 |
| --- | --- | --- |
| 硬件控制 | PLC/继电器/门锁/传感器/咖啡机/制冰机 | 实时性（毫秒级），断网必须可用 |
| 制作调度 | 制作队列、机器分配、制作步骤编排 | 本地决策，不依赖网络 |
| 取餐管理 | 取餐码校验、仓位分配、传感器检测 | 硬件交互实时性 |
| 本地订单缓存 | 接收云端订单、缓存到本地SQLite | 断网时仍可制作已下单的饮品 |
| 本地状态管理 | 制作状态、硬件状态、传感器状态 | 实时变化，高频更新 |
| 状态上报 | 定时/事件驱动上报状态到云端 | 云端需要可观测 |
| 离线降级 | 断网时本地继续运行，恢复后同步 | 保障营业不中断 |
| 扫码核销 | 扫码枪读取取餐码、本地校验 | 本地快速响应 |
| 本地打印 | 小票打印、杯贴打印 | USB直连 |
| 安全看护 | 水位监测、温度异常、设备故障本地告警 | 安全相关，不能依赖网络 |

## 三、云端详细设计

### 3.1 商品中心（Catalog Service）


```plaintext
┌──────────────────────────────────────────────┐
│              商品中心 (Catalog)                │
│                                              │
│  Product (商品)                               │
│  ├── id, name, category, images             │
│  ├── goods_type (现制/预包装)                  │
│  ├── base_recipe (基础配方，推送到工控机)       │
│  └── status (上架/下架)                       │
│                                              │
│  SKU (规格)                                   │
│  ├── id, product_id                          │
│  ├── spec_rules (大杯/燕麦奶/少冰)             │
│  └── barcode (条形码)                         │
│                                              │
│  StorePrice (门店价格)                         │
│  ├── store_id, sku_id                        │
│  ├── regular_price (日常价)                   │
│  ├── member_price (会员价)                    │
│  ├── promotion_price (促销价)                 │
│  └── effective_from / effective_to (生效区间) │
└──────────────────────────────────────────────┘
```

关键设计：商品变更推送  

```go
// 云端：商品变更后发布事件
func (s *CatalogService) UpdateProduct(ctx context.Context, product Product) error {
    db.Save(&product)
    // 发布商品变更事件 → MQTT
    mqtt.Publish("store/{storeID}/catalog/updated", ProductUpdateEvent{
        ProductID: product.ID,
        Version:   product.Version,  // 乐观锁版本号
        Timestamp: time.Now(),
    })
    return nil
}

// 工控机：收到通知后拉取最新商品
func (e *EdgeCatalog) OnProductUpdated(storeID string, event ProductUpdateEvent) {
    if event.Version <= e.localVersion {
        return  // 本地已是最新的
    }
    // 从云端拉取完整商品数据
    products := cloudAPI.FetchProducts(storeID, event.Version)
    // 写入本地 SQLite
    e.localDB.SaveProducts(products)
    e.localVersion = event.Version
}
```

### 3.2 订单中心（Order Service）

数据模型：  
- Order：id, order_no, store_id, user_id, channel(小程序/H5/POS/外卖), state, total_amount, paid_amount, discount, coupon_id, payment_id, pick_code, created_at, paid_at, completed_at, version(乐观锁)
- OrderItem：id, order_id, snapshot(商品快照JSON), sku_snapshot(规格快照JSON), unit_price, quantity, subtotal, state(独立状态，支持一单多杯)
核心流程：下单 → 支付 → 推送到工控机  

```plaintext
用户下单                          云端                          工控机
  │                               │                              │
  │── 创建订单 ──→                 │                              │
  │                 ←── 订单创建成功 (state=pending)               │
  │── 微信支付 ──→ 微信            │                              │
  │                 ←── 支付成功    │                              │
  │                  微信回调 ──→  │                              │
  │                              │── 幂等处理                     │
  │                              │── state=paid                  │
  │                              │── 分配取餐码                   │
  │                              │── 预扣库存确认                 │
  │                              │── 写 Outbox(订单下发)          │
  │                              │── MQTT 发布 ──→ order/new ──→ │
  │←── 推送 "支付成功,取餐码A42"   │                              │── 写入本地队列
  │                              │   ←── 状态上报(making) ──────│
  │←── 推送 "制作中"              │                              │
  │                              │   ←── 状态上报(ready) ──────│
  │←── 推送 "请到B2取餐"          │                              │
  │── 取走饮品 ────────────────────────────────────────────────→│
  │                              │   ←── 状态上报(done) ───────│
  │←── 推送 "取餐完成"            │                              │
```

### 3.3 设备中心（Device Service）


```go
// 工控机心跳上报（每30秒）
func (e *Edge) heartbeatLoop() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        if mqtt.IsConnected() {
            mqtt.Publish("device/{deviceID}/heartbeat", Heartbeat{
                DeviceID:        e.deviceID,
                Timestamp:       time.Now().Unix(),
                FirmwareVer:     e.firmwareVersion,
                HardwareStatus:  e.getHardwareStatus(),
                QueueLength:     e.productionQueue.Len(),
                CPU:             getCPUUsage(),
                Memory:          getMemUsage(),
                DiskSpace:       getDiskSpace(),
            })
        }
    }
}

// 云端告警规则
rules := []AlertRule{
    {Condition: "heartbeat > 90s ago",   Severity: "warning", Action: "mark_offline"},
    {Condition: "heartbeat > 5m ago",    Severity: "error",   Action: "notify_ops"},
    {Condition: "plc_error_count > 5/min", Severity: "error", Action: "notify_ops"},
    {Condition: "disk_space < 10%",      Severity: "warning", Action: "notify_ops"},
    {Condition: "queue_length > 20",     Severity: "info",    Action: "auto_scale"},
}
```

## 四、工控机详细设计

### 4.1 整体架构


```plaintext
┌─────────────────────────────────────────────────────────────────┐
│                      工控机软件架构                               │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    应用编排层                              │   │
│  │  订单接收 → 制作调度 → 硬件控制 → 状态上报 → 取餐管理     │   │
│  └────────────────────────┬────────────────────────────────┘   │
│                           │                                     │
│  ┌────────────────────────┼────────────────────────────────┐   │
│  │                    核心服务层                              │   │
│  │  ┌──────────┐  ┌───────┴──────┐  ┌──────────┐          │   │
│  │  │ MQTT通信  │  │ 本地状态管理  │  │ 硬件适配层 │          │   │
│  │  │ 云端通道  │  │ (内存+SQLite) │  │           │          │   │
│  │  └──────┬───┘  └───────┬──────┘  └─────┬────┘          │   │
│  │  ┌──────┴───┐  ┌───────┴──────┐  ┌─────┴────┐          │   │
│  │  │ 离线队列  │  │ 定时任务调度  │  │ 安全监控  │          │   │
│  │  │ (断网缓存)│  │ (超时/巡检)  │  │ (水位/温)│          │   │
│  │  └──────────┘  └──────────────┘  └──────────┘          │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│  ┌────────────────────────┼────────────────────────────────┐   │
│  │                    硬件驱动层                              │   │
│  │  ┌─────────┐  ┌────────┴───────┐  ┌──────────┐         │   │
│  │  │ Modbus  │  │ Serial(串口)    │  │ USB      │         │   │
│  │  │ TCP/RTU │  │ RS232/RS485    │  │ 打印机/印花│         │   │
│  │  └────┬────┘  └───────┬────────┘  └────┬─────┘         │   │
│  │  PLC/继电器/门锁  扫码枪/称重/POS   热敏打印/印花机      │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 本地状态管理（核心设计）

工控机内置 SQLite 数据库，保障断网时的业务连续性：  

```plaintext
本地存储的数据表：

1. LocalOrder（本地订单缓存，从云端同步）
   字段：order_no, cloud_order_id, pick_code, state, items(JSON), 
         received_at, completed_at, synced_at
   state: local_pending / local_making / local_ready / local_done / local_synced

2. LocalProduct（商品本地缓存，从云端同步）
   字段：name, price, images, recipe_data(制作配方), sync_version

3. HardwareStateLog（硬件状态历史，定时上报）
   字段：device_type, state_data(JSON), timestamp, uploaded

4. LocalOperationLog（操作日志）
   字段：order_no, action(receive/make_start/make_done/pickup), 
         timestamp, uploaded
```

### 4.3 订单接收与制作编排


```plaintext
                    工控机订单处理流程

    云端 MQTT 推送
         │
         ▼
┌─────────────────┐
│  消息接收层       │  订阅: store/{id}/order/new
│  幂等检查        │  检查 order_no 是否已存在
└────────┬────────┘
         ▼
┌─────────────────┐
│  订单入库        │  写入 LocalOrder (state=local_pending)
└────────┬────────┘
         ▼
┌─────────────────────────────────────────────────────┐
│  制作调度器 (Production Scheduler)                    │
│                                                     │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐            │
│  │ 队列位置1 │  │ 队列位置2 │  │ 队列位置3 │           │
│  │ 制作中   │  │ 排队     │  │ 排队     │            │
│  └─────────┘  └─────────┘  └─────────┘            │
│                                                     │
│  调度策略：                                          │
│  - 单机器：FIFO 先进先出                             │
│  - 多机器：按商品类型路由（咖啡→咖啡机, 冰品→冰机）    │
│  - 优先级：VIP优先 / 超时单优先                      │
└────────┬────────────────────────────────────────────┘
         ▼
┌─────────────────────────────────────────────────────┐
│  制作执行器 (Production Worker)                       │
│                                                     │
│  Step 1: 落杯 → cup_device.FallCup() → 串口          │
│  Step 2: 制作 → 按配方执行（HTTP/Modbus/Serial）      │
│  Step 3: 取餐位分配 → pickup.AssignStation()         │
│  Step 4: PLC控制 → 杯托↓→舱门开→等待取走→舱门关→杯托↑ │
│                                                     │
│  全程状态更新 → LocalOrder.State                     │
│  全程状态上报 → MQTT → 云端                          │
└─────────────────────────────────────────────────────┘
```

制作执行器代码：  

```go
func (w *ProductionWorker) Process(ctx context.Context, order *LocalOrder) error {
    traceID := order.OrderNo
    w.updateState(ctx, order, "making")
    w.reporter.Report(ctx, traceID, "making_started", nil)

    // Step 1: 落杯
    if err := w.hardware.FallCup(ctx); err != nil {
        return w.handleFailure(ctx, order, err, "fall_cup")
    }

    // Step 2: 按配方制作
    for _, item := range order.Items {
        recipe := item.Recipe  // 从本地缓存读取
        for _, step := range recipe.Steps {
            if err := w.executeStep(ctx, step); err != nil {
                return w.handleFailure(ctx, order, err, step.Name)
            }
        }
    }

    // Step 3: 分配取餐位
    station, err := w.pickup.AssignStation(ctx, order.OrderNo)
    if err != nil {
        return w.handleFailure(ctx, order, err, "assign_station")
    }

    // Step 4: 标记制作完成
    w.updateState(ctx, order, "ready")
    w.reporter.Report(ctx, traceID, "making_completed", map[string]any{
        "station": station.ID,
    })
    return nil
}

func (w *ProductionWorker) executeStep(ctx context.Context, step RecipeStep) error {
    switch step.Type {
    case StepTypeCoffee:
        return w.hardware.MakeCoffee(ctx, step.Params)         // HTTP
    case StepTypePLC:
        return w.hardware.PLCCommand(ctx, step.Risc, step.Param) // Modbus
    case StepTypeIce:
        return w.hardware.MakeIce(ctx, step.Params)             // Serial
    case StepTypeRelay:
        return w.hardware.RelayCommand(ctx, step.Addr, step.Value) // Modbus
    }
}
```

### 4.4 硬件适配层（Hardware Adapter）

统一8种协议为一个接口，按机型实现不同适配器：  

```go
// 硬件能力接口（Capability-Based）
type HardwareAdapter interface {
    // 制作类
    FallCup(ctx context.Context) error
    MakeCoffee(ctx context.Context, recipe Recipe) error
    MakeIce(ctx context.Context, params IceParams) error

    // PLC 控制
    PLCCommand(ctx context.Context, command string, param int) error

    // 取餐类
    AssignStation(ctx context.Context, orderNo string) (*Station, error)
    DetectPickup(ctx context.Context, stationID string) (<-chan bool, error)

    // 状态类
    GetStatus(ctx context.Context) (HardwareStatus, error)
    HealthCheck(ctx context.Context) (HealthReport, error)

    // 生命周期
    Connect() error
    Disconnect()
    IsConnected() bool
}

// HL100机型适配器
type HL100Adapter struct {
    plc       *ModbusClient       // Modbus TCP → 嵌入式板卡
    coffee    *HTTPClient         // HTTP → Docter 咖啡机
    cupDevice *SerialClient       // Serial → 落杯器
    printer   *USBClient          // USB → 热敏打印机
    breaker   *circuit.Breaker    // 熔断器
}

// DeskA机型适配器
type DeskAAdapter struct {
    plc        *ModbusClient      // Modbus TCP → 嵌入式板卡
    coffee     *HTTPClient        // HTTP → 咖啡机
    iceMachine *SerialClient      // Serial → 制冰机
    frontPanel *ModbusClient      // Modbus TCP → 前面板
}
```

每个硬件操作都有超时 + 熔断 + 重试：  

```go
func (a *HL100Adapter) PLCCommand(ctx context.Context, command string, param int) error {
    return a.breaker.Execute(func() error {
        ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
        defer cancel()

        var lastErr error
        for retry := 0; retry < 3; retry++ {
            err := a.plc.SendCommand(ctx, command, param)
            if err == nil {
                return nil
            }
            lastErr = err
            // 指数退避
            time.Sleep(time.Duration(math.Pow(2, float64(retry))) * 100 * time.Millisecond)
        }
        return fmt.Errorf("PLC command %s failed after 3 retries: %w", command, lastErr)
    })
}
```

### 4.5 离线降级策略（关键设计）


```plaintext
断网场景处理：

场景1：用户已下单支付，工控机断网
  → 云端订单进入"待下发"队列
  → 网络恢复后自动推送
  → 超时30分钟未推送 → 自动退款

场景2：工控机正在制作，突然断网
  → 本地继续制作（不依赖网络）
  → 状态变更写入本地 SQLite
  → 网络恢复后批量上报
  → 用户侧：WebSocket 断开 → 降级为小程序内查看订单状态

场景3：工控机断网期间POS下单
  → 本地创建订单 → 本地制作
  → 网络恢复后同步到云端

场景4：云端故障，工控机正常
  → 工控机独立运行（POS下单 → 制作 → 取餐）
  → 本地缓存所有订单和状态
  → 云端恢复后批量同步
```

离线状态机与重连代码：  

```go
type ConnectionState struct {
    State     string  // online / degraded / offline
    Since     time.Time
    RetryNext time.Time
}

func (e *Edge) onMQTTDisconnected() {
    e.connState.State = "offline"
    e.connState.Since = time.Now()
    // 1. 所有状态变更写入本地 unsynced 队列
    // 2. POS/点单屏切到"离线模式"
    // 3. 启动重连定时器
    go e.reconnectLoop()
}

func (e *Edge) reconnectLoop() {
    for {
        time.Sleep(10 * time.Second)
        if e.tryReconnect() {
            e.onReconnected()
            return
        }
    }
}

func (e *Edge) onReconnected() {
    e.connState.State = "online"
    // 1. 上报离线期间的状态变更
    e.syncPendingStates()
    // 2. 拉取离线期间的云端订单
    e.fetchPendingOrders()
    // 3. 拉取商品/价格变更
    e.syncCatalog()
    // 4. 上报设备状态
    e.reportFullStatus()
}
```

## 五、云端 ↔ 工控机通信设计

### 5.1 通信通道


```plaintext
┌──────────────┐                          ┌──────────────┐
│    云端       │                          │   工控机      │
│              │                          │              │
│  ┌────────┐  │    MQTT (双向通道)        │  ┌────────┐  │
│  │MQTT    │←─┼──────────────────────────┼─→│MQTT    │  │
│  │Broker  │  │  心跳/状态/订单/配置       │  │Client  │  │
│  └────────┘  │                          │  └────────┘  │
│              │                          │              │
│  ┌────────┐  │                          │  ┌────────┐  │
│  │HTTP    │←─┼──── HTTP (工控机拉取) ───┼─→│HTTP    │  │
│  │API     │  │  商品同步/全量上报/固件    │  │Client  │  │
│  └────────┘  │                          │  └────────┘  │
└──────────────┘                          └──────────────┘

MQTT → 实时、双向、低频消息
HTTP → 大数据量、请求-响应式同步
```

### 5.2 MQTT Topic 设计


```plaintext
Topic 命名规范：{方向}/{entity}/{action}

云端 → 工控机 (下行):
  store/{storeID}/order/new          ← 新订单（支付完成）
  store/{storeID}/order/cancel       ← 取消订单
  store/{storeID}/config/update      ← 配置变更
  store/{storeID}/catalog/updated    ← 商品/价格变更通知
  store/{storeID}/command/restart    ← 远程重启
  store/{storeID}/command/selfcheck  ← 远程自检
  device/{deviceID}/ota/notify       ← 固件升级通知

工控机 → 云端 (上行):
  device/{deviceID}/heartbeat        ← 心跳（30s）
  device/{deviceID}/status           ← 状态变更（制作开始/完成/取餐）
  device/{deviceID}/alert            ← 告警（缺水/故障/异常）
  device/{deviceID}/metrics          ← 指标（制作耗时/队列长度）
  device/{deviceID}/log              ← 关键日志

用户推送 (云端 → 用户端):
  user/{userID}/order/{orderNo}      ← 订单状态变更
```

### 5.3 MQTT QoS 策略


| Topic 方向 | QoS | 原因 |
| --- | --- | --- |
| 订单下发 | 1（至少一次） | 订单不能丢，必须送达 |
| 状态上报 | 0（至多一次） | 高频状态，丢一两条无所谓 |
| 心跳 | 0 | 30秒一次，丢了下次就来了 |
| 告警 | 1 | 告警不能丢 |
| 配置变更 | 1 + Retain | 必须送达且断线重连后能收到最新值 |

### 5.4 断网补偿：Outbox Pattern


```plaintext
云端发送订单到工控机：

  云端 Order Service
       │
       ▼
  ┌──────────────┐
  │  Outbox 表    │  ← 与订单状态变更同一事务写入
  │  (未投递消息)  │
  └──────┬───────┘
         │
  ┌──────▼───────┐
  │  MQ Sender   │  ← 后台 Worker 轮询，投递到 MQTT
  └──────┬───────┘
         │
         ▼ MQTT QoS 1
  ┌──────────────┐
  │  MQTT Broker  │
  └──────┬───────┘
         │
         ▼
  ┌──────────────┐
  │  工控机       │  ← 收到后写入本地 SQLite，回复 ACK
  │  本地队列     │
  └──────────────┘

  工控机收到 → 写本地 → MQTT ACK
  云端收到 ACK → 删除 Outbox 记录

  如果工控机离线：
    MQTT QoS 1 保证消息缓存
    工控机上线后收到所有缓存消息
    超时30分钟未收到 → 云端自动退款
```

### 5.5 关键流程：正常下单全流程时序图


```plaintext
用户          云端                    MQTT           工控机          硬件
 │              │                       │               │              │
 │── 创建订单 ──→│                       │               │              │
 │              │── 校验库存(预扣)       │               │              │
 │              │── 写 Order(state=pending)              │              │
 │←── 返回订单号 │                       │               │              │
 │── 微信支付 ──→│                       │               │              │
 │              │   ←─── 微信回调 ───────│               │              │
 │              │── 幂等检查             │               │              │
 │              │── state=paid          │               │              │
 │              │── 分配取餐码           │               │              │
 │              │── 写 Outbox           │               │              │
 │              │── MQTT 发布 ─────────→│── 投递 ──────→│              │
 │←── 推送"支付成功"                    │               │── 写本地队列 │
 │              │                       │←── ACK ───────│              │
 │              │←── 删除 Outbox         │               │── 调度制作    │
 │              │                       │               │── 落杯 ──────→│
 │              │                       │               │←── 落杯完成   │
 │              │                       │               │── 制作咖啡 ──→│
 │              │                       │               │←── 制作完成   │
 │              │                       │               │── 分配取餐位  │
 │              │                       │←── 上报(making)│              │
 │←── 推送"制作中"                      │               │              │
 │              │                       │               │── 杯托↓ ────→│
 │              │                       │               │── 舱门开 ───→│
 │              │                       │               │←── 等待取走   │
 │              │                       │←── 上报(ready)│              │
 │←── 推送"请到B2取餐"                  │               │              │
 │── 取走饮品 ────────────────────────────────────────→│←── 传感器检测 │
 │              │                       │               │── 舱门关      │
 │              │                       │               │── 杯托↑      │
 │              │                       │←── 上报(done) │              │
 │←── 推送"取餐完成"                    │               │              │
```

## 六、数据一致性保障

### 6.1 订单状态对齐策略

原则：云端是最终事实来源，工控机是临时缓存。  
状态冲突解决规则：  
1. 云端 state=paid, 工控机无此订单 → 云端重发 → 工控机接收并制作
2. 云端 state=paid, 工控机 state=done → 工控机上报 → 云端更新为 done
3. 云端 state=cancelled(退款), 工控机 state=making → 云端下发取消指令 → 工控机停止制作
4. 云端 state=paid(超时), 工控机 state=ready(已完成) → 工控机上报 done → 云端更新
定时对账（每小时）：云端遍历最近2小时内的 active 订单，与工控机上报的最新状态比对，不一致标记"待人工确认"。  

```go
// 云端定时对账任务
func (s *ReconcileService) Run(ctx context.Context) {
    devices := s.deviceRepo.GetOnlineDevices()
    for _, device := range devices {
        orders := s.orderRepo.GetActiveOrders(device.StoreID)
        for _, order := range orders {
            edgeState, err := s.edgeAPI.GetOrderState(device.ID, order.OrderNo)
            if err != nil { continue }
            if order.State != edgeState {
                s.alertConflict(order, edgeState)
                s.resolve(order, edgeState)
            }
        }
    }
}
```

## 七、安全性设计


```plaintext
┌──────────────────────────────────────────────┐
│              安全分层设计                      │
├──────────────────────────────────────────────┤
│  1. 通信安全                                  │
│     - MQTT: TLS 加密 + 设备证书双向认证        │
│     - HTTP: HTTPS + API Token                │
│     - 工控机内网: VLAN 隔离，不暴露公网        │
│                                              │
│  2. 设备身份                                  │
│     - 每台工控机有唯一 device_id + 证书        │
│     - MQTT 连接时验证证书                     │
│     - 未经注册的设备拒绝连接                  │
│                                              │
│  3. 指令安全                                  │
│     - 远程指令需要签名（HMAC）                │
│     - 关键操作（重启/退款）需要二次确认        │
│     - 指令有有效期（5分钟过期）                │
│                                              │
│  4. 数据安全                                  │
│     - 工控机本地 SQLite 加密                  │
│     - 敏感配置只在云端，不下发                │
│     - 工控机不存储用户隐私数据                │
│                                              │
│  5. 物理安全                                  │
│     - 工控机 USB 口禁用未授权设备              │
│     - SSH 只允许密钥认证                     │
│     - 远程调试通道默认关闭                    │
└──────────────────────────────────────────────┘
```

## 八、部署架构


```plaintext
┌──────────────────────────────────────────────────────────────┐
│                        云端部署                                │
│  Region: 阿里云 / AWS                                        │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ K8s 集群     │  │ MQTT Broker │  │ 对象存储     │         │
│  │             │  │ (EMQX)      │  │ (OSS/S3)    │         │
│  │ - Order Svc │  │ 集群部署     │  │ 商品图片     │         │
│  │ - Pay Svc   │  │ 持久会话     │  │ 固件包       │         │
│  │ - Catalog   │  │ 消息回溯     │  │ 日志归档     │         │
│  │ - Device    │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ PostgreSQL  │  │ Redis       │  │ Kafka       │         │
│  │ (主从)      │  │ (集群)      │  │ (事件流)    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐                                          │
│  │ 监控         │  Prometheus + Grafana + AlertManager     │
│  └─────────────┘                                          │
└────────────────────────────────┬─────────────────────────────┘
                                 │ HTTPS / MQTT / TLS
                                 │
┌────────────────────────────────┴─────────────────────────────┐
│                     门店网络                                   │
│  ┌──────────┐     ┌──────────────┐     ┌──────────────┐     │
│  │ 4G/宽带   │────→│ 门店路由器    │────→│ 工业交换机    │     │
│  │ (主+备)   │     │ (防火墙/NAT) │     │ (VLAN隔离)   │     │
│  └──────────┘     └──────────────┘     └──────┬───────┘     │
│                    ┌──────────────────────────┼─────────┐   │
│              ┌─────┴──────┐            ┌───────┴────┐   │   │
│              │ 工控机      │            │ 点单屏      │   │   │
│              │ (IPC)      │            │ (触摸屏)    │   │   │
│              └─────┬──────┘            └────────────┘   │   │
│         ┌──────────┼──────────────┐                    │   │
│    ┌────┴───┐ ┌────┴────┐ ┌──────┴───┐               │   │
│    │嵌入式   │ │ 咖啡机   │ │ 制冰机   │               │   │
│    │板卡     │ │(HTTP)   │ │(Serial) │               │   │
│    │(Modbus)│ │         │ │         │               │   │
│    └────┬───┘ └─────────┘ └─────────┘               │   │
│    ┌────┴───────────────────────┐                    │   │
│    │ PLC/继电器/门锁/传感器      │                    │   │
│    │ 杯托/舱门/水泵/水位        │                    │   │
│    │ (RS-485 总线)              │                    │   │
│    └────────────────────────────┘                    │   │
└───────────────────────────────────────────────────────────┘
```

## 九、总结对照表


| 维度 | 云端 | 工控机 |
| --- | --- | --- |
| 数据所有权 | 最终事实来源 | 本地缓存副本 |
| 商品/价格 | 定义 + 管理 | 同步缓存 + 使用 |
| 订单 | 创建 + 状态查询 | 接收 + 本地执行 |
| 支付 | 唯一处理方 | 不参与 |
| 库存 | 预扣 + 确认 | 本地库存查询 |
| 制作 | 不参与 | 唯一执行方 |
| 硬件控制 | 不参与 | 唯一控制方 |
| 取餐 | 推送通知 | 硬件交互 + 状态检测 |
| 断网行为 | 订单堆积→超时退款 | 独立运行→恢复同步 |
| 数据存储 | PostgreSQL + Redis | SQLite + 内存 |
| 通信协议 | HTTPS（面向用户） | MQTT + HTTP（面向云端） |
| 实时性要求 | 秒级 | 毫秒级 |
| 故障影响 | 全局影响 | 单店影响 |

核心原则：云端是"大脑"（决策、交易、数据），工控机是"小脑+四肢"（实时控制、反射动作、断网自治）。两者通过 MQTT 事件解耦，通过 Outbox + 本地 SQLite 保障最终一致性，通过对账任务兜底数据安全。  
## 十、设备影子（Device Shadow）

这是 IoT 系统的标配，AWS IoT 和 Azure IoT 都内置。当前文档的"状态上报"是单向推送，生产级必须有双向影子同步。  
### 10.1 概念


```plaintext
设备影子（Device Shadow）= 设备状态的缓存副本

  云端                     工控机
  ┌──────────────┐         ┌──────────────┐
  │ Desired State │         │ Reported State│
  │ (期望状态)    │ ←同步→  │ (实际状态)    │
  │              │         │              │
  │ mode=cleaning│         │ mode=idle    │  ← 有差异
  │ queue=[]     │         │ queue=[3杯]  │
  └──────────────┘         └──────────────┘

  两个状态字段：
  - desired:  云端希望设备变成什么（下发指令）
  - reported: 设备上报自己当前是什么（状态反馈）

  当 desired != reported → 设备执行变更 → 更新 reported → 两者一致 → 完成
```

### 10.2 使用场景

场景A：云端下发"清洁模式"  

```plaintext
1. 云端更新影子 desired.mode = "cleaning"
2. 工控机下次同步发现 desired.mode != reported.mode
3. 工控机执行清洁（冲洗咖啡机、清空废渣）
4. 工控机更新 reported.mode = "cleaning" → "idle"
5. 云端看到 reported == desired → 标记清洁完成
```

场景B：云端读取设备状态（设备离线）  

```plaintext
1. 用户/运维想查看设备状态
2. 设备离线 → 直接读影子的 reported（最后一次上报的状态）
3. 不需要等设备上线
```

场景C：设备恢复在线后同步  

```plaintext
1. 设备断网期间，云端更新了 desired（如配置变更）
2. 设备恢复 → 发现 desired != reported
3. 自动应用新配置 → 更新 reported
4. 不需要云端重新下发
```

### 10.3 实现设计


```go
// 影子数据结构（存储在云端 Redis + 数据库）
type DeviceShadow struct {
    DeviceID  string          `json:"device_id"`
    Version   int64           `json:"version"`    // 乐观锁
    State     ShadowState     `json:"state"`
    Metadata  ShadowMetadata  `json:"metadata"`
    Timestamp time.Time       `json:"timestamp"`
}

type ShadowState struct {
    Desired  map[string]any `json:"desired"`   // 云端期望
    Reported map[string]any `json:"reported"`  // 设备上报
}

type ShadowMetadata struct {
    Desired  map[string]time.Time `json:"desired"`
    Reported map[string]time.Time `json:"reported"`
}
```


```go
// 云端：更新 desired（下发指令）
func (s *ShadowService) UpdateDesired(ctx context.Context, deviceID string, desired map[string]any) error {
    shadow := s.getShadow(deviceID)
    shadow.State.Desired = merge(shadow.State.Desired, desired)
    shadow.Version++
    s.saveShadow(shadow)
    // MQTT 通知设备影子已更新
    mqtt.Publish("device/"+deviceID+"/shadow/delta", ShadowDelta{
        Version: shadow.Version,
        Changed: desired,
    })
    return nil
}

// 工控机：上报 reported + 检查 desired
func (e *Edge) syncShadow(ctx context.Context) {
    // 1. 上报当前状态
    reported := e.getHardwareStatus()
    resp := cloudAPI.UpdateShadowReported(e.deviceID, reported)

    // 2. 检查是否有待执行的 desired
    if resp.HasDelta {
        for key, desiredVal := range resp.Delta {
            currentVal := e.getLocalState(key)
            if desiredVal != currentVal {
                e.executeDesiredChange(key, desiredVal)
            }
        }
    }
}
```


```plaintext
影子同步频率：
  - 正常：每次状态变化时 + 每5分钟全量同步
  - 断网恢复：立即触发一次全量同步
  - 云端变更：通过 MQTT delta 消息即时通知
```

## 十一、OTA 固件升级完整流程

文档之前只提了"OTA固件升级"四个字，生产级必须有灰度、回滚、版本兼容。  
### 11.1 OTA 架构


```plaintext
┌──────────────────────────────────────────────────────────┐
│                      云端 OTA 服务                        │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ 版本管理      │  │ 灰度策略      │  │ 升级监控      │  │
│  │              │  │              │  │              │  │
│  │ 固件包存储    │  │ 分批发布规则  │  │ 成功率追踪    │  │
│  │ 版本兼容矩阵  │  │ 回滚触发条件  │  │ 异常自动暂停  │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────┬────────────────────────────────┘
                          │ MQTT + HTTPS
                          ▼
┌──────────────────────────────────────────────────────────┐
│                    工控机 OTA Agent                       │
│                                                          │
│  1. 收到升级通知 (device/{id}/ota/notify)                 │
│  2. 下载固件包 (HTTPS，支持断点续传)                       │
│  3. 校验签名 (SHA256 + RSA)                               │
│  4. 双分区写入 (A/B Partition)                            │
│  5. 切换启动分区 → 重启                                   │
│  6. 健康检查 → 成功确认 / 失败自动回滚                     │
└──────────────────────────────────────────────────────────┘
```

### 11.2 灰度发布策略


```plaintext
升级分 4 批，每批间隔 24 小时观察：

  第1批 (Day 1):   3 台机器 → 内部测试门店
      │ 观察 24h
      ▼
  第2批 (Day 2):   10% 机器 → 低流量门店
      │ 观察 24h
      ▼
  第3批 (Day 3):   50% 机器
      │ 观察 24h
      ▼
  第4批 (Day 4):   100% 全量

  自动暂停条件（任一触发）：
  - 升级失败率 > 5%
  - 升级后设备离线率 > 10%
  - 制作异常率比升级前高 20%
  - 人工暂停按钮

  自动回滚条件：
  - 升级失败率 > 20%
  - 连续 3 台设备升级后无法启动
```

### 11.3 A/B 分区与自动回滚


```plaintext
工控机存储分 A/B 两个系统分区：

  正常运行（运行A分区）：
  ┌──────────┐  ┌──────────┐
  │ A分区     │  │ B分区     │
  │ v1.2(运行)│  │ v1.1(备用)│
  └──────────┘  └──────────┘

  升级过程（写入B分区）：
  ┌──────────┐  ┌──────────┐
  │ A分区     │  │ B分区     │
  │ v1.2(运行)│  │ v1.3(写入)│  ← 下载新固件写入备用分区
  └──────────┘  └──────────┘

  重启切换（运行B分区）：
  ┌──────────┐  ┌──────────┐
  │ A分区     │  │ B分区     │
  │ v1.2(备用)│  │ v1.3(运行)│  ← Bootloader 切换到 B
  └──────────┘  └──────────┘

  健康检查（5分钟内）：
  - MQTT 能否连接？
  - 硬件能否正常通信？
  - 能否接受新订单？

  检查通过 → 确认升级，标记 B 为活跃分区
  检查失败 → Bootloader 自动切回 A（无需人工干预）
```

### 11.4 版本兼容矩阵


```go
// 云端维护版本兼容表
type FirmwareCompat struct {
    FirmwareVersion string   // 工控机固件版本
    MinCloudAPIVer  string   // 最低兼容的云端API版本
    MQTTProtocolVer int      // MQTT消息格式版本
    ConfigSchemaVer int      // 配置Schema版本
}

// 升级前检查
func (s *OTAService) CanUpgrade(deviceID, targetVer string) error {
    device := s.deviceRepo.Get(deviceID)
    compat := s.compatMatrix[targetVer]

    // 1. 当前版本不能跳太多代（v1.0 不能直接升 v3.0）
    if !isGradualUpgrade(device.FirmwareVersion, targetVer) {
        return ErrTooManyVersionSkip
    }

    // 2. 云端API版本是否兼容
    if s.cloudAPIVersion < compat.MinCloudAPIVer {
        return ErrCloudVersionTooLow
    }

    // 3. MQTT消息格式是否需要云端先升级
    if compat.MQTTProtocolVer > s.currentMQTTVer {
        return ErrMQTTProtocolIncompatible
    }

    return nil
}
```

## 十二、可观测性具体方案

文档之前只提了 Prometheus + Grafana 的名字，生产级需要完整的追踪、指标、日志、告警体系。  
### 12.1 三层可观测性


```plaintext
┌────────────────────────────────────────────────────────┐
│                    可观测性三层架构                      │
│                                                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐│
│  │ Metrics 指标  │  │ Logs 日志    │  │ Traces 追踪  ││
│  │ (Prometheus) │  │ (Loki/ELK)  │  │ (Jaeger/OTel)││
│  │              │  │              │  │              ││
│  │ 数值监控      │  │ 事件记录     │  │ 链路追踪     ││
│  │ 告警触发      │  │ 事后排查     │  │ 性能定位     ││
│  └──────────────┘  └──────────────┘  └──────────────┘│
│         │                │                │           │
│         └────────────────┼────────────────┘           │
│                          ▼                            │
│                 ┌──────────────┐                      │
│                 │ Grafana 大盘 │  统一可视化           │
│                 └──────────────┘                      │
└────────────────────────────────────────────────────────┘
```

### 12.2 分布式追踪（OpenTelemetry）


```plaintext
一个订单从创建到取餐的完整 Trace：

TraceID: a1b2c3d4e5f6

  Span 1: order.create          (云端)       0ms────15ms
    ├── Span 2: pay.notify      (云端)       15ms───120ms
    │     ├── Span 3: mqtt.publish           120ms──125ms
    │     └── Span 4: mqtt.deliver           125ms──130ms
    ├── Span 5: edge.receive   (工控机)      130ms──135ms
    │     ├── Span 6: hardware.fall_cup      135ms──155ms
    │     ├── Span 7: hardware.make_coffee   155ms──220ms
    │     └── Span 8: plc.cup_down           220ms──250ms
    └── Span 9: edge.report_done (工控机)     250ms──260ms

  每个Span包含：order_no, device_id, duration, status, error
  一个TraceID贯穿：用户端 → 云端 → MQTT → 工控机 → 硬件
```


```go
// 云端：创建订单时生成 TraceID
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderReq) (*Order, error) {
    ctx, span := tracer.Start(ctx, "order.create")
    defer span.End()
    span.SetAttributes(
        attribute.String("order_no", orderNo),
        attribute.String("store_id", storeID),
        attribute.String("user_id", userID),
    )
    // ... 业务逻辑
}

// 工控机：收到订单时继承 TraceID
func (e *Edge) onOrderReceived(orderMsg OrderMessage) {
    ctx := context.WithTraceID(orderMsg.TraceID)  // 继承云端TraceID
    ctx, span := tracer.Start(ctx, "edge.receive")
    defer span.End()
    e.productionWorker.Process(ctx, order)
}
```

### 12.3 关键指标定义（SLO）


```plaintext
订单SLO：
  - 订单成功率 > 99.9%（每月允许失败 < 0.1%）
  - 下单→制作完成 P95 < 180秒
  - 制作→取餐 P95 < 300秒
  - 支付回调延迟 P99 < 5秒

设备SLO：
  - 设备在线率 > 99.5%
  - 硬件通信成功率 > 99.8%
  - PLC响应时间 P99 < 500ms
  - OTA升级成功率 > 98%

通信SLO：
  - MQTT消息投递成功率 > 99.95%
  - 云端→工控机延迟 P95 < 2秒
  - 工控机→云端状态延迟 P95 < 3秒
```

### 12.4 指标采集


```go
// 云端关键指标
var (
    orderCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "order_created_total"},
        []string{"store_id", "channel"},
    )
    orderStateDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "order_state_duration_seconds",
            Buckets: []float64{10, 30, 60, 120, 180, 300, 600},
        },
        []string{"from_state", "to_state"},
    )
    hardwareCommandDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "hardware_command_duration_seconds",
            Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
        },
        []string{"device_type", "command"},
    )
)

// 工控机关键指标（通过 MQTT 上报，云端聚合到 Prometheus）
type EdgeMetrics struct {
    MakingDuration   float64  `json:"making_duration"`    // 制作耗时
    PickupDuration   float64  `json:"pickup_duration"`    // 取餐耗时
    QueueLength      int      `json:"queue_length"`       // 队列长度
    PLCErrorCount    int      `json:"plc_error_count"`    // PLC错误次数
    HardwareUptime   float64  `json:"hardware_uptime"`    // 硬件在线时长占比
    MQTTReconnects   int      `json:"mqtt_reconnects"`    // MQTT重连次数
}
```

### 12.5 告警分级与路由


```plaintext
告警分级路由：

  P0 紧急（影响营业）：
  - 门店所有设备离线
  - 支付系统不可用
  - MQTT Broker 宕机
  → 电话 + 短信 + 企业微信（7×24小时）

  P1 重要（影响体验）：
  - 单台设备连续离线 > 5分钟
  - 制作失败率 > 5%
  - 取餐码生成失败
  → 企业微信 + 邮件（工作时间）

  P2 警告（需关注）：
  - CPU/内存/磁盘使用率 > 80%
  - 队列积压 > 10杯
  - MQTT重连频繁
  → 企业微信（仅记录，工作时间处理）

  P3 信息（趋势分析）：
  - 制作耗时环比上升 20%
  - 某SKU销量异常波动
  → 仅入Grafana大盘，不告警
```

## 十三、设备 Provisioning（设备入网）

文档假设设备已有 device_id 和证书，生产级必须有自动化入网流程。  
### 13.1 设备生命周期


```plaintext
工厂生产          首次上电           正式运行          退役
  │                │                 │               │
  ▼                ▼                 ▼               ▼
  
  烧录             连接              连接             吊销
  Provisioning     Bootstrap         生产MQTT         证书
  证书             Server            Broker           清除设备
  
  注册设备         获取生产          正常工作          从设备
  序列号           证书+配置         定期更新证书      列表移除
```

### 13.2 首次入网流程（JITR）


```plaintext
Step 1: 工厂阶段
  ┌──────────┐
  │ 工厂产线  │
  │          │── 烧录 Provisioning 证书（所有设备共用一个）
  │          │── 烧录设备序列号 SN（唯一）
  │          │── 记录到云端设备注册表（状态=已出厂）
  └──────────┘

Step 2: 首次上电（门店安装后第一次开机）
  ┌──────────────────────────────────────────────┐
  │ 工控机 OTA Agent                              │
  │                                              │
  │ 1. 用 Provisioning 证书连接 Bootstrap Server  │
  │    POST /provision                           │
  │    Body: {                                   │
  │      serial_no: "HL100-2024-00456",         │
  │      hardware_fingerprint: "xxx",            │  ← 主板序列号等硬件指纹
  │      firmware_version: "v1.0.0",            │
  │      store_claim_code: "ABCD1234"            │  ← 门店绑定码（扫码获取）
  │    }                                         │
  └──────────────────────┬───────────────────────┘
                         ▼
  ┌──────────────────────────────────────────────┐
  │ Bootstrap Server（云端）                      │
  │                                              │
  │ 1. 验证 Provisioning 证书有效性               │
  │ 2. 验证 serial_no 在注册表中且未被使用         │
  │ 3. 验证 hardware_fingerprint 匹配             │
  │ 4. 根据 store_claim_code 绑定到门店            │
  │ 5. 生成设备唯一的生产证书（含 device_id）       │
  │ 6. 下发配置：                                 │
  │    - 生产 MQTT Broker 地址                    │
  │    - device_id: "dev_xxx"                    │
  │    - 生产证书 + 私钥                          │
  │    - 加密密钥                                 │
  │    - 门店配置                                 │
  └──────────────────────┬───────────────────────┘
                         ▼
  ┌──────────────────────────────────────────────┐
  │ 工控机                                        │
  │                                              │
  │ 1. 保存生产证书到安全存储（TPM/加密分区）       │
  │ 2. 删除 Provisioning 证书                     │
  │ 3. 用生产证书连接生产 MQTT Broker              │
  │ 4. 上报首次心跳                               │
  │ 5. 同步商品/配置                              │
  │ 6. 进入正常工作状态                            │
  └──────────────────────────────────────────────┘
```

### 13.3 证书自动轮换


```plaintext
设备证书有效期：1年

轮换流程（到期前30天自动触发）：

  云端                                工控机
    │                                    │
    │ ── MQTT通知证书即将过期 ──────────→ │
    │    device/{id}/cert/renew           │
    │                                    │
    │ ←── 申请新证书（旧证书签名） ────── │
    │                                    │
    │ ── 下发新证书 ──────────────────→ │
    │                                    │
    │                          ┌─────────┤
    │                          │ 验证新证书│
    │                          │ 切换到新证书│
    │                          │ 删除旧证书│
    │                          └─────────┤
    │                                    │
    │ ←── 新证书首次心跳 ──────────────── │
    │                                    │
    │ ── 确认轮换成功 ──────────────────→ │

  如果轮换失败（设备离线/证书错误）：
  - 旧证书仍在有效期 → 继续使用
  - 旧证书过期 → 设备用 Provisioning 证书重新入网
  - Provisioning 证书也过期 → 需要人工现场处理
```

## 十四、安全深度设计

文档之前列了5层安全框架，但缺少证书轮换、密钥管理、审计日志的具体方案。  
### 14.1 密钥管理（KMS）


```plaintext
敏感密钥分层管理：

  ┌─────────────────────────────────────────────┐
  │  L1: 支付密钥（最高安全级别）                  │
  │  - 微信支付 API Key / 证书                    │
  │  - 支付宝私钥                                │
  │  存储：云端 KMS（AWS KMS / 阿里云 KMS）       │
  │  规则：绝不下发到工控机，只在云端使用           │
  └─────────────────────────────────────────────┘
  ┌─────────────────────────────────────────────┐
  │  L2: 设备证书（中等安全级别）                  │
  │  - MQTT TLS 客户端证书                       │
  │  - 设备签名密钥                              │
  │  存储：云端 KMS 签发，工控机 TPM/加密分区存储  │
  │  规则：私钥永不离开设备，1年自动轮换           │
  └─────────────────────────────────────────────┘
  ┌─────────────────────────────────────────────┐
  │  L3: 业务加密密钥                             │
  │  - 本地 SQLite 加密密钥                      │
  │  - MQTT 消息签名密钥                         │
  │  存储：云端下发，工控机安全存储                │
  │  规则：每季度轮换，设备维度独立                │
  └─────────────────────────────────────────────┘
```

### 14.2 零信任网络


```plaintext
传统网络（不安全）：
  工控机 ←→ 工控机    ← 可以互相访问
  工控机 ←→ 咖啡机    ← 可以访问任意设备

零信任网络：
  规则：每个连接都必须验证身份，默认拒绝一切

  工控机A ──✕──→ 工控机B（禁止，不需要互相访问）
  工控机A ──→ MQTT Broker（允许，需证书验证）
  工控机A ──→ 云端API（允许，需Token验证）
  工控机A ──→ 本门店咖啡机（允许，VLAN隔离）
  工控机A ──✕──→ 其他门店设备（禁止，跨VLAN阻断）

实现方式：
  - 门店路由器：ACL规则，只允许工控机→云端 + 工控机→本店硬件
  - 工控机防火墙：iptables 白名单，只开放需要的端口
  - MQTT Broker：每台设备独立Topic权限（ACL）
  - 云端API：每台设备独立API Token + IP白名单
```

### 14.3 审计日志


```go
// 所有远程操作必须记录审计日志
type AuditLog struct {
    ID          string    `json:"id"`
    Timestamp   time.Time `json:"timestamp"`
    Actor       string    `json:"actor"`         // 操作者（用户ID/系统）
    ActorIP     string    `json:"actor_ip"`      // 操作来源IP
    ActionType  string    `json:"action_type"`   // ota_update / config_change / restart / refund
    TargetType  string    `json:"target_type"`   // device / store / order
    TargetID    string    `json:"target_id"`     // 目标ID
    BeforeState string    `json:"before_state"`  // 变更前状态
    AfterState  string    `json:"after_state"`   // 变更后状态
    Reason      string    `json:"reason"`        // 操作原因
    ApprovalID  string    `json:"approval_id"`   // 审批单号（高危操作）
}

// 高危操作需要二次审批
var HighRiskActions = map[string]bool{
    "ota_update":      true,   // 固件升级
    "device_restart":  true,   // 远程重启
    "refund":          true,   // 退款
    "config_override": true,   // 覆盖配置
    "factory_reset":   true,   // 恢复出厂
}

func (s *AuditService) Log(ctx context.Context, log AuditLog) error {
    if HighRiskActions[log.ActionType] && log.ApprovalID == "" {
        return ErrApprovalRequired  // 高危操作必须先审批
    }
    return s.auditRepo.Create(ctx, &log)
}
```

### 14.4 MQTT 安全加固


```plaintext
MQTT Broker (EMQX) 安全配置：

  1. TLS 双向认证
     - Broker 出示服务器证书给设备验证
     - 设备出示客户端证书给 Broker 验证
     - 证书链验证（CA → 中间CA → 设备证书）

  2. 设备 Topic ACL（每台设备只能操作自己的Topic）
     设备 dev_001 的权限：
       允许订阅: device/dev_001/heartbeat（自己上报）
       允许订阅: device/dev_001/shadow/delta（自己接收指令）
       允订订阅: store/store_042/order/new（自己门店的订单）
       禁止订阅: device/dev_002/#（其他设备）
       禁止订阅: store/store_999/#（其他门店）

  3. 消息签名验证（防伪造）
     - 关键消息（订单下发/取消/退款）附加 HMAC 签名
     - 工控机验证签名后才执行
     - 签名密钥每设备独立，定期轮换

  4. 速率限制
     - 单设备消息频率限制（防DDoS）
     - 异常消息模式检测（短时间内大量订单=异常）
```

## 十五、灾备与多活

文档之前没有灾备设计，生产级不能有单点故障。  
### 15.1 云端高可用架构


```plaintext
┌─────────────────────────────────────────────────────────┐
│                  云端多可用区部署                         │
│                                                         │
│  ┌───────────── 可用区 A ──────────────┐                │
│  │  K8s Node × 3                      │                │
│  │  PostgreSQL Primary                │                │
│  │  Redis Primary                     │                │
│  │  EMQX Node × 2                    │                │
│  └────────────────────────────────────┘                │
│         │                │                              │
│      同步复制          同步复制                          │
│         │                │                              │
│  ┌───────────── 可用区 B ──────────────┐                │
│  │  K8s Node × 3                      │                │
│  │  PostgreSQL Standby                │                │
│  │  Redis Replica                     │                │
│  │  EMQX Node × 2                    │                │
│  └────────────────────────────────────┘                │
│                                                         │
│  Failover 规则：                                        │
│  - 可用区A故障 → 自动切换到B（< 30秒）                   │
│  - DNS健康检查 → 故障时自动切换流量                      │
│  - 数据库主从切换 → Patroni 自动 Failover               │
└─────────────────────────────────────────────────────────┘
```

### 15.2 MQTT 高可用


```plaintext
EMQX 集群（4节点，跨可用区）：

  ┌──── AZ-A ────┐        ┌──── AZ-B ────┐
  │  EMQX Node1  │←─共享─→│  EMQX Node3  │
  │  EMQX Node2  │←─状态─→│  EMQX Node4  │
  └──────────────┘  Redis └──────────────┘
                    共享

  特性：
  - 任意节点宕机 → 设备自动重连到其他节点
  - 消息持久化到共享数据库 → 不丢消息
  - 会话跨节点共享 → 设备重连后恢复订阅
  - 负载均衡 → 按设备数量自动分配

  工控机 MQTT 客户端配置：
  - 连接地址：mqtt.example.com（DNS 轮询多个节点）
  - 自动重连：指数退避（1s → 2s → 4s → ... → 60s）
  - Clean Session: false（重连后恢复订阅）
  - LWT (Last Will Testament)：设备意外断开时自动上报离线
```

### 15.3 工控机极端灾备（云端完全宕机）


```plaintext
极端场景：云端完全不可用（72小时以上）

  工控机独立运行能力评估：

  ┌────────────┬──────────────┬──────────────────┐
  │ 功能       │ 云端宕机是否可用│ 说明              │
  ├────────────┼──────────────┼──────────────────┤
  │ POS下单    │  可用       │ 本地创建订单      │
  │ 制作       │  可用       │ 本地调度+硬件控制 │
  │ 取餐       │  可用       │ 本地传感器检测    │
  │ 小程序下单 │  不可用     │ 依赖云端          │
  │ 在线支付   │  不可用     │ 依赖云端          │
  │ 离线支付   │  可用       │ 现金/刷卡         │
  │ 打印小票   │  可用       │ USB直连          │
  │ 商品价格   │  可用       │ 本地SQLite缓存   │
  │ 会员优惠   │  降级      │ 不记名折扣        │
  └────────────┴──────────────┴──────────────────┘

  关键设计：
  1. 工控机本地SQLite能存至少7天订单数据
  2. POS端有"离线模式"切换开关
  3. 支持离线收款（现金/刷卡），恢复后对账
  4. 商品价格/配方本地有完整缓存
  5. 网络恢复后自动同步所有积压数据

  数据同步优先级（恢复后）：
  P0: 上传已完成订单（财务相关）
  P1: 上传退款/取消记录
  P2: 同步库存变更
  P3: 上传硬件日志
  P4: 同步商品/配置变更
```

### 15.4 数据备份策略


```plaintext
备份策略：

  PostgreSQL：
  - 实时：WAL 流复制到备用库
  - 每小时：增量快照到对象存储
  - 每天：全量备份到异地对象存储
  - 保留：每日备份保留30天，月末备份保留1年

  Redis：
  - 每5分钟：RDB 快照
  - AOF 持久化（每秒 fsync）

  EMQX 消息：
  - 持久化消息保留7天
  - 重要消息（订单/支付）转存到 Kafka → 对象存储

  工控机本地数据：
  - 每日定时上传到云端
  - 云端验证完整性

  恢复演练：
  - 每月进行一次灾备切换演练
  - 每季度进行一次"云端完全宕机"演练
  - RTO（恢复时间目标）< 30分钟
  - RPO（数据恢复点目标）< 5分钟
```

## 十六、更新后的总结

### 完整最佳实践清单


|  | 实践 | 状态 | 章节 |
| --- | --- | --- | --- |
| 1 | 云端/工控机分层架构 | ✅ | 二 |
| 2 | MQTT + Outbox Pattern | ✅ | 五 |
| 3 | 硬件适配层 + 熔断重试 | ✅ | 四 |
| 4 | 离线降级 + 自动恢复 | ✅ | 四 |
| 5 | 设备影子（Device Shadow） | ✅ 补充 | 十 |
| 6 | OTA 灰度 + A/B分区 + 回滚 | ✅ 补充 | 十一 |
| 7 | 分布式追踪 + SLO + 告警分级 | ✅ 补充 | 十二 |
| 8 | 设备 JITR 自动入网 | ✅ 补充 | 十三 |
| 9 | KMS密钥管理 + 零信任 + 审计 | ✅ 补充 | 十四 |
| 10 | 多可用区 + MQTT集群 + 极端灾备 | ✅ 补充 | 十五 |

### 与业界标杆对标


| 维度 | 本文档(补充后) | AWS IoT | 瑞幸/星巴克 |
| --- | --- | --- | --- |
| 云边分层 | ✅ | ✅ | ✅ |
| MQTT | ✅ | ✅ | ✅ |
| 设备影子 | ✅ | ✅ Device Shadow | ✅ |
| OTA灰度+回滚 | ✅ | ✅ Jobs | ✅ |
| JITR入网 | ✅ | ✅ | ✅ |
| 可观测性 | ✅ OTel+SLO | ✅ CloudWatch | ✅ |
| 安全 | ✅ KMS+零信任 | ✅ X.509+IAM | ✅ |
| 灾备多活 | ✅ 多AZ | ✅ 多Region | ✅ |
| 版本管理 | ✅ 兼容矩阵 | ✅ | ✅ |

补充后的文档已覆盖 IoT 零售系统完整的生产级最佳实践。  
## 十七、当前代码与设计文档差距分析

本章基于对 hooloo 项目实际代码库的审计，对比设计文档（第一至十六章）与真实代码的差距。设计文档是架构愿景，代码是当前现实。  
### 17.1 差距清单

#### 差距 1：MQTT Connect() 被注释（🔴 高）

- 文件：common/pkg/mqtt/mqtt.go:442-445
- 现状：Connect() 调用被完全注释，GetMQTTSubClientNew 返回一个从未连接的客户端
- 代码证据：

```go
// if token := c.Connect(); token.Wait() && token.Error() != nil {
//     panic(token.Error())
// }
```

- 影响：MQTT 订阅客户端可能静默失败，消息无法接收
- 建议：取消注释，增加指数退避重连策略，移除 panic 改为 error 返回
#### 差距 2：硬编码 X.509 根证书（🟡 中）

- 文件：common/pkg/mqtt/mqtt.go:40, 453
- 现状：X509RootPem 证书内容直接硬编码在 Go 源码中
- 影响：CA 证书过期后需要重新编译部署，无法热更新
- 建议：移到外部文件（如 /etc/hooloo/certs/ca.pem），启动时读取
#### 差距 3：SetOrderMatters 被注释（🟡 中）

- 文件：common/pkg/mqtt/mqtt.go:432
- 现状：MQTT v5 的 SetOrderMatters 被注释掉
- 影响：消息可能乱序处理，导致状态机跳转异常
- 建议：启用 SetOrderMatters(true)，或在业务层做版本号校验
#### 差距 4：50 秒读取超时，无写入超时（🔴 高）

- 文件：iot/hardware/embedded/huijin_v2_tcp.go:204
- 现状：TCP 读取设置 50 秒超时，但写入操作（line 410 附近）没有设置超时
- 影响：写入阻塞时 goroutine 永久挂起，导致 PLC 通信线程泄漏
- 建议：写入也设置超时（3 秒），并用 context 传播取消信号
#### 差距 5：缺少写入 ACK 验证（🟡 中）

- 文件：iot/hardware/embedded/huijin_v2_tcp.go:242
- 现状：代码中有 TODO 注释 "验证是否收到发送的指令"，但从未实现
- 影响：PLC 指令发送后无法确认硬件是否执行成功
- 建议：实现写入后的状态读取确认，增加重试机制
#### 差距 6：PumpStateMap 无互斥锁（🔴 高）

- 文件：iot/hardware/plc/interface.go:36
- 现状：PumpStateMap 是普通 map[int]int，被多个 goroutine 并发读写
- 影响：Go 运行时检测到并发 map 写入会 panic（fatal error: concurrent map writes）
- 建议：改为 sync.Map 或用 sync.Mutex 保护
#### 差距 7：SQLite 无 WAL 模式 + panic 崩溃（🟡 中）

- 文件：common/drives/sqlite.go:42, 45
- 现状：gorm.Open 使用空配置（无 WAL、无 busy_timeout），错误时 panic
- 影响：断电时数据库可能损坏；panic 会导致工控机进入重启循环
- 建议：开启 WAL 模式、设置 busy_timeout、错误返回而非 panic
#### 差距 8：本地 MQTT Broker 无认证（🟡 中）

- 文件：iot/hardware/mochi/mqtt_local_service.go:69
- 现状：AllowHook 允许匿名连接，无任何认证
- 影响：门店网络内任何设备都可以连接本地 MQTT Broker 发送伪造消息
- 建议：添加 Token 或证书认证，或至少限制为 localhost 连接
### 17.2 差距优先级矩阵


| 差距 | 严重度 | 业务影响 | 修复难度 | 优先级 |
| --- | --- | --- | --- | --- |
| 1.MQTT Connect 注释 | 🔴 高 | 消息静默丢失 | 低（取消注释） | P0 |
| 1.无写入超时 | 🔴 高 | goroutine 泄漏 | 低（加 SetWriteDeadline） | P0 |
| 1.PumpStateMap 竞争 | 🔴 高 | 运行时 panic | 低（改 sync.Map） | P0 |
| 1.硬编码证书 | 🟡 中 | 无法热更新 CA | 中（改配置加载） | P1 |
| 1.SetOrderMatters | 🟡 中 | 消息乱序 | 低（取消注释） | P1 |
| 1.无写入 ACK | 🟡 中 | 指令不可靠 | 中（实现确认逻辑） | P1 |
| 1.SQLite 配置 | 🟡 中 | 断电损坏 | 低（加 WAL 配置） | P1 |
| 1.MQTT 无认证 | 🟡 中 | 伪造消息 | 中（加认证逻辑） | P2 |

## 十八、Go 生态参考项目附录

本章基于 GitHub 深度调研，列出验证或改进本设计方案的开源 Go 项目。所有项目均经过交叉审查确认。  
### 18.1 核心依赖推荐


| 领域 | 推荐项目 | Stars | 状态 | 选型理由 |
| --- | --- | --- | --- | --- |
| MQTT 客户端 | eclipse/paho.golang | 1.3k | 活跃 | MQTT v5 官方 Go 实现，支持 QoS 0/1/2 |
| MQTT 边缘 Broker | mochi-mqtt/server | 1.8k | 活跃 | 已在项目中集成（mqtt_local_service.go），嵌入式 MQTT v5 |
| MQTT 云端 Broker | EMQX | 16k | 活跃 | 100M 连接，内置规则引擎，OTel 追踪传播 |
| 状态机 | looplab/fsm | 2.8k | 活跃 | 最流行的 Go FSM，适合订单状态机 |
| 状态机（替代） | machina | 282 | 活跃 | 泛型编译期类型安全，适合严格状态约束 |
| 熔断器 | sony/gobreaker | 2.8k | 活跃 v2.4.0 | 替代已归档的 hystrix-go |
| 韧性套件 | resilience4go | 41 | 活跃 | 重试 + 熔断 + 限流一体化 |
| 分布式追踪 | OpenTelemetry Go SDK | — | 活跃 | OTel + EMQX MQTT v5 用户属性传播 TraceID |
| Modbus | 自有实现 | — | — | 项目已有 CRC（common/tool/crc_modbus.go），不迁移 |

### 18.2 参考架构项目


| 项目 | Stars | 语言 | 学习要点 |
| --- | --- | --- | --- |
| EdgeX Foundry | 1.5k | Go | LF Edge IoT 框架，设备服务模式（device service pattern），配置驱动硬件抽象 |
| Magistrala (MainfluxLabs) | 3k | Go | IoT 平台，Device Shadow 设计（issue \#1256 独立验证了我们的设备影子方案） |
| thin-edge.io | 279 | Rust | OTA 灰度策略 + 硬件安全模块（HSM）架构参考，语言无关的模式 |

### 18.3 技术选型决策矩阵


| 领域 | 推荐方案 | 当前项目状态 | 迁移优先级 | 备注 |
| --- | --- | --- | --- | --- |
| MQTT 客户端 | paho.golang | 自封装客户端 | P2 | 功能可用但缺少 v5 特性 |
| MQTT 边缘 Broker | mochi-mqtt | ✅ 已集成 | — | 仅需修复认证配置 |
| 状态机 | looplab/fsm | 整数枚举手动流转 | P1 | 集中状态管理，防非法转换 |
| 熔断器 | sony/gobreaker | 手动 channel 互斥 | P1 | 替换脆弱的并发控制 |
| Modbus 库 | 自有实现 | ✅ 已有 CRC | — | 不迁移 goburrow/modbus |
| 分布式追踪 | OTel SDK | logger.Println | P2 | 全链路追踪，SLO 监控 |

### 18.4 明确不推荐


| 项目/方案 | 原因 |
| --- | --- |
| goburrow/modbus | 项目未使用此库（有自有 CRC 实现）；Modbus 协议自 1979 年未变，"无更新" ≠ "不可用" |
| hystrix-go | 已归档（archived），5 年无安全补丁，无 Go modules 支持 |
| EST 协议 (RFC 7030) | 对于 ~50 台设备的封闭车队属于过度工程化；手动证书轮换是合理的选择 |
| 自建 OTA 框架 | OTA 由基础设施层（k3s/systemd/Ansible）处理是正确的架构选择，无需在 Go 应用层实现 |

### 18.5 扩展参考项目（深度调研）

以下项目通过 GitHub 代码搜索 + Web 搜索发现，均经过与 hooloo 架构的匹配度评估。按相关度排序。  
#### 高度相关（强烈推荐研究）

1. [fdi-iiot-gateway](https://github.com/sophie-nguyenthuthuy/fdi-iiot-gateway) — IIoT 边缘到云端遥测系统  
与 hooloo 架构几乎一模一样的生产级实现：  

| 维度 | fdi-iiot-gateway | hooloo |
| --- | --- | --- |
| 边缘端 | Go 二进制（~25MB），systemd 部署 | Go 工控机程序 |
| 南向协议 | OPC-UA + Modbus TCP/RTU（读 PLC） | Modbus + Serial + USB |
| 北向协议 | MQTT 5.0 + mTLS + Sparkplug B | MQTT + HTTP |
| 断网缓冲 | BoltDB store-and-forward 队列 | SQLite 本地缓存（设计文档） |
| OTA | 命令 Topic 监听 + 配置重载 + 远程诊断 | k3s/systemd 基础设施层 |
| 心跳上报 | 每 30s 网关健康（队列深度/CPU/PLC 可达性） | 每 30s 心跳（设计一致） |
| 安全 | mTLS + 每租户 ACL + 仅出站连接 | TLS + 证书认证 |

技术栈：Go + chi + EMQX + TimescaleDB + Helm/K8s。核心学习点：BoltDB store-and-forward 队列的断网缓冲实现，可直接参考替代我们的 SQLite 方案。  
2. [anviod/edgex](https://github.com/anviod/edgex) — 工业边缘计算网关  
Go 1.25+ 实现的工业边缘网关，南向/北向分离架构与 hooloo 完全一致：  
- 南向：Modbus / BACnet / OPC-UA / CAN / PLC(S7) / EtherNet-IP
- 北向：MQTT / OPC-UA Server / Sparkplug-B
- 边缘规则引擎：expr-lang 表达式，支持 Check / Fail Action / Rollback
- 使用 [simonvetter/modbus](https://github.com/simonvetter/modbus)（而非 goburrow/modbus）—— Modbus 库的活跃替代
- 使用 [gopcua/opcua](https://github.com/gopcua/opcua) —— OPC UA Go 实现
- 技术栈：Go + Fiber + Paho MQTT + Vue 3
- 配置：YAML 驱动（channels.yaml / northbound.yaml / edge_rules.yaml）
核心学习点：边缘规则引擎的 Check/Rollback 模式（制作失败自动回滚），以及 simonvetter/modbus 作为 goburrow 的替代方案。  
3. [unitedrhino/things（联犀）](https://github.com/unitedrhino/things) — SaaS 云原生物联网平台  
Stars: 617 | 框架: go-zero 微服务 + gRPC  
- 协议：MQTT / CoAP / HTTP / TCP / Modbus / 阿里云/腾讯云/涂鸦云对接
- 核心能力：OTA 升级、物模型（通用/品类/产品/设备四级）、规则引擎、场景自动化、告警管理、多租户
- 存储：MySQL + TDengine（时序） + Redis + MinIO | 消息：NATS + EMQX
- 部署：单体/微服务/集群三模式，最低 2G 内存
核心学习点：如果 hooloo 要从"单店系统"升级到"多门店 SaaS 平台"，联犀是最直接的参考。物模型设计、多租户权限、OTA 全流程、规则引擎都完整实现了。[官方文档](https://doc.unitedrhino.com/)  
4. [sagoo-cloud/sagooiot](https://github.com/sagoo-cloud/sagooiot) — 企业级物联网平台  
Stars: 834 | 框架: GoFrame 2.9 + Vue 3  
- 协议：TCP / MQTT / UDP / CoAP / HTTP / Modbus / OPC UA / IEC104
- 核心能力：物模型管理、设备全生命周期、热插拔插件系统（C/Python/Go）、规则引擎、边缘计算
- 存储：MySQL + TDengine/InfluxDB + Redis
核心学习点：热插拔插件系统——hoolloo 目前每加一种新硬件协议就要改 Go 代码重新编译，SagooIOT 的插件化设计（跨进程 gRPC 通信 + 热更新）可以解决这个问题。[官方文档](https://iotdoc.sagoo.cn/)  
5. [Edgenesis/shifu](https://github.com/Edgenesis/shifu) — Kubernetes 原生 IoT 网关  
CNCF Landscape 项目 | 核心概念：DeviceShifu = 设备的数字孪生（K8s Pod）  
- 协议：HTTP / MQTT / RTSP / Siemens S7 / TCP socket / OPC UA
- 每个物理设备对应一个 K8s Pod，提供高层抽象 API
- 协议无关、即插即用、CNCF 认可
核心学习点：如果 hooloo 未来迁移到 K8s 部署，Shifu 的 DeviceShifu 模式可以让每个 PLC/咖啡机/制冰机变成一个 K8s CRD 资源。[官方文档](https://shifu.dev/)  
#### 中度相关（值得学习特定模块）

6. [maomao94/zero-service](https://github.com/maomao94/zero-service) — go-zero 工业微服务脚手架  
IEC 104 / Modbus TCP/RTU / MQTT / gRPC / HTTP 多协议。核心服务 bridgemodbus（Modbus 桥接）+ bridgemqtt（MQTT 桥接）+ SocketIO 实时推送。bridgemodbus/bridgemqtt 的 Go 代码结构可直接参考。  
7. [lf-edge/ekuiper](https://github.com/lf-edge/ekuiper) — LF Edge 流处理引擎  
边缘端轻量级流处理规则引擎（SQL 式规则查询），使用 paho.golang autopaho 自动重连。适合工控机端做本地告警逻辑。  
8. [vogler75/monster-mq-edge](https://github.com/vogler75/monster-mq-edge) — 边缘 MQTT Broker  
基于 [mochi-mqtt](https://github.com/mochi-mqtt/server)（和 hooloo 一样），增加 SQLite/PostgreSQL 存储 + GraphQL API + MQTT Bridge + Users/ACL。展示如何在 mochi-mqtt 基础上构建完整边缘 Broker。  
9. [lgustavopalmieri/microbroker-mqtt-edge](https://github.com/lgustavopalmieri/microbroker-mqtt-edge) — 超轻量边缘 MQTT Broker  
从零实现 MQTT 3.1.1 + SQLite 持久化 + FIFO 队列，内置 OEE（设备可用性）实时计算引擎。适合资源极度受限的边缘设备（~50MB RAM）。  
10. [simonvetter/modbus](https://github.com/simonvetter/modbus) — Modbus 库（goburrow 替代）  
被 anviod/edgex 工业网关选用的活跃 Modbus 库。如果 hooloo 未来需要标准化 Modbus 库（目前是自有实现），这是推荐替代。  
11. [gopcua/opcua](https://github.com/gopcua/opcua) — Go OPC UA 实现  
OPC UA Client + Server 完整实现。如果 hooloo 未来接入更高端工业设备（支持 OPC UA 的咖啡机/制冰机）。  
12. [goingforstudying-ctrl/comqtt](https://github.com/goingforstudying-ctrl/comqtt) — 分布式 MQTT Broker  
mochi-mqtt 的 fork，增加集群支持（Gossip + Raft）。如果需要在云端自建 MQTT Broker 集群（替代 EMQX 商业版）。  
13. [amenzhinsky/iothub](https://github.com/amenzhinsky/iothub) — Azure IoT Hub Go 客户端  
设备 provisioning、X.509 证书认证、Edge 模块。JITR 设备入网流程参考，X.509 证书管理实现。  
14. [junedo/fluxmq](https://github.com/junedo/fluxmq) — 高性能多协议消息 Broker  
MQTT 3.1.1/5.0 + AMQP 多协议 + 内嵌 etcd 集群 + gRPC。Magistrala 底层 Broker，跨协议 durable queues 模式。  
15. [jaab-tech/fluxrig](https://github.com/jaab-tech/fluxrig) — 边缘到云端编排引擎  
协议翻译 + 业务逻辑编排（ISO8583 / Modbus / MQTT / JSON）。架构：Mixer（控制面）+ Rack（边缘代理）+ Gear（处理模块）。边缘自治 + OpenTelemetry 原生，适合支付+IoT 混合场景。  
#### 按用途分类速查


| 需求 | 推荐项目 | 链接 |
| --- | --- | --- |
| 完整边缘网关参考 | fdi-iiot-gateway | [GitHub](https://github.com/sophie-nguyenthuthuy/fdi-iiot-gateway) |
| 中国 IoT 平台参考 | unitedrhino/things | [GitHub](https://github.com/unitedrhino/things) |
| 插件化协议扩展 | sagoo-cloud/sagooiot | [GitHub](https://github.com/sagoo-cloud/sagooiot) |
| K8s 原生设备管理 | Edgenesis/shifu | [GitHub](https://github.com/Edgenesis/shifu) |
| Modbus 桥接封装 | maomao94/zero-service | [GitHub](https://github.com/maomao94/zero-service) |
| 边缘端规则引擎 | lf-edge/ekuiper | [GitHub](https://github.com/lf-edge/ekuiper) |
| mochi-mqtt 增强 | monster-mq-edge | [GitHub](https://github.com/vogler75/monster-mq-edge) |
| Modbus 库替代 | simonvetter/modbus | [GitHub](https://github.com/simonvetter/modbus) |
| 自建 MQTT 集群 | comqtt | [GitHub](https://github.com/goingforstudying-ctrl/comqtt) |
| 设备入网（JITR） | amenzhinsky/iothub | [GitHub](https://github.com/amenzhinsky/iothub) |

## 附录：修订记录（Revision Notes）

本节记录对文档已有章节（第一至十六章）的勘误和补充，保留原始内容以维护版本历史。经 4 人对抗性评审团队（ultrabrain / unspecified-high / unspecified-low / artistry）三轮交叉审查后确认。  
### 修订 R1：第 5.2 节 — MQTT Topic 设计修正

#### R1.1 PII 暴露风险

- 原文：Topic 列表包含 user/{userID}/order/{orderNo}
- 问题：AWS IoT Core 明确规定 Topic 不得包含 PII（个人身份信息）；user/{userID} 暴露用户标识在 Topic 路径中
- 修正：用户通知应使用 WebSocket/SSE 服务端推送，而非包含用户 ID 的 MQTT Topic
- 替代方案：如需 MQTT 推送，使用设备维度 Topic（device/{deviceID}/notify），由设备端转发给用户
#### R1.2 缺少 Shadow Topic

- 问题：第十章添加了 Device Shadow 设计，但第 5.2 节的 Topic 列表未同步更新
- 补充 Topic：
    - device/{deviceID}/shadow/delta — 云端下发期望状态变更
    - device/{deviceID}/shadow/reported — 设备上报实际状态
    - device/{deviceID}/shadow/get — 设备请求当前影子
    - device/{deviceID}/shadow/get/accepted — 云端返回影子内容
#### R1.3 缺少共享订阅前缀

- 问题：IPC 集群场景下多个工控机需要负载均衡消费订单
- 补充：增加共享订阅 Topic
    - $share/<group>/store/{storeID}/order/new — 同组工控机竞争消费
    - $share/<group>/store/{storeID}/command/restart — 运维指令负载均衡
### 修订 R2：第十四章 — X.509 证书轮换补充

#### R2.1 当前方案评估

- 当前文档：证书有效期 1 年，到期前 30 天自动轮换
- 评估结论：对 ~50 台设备的封闭车队，1 年有效期 + 手动轮换是可接受的
- 补充建议：增加 CRL（证书吊销列表）或短证书隐式吊销机制，用于设备被盗/密钥泄露场景
#### R2.2 扩展建议（非当前优先级）

- 当车队规模超过 500 台时，考虑引入：
    - EST 协议（RFC 7030）自动化证书 enrollment
    - 90 天短生命周期证书
    - TPM / Secure Element 硬件密钥存储
    - IDevID（工厂身份）→ LDevID（运行身份）分层模型
- 当前规模无需引入这些复杂度
### 修订 R3：误报排除记录

以下问题在初轮审查中被提出，经交叉审查后确认为误报，特此记录以避免重复。  

| 初轮发现 | 交叉审查结论 | 排除理由 |
| --- | --- | --- |
| goburrow/modbus "未维护"，风险极高 | 误报 | 项目未使用此库；有自有 CRC 实现（common/tool/crc_modbus.go）；Modbus 协议自 1979 年未变 |
| "代码中完全没有 OTA 实现" | 误报 | OTA 在基础设施层处理（k3s 容器编排 / systemd / Ansible），Go 应用层无需实现 OTA |
| "必须使用 EST 协议做证书轮换" | 误报 | 对 50 台设备用 EST 属于过度工程化；手动轮换 + CRL 是合理选择 |

### 修订 R4：Modbus 帧结构补充（第 5.1 节关联）

原架构文档未展开 Modbus 帧结构，补充参考：  

```plaintext
Modbus RTU 串行帧：
┌──────────┬──────────┬───────────────┬──────────┐
│ 从站地址  │ 功能码    │ 数据区        │ CRC16    │
│ 1 byte   │ 1 byte   │ 0-252 bytes  │ 2 bytes  │
└──────────┴──────────┴───────────────┴──────────┘
帧间隔：≥ 3.5 字符静默时间

Modbus TCP 以太网帧：
┌─────────────────────── MBAP Header ──────────────────────┐
│ Transaction ID (2B) │ Protocol ID (2B) │ Length (2B) │ Unit ID (1B) │
└───────────────────────────────────────────────────────────┘
│                     PDU（功能码 + 数据）                     │
端口：502
无 CRC（TCP 自带校验）

功能码：
  0x03 — 读保持寄存器（状态查询）
  0x06 — 写单个寄存器（控制指令）
```
