package service

import (
	"encoding/json"

	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"go.etcd.io/bbolt"
	"k8s.io/apimachinery/pkg/util/uuid"
)

func AddRconCommand(db *bbolt.DB, rcon database.RconCommand) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为 "rcons" 的 bucket
		b := tx.Bucket([]byte("rcons"))

		// 将 rcon 对象序列化为 JSON 格式
		v, err := json.Marshal(rcon)
		if err != nil {
			// 如果序列化出错，返回错误
			return err
		}

		// 生成一个新的 UUID
		uuid := uuid.NewUUID()

		// 将序列化的 JSON 数据和对应的 UUID 存储到 bucket 中
		return b.Put([]byte(uuid), v)
	})
}

func PutRconCommand(db *bbolt.DB, uuid string, rcon database.RconCommand) error {
	// 使用db.Update方法进行事务操作
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为"rcons"的Bucket
		b := tx.Bucket([]byte("rcons"))
		// 将rcon对象序列化为JSON格式的字节数组
		v, err := json.Marshal(rcon)
		if err != nil {
			// 如果序列化失败，则返回错误
			return err
		}
		// 将序列化后的JSON字节数组存入Bucket中，以uuid为键
		return b.Put([]byte(uuid), v)
	})
}

func ListRconCommands(db *bbolt.DB) ([]database.RconCommandList, error) {
	// 初始化一个空的 RconCommandList 切片
	rcons := make([]database.RconCommandList, 0)

	// 以只读模式打开数据库事务
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为 "rcons" 的 bucket
		b := tx.Bucket([]byte("rcons"))

		// 遍历 bucket 中的每个键值对
		return b.ForEach(func(k, v []byte) error {
			// 定义一个 RconCommand 变量
			var rcon database.RconCommand

			// 将字节切片 v 反序列化为 RconCommand 结构体
			if err := json.Unmarshal(v, &rcon); err != nil {
				// 如果反序列化失败，则返回错误
				return err
			}

			// 将 RconCommand 结构体和对应的键（UUID）添加到 RconCommandList 切片中
			rcons = append(rcons, database.RconCommandList{
				UUID:        string(k),
				RconCommand: rcon,
			})
			return nil
		})
	})

	// 如果出现错误，则返回错误
	if err != nil {
		return nil, err
	}

	// 返回 RconCommandList 切片和 nil 错误
	return rcons, nil
}

func GetRconCommand(db *bbolt.DB, uuid string) (database.RconCommand, error) {
	var rcon database.RconCommand
	// 使用只读事务访问数据库
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为"rcons"的bucket
		b := tx.Bucket([]byte("rcons"))
		// 根据uuid获取对应的值
		v := b.Get([]byte(uuid))
		// 如果没有找到对应的记录，则返回ErrNoRecord错误
		if v == nil {
			return ErrNoRecord
		}
		// 将获取到的值解析为RconCommand结构体
		return json.Unmarshal(v, &rcon)
	})
	return rcon, err
}

func RemoveRconCommand(db *bbolt.DB, uuid string) error {
	// 使用db.Update方法更新数据库
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为"rcons"的bucket
		b := tx.Bucket([]byte("rcons"))
		// 从bucket中删除指定uuid的键值对
		return b.Delete([]byte(uuid))
	})
}
