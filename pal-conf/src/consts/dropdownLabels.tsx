// 定义死亡惩罚标签的常量数组
// 这些标签描述了玩家死亡时可能失去的物品类型
export const DeathPenaltyLabels = [
    {
        name: "None",        // 无惩罚，不失去任何物品
        desc: "No lost",     // 描述：不丢失任何物品
    },
    {
        name: "Item",        // 物品惩罚，失去非装备物品
        desc: "Lost item without equipment", // 描述：丢失物品（不包括装备）
    },
    {
        name: "ItemAndEquipment", // 物品和装备惩罚，失去所有物品和装备
        desc: "Lost item and equipment",    // 描述：丢失物品和装备
    },
    {
        name: "All",         // 全部惩罚，失去所有物品、装备和背包中的pal
        desc: "Lost All item, equipment, pal(in inventory)", // 描述：丢失所有物品、装备以及背包中的pal
    },
] as const; // 使用as const断言这些值在初始化后不会被修改

// 定义允许连接的平台标签的常量数组
// 这些标签描述了允许哪些平台连接到游戏服务器
export const AllowConnectPlatformLabels = [
    {
        name: "Steam",       // Steam平台
        desc: "Only allow Steam to connect", // 描述：只允许Steam平台连接
    },
    {
        name: "Xbox",        // Xbox平台
        desc: "Only allow Xbox to connect",  // 描述：只允许Xbox平台连接
    },
] as const; // 使用as const断言这些值在初始化后不会被修改

// 定义日志格式类型标签的常量数组
// 这些标签描述了日志记录的格式
export const LogFormatTypeLabels = [
    {
        name: "Text",        // 文本格式
        desc: "Use Text format to log",      // 描述：使用文本格式记录日志
    },
    {
        name: "Json",        // JSON格式
        desc: "Use Json format to log",      // 描述：使用JSON格式记录日志
    },
] as const; // 使用as const断言这些值在初始化后不会被修改

// 定义随机化类型标签的常量数组
// 这些标签描述了游戏内的随机化设置
export const RandomizerTypeLabels = [
    {
        name: "None",        // 无随机化
        desc: "No randomizer set",          // 描述：未设置随机化
    },
    {
        name: "Region",      // 区域随机化
        desc: "Set to randomize region",     // 描述：设置为随机化区域
    },
    {
        name: "All",         // 全部随机化
        desc: "Set to randomize all",        // 描述：设置为随机化所有内容
    },
] as const; // 使用as const断言这些值在初始化后不会被修改
