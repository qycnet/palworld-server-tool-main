package service

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"go.etcd.io/bbolt"
)

func PutPlayers(db *bbolt.DB, players []database.Player) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为 "players" 的 bucket
		b := tx.Bucket([]byte("players"))
		for _, p := range players {
			// 获取当前玩家的数据
			existingPlayerData := b.Get([]byte(p.PlayerUid))
			if existingPlayerData != nil {
				var existingPlayer database.Player
				// 将现有的玩家数据反序列化为结构体
				if err := json.Unmarshal(existingPlayerData, &existingPlayer); err != nil {
					return err
				}
				// 如果现有玩家的 SteamId 不为空，则保留其 SteamId
				if existingPlayer.SteamId != "" {
					p.SteamId = existingPlayer.SteamId
				}
				// 保留现有玩家的 IP 和 Ping
				p.Ip = existingPlayer.Ip
				p.Ping = existingPlayer.Ping
			}
			// 如果玩家需要保存最后在线时间，则解析时间字符串
			if p.SaveLastOnline != "" {
				p.LastOnline, _ = time.Parse(time.RFC3339, p.SaveLastOnline)
			}
			// 将玩家数据序列化为 JSON
			v, err := json.Marshal(p)
			if err != nil {
				return err
			}
			// 将序列化后的玩家数据保存到 bucket 中
			if err := b.Put([]byte(p.PlayerUid), v); err != nil {
				return err
			}
		}
		return nil
	})
}

func PutPlayersOnline(db *bbolt.DB, players []database.OnlinePlayer) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为 "players" 的 bucket
		b := tx.Bucket([]byte("players"))
		for _, p := range players {
			// 从 bucket 中获取玩家的数据
			existingPlayerData := b.Get([]byte(p.PlayerUid))
			var player database.Player
			// 如果玩家数据不存在，说明玩家在线但不在数据库中
			if existingPlayerData == nil {
				// player online but not in database
				player.PlayerUid = p.PlayerUid
				player.SteamId = p.SteamId
				player.Nickname = p.Nickname
			} else {
				// 反序列化现有玩家数据
				if err := json.Unmarshal(existingPlayerData, &player); err != nil {
					return err
				}
				// 如果玩家的 SteamId 为空或包含 "000000"，则更新 SteamId
				if player.SteamId == "" || strings.Contains(player.SteamId, "000000") {
					player.SteamId = p.SteamId
				}
			}
			// 更新玩家的其他信息
			player.Ip = p.Ip
			player.Ping = p.Ping
			player.LocationX = p.LocationX
			player.LocationY = p.LocationY
			player.Level = p.Level
			player.LastOnline = time.Now()

			// 序列化玩家数据
			v, err := json.Marshal(player)
			if err != nil {
				return err
			}
			// 将玩家数据存入 bucket
			if err := b.Put([]byte(p.PlayerUid), v); err != nil {
				return err
			}
		}
		return nil
	})
}

func ListPlayers(db *bbolt.DB) ([]database.TersePlayer, error) {
	players := make([]database.TersePlayer, 0)
	// 开启数据库只读事务
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为"players"的桶
		b := tx.Bucket([]byte("players"))
		// 遍历桶中的每个键值对
		return b.ForEach(func(k, v []byte) error {
			// 如果键包含"000000"，则跳过该键值对
			if strings.Contains(string(k), "000000") {
				return nil
			}
			// 定义一个TersePlayer类型的变量
			var player database.TersePlayer
			// 将值从JSON格式反序列化为TersePlayer类型
			if err := json.Unmarshal(v, &player); err != nil {
				return err
			}
			// 将反序列化后的TersePlayer对象添加到players切片中
			players = append(players, player)
			return nil
		})
	})
	// 如果发生错误，则返回错误
	if err != nil {
		return nil, err
	}
	// 返回players切片和nil错误
	return players, nil
}

func GetPlayer(db *bbolt.DB, playerUid string) (database.Player, error) {
	var player database.Player

	// 开启一个只读事务
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为"players"的Bucket
		b := tx.Bucket([]byte("players"))
		// 根据playerUid获取对应的值
		v := b.Get([]byte(playerUid))
		// 如果没有找到对应的值，则返回ErrNoRecord错误
		if v == nil {
			return ErrNoRecord
		}
		// 将字节数组解析为Player对象
		if err := json.Unmarshal(v, &player); err != nil {
			return err
		}
		return nil
	})

	// 如果事务执行过程中发生错误，则返回错误
	if err != nil {
		return database.Player{}, err
	}

	// 返回找到的Player对象
	return player, nil
}

func AddWhitelist(db *bbolt.DB, player database.PlayerW) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取或创建白名单bucket
		b, err := tx.CreateBucketIfNotExists([]byte("whitelist"))
		if err != nil {
			return err
		}

		// 序列化玩家数据为JSON
		playerData, err := json.Marshal(player)
		if err != nil {
			return err
		}

		// 使用 findPlayerKey 检查玩家是否已经在白名单中
		key, err := findPlayerKey(b, player)
		if err != nil {
			return err
		}

		// 如果玩家已存在，更新其信息；如果不存在，创建新的键
		if key != nil {
			// 玩家已存在，更新其信息
			if err := b.Put(key, playerData); err != nil {
				return err
			}
		} else {
			// 玩家不存在，添加新玩家
			// 生成新玩家的唯一键
			newPlayerKey := []byte(player.Name + "|" + player.SteamID + "|" + player.PlayerUID)
			if err := b.Put(newPlayerKey, playerData); err != nil {
				return err
			}
		}

		return nil
	})
}

func ListWhitelist(db *bbolt.DB) ([]database.PlayerW, error) {
	var players []database.PlayerW

	// 开启只读事务
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为 "whitelist" 的 bucket
		b := tx.Bucket([]byte("whitelist"))
		if b == nil {
			// 如果 bucket 不存在，返回 nil，表示没有错误，只是一个空列表
			return nil // No error, just an empty list if the bucket doesn't exist.
		}

		// 遍历 bucket 中的所有键值对
		return b.ForEach(func(k, v []byte) error {
			var player database.PlayerW
			// 将值 v 从 JSON 格式反序列化为 database.PlayerW 结构体
			if err := json.Unmarshal(v, &player); err != nil {
				return err
			}
			// 将解析后的 player 添加到 players 列表中
			players = append(players, player)
			return nil
		})
	})

	return players, err
}

// findPlayerKey tries to find a player in the whitelist and returns the key if found.
func findPlayerKey(b *bbolt.Bucket, player database.PlayerW) ([]byte, error) {
	var keyFound []byte
	err := b.ForEach(func(k, v []byte) error {
		// 定义一个变量 existingPlayer 来存储从 bucket 中读取的玩家信息
		var existingPlayer database.PlayerW
		// 将 bucket 中的值反序列化为 existingPlayer
		if err := json.Unmarshal(v, &existingPlayer); err != nil {
			return err
		}
		// 判断 existingPlayer 是否满足给定的条件
		if matchesCriteria(existingPlayer, player) {
			// 将找到的键复制到一个新的切片中
			keyFound = append([]byte(nil), k...) // Make a copy of the key
			// 抛出一个错误来提前退出迭代
			return errors.New("player found")    // Use an error to break out of the iteration early.
		}
		return nil
	})

	// 判断错误是否为 "player found"
	if err != nil && err.Error() == "player found" {
		return keyFound, nil
	}

	return nil, err
}

// RemoveWhitelist removes a player from the whitelist.
func RemoveWhitelist(db *bbolt.DB, player database.PlayerW) error {
	// 使用db.Update方法执行事务
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为"whitelist"的bucket
		b := tx.Bucket([]byte("whitelist"))
		// 如果bucket不存在，返回错误
		if b == nil {
			// bucket不存在
			return errors.New("whitelist bucket does not exist")
		}

		// 在bucket中查找玩家对应的key
		key, err := findPlayerKey(b, player)
		// 如果查找过程中出现错误，返回错误
		if err != nil {
			return err
		}
		// 如果未找到对应的key，返回错误
		if key == nil {
			// 玩家未在白名单中找到
			return errors.New("player not found in whitelist")
		}

		// 从bucket中删除对应的key
		return b.Delete(key)
	})
}

// matchesCriteria checks if the given player matches the criteria.
func matchesCriteria(existingPlayer, player database.PlayerW) bool {
	// 如果PlayerUID非空且匹配，认为是同一个玩家
	if player.PlayerUID != "" && existingPlayer.PlayerUID == player.PlayerUID {
		return true
	}
	// 如果Name非空且匹配，认为是同一个玩家
	if player.Name != "" && existingPlayer.Name == player.Name {
		return true
	}
	// 如果SteamID非空且匹配，认为是同一个玩家
	if player.SteamID != "" && existingPlayer.SteamID == player.SteamID {
		return true
	}
	// 如果没有任何字段匹配，返回false
	return false
}

func PutWhitelist(db *bbolt.DB, players []database.PlayerW) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取或创建白名单bucket
		b, err := tx.CreateBucketIfNotExists([]byte("whitelist"))
		if err != nil {
			return err
		}

		// 清空现有的白名单
		err = b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})
		if err != nil {
			return err
		}

		// 遍历并添加新的玩家数据到白名单
		for _, player := range players {
			playerData, err := json.Marshal(player)
			if err != nil {
				return err
			}
			identifier := player.PlayerUID
			if identifier == "" {
				if identifier = player.SteamID; identifier == "" {
					continue
				}
			}
			if err := b.Put([]byte(identifier), playerData); err != nil {
				return err
			}
		}

		return nil
	})
}
