package service

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"go.etcd.io/bbolt"
)

func AddBackup(db *bbolt.DB, backup database.Backup) error {
	// 使用db.Update方法启动一个事务
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取名为"backups"的Bucket
		b := tx.Bucket([]byte("backups"))
		// 将backup对象序列化为JSON格式的字节切片
		v, err := json.Marshal(backup)
		if err != nil {
			// 如果序列化出错，则返回错误
			return err
		}
		// 将序列化后的数据存入Bucket中，键为backup.BackupId
		if err := b.Put([]byte(backup.BackupId), v); err != nil {
			// 如果存入数据出错，则返回错误
			return err
		}
		// 返回nil表示操作成功
		return nil
	})
}

func GetBackup(db *bbolt.DB, backupId string) (database.Backup, error) {
	var backup database.Backup
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为 "backups" 的 bucket
		b := tx.Bucket([]byte("backups"))
		// 根据 backupId 获取对应的值
		v := b.Get([]byte(backupId))
		// 如果没有找到对应的记录，则返回 ErrNoRecord 错误
		if v == nil {
			return ErrNoRecord
		}
		// 将获取到的值解析到 backup 变量中
		return json.Unmarshal(v, &backup)
	})
	// 如果出现错误，则进行处理
	if err != nil {
		// 如果错误是 ErrNoRecord，则返回 ErrNoRecord 错误
		if err == ErrNoRecord {
			return backup, ErrNoRecord
		}
		// 如果不是 ErrNoRecord 错误，则返回原始错误
		return backup, err
	}
	// 如果没有错误，则返回 backup 和 nil
	return backup, nil
}

func DeleteBackup(db *bbolt.DB, backupId string) error {
	// 使用db.Update进行数据库事务更新
	return db.Update(func(tx *bbolt.Tx) error {
		// 从事务中获取名为"backups"的Bucket
		b := tx.Bucket([]byte("backups"))
		// 从Bucket中删除指定backupId的备份
		return b.Delete([]byte(backupId))
	})
}

func ListBackups(db *bbolt.DB, startTime, endTime time.Time) ([]database.Backup, error) {
	backups := make([]database.Backup, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		// 获取名为 "backups" 的 bucket
		b := tx.Bucket([]byte("backups"))
		return b.ForEach(func(k, v []byte) error {
			var backup database.Backup
			// 将字节数据反序列化为 Backup 对象
			if err := json.Unmarshal(v, &backup); err != nil {
				return err
			}
			// 根据时间筛选
			// 如果 startTime 为零或者 backup 的保存时间晚于 startTime
			// 并且 endTime 为零或者 backup 的保存时间早于 endTime
			if (startTime.IsZero() || backup.SaveTime.After(startTime)) &&
				(endTime.IsZero() || backup.SaveTime.Before(endTime)) {
				backups = append(backups, backup)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	// 根据保存时间对 backups 进行排序
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].SaveTime.Before(backups[j].SaveTime)
	})
	return backups, nil
}
