import argparse
import asyncio
import json
import os

import aiohttp
import requests
from tqdm import tqdm

# 重试次数常量
RETRY_TIMES = 3

# 设置图片下载 URL 模板
base_url = "https://palworld.gg/images/tiles/{z}/{x}/{y}.png"

# 本地保存文件的根目录
save_dir = "./map"

# 为图片请求设置 headers
headers = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
}

# 根据z的值设置x和y的最大范围
z_to_range = {
    0: (0, 0),
    1: (1, 1),
    2: (3, 3),
    3: (7, 7),
    4: (15, 15),
    5: (31, 31),
    6: (63, 63),
}


async def download_image(session, url, file_path, custom_headers, progress_bar, redown):
    # 如果文件已存在且未启用重新下载，则更新进度条并返回
    if os.path.exists(file_path) and not redown:
        progress_bar.update(1)
        return

    attempt = 0
    success = False
    while attempt < RETRY_TIMES:
        try:
            # 发起GET请求并设置自定义请求头
            async with session.get(url, headers=custom_headers) as response:
                # 检查响应状态码
                if response.status == 200:
                    # 如果响应状态码为200，则打开文件并写入数据
                    with open(file_path, "wb") as f:
                        f.write(await response.read())
                    success = True
                    break
                elif response.status == 404:
                    # 如果响应状态码为404，则打印提示信息并退出循环
                    print(f"Skipped {url} - Not Found (404)")
                    break
                elif response.status == 403:
                    # 如果响应状态码为403，则打印提示信息并退出循环
                    print(f"Skipped {url} - Forbidden (403)")
                    break
                else:
                    # 如果响应状态码不是200，则打印错误信息并退出循环
                    print(f"Failed to download {url} (status code: {response.status})")
                    break
        except aiohttp.ClientError as e:
            # 捕获aiohttp.ClientError异常并打印错误信息
            print(f"HTTP error downloading {url}: {e}")
        except Exception as e:
            # 捕获其他异常并打印错误信息
            print(f"Unexpected error downloading {url}: {e}")
        attempt += 1
        # 如果重试次数未达到最大限制，则打印重试信息并暂停1秒
        if attempt < RETRY_TIMES:
            print(f"Retrying ({attempt}/{RETRY_TIMES}) for {url}")
            await asyncio.sleep(1)

    # 如果下载成功，则更新进度条
    if success:
        progress_bar.update(1)
    else:
        # 如果下载失败，则打印失败信息并更新进度条
        print(f"Failed to download {url} after {RETRY_TIMES} attempts")
        progress_bar.update(1)


async def download_images_async(redown=False):
    # 计算总图片数量
    total_images = sum((x_max + 1) * (y_max + 1) for x_max, y_max in z_to_range.values())
    # 创建进度条
    progress_bar = tqdm(total=total_images, desc="Downloading images", unit="img")

    async with aiohttp.ClientSession() as session:
        tasks = []
        for z, (x_max, y_max) in z_to_range.items():
            for x in range(0, x_max + 1):
                for y in range(0, y_max + 1):
                    # 生成图片URL
                    url = base_url.format(z=z, x=x, y=y)
                    # 生成保存路径
                    save_path = os.path.join(save_dir, str(z), str(x))
                    file_name = f"{y}.png"
                    file_path = os.path.join(save_path, file_name)
                    # 创建保存路径文件夹
                    os.makedirs(save_path, exist_ok=True)
                    # 创建下载任务
                    task = download_image(session, url, file_path, headers, progress_bar, redown)
                    # 将任务添加到任务列表中
                    tasks.append(task)
        # 并发执行所有下载任务
        await asyncio.gather(*tasks)

    # 关闭进度条
    progress_bar.close()


def parse_js_file():
    url = "https://paldb.cc/js/map_data_cn.js"
    try:
        # 使用requests库下载JS文件
        with requests.get(url, timeout=10) as response:
            # 将下载的文件内容写入本地文件
            with open("map_data_cn.js", "wb") as file:
                file.write(response.content)

        # 以utf-8编码读取本地JS文件内容
        with open("map_data_cn.js", "r", encoding="utf-8") as file:
            js_content = file.read()

        # 使用js2py库执行JS代码
        context = js2py.EvalJs()
        context.execute(js_content)

        # 将fixedDungeon对象转换为字典，并获取其值列表
        fixed_dungeon_obj = context.fixedDungeon.to_dict()
        fixed_dungeon = list(fixed_dungeon_obj.values())

        # 初始化结果字典
        result = {"boss_tower": [], "fast_travel": []}
        # 遍历fixedDungeon列表，根据item类型将坐标添加到相应列表中
        for item in fixed_dungeon:
            if item["type"] == "Tower":
                result["boss_tower"].append([float(item["pos"]["X"]), float(item["pos"]["Y"])])
            elif item["type"] == "Fast Travel":
                result["fast_travel"].append([float(item["pos"]["X"]), float(item["pos"]["Y"])])

        # 将结果字典写入JSON文件
        with open("web/src/assets/map/points.json", "w", encoding="utf-8") as json_file:
            json.dump(result, json_file, ensure_ascii=False, indent=4)
    except requests.RequestException as e:
        # 捕获requests异常并打印错误信息
        print(f"Error downloading JS file: {e}")
    except Exception as e:
        # 捕获其他异常并打印错误信息
        print(f"Error processing JS file: {e}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Download tiles from palworld.gg")
    parser.add_argument("--redown", action="store_true", help="Redownload existing files")
    args = parser.parse_args()

    asyncio.run(download_images_async(args.redown))
    # parse_js_file()