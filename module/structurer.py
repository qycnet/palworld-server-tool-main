import copy
import os
import sys
import zlib
import json
import time
from typing import Any

from palworld_save_tools.gvas import GvasFile
from palworld_save_tools.palsav import decompress_sav_to_gvas
from palworld_save_tools.paltypes import PALWORLD_CUSTOM_PROPERTIES, PALWORLD_TYPE_HINTS
from palworld_save_tools.archive import FArchiveReader, FArchiveWriter
import item_container_slots
import base_camp
import group

from world_types import Player, Pal, Guild, BaseCamp
from logger import log, redirect_stdout_stderr

PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.ItemContainerSaveData.Value.Slots.Slots.RawData"
] = (
    item_container_slots.decode,
    item_container_slots.encode,
)
PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.BaseCampSaveData.Value.WorkerDirector.RawData"
] = (
    base_camp.decode,
    base_camp.encode,
)
PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.GroupSaveDataMap"] = (
    group.decode,
    group.encode,
)


wsd = None
gvas_file = None


def skip_decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    # 如果是ArrayProperty类型
    if type_name == "ArrayProperty":
        # 读取数组类型
        array_type = reader.fstring()
        # 构建返回字典
        value = {
            "skip_type": type_name,
            "array_type": array_type,
            "id": reader.optional_guid(),  # 读取可选GUID
            "value": reader.read(size),  # 读取指定大小的数据
        }
    # 如果是MapProperty类型
    elif type_name == "MapProperty":
        # 读取键类型
        key_type = reader.fstring()
        # 读取值类型
        value_type = reader.fstring()
        # 读取GUID
        _id = reader.optional_guid()
        # 构建返回字典
        value = {
            "skip_type": type_name,
            "key_type": key_type,
            "value_type": value_type,
            "id": _id,
            "value": reader.read(size),
        }
    # 如果是StructProperty类型
    elif type_name == "StructProperty":
        # 构建返回字典
        value = {
            "skip_type": type_name,
            "struct_type": reader.fstring(),  # 读取结构体类型
            "struct_id": reader.guid(),  # 读取结构体GUID
            "id": reader.optional_guid(),  # 读取可选GUID
            "value": reader.read(size),  # 读取指定大小的数据
        }
    # 如果不是上述三种类型，抛出异常
    else:
        raise Exception(
            f"Expected ArrayProperty or MapProperty or StructProperty, got {type_name} in {path}"
        )
    return value


def skip_encode(
    writer: FArchiveWriter, property_type: str, properties: dict[str, Any]
) -> int:
    # 检查属性中是否包含"skip_type"
    if "skip_type" not in properties:
        # 检查"custom_type"是否在PALWORLD_CUSTOM_PROPERTIES中
        if properties["custom_type"] in PALWORLD_CUSTOM_PROPERTIES is not None:
            # 调用PALWORLD_CUSTOM_PROPERTIES中对应的函数
            return PALWORLD_CUSTOM_PROPERTIES[properties["custom_type"]][1](
                writer, property_type, properties
            )

    # 如果属性类型为"ArrayProperty"
    if property_type == "ArrayProperty":
        # 删除"custom_type"和"skip_type"属性
        del properties["custom_type"]
        del properties["skip_type"]
        # 写入数组类型
        writer.fstring(properties["array_type"])
        # 写入可选的GUID
        writer.optional_guid(properties.get("id", None))
        # 写入数组值
        writer.write(properties["value"])
        # 返回数组长度
        return len(properties["value"])

    # 如果属性类型为"MapProperty"
    elif property_type == "MapProperty":
        # 删除"custom_type"和"skip_type"属性
        del properties["custom_type"]
        del properties["skip_type"]
        # 写入键类型
        writer.fstring(properties["key_type"])
        # 写入值类型
        writer.fstring(properties["value_type"])
        # 写入可选的GUID
        writer.optional_guid(properties.get("id", None))
        # 写入映射值
        writer.write(properties["value"])
        # 返回映射值长度
        return len(properties["value"])

    # 如果属性类型为"StructProperty"
    elif property_type == "StructProperty":
        # 删除"custom_type"和"skip_type"属性
        del properties["custom_type"]
        del properties["skip_type"]
        # 写入结构体类型
        writer.fstring(properties["struct_type"])
        # 写入结构体ID
        writer.guid(properties["struct_id"])
        # 写入可选的GUID
        writer.optional_guid(properties.get("id", None))
        # 写入结构体值
        writer.write(properties["value"])
        # 返回结构体值长度
        return len(properties["value"])

    # 如果属性类型不是ArrayProperty、MapProperty或StructProperty
    else:
        # 抛出异常
        raise Exception(
            f"Expected ArrayProperty or MapProperty or StructProperty, got {property_type}"
        )


def load_skiped_decode(wsd, skip_paths, recursive=True):
    # 如果 skip_paths 是字符串类型，则将其转换为列表
    if isinstance(skip_paths, str):
        skip_paths = [skip_paths]

    # 遍历 skip_paths 列表中的每个路径
    for skip_path in skip_paths:
        # 获取当前路径对应的属性
        properties = wsd[skip_path]

        # 如果属性中不包含 "skip_type"，则跳过当前循环
        if "skip_type" not in properties:
            continue

        # 解析跳过的项目
        parse_skiped_item(properties, skip_path, recursive)

        # 如果当前路径在 SKP_PALWORLD_CUSTOM_PROPERTIES 中存在，则删除该路径对应的属性
        if ".worldSaveData.%s" % skip_path in SKP_PALWORLD_CUSTOM_PROPERTIES:
            del SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.%s" % skip_path]


SKP_PALWORLD_CUSTOM_PROPERTIES = copy.deepcopy(PALWORLD_CUSTOM_PROPERTIES)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.MapObjectSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.FoliageGridSaveDataMap"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.MapObjectSpawnerInStageSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.DynamicItemSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.CharacterContainerSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.CharacterContainerSaveData.Value.Slots"
] = (skip_decode, skip_encode)
SKP_PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.CharacterContainerSaveData.Value.RawData"
] = (skip_decode, skip_encode)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.ItemContainerSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.ItemContainerSaveData.Value.BelongInfo"
] = (skip_decode, skip_encode)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.ItemContainerSaveData.Value.Slots"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.ItemContainerSaveData.Value.RawData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.RandomizerSaveData"] = (
    skip_decode,
    skip_encode,
)


def convert_sav(file):
    global gvas_file, wsd

    # 判断文件是否以".sav.json"结尾
    if file.endswith(".sav.json"):
        # 打印日志，表示正在加载
        log("加载...")
        with open(file, "r", encoding="utf-8") as f:
            # 返回文件内容
            return f.read()

    # 打印日志，表示正在转换
    log("转换...")
    with redirect_stdout_stderr():
        try:
            # 以二进制模式打开文件
            with open(file, "rb") as f:
                # 读取文件内容
                data = f.read()
                # 解压缩sav文件到gvas格式
                raw_gvas, _ = decompress_sav_to_gvas(data)
                # 读取gvas文件
                gvas_file = GvasFile.read(
                    raw_gvas, PALWORLD_TYPE_HINTS, SKP_PALWORLD_CUSTOM_PROPERTIES
                )
        except zlib.error:
            # 打印日志，表示sav文件已损坏，并退出程序
            log("此 .sav 文件已损坏. :(", "ERROR")
            sys.exit(1)

    # 将gvas文件的内容转换为json字符串并返回（此行代码被注释掉）
    # return json.dumps(gvas_file.dump(), cls=CustomEncoder)

    # 从gvas文件中获取worldSaveData的值
    wsd = gvas_file.properties["worldSaveData"]["value"]


def structure_player(dir_path, data_source=None, filetime: int = -1):
    # 记录日志，提示正在处理玩家结构
    log("构建玩家结构...")
    # 声明全局变量 wsd
    global wsd
    # 如果数据源为空，则使用全局变量 wsd 作为数据源
    if data_source is None:
        data_source = wsd
    # 如果数据源中没有 "CharacterSaveParameterMap" 键，则返回空列表
    if not data_source.get("CharacterSaveParameterMap"):
        return []
    # 提取玩家 UID 和玩家保存参数
    uid_character = (
        # 遍历 "CharacterSaveParameterMap" 中的每个条目
        (
            c["key"]["PlayerUId"]["value"],  # 玩家 UID
            c["value"]["RawData"]["value"]["object"]["SaveParameter"]["value"],  # 玩家保存参数
        )
        for c in wsd["CharacterSaveParameterMap"]["value"]
    )

    players = []  # 初始化玩家列表
    pals = []  # 初始化伙伴列表
    ticks = wsd["GameTimeSaveData"]["value"]["RealDateTimeTicks"]["value"]
    # 遍历 UID 和玩家保存参数的元组
    for uid, c in uid_character:
        # 如果玩家保存参数中包含 "IsPlayer" 键且值为真
        if c.get("IsPlayer") and c["IsPlayer"]["value"]:
            # 获取玩家物品，并添加到玩家保存参数中
            c["Items"] = getPlayerItems(uid, dir_path)
            # 创建玩家对象，并转换为字典后添加到玩家列表中
            players.append(Player(uid, c).to_dict())
        else:
            # 如果玩家保存参数中没有 "OwnerPlayerUId" 键，则跳过当前循环
            if not c.get("OwnerPlayerUId"):
                continue
            # 创建伙伴对象，并转换为字典后添加到伙伴列表中
            pals.append(Pal(c, ticks, filetime).to_dict())

    unique_players_dict = {}  # 初始化唯一玩家字典
    # 遍历玩家列表，处理重复玩家
    for player in players:
        player_uid = player["player_uid"]  # 玩家 UID
        if player_uid in unique_players_dict:  # 如果玩家 UID 已存在于字典中
            existing_player = unique_players_dict[player_uid]  # 获取已存在的玩家
            if player["level"] > existing_player["level"]:  # 如果当前玩家等级更高
                unique_players_dict[player_uid] = player  # 更新字典中的玩家
        else:
            unique_players_dict[player_uid] = player  # 将玩家添加到字典中

    unique_players = list(unique_players_dict.values())  # 将字典转换为玩家列表
    # 遍历伙伴列表，将伙伴添加到对应玩家的 "pals" 键中
    for pal in pals:
        for player in unique_players:
            if player["player_uid"] == pal["owner"]:  # 如果伙伴的 "owner" 键与玩家的 UID 匹配
                pal.pop("owner")  # 从伙伴中移除 "owner" 键
                player["pals"].append(pal)  # 将伙伴添加到玩家的 "pals" 键中
                break  # 跳出内层循环

    sorted_players = sorted(unique_players, key=lambda p: p["level"], reverse=True)  # 按等级降序排序玩家列表

    return sorted_players  # 返回排序后的玩家列表


def parse_skiped_item(properties, skip_path, recursive=True):
    # 检查属性中是否包含"skip_type"
    if "skip_type" not in properties:
        return properties

    # 使用FArchiveReader读取属性
    with FArchiveReader(
        properties["value"],
        PALWORLD_TYPE_HINTS,
        (
            SKP_PALWORLD_CUSTOM_PROPERTIES
            if recursive == False
            else PALWORLD_CUSTOM_PROPERTIES
        ),
    ) as reader:
        # 根据skip_type的值处理不同的属性类型
        if properties["skip_type"] == "ArrayProperty":
            # hack: 0.3.7 later version has a bug that the array type include bytes
            # current use custom item_container_slots.decode to fix it
            properties["value"] = reader.array_property(
                properties["array_type"],
                len(properties["value"]) - 4,
                ".worldSaveData.%s" % skip_path,
            )
        elif properties["skip_type"] == "StructProperty":
            properties["value"] = reader.struct_value(
                properties["struct_type"], ".worldSaveData.%s" % skip_path
            )
        elif properties["skip_type"] == "MapProperty":
            reader.u32()
            count = reader.u32()
            path = ".worldSaveData.%s" % skip_path
            key_path = path + ".Key"
            key_type = properties["key_type"]
            value_type = properties["value_type"]
            # 处理键的类型
            if key_type == "StructProperty":
                key_struct_type = reader.get_type_or(key_path, "Guid")
            else:
                key_struct_type = None
            value_path = path + ".Value"
            # 处理值的类型
            if value_type == "StructProperty":
                value_struct_type = reader.get_type_or(value_path, "StructProperty")
            else:
                value_struct_type = None
            # 存储键值对
            values: list[dict[str, Any]] = []
            for _ in range(count):
                key = reader.prop_value(key_type, key_struct_type, key_path)
                value = reader.prop_value(value_type, value_struct_type, value_path)
                values.append(
                    {
                        "key": key,
                        "value": value,
                    }
                )
            properties["key_struct_type"] = key_struct_type
            properties["value_struct_type"] = value_struct_type
            properties["value"] = values
        # 删除自定义类型和skip_type属性
        del properties["custom_type"]
        del properties["skip_type"]
    return properties


def parse_item(properties, skip_path):
    # 判断properties是否为字典类型
    if isinstance(properties, dict):
        # 遍历字典中的每个键值对
        for key in properties:
            # 构建新的路径
            call_skip_path = skip_path + "." + key[0].upper() + key[1:]
            # 判断当前值是否为字典类型，并且包含"type"键，且"type"的值在["StructProperty", "ArrayProperty", "MapProperty"]中
            if (
                isinstance(properties[key], dict)
                and "type" in properties[key]
                and properties[key]["type"]
                in ["StructProperty", "ArrayProperty", "MapProperty"]
            ):
                # 判断当前值是否包含"skip_type"键
                if "skip_type" in properties[key]:
                    # print("Parsing worldSaveData.%s..." % call_skip_path, end="", flush=True)
                    # 调用parse_skiped_item函数处理当前值
                    properties[key] = parse_skiped_item(
                        properties[key], call_skip_path, True
                    )
                    # print("Done")
                else:
                    # 如果不包含"skip_type"键，则递归调用parse_item函数处理当前值的"value"字段
                    properties[key]["value"] = parse_item(
                        properties[key]["value"], call_skip_path
                    )
            else:
                # 如果当前值不满足上述条件，则递归调用parse_item函数处理当前值
                properties[key] = parse_item(properties[key], call_skip_path)
    # 判断properties是否为列表类型
    elif isinstance(properties, list):
        # 获取顶层的skip_path
        top_skip_path = ".".join(skip_path.split(".")[:-1])
        # 遍历列表中的每个元素
        for idx, item in enumerate(properties):
            # 递归调用parse_item函数处理当前元素
            properties[idx] = parse_item(item, top_skip_path)
    # 返回处理后的properties
    return properties


def getPlayerItems(player_uid, dir_path):
    # 加载跳过的解码数据
    load_skiped_decode(wsd, ["ItemContainerSaveData"], False)

    # 初始化物品容器字典
    item_containers = {}
    for item_container in wsd["ItemContainerSaveData"]["value"]:
        # 将物品容器ID作为键，物品容器作为值存入字典
        item_containers[str(item_container["key"]["ID"]["value"])] = item_container

    # 构造玩家存档文件的路径
    player_sav_file = os.path.join(
        dir_path, str(player_uid).upper().replace("-", "") + ".sav"
    )

    # 检查玩家存档文件是否存在
    if not os.path.exists(player_sav_file):
        # log("Player Sav file Not exists: %s" % player_sav_file)
        return
    else:
        # 重定向标准输出和标准错误输出
        with redirect_stdout_stderr():
            try:
                # 打开玩家存档文件并读取
                with open(player_sav_file, "rb") as f:
                    raw_gvas, _ = decompress_sav_to_gvas(f.read())
                    # 读取并解析Gvas文件
                    player_gvas_file = GvasFile.read(
                        raw_gvas, PALWORLD_TYPE_HINTS, PALWORLD_CUSTOM_PROPERTIES
                    )
                player_gvas = player_gvas_file.properties["SaveData"]["value"]
            except Exception as e:
                # 记录错误信息并返回
                log(
                    f"玩家 Sav 文件已损坏: {os.path.basename(player_sav_file)}: {str(e)}",
                    "ERROR",
                )
                return

    # 初始化容器数据字典
    containers_data = {
        "CommonContainerId": [],
        "DropSlotContainerId": [],
        "EssentialContainerId": [],
        "FoodEquipContainerId": [],
        "PlayerEquipArmorContainerId": [],
        "WeaponLoadOutContainerId": [],
    }

    # 遍历容器数据字典的键
    for idx_key in containers_data.keys():
        # 检查玩家Gvas数据中的InventoryInfo是否存在
        if player_gvas.get("InventoryInfo") is None:
            continue
        # 检查InventoryInfo中对应键的值是否存在
        if player_gvas["InventoryInfo"]["value"].get(idx_key) is None:
            continue
        # 获取容器ID
        container_id = str(
            player_gvas["InventoryInfo"]["value"][idx_key]["value"]["ID"]["value"]
        )
        # 检查容器ID是否在物品容器字典中
        if container_id in item_containers:
            # 解析对应的物品容器数据
            item_container = parse_item(
                item_containers[container_id], "ItemContainerSaveData"
            )

            # 提取每个物品的相关数据并保存到字典中
            containers_data[idx_key] = [
                {
                    "SlotIndex": item["RawData"]["value"]["permission"]["type_a"],
                    "ItemId": item["RawData"]["value"]["permission"][
                        "item_static_id"
                    ].lower(),
                    "StackCount": item["RawData"]["value"]["permission"]["type_b"],
                }
                for item in item_container["value"]["Slots"]["value"]["values"]
                if item["RawData"]["value"]["permission"]["item_static_id"].lower()
                != "none"
            ]

    # 返回容器数据字典
    return containers_data


def structure_base_camp():
    # 记录日志：正在构建基地营...
    log("构建基地营...")
    # 如果BaseCampSaveData不存在，则返回空列表
    if not wsd.get("BaseCampSaveData"):
        return []
    # 从BaseCampSaveData中提取基地营数据
    base_camps = (
        # 提取每个基地营的RawData值
        b["value"]["RawData"]["value"] for b in wsd["BaseCampSaveData"]["value"]
    )
    # 将每个基地营数据转换为BaseCamp对象，并转换为字典形式
    base_camps_generator = (BaseCamp(b).to_dict() for b in base_camps)
    # 将生成器转换为列表并返回
    return list(base_camps_generator)


def structure_guild(filetime: int = -1):
    # 输出日志信息
    log("构建公会...")

    # 检查是否存在"GroupSaveDataMap"
    if not wsd.get("GroupSaveDataMap"):
        return []

    # 调用函数获取基本营地信息
    base_camps = structure_base_camp()

    # 过滤出类型为"EPalGroupType::Guild"的群组
    groups = (
        g["value"]["RawData"]["value"]
        for g in wsd["GroupSaveDataMap"]["value"]
        if g["value"]["GroupType"]["value"]["value"] == "EPalGroupType::Guild"
    )

    # 获取游戏时间信息
    Ticks = wsd["GameTimeSaveData"]["value"]["RealDateTimeTicks"]["value"]

    # 生成公会信息的生成器
    guilds_generator = (Guild(g, Ticks, filetime).to_dict() for g in groups)

    # 根据基本营地等级对公会信息进行排序
    sorted_guilds = sorted(
        guilds_generator, key=lambda g: g["base_camp_level"], reverse=True
    )

    # 遍历排序后的公会信息
    for guild in sorted_guilds:
        # 遍历基本营地信息
        for camp in base_camps:
            # 如果基本营地的ID在公会的基地ID列表中
            if camp["id"] in guild["base_ids"]:
                # 将基本营地信息添加到公会的基地列表中
                guild["base_camp"].append(
                    {
                        "id": camp["id"],
                        "area": camp["area_range"],
                        "location_x": camp["transform"]["x"],
                        "location_y": camp["transform"]["y"],
                    }
                )

    # 返回排序后的公会信息列表
    return list(sorted_guilds)


if __name__ == "__main__":
    import time

    start = time.time()
    file = "./Level.sav"
    converted = convert_sav(file)
    players = structure_player(converted)
    log("储蓄玩家...")
    with open("players.json", "w", encoding="utf-8") as f:
        json.dump(players, f, indent=4, ensure_ascii=False)
    guilds = structure_guild(converted)
    log("储蓄工会...")
    with open("guilds.json", "w", encoding="utf-8") as f:
        json.dump(guilds, f, indent=4, ensure_ascii=False)
    log(f"完成时间 {time.time() - start}s")