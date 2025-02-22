import Service from "./service";

class ApiService extends Service {
  async login(param) {
    // 接收传入的参数
    let data = param;
    // 发起POST请求，请求地址为/api/login，请求体为data
    return this.fetch(`/api/login`).post(data).json();
  }

  async getServerToolInfo() {
    // 调用fetch函数，传入API地址
    return this.fetch(`/api/server/tool`)
      // 发起GET请求
      .get()
      // 解析返回的数据为JSON格式
      .json();
  }
  async getServerInfo() {
    // 使用fetch API向服务器发送GET请求
    return this.fetch(`/api/server`).get().json();
  }
  async getServerMetrics() {
    // 调用fetch方法获取/api/server/metrics接口的数据
    return this.fetch(`/api/server/metrics`)
      .get() // 使用get方法发送请求
      .json(); // 将响应体解析为JSON格式
  }
  async sendBroadcast(param) {
    // 将传入的参数赋值给变量 data
    let data = param;

    // 使用 fetch 方法向服务器发送 POST 请求
    // 请求的路径为 /api/server/broadcast
    // 将变量 data 作为请求体发送
    // 并返回请求的响应体，解析为 JSON 格式
    return this.fetch(`/api/server/broadcast`).post(data).json();
  }
  async shutdownServer(param) {
    // 将传入的参数赋值给局部变量data
    let data = param;

    // 使用fetch方法向指定的API地址发送POST请求，并传递data作为请求体
    // 然后调用.json()方法将响应体解析为JSON格式
    return this.fetch(`/api/server/shutdown`).post(data).json();
  }

  async getPlayerList(param) {
    // 生成查询字符串
    const query = this.generateQuery(param);

    // 发送HTTP GET请求到指定的API端点
    // `/api/player?` 后面拼接查询字符串
    return this.fetch(`/api/player?${query}`).get().json();
  }
  async getOnlinePlayerList() {
    // 调用fetch方法获取在线玩家列表
    return this.fetch(`/api/online_player`)
      .get()  // 发送GET请求
      .json();  // 解析JSON格式的响应数据
  }
  async getPlayer(param) {
    // 解构参数对象，获取 playerUid
    const { playerUid } = param;
    // 使用 fetch 方法请求 `/api/player/${playerUid}`
    return this.fetch(`/api/player/${playerUid}`).get().json();
  }
  async kickPlayer(param) {
    // 从参数中获取玩家UID
    const { playerUid } = param;

    // 调用fetch方法，向指定的URL发送POST请求
    // URL格式为 /api/player/{playerUid}/kick
    return this.fetch(`/api/player/${playerUid}/kick`).post().json();
  }
  async banPlayer(param) {
    // 从参数中解构出 playerUid
    const { playerUid } = param;

    // 发送 POST 请求到 /api/player/${playerUid}/ban
    return this.fetch(`/api/player/${playerUid}/ban`).post().json();
  }
  async unbanPlayer(param) {
    // 从参数中解构出 playerUid
    const { playerUid } = param;

    // 发起 POST 请求到 '/api/player/${playerUid}/unban'
    return this.fetch(`/api/player/${playerUid}/unban`).post().json();
  }

  async getGuildList() {
    // 调用 fetch 方法，向 '/api/guild' 发送 GET 请求
    return this.fetch(`/api/guild`).get().json();
  }
  async getGuild(param) {
    // 解构参数对象，获取 adminPlayerUid
    const { adminPlayerUid } = param;

    // 调用 fetch 方法，传入 API 路径
    return this.fetch(`/api/guild/${adminPlayerUid}`).get().json();
    // 获取数据并返回 JSON 格式的结果
  }

  async getWhitelist() {
    // 发起GET请求获取白名单数据
    return this.fetch(`/api/whitelist`).get().json();
  }

  async addWhitelist(param) {
    // 将传入的参数赋值给变量data
    let data = param;

    // 调用fetch方法获取/api/whitelist的URL，然后使用post方法发送data数据，并将结果转换为json格式
    return this.fetch(`/api/whitelist`).post(data).json();
  }

  async removeWhitelist(param) {
    // 将传入的参数赋值给变量 data
    let data = param;

    // 使用 fetch 发起请求，请求路径为 '/api/whitelist'
    // 然后调用 delete 方法，传入 data 参数，最后调用 json 方法解析响应数据
    return this.fetch(`/api/whitelist`).delete(data).json();
  }

  async putWhitelist(param) {
    // 接收传入的参数
    let data = param;

    // 使用fetch函数向/api/whitelist路径发送PUT请求，并传递data作为请求体
    // 然后将响应结果转换为JSON格式
    return this.fetch(`/api/whitelist`).put(data).json();
  }

  async getRconCommands() {
    // 调用 fetch 方法，获取 RCON 命令
    // 访问路径为 /api/rcon
    return this.fetch(`/api/rcon`).get().json();
  }

  async sendRconCommand(param) {
    // 将传入的参数赋值给data变量
    let data = param;
    // 调用fetch函数，向/api/rcon/send路径发送POST请求，请求体为data
    // 并使用.post()方法设置请求方法为POST
    // 使用.json()方法将响应解析为JSON格式
    return this.fetch(`/api/rcon/send`).post(data).json();
  }

  async addRconCommand(param) {
    // 将传入的参数赋值给变量data
    let data = param;
    // 调用fetch函数，访问'/api/rcon'路径，并发送POST请求，请求体为data
    // 然后将响应体转换为JSON格式
    return this.fetch(`/api/rcon`).post(data).json();
  }

  async putRconCommand(uuid, param) {
    // 将参数赋值给data变量
    let data = param;
    // 调用fetch方法，传入指定的URL和参数，然后使用put方法发送数据，最后调用json方法解析响应
    return this.fetch(`/api/rcon/${uuid}`).put(data).json();
  }

  async removeRconCommand(uuid) {
    // 发送HTTP DELETE请求到指定的API端点
    return this.fetch(`/api/rcon/${uuid}`).delete().json(); // 解析返回的JSON数据
  }

  async getBackupList(param) {
    // 生成查询字符串
    const query = this.generateQuery(param);
    // 发送HTTP GET请求到指定的API端点
    // 并解析返回的JSON数据
    return this.fetch(`/api/backup?${query}`).get().json();
  }

  async removeBackup(uuid) {
    // 发送删除请求到指定的API路径
    // 例如：'/api/backup/12345'
    return this.fetch(`/api/backup/${uuid}`).delete().json();
  }

  async downloadBackup(uuid) {
    // 调用 fetch 方法获取备份数据
    return this.fetch(`/api/backup/${uuid}`).get().blob();
    // 获取数据并作为 Blob 对象返回
  }
}

export default ApiService;