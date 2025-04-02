package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
	"github.com/qycnet/palworld-server-tool-main/internal/system"

	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"github.com/qycnet/palworld-server-tool-main/internal/tool"
	"github.com/qycnet/palworld-server-tool-main/service"
	"go.etcd.io/bbolt"
)

var s gocron.Scheduler

func BackupTask(db *bbolt.DB) {
	// 记录日志，提示开始安排备份
	logger.Info("开始备份...\n")

	// 调用备份工具进行备份，并获取备份路径
	path, err := tool.Backup()
	if err != nil {
		// 如果备份过程中出现错误，记录错误日志并返回
		logger.Errorf("%v\n", err)
		return
	}

	// 将备份信息添加到数据库中
	err = service.AddBackup(db, database.Backup{
		BackupId: uuid.New().String(), // 生成新的备份ID
		Path:     path,                // 设置备份路径
		SaveTime: time.Now(),          // 设置备份时间
	})
	if err != nil {
		// 如果向数据库添加备份信息时出现错误，记录错误日志并返回
		logger.Errorf("%v\n", err)
		return
	}

	// 记录日志，提示备份完成并显示备份路径
	logger.Infof("自动备份到backups目录 %s\n", path)
	
		keepDays := viper.GetInt("save.backup_keep_days")
	if keepDays == 0 {
		keepDays = 7
	}
	err = tool.CleanOldBackups(db, keepDays)
	if err != nil {
		logger.Errorf("无法清理旧备份: %v\n", err)
	}
}

func PlayerSync(db *bbolt.DB) {
	logger.Info("玩家信息同步...\n")
	// 获取在线玩家列表
	onlinePlayers, err := tool.ShowPlayers()
	if err != nil {
		// 如果获取在线玩家列表出错，记录错误日志
		logger.Errorf("获取在线玩家列表出错 %v\n", err)
	}
	// 将在线玩家列表存入数据库
	err = service.PutPlayersOnline(db, onlinePlayers)
	if err != nil {
		// 如果将在线玩家列表存入数据库出错，记录错误日志
		logger.Errorf("%v\n", err)
	}
	logger.Info("玩家信息同步完成\n")

	// 获取玩家日志记录的配置项
	playerLogging := viper.GetBool("task.player_logging")
	if playerLogging {
		// 如果开启玩家日志记录，启动玩家日志记录协程
		go PlayerLogging(onlinePlayers)
	}

	// 获取是否踢除非白名单玩家的配置项
	kickInterval := viper.GetBool("manage.kick_non_whitelist")
	if kickInterval {
		// 如果开启踢除非白名单玩家的功能，启动检查并踢除玩家的协程
		go CheckAndKickPlayers(db, onlinePlayers)
	}
}

func isPlayerWhitelisted(player database.OnlinePlayer, whitelist []database.PlayerW) bool {
	// 遍历白名单中的每个玩家
	for _, whitelistedPlayer := range whitelist {
		// 如果玩家的UID不为空且等于白名单中的玩家UID，或者玩家的SteamID不为空且等于白名单中的玩家SteamID
		if (player.PlayerUid != "" && player.PlayerUid == whitelistedPlayer.PlayerUID) ||
			(player.SteamId != "" && player.SteamId == whitelistedPlayer.SteamID) {
			// 返回true，表示玩家在白名单中
			return true
		}
	}
	// 如果循环结束后仍未找到匹配的玩家，则返回false
	return false
}

var playerCache map[string]string
var firstPoll = true

func PlayerLogging(players []database.OnlinePlayer) {
	// 获取配置文件中定义的登录和登出消息
	loginMsg := viper.GetString("task.player_login_message")
	logoutMsg := viper.GetString("task.player_logout_message")

	// 创建一个临时map，用于存储玩家UID和昵称的映射关系
	tmp := make(map[string]string, len(players))
	// 遍历玩家列表，将玩家UID和昵称存储到临时map中
	for _, player := range players {
		// 如果玩家UID不为空，则存储到临时map中
		if player.PlayerUid != "" {
			tmp[player.PlayerUid] = player.Nickname
		}
	}

	// 如果不是第一次轮询
	if !firstPoll {
		// 遍历临时map，如果玩家在playerCache中不存在，则广播登录消息
		for id, name := range tmp {
			if _, ok := playerCache[id]; !ok {
				// 广播登录消息
				BroadcastVariableMessage(loginMsg, name, len(players))
			}
		}
		// 遍历playerCache，如果玩家在临时map中不存在，则广播登出消息
		for id, name := range playerCache {
			if _, ok := tmp[id]; !ok {
				// 广播登出消息
				BroadcastVariableMessage(logoutMsg, name, len(players))
			}
		}
	}
	// 设置firstPoll为false，表示已经进行过一次轮询
	firstPoll = false
	// 更新playerCache为最新的临时map
	playerCache = tmp
}

func BroadcastVariableMessage(message string, username string, onlineNum int) {
	// 将消息中的"{username}"替换为实际的用户名
	message = strings.ReplaceAll(message, "{username}", username)
	// 将消息中的"{online_num}"替换为实际的在线人数
	message = strings.ReplaceAll(message, "{online_num}", strconv.Itoa(onlineNum))
	// 按换行符分割消息
	arr := strings.Split(message, "\n")
	for _, msg := range arr {
		// 广播消息
		err := tool.Broadcast(msg)
		if err != nil {
			// 如果广播失败，记录警告日志
			logger.Warnf("广播失败, %s \n", err)
		}
		// 连续发送不知道为啥行会错乱, 只能加点延迟
		// 延迟1秒
		time.Sleep(1000 * time.Millisecond)
	}
}

func CheckAndKickPlayers(db *bbolt.DB, players []database.OnlinePlayer) {
	// 获取白名单
	whitelist, err := service.ListWhitelist(db)
	if err != nil {
		logger.Errorf("%v\n", err)
	}

	// 遍历所有在线玩家
	for _, player := range players {
		// 如果玩家不在白名单中
		if !isPlayerWhitelisted(player, whitelist) {
			// 获取玩家的SteamId
			identifier := player.SteamId
			// 如果SteamId为空
			if identifier == "" {
				// 日志记录：踢出失败，SteamId为空
				logger.Warnf("踢 %s 失败, SteamId 为空 \n", player.Nickname)
				continue
			}
			// 踢出玩家
			err := tool.KickPlayer(fmt.Sprintf("steam_%s", identifier))
			// 如果踢出失败
			if err != nil {
				// 日志记录：踢出失败，记录错误信息
				logger.Warnf("踢 %s 失败, %s \n", player.Nickname, err)
				continue
			}
			// 日志记录：踢出成功
			logger.Warnf("踢 %s 成功 \n", player.Nickname)
		}
	}
	// 日志记录：白名单检查完成
	logger.Info("检查白名单完成\n")
}

func SavSync() {
	// 记录日志：调度Sav同步...
	logger.Info("调度Sav同步...\n")

	// 解码viper配置文件中的保存路径
	err := tool.Decode(viper.GetString("save.path"))
	if err != nil {
		// 记录错误日志
		logger.Errorf("%v\n", err)
	}

	// 记录日志：Sav同步完成
	logger.Info("Sav同步完成\n")
}

func Schedule(db *bbolt.DB) {
	// 获取调度器实例
	s := getScheduler()

	// 从配置文件中获取玩家同步间隔时间
	playerSyncInterval := time.Duration(viper.GetInt("task.sync_interval"))
	// 从配置文件中获取保存同步间隔时间
	savSyncInterval := time.Duration(viper.GetInt("save.sync_interval"))
	// 从配置文件中获取备份间隔时间
	backupInterval := time.Duration(viper.GetInt("save.backup_interval"))

	// 如果玩家同步间隔时间大于0
	if playerSyncInterval > 0 {
		// 启动玩家同步任务
		go PlayerSync(db)
		// 创建玩家同步任务
		_, err := s.NewJob(
			gocron.DurationJob(playerSyncInterval*time.Second),
			gocron.NewTask(PlayerSync, db),
		)
		if err != nil {
			// 记录错误日志
			logger.Errorf("%v\n", err)
		}
	}

	// 如果保存同步间隔时间大于0
	if savSyncInterval > 0 {
		// 启动保存同步任务
		go SavSync()
		// 创建保存同步任务
		_, err := s.NewJob(
			gocron.DurationJob(savSyncInterval*time.Second),
			gocron.NewTask(SavSync),
		)
		if err != nil {
			// 记录错误日志
			logger.Errorf("%v\n", err)
		}
	}

	// 如果备份间隔时间大于0
	if backupInterval > 0 {
		// 启动备份任务
		go BackupTask(db)
		// 创建备份任务
		_, err := s.NewJob(
			gocron.DurationJob(backupInterval*time.Second),
			gocron.NewTask(BackupTask, db),
		)
		if err != nil {
			// 记录错误日志
			logger.Error(err)
		}
	}

	// 创建限制缓存目录大小的任务
	_, err := s.NewJob(
		gocron.DurationJob(300*time.Second),
		gocron.NewTask(system.LimitCacheDir, filepath.Join(os.TempDir(), "palworldsav-"), 5),
	)
	if err != nil {
		// 记录错误日志
		logger.Errorf("%v\n", err)
	}

	// 启动调度器
	s.Start()
}

func Shutdown() {
	// 获取调度器实例
	s := getScheduler()
	// 关闭调度器
	err := s.Shutdown()
	// 如果关闭过程中出现错误
	if err != nil {
		// 记录错误信息
		logger.Errorf("%v\n", err)
	}
}

func initScheduler() gocron.Scheduler {
	// 创建一个新的调度器实例
	s, err := gocron.NewScheduler()
	if err != nil {
		// 如果创建调度器实例时出错，记录错误信息
		logger.Errorf("%v\n", err)
	}
	// 返回创建的调度器实例
	return s
}

func getScheduler() gocron.Scheduler {
	// 如果 s 为 nil
	if s == nil {
		// 初始化调度器
		return initScheduler()
	}
	// 返回 s
	return s
}
