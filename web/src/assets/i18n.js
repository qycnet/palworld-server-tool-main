import { createI18n } from "vue-i18n";

const messages = {
  en: {
    title: "PalWorld Server Tool",
    modal: {
      customPal: "Custom Palu",
      auth: "Auth",
      broadcast: "Publish Broadcast",
      whitelist: "Whitelist",
      addWhitelist: "Add Whitelist",
      rcon: "Customize RCON",
      backup: "Backup File",
    },
    status: {
      online: "Online",
      offline: "Offline",
      last_online: "Last Online",
      authenticated: "Authenticated",
      master: "Master",
      player_number: "Player: {number}",
      online_number: "Online: {number}",
      whitelist: "Whitelist",
    },
    message: {
      customPalDesc: "Go to Custom Palu",
      warn: "Warning",
      success: "Success",
      fail: "Fail",
      changelanguage: "切换语言中...",
      authdesc: "You need to log in before you can operate.",
      copyfail: "Your browser does not support this feature!",
      copysuccess: "Copy success!",
      copyerr: "Failed to copy: {err}",
      copyempty: "Copy content is empty!",
      autherr: "Password error!",
      authsuccess: "Auth success!",
      requireauth: "Please auth first!",
      bantitle: "Ban Player",
      banwarn: "Are you sure to ban this player?",
      bansuccess: "Ban success!",
      banfail: "Ban fail: {err}",
      unbantitle: "Unban Player",
      unbanwarn: "Are you sure to unban this player?",
      unbansuccess: "Unban success!",
      unbanfail: "Unban fail: {err}",
      kicktitle: "Kick Player",
      kickwarn: "Are you sure to kick this player?",
      kicksuccess: "Kick success!",
      kickfail: "Kick fail: {err}",
      broadcastsuccess: "Broadcast success!",
      shutdownsuccess: "Shutdown success!",
      broadcastfail: "Broadcast fail: {err}",
      shutdownfail: "Shutdown fail: {err}",
      shutdowntip:
        "This will shut down the server after 60 seconds and send a broadcast.",
      selectVerify: "Select verification method",
      addwhitesuccess: "Add whitelist success!",
      addwhitefail: "Add whitelist fail: {err}",
      removewhitesuccess: "Remove whitelist success!",
      removewhitefail: "Remove whitelist fail: {err}",
      rconfail: "RCON execution failed: {err}",
      addrconsuccess: "Add RCON command success!",
      addrconfail: "Add RCON command fail: {err}",
      removerconsuccess: "Remove RCON command success!",
      removerconfail: "Remove RCON command fail: {err}",
      importRconTitle: "Click or drag file containing RCON commands to upload",
      importRconDesc:
        "Make sure your file content is in the format of command,remark,placeholder per line",
      importRconSuccess: "Import RCON commands success!",
      importRconFail: "Import RCON commands fail: {err}",
      loading: "Loading...",
      removebackupsuccess: "Remove backup success!",
      removebackupfail: "Remove backup fail: {err}",
      downloadsuccess: "Download success!",
    },
    button: {
      customPal: "Custom Pal",
      auth: "Admin Mode",
      players: "Players",
      guilds: "Guilds",
      broadcast: "Broadcast",
      shutdown: "Shutdown",
      rcon: "RCON",
      execute: "Execute",
      addRcon: "Add new Command",
      ban: "Ban",
      unban: "Unban",
      kick: "Kick",
      detail: "Detail",
      confirm: "Confirm",
      cancel: "Cancel",
      close: "Close",
      remove: "Remove",
      import: "Import",
      add: "Add",
      viewGuild: "View Guild",
      viewPlayer: "View Player",
      whitelist: "Whitelist",
      joinWhitelist: "Join Whitelist",
      removeWhitelist: "Remove Whitelist",
      search: "Search",
      addNew: "Add New",
      save: "Save",
      palconf: "PalConf",
      controlCenter: "Control Center",
      download: "Download",
      backup: "Backup",
      copysid: "Copy Steam64",
      copypid: "Copy Player ID",
      copyitem: "Copy Item ID",
      copypal: "Copy Pal ID",
      map: "Map",
    },
    pal: {
      type: "Type",
      level: "Level",
      skills: "Skills",
      ranged: "ATK(potential)",
      defense: "Defense(potential)",
      melee: "Blood(potential)",
      rank: "Rank",
      tower: "Tower",
      lucky: "Lucky",
      rank_boost: "Rank Boost",
      rank_attack: "Attack Boost",
      rank_defence: "Defense Boost",
      rank_craftspeed: "Work Speed Boost",
    },
    item: {
      palList: "Pal List",
      itemList: "Item List",
      commonContainer: "Common Container",
      essentialContainer: "Essential Container",
      weaponContainer: "Weapon Container",
      armorContainer: "Armor Container",
      name: "Name",
      description: "Description",
      count: "Count",
      time: "Time",
      serverFps: "ServerFPS",
      serverUptime: "ServerUptime",
      serverDays: "ServerDays",
      serverFrameTime: "ServerFrameTime",
      maxPlayerNum: "maxPlayerNum",
    },
    input: {
      searchPlaceholder: "Search for pal type and skills",
      nickname: "Nickname",
      player_uid: "Player UID",
      steam_id: "Steam64",
      rcon: "RCON Command",
      placeholder: "RCON Placeholder",
      remark: "Remark",
      selectPlayer: "Select Player",
      selectItem: "Select Item",
      selectPal: "Select Pal",
    },
    map: {
      baseCampTitle: 'The base of the "{name}" guild',
      guildMember: "Member: ",
      showFastTravel: "Show fast travel",
      showBossTower: "Show boss tower",
      showPlayer: "Show online player",
      showBaseCamp: "Show basecamp",
    },
  },
  zh: {
    title: "幻兽帕鲁服务器工具",
    modal: {
      customPal: "定制帕鲁",
      auth: "管理认证",
      broadcast: "发布游戏内广播",
      whitelist: "白名单管理",
      addWhitelist: "添加白名单",
      rcon: "自定义RCON命令",
      backup: "备份存档管理",
    },
    status: {
      online: "在线",
      offline: "离线",
      last_online: "最后在线",
      authenticated: "已进入管理模式",
      master: "会长",
      player_number: "玩家: {number}",
      online_number: "在线: {number}",
      whitelist: "白名单",
    },
    message: {
      customPalDesc: "前往定制帕鲁",
      warn: "警告",
      success: "成功",
      fail: "失败",
      changelanguage: "Changing language...",
      authdesc: "您需要先进入管理模式才能操作。",
      copyfail: "您的浏览器不支持此功能!",
      copysuccess: "复制成功!",
      copyerr: "复制失败: {err}",
      copyempty: "复制内容为空!",
      autherr: "密码错误!",
      authsuccess: "认证成功!",
      requireauth: "请先进入管理模式!",
      bantitle: "封禁玩家",
      banwarn: "您确定要封禁此玩家吗?",
      bansuccess: "封禁成功!",
      banfail: "封禁失败: {err}",
      unbantitle: "解封玩家",
      unbanwarn: "您确定要解封此玩家吗?",
      unbansuccess: "解封成功!",
      unbanfail: "解封失败: {err}",
      kicktitle: "踢出玩家",
      kickwarn: "您确定要踢出此玩家吗?",
      kicksuccess: "踢出成功!",
      kickfail: "踢出失败: {err}",
      broadcastsuccess: "广播成功!",
      shutdownsuccess: "关闭成功!",
      broadcastfail: "广播失败: {err}",
      shutdownfail: "关闭失败: {err}",
	  quitpstsuccess: "退出成功!",
	  quitpstfail: "退出失败: {err}",
      shutdowntip: "此操作将在60秒后关闭服务器并发送广播。",
	  Signoutprompst: "此操作将退出登录pst。",
      selectVerify: "请选择验证方式",
      addwhitesuccess: "添加白名单成功!",
      addwhitefail: "添加白名单失败: {err}",
      removewhitesuccess: "移除白名单成功!",
      removewhitefail: "移除白名单失败: {err}",
      rconfail: "RCON执行失败: {err}",
      addrconsuccess: "添加RCON命令成功!",
      addrconfail: "添加RCON命令失败: {err}",
      removerconsuccess: "移除RCON命令成功!",
      removerconfail: "移除RCON命令失败: {err}",
      importRconTitle: "点击或者拖动包含RCON命令的文件上传",
      importRconDesc:
        "请确保你的文件内容每一行是 command,remark,placeholder 的格式",
      importRconSuccess: "导入RCON命令成功!",
      importRconFail: "导入RCON命令失败: {err}",
      loading: "加载中...",
      removebackupsuccess: "删除备份成功!",
      removebackupfail: "删除备份失败: {err}",
      downloadsuccess: "下载成功!",
    },
    button: {
      customPal: "定制帕鲁",
      auth: "管理模式",
      players: "玩家",
      guilds: "公会",
      broadcast: "游戏内广播",
      shutdown: "关闭服务器",
	  quitpst: "退出管理员登录",
      rcon: "RCON命令",
      execute: "执行",
      addRcon: "新增命令",
      ban: "封禁",
      unban: "解封",
      kick: "踢出",
      detail: "详情",
      confirm: "确认",
      cancel: "取消",
      close: "关闭",
      remove: "删除",
      import: "导入",
      add: "新增",
      viewGuild: "查看公会",
      viewPlayer: "查看玩家",
      whitelist: "白名单管理",
      joinWhitelist: "加入白名单",
      removeWhitelist: " 移除白名单",
      search: "搜索",
      addNew: "新增",
      save: "保存",
      palconf: "配置生成器",
      controlCenter: "控制中心",
      download: "下载",
      backup: "存档管理",
      copysid: "复制Steam64",
      copypid: "复制玩家ID",
      copyitem: "复制物品ID",
      copypal: "复制帕鲁ID",
      map: "地图",
    },
    pal: {
      type: "类型",
      level: "等级",
      skills: "被动词条",
      ranged: "攻击(潜力)",
      defense: "防御(潜力）",
      melee: "血量(潜力)",
      rank: "星级",
      tower: "塔主",
      lucky: "稀有",
      rank_boost: "强化属性",
      rank_attack: "攻击强化",
      rank_defence: "防御强化",
      rank_craftspeed: "作业速度强化",
    },
    item: {
      palList: "幻兽列表",
      itemList: "物品列表",
      commonContainer: "物品栏",
      essentialContainer: "重要物品",
      weaponContainer: "武器栏",
      armorContainer: "防具栏",
      name: "名字",
      description: "描述",
      count: "数量",
      time: "时间",
      serverFps: "服务器FPS",
      serverUptime: "服务器运行时间",
      serverDays: "服务器天数",
      serverFrameTime: "服务器帧时间",
      maxPlayerNum: "最大玩家数",
    },
    input: {
      searchPlaceholder: "搜索帕鲁类型、技能",
      nickname: "昵称",
      player_uid: "玩家 UID",
      steam_id: "Steam64",
      rcon: "RCON命令",
      placeholder: "占位描述(如玩家ID+物品ID)",
      remark: "备注",
      selectPlayer: "选择玩家",
      selectItem: "选择物品",
      selectPal: "选择帕鲁",
    },
    map: {
      baseCampTitle: '"{name}" 公会的据点',
      guildMember: "成员: ",
      showFastTravel: "显示传送点",
      showBossTower: "显示高塔",
      showPlayer: "显示在线玩家",
      showBaseCamp: "显示据点",
    },
  },
  ja: {
    title: "パルワールドサーバーツール",
    modal: {
      customPal: "カスタムパル",
      auth: "ログイン",
      broadcast: "全体メッセージ配信",
      whitelist: "ホワイトリスト",
      addWhitelist: "ホワイトリストに追加",
      rcon: "カスタムRCONです",
      backup: "バックアップファイル",
    },
    status: {
      online: "オンライン",
      offline: "オフライン",
      last_online: "最終オンライン",
      authenticated: "認証済み",
      master: "ギルドマスター",
      player_number: "プレイヤー: {number}",
      online_number: "オンライン: {number}",
    },
    message: {
      customPalDesc: "カスタマイズされたパルに行く",
      warn: "警告",
      success: "成功",
      fail: "失敗",
      changelanguage: "Changing language...",
      authdesc: "この操作を行う前にログインしてください",
      copyfail: "お使いのブラウザーではこの機能は使えません",
      copysuccess: "コピーに成功しました！",
      copyerr: "コピーに失敗しました: {err}",
      copyempty: "コピー内容が空です！",
      autherr: "パスワードが正しくありません！",
      authsuccess: "ログインに成功しました！",
      requireauth: "先にログインをしてください！",
      bantitle: "プレイヤーをBANする",
      banwarn: "このプレイヤーを本当にBANしますか？",
      bansuccess: "BANに成功しました！",
      banfail: "BANに失敗しました: {err}",
      unbantitle: "プレイヤーをUNBANする",
      unbanwarn: "このプレイヤーを本当にUNBANしますか？",
      unbansuccess: "UNBANに成功しました！",
      unbanfail: "UNBANに失敗しました: {err}",
      kicktitle: "プレイヤーをキックする",
      kickwarn: "本当にこのプレイヤーをキックしますか？",
      kicksuccess: "キックに成功しました！",
      kickfail: "キックに失敗しました: {err}",
      broadcastsuccess: "全体メッセージを配信しました！",
      shutdownsuccess: "シャットダウンに成功しました！",
      broadcastfail: "全体メッセージ配信に失敗しました: {err}",
      shutdownfail: "シャットダウンに失敗しました: {err}",
      shutdowntip:
        "全体メッセージを配信し、60秒後にサーバーをシャットダウンします",
      selectVerify: "認証方法を選択してください",
      addwhitesuccess: "ホワイトリストに追加しました！",
      addwhitefail: "ホワイトリスト追加に失敗しました: {err}",
      removewhitesuccess: "ホワイトリストから削除しました",
      removewhitefail: "ホワイトリスト削除に失敗しました: {err}",
      rconfail: "RCONの実行に失敗しました: {err}",
      addrconsuccess: "RCONコマンドを追加しました！",
      addrconfail: "RCONコマンド追加に失敗しました: {err}",
      removerconsuccess: "RCONコマンドを削除しました！",
      removerconfail: "RCONコマンド削除に失敗しました: {err}",
      importRconTitle:
        "RCONコマンドを含むファイルをクリックまたはドラッグしてアップロードしてください",
      importRconDesc:
        "ファイルの内容が1行につき command,remark,placeholder の形式であることを確認してください",
      importRconSuccess: "RCONコマンドを導入しました！",
      importRconFail: "RCONコマンド導入に失敗しました: {err}",
      loading: "読み込み中...",
      removebackupsuccess: "バックアップを削除しました",
      removebackupfail: "バックアップ削除に失敗しました: {err}",
      downloadsuccess: "ダウンロードに成功しました！",
    },
    button: {
      customPal: "カスタムパル",
      auth: "管理者モード",
      players: "プレイヤー",
      guilds: "ギルド",
      broadcast: "全体メッセージ",
      shutdown: "シャットダウン",
      rcon: "RCON",
      addRcon: "新しいコマンドを追加",
      execute: "実行",
      ban: "BANする",
      unban: "UNBANする",
      kick: "キックする",
      detail: "詳細",
      confirm: "確定",
      cancel: "キャンセル",
      close: "閉じる",
      remove: "削除",
      import: "導入",
      add: "追加",
      viewGuild: "ギルドを見る",
      viewPlayer: "プレイヤーを見る",
      whitelist: "ホワイトリスト",
      joinWhitelist: "ホワイトリストに追加",
      removeWhitelist: "ホワイトリストを削除します",
      search: "検索",
      palconf: "パルの設定",
      controlCenter: "コントロールセンター",
      download: "ダウンロード",
      backup: "バックアップ",
      copysid: "Steam64をコピー",
      copypid: "プレイヤーIDをコピー",
      copyitem: "アイテムIDをコピー",
      copypal: "パルIDをコピー",
      map: "地図",
    },
    item: {
      palList: "幻獣リスト",
      itemList: "アイテムリスト",
      commonContainer: "一般アイテム",
      essentialContainer: "重要アイテム",
      weaponContainer: "武器アイテム",
      armorContainer: "防具アイテム",
      name: "名前",
      description: "説明",
      count: "数",
      time: "時間",
      serverFps: "サーバーFPS",
      serverUptime: "サーバー稼働時間",
      serverDays: "サーバー天数",
      serverFrameTime: "サーバーフレーム時間",
      maxPlayerNum: "最大プレイヤー数",
    },
    input: {
      searchPlaceholder: "パルのタイプやスキルで検索してください",
      nickname: "ニックネーム",
      player_uid: "プレイヤー UID",
      steam_id: "Steam64",
      rcon: "RCONコマンド",
      placeholder: "プレースホルダー",
      remark: "備考",
      selectPlayer: "プレイヤーを選択",
      selectItem: "アイテムを選択",
      selectPal: "パルを選択",
    },
    pal: {
      type: "タイプ",
      level: "レベル",
      skills: "スキル",
      ranged: "攻撃(潜在力)です",
      defense: "防御(潜在力)です",
      melee: "血の量(ポテンシャル)です",
      rank: "ランク",
      tower: "タワーボス",
      lucky: "ラッキーパル",
      rank_boost: "強化ステータス",
      rank_attack: "攻撃力強化",
      rank_defence: "防御力強化",
      rank_craftspeed: "作業速度強化",
    },
    map: {
      baseCampTitle: "「{name}」ギルドの基地",
      guildMember: "メンバー: ",
      showFastTravel: "Show fast travel",
      showBossTower: "Show boss tower",
      showPlayer: "Show online player",
      showBaseCamp: "Show basecamp",
    },
  },
};

const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackLocale: "zh",
  messages,
});

export default i18n;
