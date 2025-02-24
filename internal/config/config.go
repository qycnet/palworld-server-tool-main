package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
)

type Config struct {
	Web struct {
		Password  string `mapstructure:"password"`
		Port      int    `mapstructure:"port"`
		Tls       bool   `mapstructure:"tls"`
		CertPath  string `mapstructure:"cert_path"`
		KeyPath   string `mapstructure:"key_path"`
		PublicUrl string `mapstructure:"public_url"`
	} `mapstructure:"web"`
	Task struct {
		SyncInterval        int    `mapstructure:"sync_interval"`
		PlayerLogging       bool   `mapstructure:"player_logging"`
		PlayerLoginMessage  string `mapstructure:"player_login_message"`
		PlayerLogoutMessage string `mapstructure:"player_logout_message"`
	} `mapstructure:"task"`
	Rcon struct {
		Address   string `mapstructure:"address"`
		Password  string `mapstructure:"password"`
		UseBase64 bool   `mapstructure:"use_base64"`
		Timeout   int    `mapstructure:"timeout"`
	} `mapstructure:"rcon"`
	Rest struct {
		Address  string `mapstructure:"address"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Timeout  int    `mapstructure:"timeout"`
	} `mapstructure:"rest"`
	Save struct {
		Path           string `mapstructure:"path"`
		DecodePath     string `mapstructure:"decode_path"`
		SyncInterval   int    `mapstructure:"sync_interval"`
		BackupInterval int    `mapstructure:"backup_interval"`
		BackupKeepDays int    `mapstructure:"backup_keep_days"`
	} `mapstructure:"save"`
	Manage struct {
		KickNonWhitelist bool `mapstructure:"kick_non_whitelist"`
	}
}

func Init(cfgFile string, conf *Config) {
	// 如果指定了配置文件路径
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	} else {
		// 如果没有指定配置文件路径，则使用默认配置
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		// 如果配置文件未找到，记录警告日志
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//logger.Warn("config file not found, try to read from env\n")
			logger.Warn("找不到配置文件，请尝试从 env 读取\n")
		} else {
			// 如果配置文件找到但出现其他错误，记录panic日志
			//logger.Panic("config file was found but another error was produced\n")
			logger.Panic("找到配置文件，但产生另一个错误\n")
		}
	}

	// 设置默认配置
	viper.SetDefault("web.port", 8080)

	viper.SetDefault("task.sync_interval", 60)

	viper.SetDefault("rcon.timeout", 5)
	viper.SetDefault("rcon.use_base64", false)

	viper.SetDefault("rest.username", "admin")
	viper.SetDefault("rest.timeout", 5)

	viper.SetDefault("save.sync_interval", 600)
	viper.SetDefault("save.backup_interval", 14400)
	viper.SetDefault("save.backup_keep_days", 7)

	// 设置环境变量前缀和替换器
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	// 将配置解析到结构体中
	err = viper.Unmarshal(conf)
	if err != nil {
		// 如果解析失败，记录panic日志
		//logger.Panicf("Unable to decode config into struct, %s", err)
		logger.Panicf("无法将配置解码到结构体, %s", err)
	}
}
