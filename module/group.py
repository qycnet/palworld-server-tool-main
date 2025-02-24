from typing import Sequence

from palworld_save_tools.archive import *


def decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    # 检查类型名称是否为"MapProperty"
    if type_name != "MapProperty":
        raise Exception(f"预期的 MapProperty(地图属性), got {type_name}")
    # 从reader中读取属性值
    value = reader.property(type_name, size, path, nested_caller_path=path)
    # Decode the raw bytes and replace the raw data
    # 获取值中的group_map
    group_map = value["value"]
    # 遍历group_map中的每个group
    for group in group_map:
        # 获取group的类型
        group_type = group["value"]["GroupType"]["value"]["value"]
        # 获取group的原始字节数据
        group_bytes = group["value"]["RawData"]["value"]["values"]
        # 解码原始字节数据并替换原始数据
        group["value"]["RawData"]["value"] = decode_bytes(
            reader, group_bytes, group_type
        )
    # 返回处理后的值
    return value


def decode_bytes(
    parent_reader: FArchiveReader, group_bytes: Sequence[int], group_type: str
) -> dict[str, Any]:
    # 创建reader的副本，读取group_bytes数据
    reader = parent_reader.internal_copy(bytes(group_bytes), debug=False)
    # 初始化group_data字典，包含group的基本信息
    group_data = {
        "group_type": group_type,
        "group_id": reader.guid(),
        "group_name": reader.fstring(),
        "individual_character_handle_ids": reader.tarray(instance_id_reader),
    }

    # 如果是公会或独立公会或组织
    if group_type in [
        "EPalGroupType::Guild",
        "EPalGroupType::IndependentGuild",
        "EPalGroupType::Organization",
    ]:
        # 初始化org字典，包含组织类型和基础ID数组
        org = {
            "org_type": reader.byte(),
            "base_ids": reader.tarray(uuid_reader),
        }
        # 将org字典合并到group_data中
        group_data |= org

    # 如果是公会或独立公会
    if group_type in ["EPalGroupType::Guild", "EPalGroupType::IndependentGuild"]:
        # 初始化guild字典，包含公会的基础营地等级、基础营地点对象实例ID数组和公会名称
        guild: dict[str, Any] = {
            "base_camp_level": reader.i32(),
            "map_object_instance_ids_base_camp_points": reader.tarray(uuid_reader),
            "guild_name": reader.fstring(),
        }
        # 将guild字典合并到group_data中
        group_data |= guild

    # 如果是独立公会
    if group_type == "EPalGroupType::IndependentGuild":
        # 初始化indie字典，包含玩家UID、公会名称2和玩家信息
        indie = {
            "player_uid": reader.guid(),
            "guild_name_2": reader.fstring(),
            "player_info": {
                "last_online_real_time": reader.i64(),
                "player_name": reader.fstring(),
            },
        }
        # 将indie字典合并到group_data中
        group_data |= indie

    # 如果是公会
    if group_type == "EPalGroupType::Guild":
        # 初始化guild字典，包含未知值u1和u2、管理员玩家UID和玩家列表
        guild = {
            # these are unknown values that don't seem to have any meaning
            "u1": reader.i64(), # perhaps like history_admin_player_uid
            "u2": reader.i64(), # always 0, not sure what this is
            "admin_player_uid": reader.guid(),
            "players": [],
        }
        # 读取玩家数量
        player_count = reader.i32()
        # 遍历每个玩家，读取玩家UID和玩家信息，并添加到玩家列表中
        for _ in range(player_count):
            player = {
                "player_uid": reader.guid(),
                "player_info": {
                    "last_online_real_time": reader.i64(),
                    "player_name": reader.fstring(),
                },
            }
            guild["players"].append(player)
        # 将guild字典合并到group_data中
        group_data |= guild

    # 如果reader没有到达文件末尾，则抛出异常
    if not reader.eof():
        raise Exception("警告：未达到 EOF")

    # 返回group_data字典
    return group_data


def encode(
    writer: FArchiveWriter, property_type: str, properties: dict[str, Any]
) -> int:
    # 检查property_type是否为"MapProperty"
    if property_type != "MapProperty":
        # 如果不是，抛出异常
        raise Exception(f"预期的 MapProperty(地图属性), got {property_type}")
    # 删除properties字典中的"custom_type"键
    del properties["custom_type"]
    # 获取properties字典中"value"键对应的值
    group_map = properties["value"]
    # 遍历group_map中的每个元素
    for group in group_map:
        # 如果group的"value"键对应的字典中的"RawData"键对应的字典中的"value"键对应的字典中包含"values"键
        if "values" in group["value"]["RawData"]["value"]:
            # 则跳过当前循环
            continue
        # 否则，获取group的"value"键对应的字典中的"RawData"键对应的字典中的"value"键对应的值
        p = group["value"]["RawData"]["value"]
        # 对p进行编码
        encoded_bytes = encode_bytes(p)
        # 将编码后的字节列表赋值给group的"value"键对应的字典中的"RawData"键对应的字典中的"value"键
        group["value"]["RawData"]["value"] = {"values": [b for b in encoded_bytes]}
    # 调用writer的property_inner方法，并返回结果
    return writer.property_inner(property_type, properties)


def encode_bytes(p: dict[str, Any]) -> bytes:
    # 创建一个FArchiveWriter对象
    writer = FArchiveWriter()
    # 写入组ID
    writer.guid(p["group_id"])
    # 写入组名
    writer.fstring(p["group_name"])
    # 写入角色ID列表
    writer.tarray(instance_id_writer, p["individual_character_handle_ids"])

    # 如果组类型为特定类型之一
    if p["group_type"] in [
        "EPalGroupType::Guild",
        "EPalGroupType::IndependentGuild",
        "EPalGroupType::Organization",
    ]:
        # 写入组织类型
        writer.byte(p["org_type"])
        # 写入基础ID列表
        writer.tarray(uuid_writer, p["base_ids"])

    # 如果组类型为特定类型之一
    if p["group_type"] in ["EPalGroupType::Guild", "EPalGroupType::IndependentGuild"]:
        # 写入基础营地等级
        writer.i32(p["base_camp_level"])
        # 写入基础营地点ID列表
        writer.tarray(uuid_writer, p["map_object_instance_ids_base_camp_points"])
        # 写入公会名
        writer.fstring(p["guild_name"])

    # 如果组类型为独立公会
    if p["group_type"] == "EPalGroupType::IndependentGuild":
        # 写入玩家UID
        writer.guid(p["player_uid"])
        # 写入第二个公会名
        writer.fstring(p["guild_name_2"])
        # 写入玩家最后在线时间
        writer.i64(p["player_info"]["last_online_real_time"])
        # 写入玩家名
        writer.fstring(p["player_info"]["player_name"])

    # 如果组类型为公会
    if p["group_type"] == "EPalGroupType::Guild":
        # 写入管理员玩家UID
        writer.guid(p["admin_player_uid"])
        # 写入玩家数量
        writer.i32(len(p["players"]))
        # 遍历玩家列表
        for i in range(len(p["players"])):
            # 写入玩家UID
            writer.guid(p["players"][i]["player_uid"])
            # 写入玩家最后在线时间
            writer.i64(p["players"][i]["player_info"]["last_online_real_time"])
            # 写入玩家名
            writer.fstring(p["players"][i]["player_info"]["player_name"])

    # 获取编码后的字节数据
    encoded_bytes = writer.bytes()
    return encoded_bytes