package database

import (
	"sync"
	"time"

	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB
var once sync.Once

func InitDB() *bbolt.DB {
	// 打开或创建数据库文件
	db_, err := bbolt.Open("pst.db", 0600, &bbolt.Options{Timeout: 1 * time.Minute})
	if err != nil {
		logger.Panic(err)
	}

	// 创建"players"桶
	// players
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("players"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}

	// 创建"guilds"桶
	// guilds
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}

	// 创建"rcons"桶
	// rcons
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("rcons"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}

	// 创建"backups"桶
	// backups
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("backups"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}

	return db_
}

func GetDB() *bbolt.DB {
	// 确保db的初始化只执行一次
	once.Do(func() {
		// 初始化数据库
		db = InitDB()
	})
	// 返回初始化后的数据库实例
	return db
}
