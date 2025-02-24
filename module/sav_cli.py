import sys
import json
import shutil
import time
import argparse
from urllib.parse import urljoin
import requests

from structurer import convert_sav, structure_player, structure_guild
from logger import log

if __name__ == "__main__":
    start = time.time()
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--file", "-f", help="File to convert", type=str, default="Level.sav"
    )
    parser.add_argument("--clear", "-c", help="Clear input file", action="store_true")
    parser.add_argument(
        "--output", "-o", help="Output file", type=str, default="structure.json"
    )
    parser.add_argument("--request", "-r", help="Request", type=str, default="")
    parser.add_argument("--token", "-t", help="Request token", type=str, default="")
    args = parser.parse_args()

    if args.request == "":
        output = args.output
        if not args.output.endswith(".json"):
            output = args.output + ".json"
    if not os.path.exists(args.file):
        log(f"文件不存在: {args.file}", "ERROR")
        sys.exit(1)

    convert_sav(args.file)
    filetime = os.stat(args.file).st_mtime

    # 同路径下的Players文件夹
    dir_path = os.path.join(os.path.dirname(args.file), "Players")

    players = structure_player(dir_path, filetime=filetime)
    guilds = structure_guild(filetime)

    # Add last_online to players
    for player in players:
        for guild in guilds:
            guild_players = guild["players"]
            for guild_player in guild_players:
                if player["player_uid"] == guild_player["player_uid"]:
                    player["save_last_online"] = guild_player["last_online"]
                    break

    if args.request == "":
        with open(output, "w", encoding="utf-8") as f:
            json.dump(
                {"players": players, "guilds": guilds}, f, indent=4, ensure_ascii=False
            )
        log(f"Players: {len(players)}")
        log(f"Guilds: {len(guilds)}")
    else:
        player_url = urljoin(args.request, "player")
        guild_url = urljoin(args.request, "guild")
        # 记录将玩家信息发送到指定URL的操作，并记录玩家的数量
        log(f"将玩家信息发送到 {player_url} 并记录玩家数量: {len(players)}")
        player_res = requests.put(
            player_url,
            headers={"Authorization": f"Bearer {args.token}"},
            json=players,
            timeout=10,
        )
        if player_res.status_code != 200:
            # 记录发送玩家数据时的错误信息
            log(f"发送玩家数据错误: {player_res.text}")

        # 记录将工会信息发送到指定URL的操作，并记录工会数量
        log(f"将工会信息发送到 {guild_url} 并记录工会数量: {len(guilds)}")
        guild_res = requests.put(
            guild_url,
            headers={"Authorization": f"Bearer {args.token}"},
            json=guilds,
            timeout=10,
        )
        if guild_res.status_code != 200:
            # 记录发送工会数据时的错误信息
            log(f"发送工会数据错误: {guild_res.text}")

    try:
        if args.clear:
            os.remove(args.file)

            if os.path.exists(dir_path):
                # 删除Players文件夹
                shutil.rmtree(dir_path)
    except FileNotFoundError:
        pass

    # 记录操作的完成时间，并计算操作耗时（以秒为单位）
    log(f"完成时间 {round(time.time() - start, 3)}s")
