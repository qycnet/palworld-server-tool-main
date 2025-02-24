from typing import Any, Sequence, Dict

# 假设FArchiveReader和FArchiveWriter是从palworld_save_tools.archive模块导入的类
# 用于读取和写入SAV文件中的数据
from palworld_save_tools.archive import *

# 解码函数，用于解码ArrayProperty类型的数据
def decode(
    reader: FArchiveReader,  # FArchiveReader对象，用于读取SAV文件中的数据
    type_name: str,  # 数据的类型名称，期望为"ArrayProperty"
    size: int,  # 数据的大小（可能表示数组中元素的数量）
    path: str  # 数据在SAV文件中的路径或标识
) -> Dict[str, Any]:  # 返回一个字典，包含解码后的数据
    """
    解码ArrayProperty类型的数据。
    
    参数:
    reader (FArchiveReader): 用于读取SAV文件数据的对象。
    type_name (str): 数据的类型名称，期望为"ArrayProperty"。
    size (int): 数据的大小（可能表示数组中元素的数量）。
    path (str): 数据在SAV文件中的路径或标识。
    
    返回:
    Dict[str, Any]: 包含解码后的数据的字典。
    
    异常:
    Exception: 如果type_name不是"ArrayProperty"，则抛出异常。
    """
    if type_name != "ArrayProperty":
        raise Exception(f"预期的ArrayProperty(数组属性), got {type_name}")  # 类型不匹配时抛出异常
    
    # 读取ArrayProperty类型的数据
    value = reader.property(type_name, size, path, nested_caller_path=path)
    
    # 获取数据中的字节值
    data_bytes = value["value"]["values"]
    
    # 解码字节值
    value["value"] = decode_bytes(reader, data_bytes)  # 调用decode_bytes函数解码字节值
    
    return value  # 返回解码后的数据

# 解码字节数据的函数
def decode_bytes(
    parent_reader: FArchiveReader,  # FArchiveReader对象，用于读取字节数据
    b_bytes: Sequence[int]  # 要解码的字节数据（整数序列）
) -> Dict[str, Any]:  # 返回一个字典，包含解码后的数据
    """
    解码字节数据。
    
    参数:
    parent_reader (FArchiveReader): 用于读取字节数据的对象。
    b_bytes (Sequence[int]): 要解码的字节数据（整数序列）。
    
    返回:
    Dict[str, Any]: 包含解码后的数据的字典。
    
    异常:
    Exception: 如果在读取数据时未达到文件末尾（EOF），则抛出异常。
    """
    reader = parent_reader.internal_copy(bytes(b_bytes), debug=False)  # 创建内部复制的reader对象
    
    # 解码数据并存储到字典中
    data = {
        "id": reader.guid(),  # 读取GUID作为ID
        "name": reader.fstring(),  # 读取FString作为名称
        "state": reader.byte(),  # 读取一个字节作为状态
        "transform": reader.ftransform(),  # 读取FTransform作为变换信息
        "area_range": reader.float(),  # 读取一个浮点数作为区域范围
        "group_id_belong_to": reader.guid(),  # 读取GUID作为所属组的ID
        # "fast_travel_local_transform": reader.ftransform(),  # 注释掉的代码：读取FTransform作为快速旅行的本地变换信息（可能不需要）
        "owner_map_object_instance_id": reader.guid(),  # 读取GUID作为拥有者地图对象实例的ID
    }
    
    # 检查是否已达到文件末尾（EOF）
    if not reader.eof():
        raise Exception("警告：未达到 EOF")  # 未达到EOF时抛出异常
    
    return data  # 返回解码后的数据

# 编码函数，用于编码ArrayProperty类型的数据
def encode(
    writer: FArchiveWriter,  # FArchiveWriter对象，用于写入编码后的数据
    property_type: str,  # 数据的类型名称，期望为"ArrayProperty"
    properties: Dict[str, Any]  # 要编码的数据属性字典
) -> int:  # 返回写入的数据大小（可能是字节数）
    """
    编码ArrayProperty类型的数据。
    
    参数:
    writer (FArchiveWriter): 用于写入编码后数据的对象。
    property_type (str): 数据的类型名称，期望为"ArrayProperty"。
    properties (Dict[str, Any]): 要编码的数据属性字典。
    
    返回:
    int: 写入的数据大小（可能是字节数）。
    
    异常:
    Exception: 如果property_type不是"ArrayProperty"，则抛出异常。
    """
    if property_type != "ArrayProperty":
        raise Exception(f"预期的ArrayProperty(数组属性), got {property_type}")  # 类型不匹配时抛出异常
    
    # 删除自定义类型属性（可能不是必需的或不被支持）
    del properties["custom_type"]
    
    # 编码属性中的值数据
    encoded_bytes = encode_bytes(properties["value"])  # 调用encode_bytes函数编码值数据
    
    # 将编码后的字节数据设置为属性中的值（以列表形式，这里可能不需要列表推导）
    properties["value"] = {"values": encoded_bytes}  # 直接使用encoded_bytes，因为它已经是一个字节序列
    
    # 写入属性数据并返回写入的大小
    return writer.property_inner(property_type, properties)  # 调用writer的property_inner方法写入数据并返回大小

# 编码字节数据的函数
def encode_bytes(p: Dict[str, Any]) -> bytes:  # p是包含要编码数据的字典
    """
    编码字节数据。
    
    参数:
    p (Dict[str, Any]): 包含要编码数据的字典。
    
    返回:
    bytes: 编码后的字节数据。
    """
    writer = FArchiveWriter()  # 创建一个FArchiveWriter对象用于写入数据
    
    # 写入数据到writer对象
    writer.guid(p["id"])  # 写入GUID作为ID
    writer.fstring(p["name"])  # 写入FString作为名称
    writer.byte(p["state"])  # 写入一个字节作为状态
    writer.ftransform(p["transform"])  # 写入FTransform作为变换信息
    writer.float(p["area_range"])  # 写入一个浮点数作为区域范围
    writer.guid(p["group_id_belong_to"])  # 写入GUID作为所属组的ID
    # writer.ftransform(p["fast_travel_local_transform"])  # 注释掉的代码：写入FTransform作为快速旅行的本地变换信息（可能不需要）
    writer.guid(p["owner_map_object_instance_id"])  # 写入GUID作为拥有者地图对象实例的ID
    
    # 获取编码后的字节数据
    encoded_bytes = writer.bytes()  # 调用writer的bytes方法获取编码后的字节数据
    
    return encoded_bytes  # 返回编码后的字节数据
