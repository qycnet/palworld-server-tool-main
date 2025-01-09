package service

import (
	"encoding/json"

	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"go.etcd.io/bbolt"
)

func PutGuilds(db *bbolt.DB, guilds []database.Guild) error {
	// 使用db的Update方法执行数据库操作
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为"guilds"的bucket
		b := tx.Bucket([]byte("guilds"))
		// 遍历guilds数组
		for _, g := range guilds {
			// 将guild对象序列化为JSON格式的字节数组
			v, err := json.Marshal(g)
			if err != nil {
				// 如果序列化出错，则返回错误
				return err
			}
			// 将序列化后的字节数组存入bucket中，键为guild的AdminPlayerUid
			if err := b.Put([]byte(g.AdminPlayerUid), v); err != nil {
				// 如果存入bucket出错，则返回错误
				return err
			}
		}
		// 所有操作成功完成，返回nil
		return nil
	})
}

func ListGuilds(db *bbolt.DB) ([]database.Guild, error) {
	// 初始化一个空的 Guild 切片
	guilds := make([]database.Guild, 0)
	// 使用 db.View 方法以只读模式打开数据库事务
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为 "guilds" 的 Bucket
		b := tx.Bucket([]byte("guilds"))
		// 遍历 Bucket 中的所有键值对
		return b.ForEach(func(k, v []byte) error {
			// 初始化一个 Guild 结构体变量
			var guild database.Guild
			// 将字节切片 v 解码到 guild 结构体中
			if err := json.Unmarshal(v, &guild); err != nil {
				// 如果解码失败，返回错误
				return err
			}
			// 将解码后的 guild 结构体添加到 guilds 切片中
			guilds = append(guilds, guild)
			// 返回 nil 表示处理成功
			return nil
		})
	})
	// 如果处理过程中发生错误，返回错误
	if err != nil {
		return nil, err
	}
	// 返回 guilds 切片和 nil 表示成功
	return guilds, nil
}

func GetGuild(db *bbolt.DB, playerUID string) (database.Guild, error) {
	var guild database.Guild
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为"guilds"的bucket
		b := tx.Bucket([]byte("guilds"))

		// 遍历bucket中的所有guild
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// 解析出当前遍历到的guild对象
			var g database.Guild
			if err := json.Unmarshal(v, &g); err != nil {
				return err
			}

			// 检查当前guild的players是否包含指定的player_uid
			for _, player := range g.Players {
				if player.PlayerUid == playerUID {
					// 如果找到，将找到的guild赋值给外部的guild变量
					guild = g
					return nil
				}
			}
		}
		// 如果没有找到匹配的guild，返回ErrNoRecord错误
		return ErrNoRecord
	})
	if err != nil {
		return database.Guild{}, err
	}
	return guild, nil
}
