package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/qycnet/palworld-server-tool-main/internal/logger"

	"github.com/spf13/viper"
	"github.com/qycnet/palworld-server-tool-main/internal/database"
)

var client = &http.Client{}

func callApi(method string, api string, param []byte) ([]byte, error) {

	// 从配置文件中获取REST服务的地址
	addr := viper.GetString("rest.address")
	// 从配置文件中获取REST服务的用户名
	user := viper.GetString("rest.username")
	// 从配置文件中获取REST服务的密码
	pass := viper.GetString("rest.password")
	// 从配置文件中获取REST服务的超时时间
	timeout := viper.GetInt("rest.timeout")

	// 将API地址和基地址拼接成完整的URL
	api, err := url.JoinPath(addr, api)
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, _ := http.NewRequest(method, api, bytes.NewReader(param))
	// 设置HTTP请求的基本认证信息
	req.SetBasicAuth(user, pass)

	// 设置HTTP客户端的超时时间
	client.Timeout = time.Duration(timeout) * time.Second
	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 确保在函数返回前关闭响应体
	defer resp.Body.Close()

	// 读取响应体的内容
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 检查HTTP响应状态码是否为200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("休息: %d %s", resp.StatusCode, b)
	}
	return b, nil
}

type ResponseInfo struct {
	Version     string `json:"version"`
	ServerName  string `json:"servername"`
	Description string `json:"description"`
	Worldguid string `json:"worldguid"`
}

func Info() (map[string]string, error) {
	// 调用API获取信息
	resp, err := callApi("GET", "/v1/api/info", nil)
	if err != nil {
		// 如果调用API出错，则返回错误
		return nil, err
	}

	// 定义ResponseInfo结构体变量
	var data ResponseInfo
	// 将响应的JSON数据解析到ResponseInfo结构体中
	err = json.Unmarshal(resp, &data)
	if err != nil {
		// 如果解析出错，则返回错误
		return nil, err
	}

	// 创建结果map
	result := map[string]string{
		// 将API返回的版本信息添加到结果map中
		"version": data.Version,
		// 将API返回的服务器名称添加到结果map中
		"name":    data.ServerName,
		// 将API返回的世界指南添加到结果map中
		"worldguid":    data.Worldguid,
	}
	// 返回结果map和nil错误
	return result, nil
}

type ResponseMetrics struct {
	ServerFps        int     `json:"serverfps"`
	CurrentPlayerNum int     `json:"currentplayernum"`
	ServerFrameTime  float64 `json:"serverframetime"`
	MaxPlayerNum     int     `json:"maxplayernum"`
	Uptime           int     `json:"uptime"`
	Days             int     `json:"days"`
}

func Metrics() (map[string]interface{}, error) {
	// 调用API获取指标数据
	resp, err := callApi("GET", "/v1/api/metrics", nil)
	if err != nil {
		// 如果调用API出错，则返回错误
		return nil, err
	}

	// 定义用于存储响应数据的结构体
	var data ResponseMetrics
	// 将响应数据解析为结构体
	err = json.Unmarshal(resp, &data)
	if err != nil {
		// 如果解析数据出错，则返回错误
		return nil, err
	}

	// 构建结果映射
	result := map[string]interface{}{
		// 服务器帧率
		"server_fps":         data.ServerFps,
		// 当前玩家数量
		"current_player_num": data.CurrentPlayerNum,
		// 服务器帧时间
		"server_frame_time":  float64(int64(data.ServerFrameTime*100+0.5)) / 100,
		// 最大玩家数量
		"max_player_num":     data.MaxPlayerNum,
		// 服务器运行时间
		"uptime":             data.Uptime,
		// 服务器运行天数
		"days":               data.Days,
	}
	return result, nil
}

type ResponsePlayer struct {
	Name      string  `json:"name"`
	PlayerId  string  `json:"playerId"`
	UserId    string  `json:"userId"`
	Ip        string  `json:"ip"`
	Ping      float64 `json:"ping"`
	LocationX float64 `json:"location_x"`
	LocationY float64 `json:"location_y"`
	Level     int     `json:"level"`
}

type ResponsePlayers struct {
	Players []ResponsePlayer `json:"players"`
}

func ShowPlayers() ([]database.OnlinePlayer, error) {
	// 调用API获取玩家数据
	resp, err := callApi("GET", "/v1/api/players", nil)
	if err != nil {
		return nil, err
	}

	// 定义响应数据结构
	var data ResponsePlayers
	// 解析响应数据
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	// 初始化在线玩家列表
	onlinePlayers := make([]database.OnlinePlayer, 0)
	// 遍历玩家数据
	for _, player := range data.Players {
		// 创建在线玩家对象
		onlinePlayer := database.OnlinePlayer{
			PlayerUid:  getPlayerUid(player.PlayerId),
			SteamId:    getSteamId(player.UserId),
			Nickname:   player.Name,
			Ip:         player.Ip,
			Ping:       player.Ping,
			LocationX:  player.LocationX,
			LocationY:  player.LocationY,
			Level:      int32(player.Level),
			LastOnline: time.Now(),
		}
		// 将在线玩家添加到列表中
		onlinePlayers = append(onlinePlayers, onlinePlayer)
	}
	return onlinePlayers, nil
}

func getSteamId(userId string) string {
	// 如果userId不为空且以"steam_"开头
	if userId != "" && strings.HasPrefix(userId, "steam_") {
		// 去除userId的"steam_"前缀
		return strings.TrimPrefix(userId, "steam_")
	}
	// 如果不满足条件，返回空字符串
	return ""
}

func getPlayerUid(playerId string) string {
	// 如果playerId的长度小于8
	if len(playerId) < 8 {
		// 记录错误日志
		logger.Errorf("解析 PlayerId 失败: %s\n", playerId)
		// 返回空字符串
		return ""
	}
	// 截取playerId的前8个字符
	hexPart := playerId[:8]
	// 将截取的前8个字符从16进制转换为无符号32位整数
	decimalValue, err := strconv.ParseUint(hexPart, 16, 32)
	// 如果转换失败
	if err != nil {
		// 记录错误日志
		logger.Errorf("解析 PlayerId 失败: %s\n", err)
		// 返回空字符串
		return ""
	}
	// 将无符号32位整数转换为10进制字符串
	return strconv.FormatUint(decimalValue, 10)
}

type RequestUserId struct {
	UserId string `json:"userid"`
}

func KickPlayer(steamId string) error {
	// 将steamId封装成json格式的数据
	b, err := json.Marshal(RequestUserId{
		UserId: steamId,
	})
	if err != nil {
		// 如果出现错误，则返回错误
		return err
	}

	// 调用API接口
	_, err = callApi("POST", "/v1/api/kick", b)
	if err != nil {
		// 如果出现错误，则返回错误
		return err
	}

	return nil
}

func BanPlayer(steamId string) error {
	// 将steamId封装到RequestUserId结构体中，并序列化为JSON格式的字节切片
	b, err := json.Marshal(RequestUserId{
		UserId: steamId,
	})
	if err != nil {
		// 如果序列化失败，则返回错误
		return err
	}
	// 调用callApi函数发送POST请求到"/v1/api/ban"接口，并传入序列化后的数据
	_, err = callApi("POST", "/v1/api/ban", b)
	if err != nil {
		// 如果请求失败，则返回错误
		return err
	}
	return nil
}

func UnBanPlayer(steamId string) error {
	// 将steamId封装到RequestUserId结构体中，并序列化为JSON格式
	b, err := json.Marshal(RequestUserId{
		UserId: steamId,
	})
	if err != nil {
		// 如果序列化失败，则返回错误
		return err
	}

	// 调用API进行POST请求，解禁玩家
	_, err = callApi("POST", "/v1/api/unban", b)
	if err != nil {
		// 如果调用API失败，则返回错误
		return err
	}

	return nil
}

type RequestBroadcast struct {
	Message string `json:"message"`
}

func Broadcast(message string) error {
	// 将RequestBroadcast结构体序列化为JSON格式的字节数组
	b, err := json.Marshal(RequestBroadcast{
		Message: message,
	})
	if err != nil {
		// 如果序列化失败，则返回错误
		return err
	}

	// 调用API接口发送POST请求
	_, err = callApi("POST", "/v1/api/announce", b)
	if err != nil {
		// 如果请求发送失败，则返回错误
		return err
	}

	// 如果没有错误发生，则返回nil
	return nil
}

type RequestShutdown struct {
	Waittime int    `json:"waittime"`
	Message  string `json:"message"`
}

func Shutdown(seconds int, message string) error {
	// 将RequestShutdown结构体序列化为JSON格式
	b, err := json.Marshal(RequestShutdown{
		Waittime: seconds,
		Message:  message,
	})
	if err != nil {
		// 如果序列化出错，则返回错误
		return err
	}

	// 调用API接口，发送POST请求
	_, err = callApi("POST", "/v1/api/shutdown", b)
	if err != nil {
		// 如果API调用出错，则返回错误
		return err
	}

	// 返回nil，表示成功
	return nil
}

func DoExit() error {
	// 调用API，使用POST方法向"/v1/api/stop"路径发送请求，携带的数据为nil
	_, err := callApi("POST", "/v1/api/stop", nil)
	if err != nil {
		// 如果调用API时发生错误，则返回错误
		return err
	}
	return nil
}
