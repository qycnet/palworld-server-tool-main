from typing import Sequence

from palworld_save_tools.archive import *


def decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    # 检查类型名是否为"MapProperty"
    if type_name != "MapProperty":
        raise Exception(f"Expected MapProperty, got {type_name}")

    # 从reader中获取属性值
    value = reader.property(type_name, size, path, nested_caller_path=path)

    # Decode the raw bytes and replace the raw data
    # 解码原始字节并替换原始数据
    group_map = value["value"]
    for group in group_map:
        # 获取组类型
        group_type = group["value"]["GroupType"]["value"]["value"]
        # 获取原始数据字节
        group_bytes = group["value"]["RawData"]["value"]["values"]
        # 解码原始数据字节
        group["value"]["RawData"]["value"] = decode_bytes(
            reader, group_bytes, group_type
        )

    return value


def decode_bytes(
    parent_reader: FArchiveReader, group_bytes: Sequence[int], group_type: str
) -> dict[str, Any]:
    # 创建parent_reader的内部副本，用于读取group_bytes
    reader = parent_reader.internal_copy(bytes(group_bytes), debug=False)

    # 初始化group_data字典
    group_data = {
        # 设置group_type
        "group_type": group_type,
        # 读取并设置group_id
        "group_id": reader.guid(),
        # 读取并设置group_name
        "group_name": reader.fstring(),
        # 读取并设置individual_character_handle_ids
        "individual_character_handle_ids": reader.tarray(instance_id_reader),
    }

    # 如果group_type为"Guild"、"IndependentGuild"或"Organization"
    if group_type in [
        "EPalGroupType::Guild",
        "EPalGroupType::IndependentGuild",
        "EPalGroupType::Organization",
    ]:
        # 设置org字典
        org = {
            # 读取并设置org_type
            "org_type": reader.byte(),
            # 读取并设置base_ids
            "base_ids": reader.tarray(uuid_reader),
        }
        # 将org字典合并到group_data中
        group_data |= org

    # 如果group_type为"Guild"或"IndependentGuild"
    if group_type in ["EPalGroupType::Guild", "EPalGroupType::IndependentGuild"]:
        # 设置guild字典
        guild: dict[str, Any] = {
            # 读取并设置base_camp_level
            "base_camp_level": reader.i32(),
            # 读取并设置map_object_instance_ids_base_camp_points
            "map_object_instance_ids_base_camp_points": reader.tarray(uuid_reader),
            # 读取并设置guild_name
            "guild_name": reader.fstring(),
        }
        # 将guild字典合并到group_data中
        group_data |= guild

    # 如果group_type为"IndependentGuild"
    if group_type == "EPalGroupType::IndependentGuild":
        # 设置indie字典
        indie = {
            # 读取并设置player_uid
            "player_uid": reader.guid(),
            # 读取并设置guild_name_2
            "guild_name_2": reader.fstring(),
            # 设置player_info字典
            "player_info": {
                # 读取并设置last_online_real_time
                "last_online_real_time": reader.i64(),
                # 读取并设置player_name
                "player_name": reader.fstring(),
            },
        }
        # 将indie字典合并到group_data中
        group_data |= indie

    # 如果group_type为"Guild"
    if group_type == "EPalGroupType::Guild":
        # 设置guild字典
        guild = {
            # 读取并设置unknown_bytes
            "unknown_bytes": reader.byte_list(16),
            # 读取并设置admin_player_uid
            "admin_player_uid": reader.guid(),
            # 初始化players列表
            "players": [],
        }
        # 读取并设置player_count
        player_count = reader.i32()
        # 循环读取players信息
        for _ in range(player_count):
            player = {
                # 读取并设置player_uid
                "player_uid": reader.guid(),
                # 设置player_info字典
                "player_info": {
                    # 读取并设置last_online_real_time
                    "last_online_real_time": reader.i64(),
                    # 读取并设置player_name
                    "player_name": reader.fstring(),
                },
            }
            # 将player添加到players列表中
            guild["players"].append(player)
        # 将guild字典合并到group_data中
        group_data |= guild

    # 如果未到达文件末尾
    if not reader.eof():
        raise Exception("Warning: EOF not reached")

    # 返回group_data字典
    return group_data


def encode(
    writer: FArchiveWriter, property_type: str, properties: dict[str, Any]
) -> int:
    # 检查property_type是否为"MapProperty"，如果不是则抛出异常
    if property_type != "MapProperty":
        raise Exception(f"Expected MapProperty, got {property_type}")
    
    # 删除properties中的"custom_type"键
    del properties["custom_type"]
    
    # 获取properties中的"value"键对应的值，赋给group_map
    group_map = properties["value"]
    
    # 遍历group_map中的每个元素
    for group in group_map:
        # 如果group中的"value"键对应的字典中的"RawData"键对应的字典中的"value"键包含"values"键，则跳过当前循环
        if "values" in group["value"]["RawData"]["value"]:
            continue
        
        # 获取group中的"value"键对应的字典中的"RawData"键对应的字典中的"value"键对应的值，赋给p
        p = group["value"]["RawData"]["value"]
        
        # 对p进行编码，得到encoded_bytes
        encoded_bytes = encode_bytes(p)
        
        # 将encoded_bytes中的每个字节放入一个列表中，并将这个列表赋给group中的"value"键对应的字典中的"RawData"键对应的字典中的"value"键
        group["value"]["RawData"]["value"] = {"values": [b for b in encoded_bytes]}
    
    # 调用writer的property_inner方法，传入property_type和properties，并返回结果
    return writer.property_inner(property_type, properties)


def encode_bytes(p: dict[str, Any]) -> bytes:
    # 创建一个FArchiveWriter对象
    writer = FArchiveWriter()
    # 写入组ID
    writer.guid(p["group_id"])
    # 写入组名
    writer.fstring(p["group_name"])
    # 写入个体角色句柄ID数组
    writer.tarray(instance_id_writer, p["individual_character_handle_ids"])
    
    # 判断组类型，如果为工会、独立工会或组织类型，则执行以下操作
    if p["group_type"] in [
        "EPalGroupType::Guild",
        "EPalGroupType::IndependentGuild",
        "EPalGroupType::Organization",
    ]:
        # 写入组织类型
        writer.byte(p["org_type"])
        # 写入基础ID数组
        writer.tarray(uuid_writer, p["base_ids"])
    
    # 判断组类型，如果为工会或独立工会类型，则执行以下操作
    if p["group_type"] in ["EPalGroupType::Guild", "EPalGroupType::IndependentGuild"]:
        # 写入基础营地等级
        writer.i32(p["base_camp_level"])
        # 写入基础营地点映射对象实例ID数组
        writer.tarray(uuid_writer, p["map_object_instance_ids_base_camp_points"])
        # 写入工会名
        writer.fstring(p["guild_name"])
    
    # 判断组类型，如果为独立工会类型，则执行以下操作
    if p["group_type"] == "EPalGroupType::IndependentGuild":
        # 写入玩家UID
        writer.guid(p["player_uid"])
        # 写入第二个工会名
        writer.fstring(p["guild_name_2"])
        # 写入玩家信息中的最后在线时间
        writer.i64(p["player_info"]["last_online_real_time"])
        # 写入玩家名
        writer.fstring(p["player_info"]["player_name"])
    
    # 判断组类型，如果为工会类型，则执行以下操作
    if p["group_type"] == "EPalGroupType::Guild":
        # 写入未知字节
        writer.write(bytes(p["unknown_bytes"]))
        # 写入管理员玩家UID
        writer.guid(p["admin_player_uid"])
        # 写入玩家数量
        writer.i32(len(p["players"]))
        # 遍历玩家列表，写入每个玩家的信息
        for i in range(len(p["players"])):
            writer.guid(p["players"][i]["player_uid"])
            writer.i64(p["players"][i]["player_info"]["last_online_real_time"])
            writer.fstring(p["players"][i]["player_info"]["player_name"])
    
    # 获取编码后的字节数据
    encoded_bytes = writer.bytes()
    return encoded_bytes