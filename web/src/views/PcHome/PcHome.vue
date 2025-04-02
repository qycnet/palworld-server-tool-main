<script setup>
// 从@vicons/material导入图标组件
import {
  AdminPanelSettingsOutlined,
  SupervisedUserCircleRound,
  SettingsPowerRound,
  DeleteOutlineTwotone,
  RemoveRedEyeTwotone,
  DeleteFilled,
  ArchiveOutlined,
  CloudDownloadOutlined,
  PublicRound,
} from "@vicons/material";
// 从@vicons/ionicons5导入图标组件
import {
  GameController,
  LanguageSharp,
  ShieldCheckmarkSharp,
  Terminal,
  ArchiveOutline,
  Settings,
} from "@vicons/ionicons5";
// 从@vicons/carbon导入图标组件
import { GuiManagement } from "@vicons/carbon";
// 从@vicons/fa导入图标组件
import { BroadcastTower } from "@vicons/fa";
// 从vue导入响应式API和生命周期钩子
import { computed, onMounted, ref, watch, h } from "vue";

// 从naive-ui导入UI组件
import { NTag, NButton, NIcon, useMessage, useDialog } from "naive-ui";

// 从vue-i18n导入国际化API
import { useI18n } from "vue-i18n";

// 导入自定义的API服务
import ApiService from "@/service/api";

// 导入Vuex或Pinia状态管理库的store
import pageStore from "@/stores/model/page.js";
import whitelistStore from "@/stores/model/whitelist";
import playerToGuildStore from "@/stores/model/playerToGuild";
import userStore from "@/stores/model/user";

// 导入第三方库
import dayjs from "dayjs";

// 导入静态资源
import palMap from "@/assets/pal.json";
import itemMap from "@/assets/items.json";
import skillMap from "@/assets/skill.json";

// 导入Vue组件
import PlayerList from "./component/PlayerList.vue";
import GuildList from "./component/GuildList.vue";
import MapView from "./component/MapView.vue";
// 使用国际化API获取当前的语言环境和翻译函数
const { t, locale } = useI18n();

// 使用naive-ui的消息和对话框API
const message = useMessage();
const dialog = useDialog();

// 定义全局常量
const PALWORLD_TOKEN = "palworld_token";

// 计算属性：获取页面宽度
const pageWidth = computed(() => pageStore().getScreenWidth());

// 计算属性：判断是否为小屏幕
const smallScreen = computed(() => pageWidth.value < 1024);

// 响应式数据：加载状态
const loading = ref(false);

// 响应式数据：服务器信息
const serverInfo = ref({});

// 响应式数据：服务器指标
const serverMetrics = ref({});

// 响应式数据：当前显示的内容 玩家/工会/地图（players或guilds或map）
const currentDisplay = ref("players");

// 响应式数据：玩家列表
const playerList = ref([]);

// 响应式数据：在线玩家列表
const onlinePlayerList = ref([]);

// 响应式数据：工会列表
const guildList = ref([]);

// 响应式数据：玩家信息
const playerInfo = ref({});

// 响应式数据：玩家的PAL列表
const playerPalsList = ref([]);

// 响应式数据：工会信息
const guildInfo = ref({});

// 响应式数据：技能类型列表
const skillTypeList = ref([]);

// 响应式数据：语言选项列表
const languageOptions = ref([]);

// 响应式数据：登录状态
const isLogin = ref(false);

// 响应式数据：认证令牌
const authToken = ref("");

// 响应式数据：是否开启暗色模式
const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches
);

// 更新暗色模式的函数
const updateDarkMode = (e) => {
  isDarkMode.value = e.matches;
};

// 获取暗色模式颜色的函数
const getDarkModeColor = () => {
  return isDarkMode.value ? "#fff" : "#000";
};

// 获取用户头像URL的函数
const getUserAvatar = () => {
  return new URL("../../assets/avatar.webp", import.meta.url).href;
};

const handleSelectLanguage = (key) => {
  message.info(t("message.changelanguage"));
  if (key === "zh") {
    localStorage.setItem("locale", "zh");
    // locale.value = "zh";
  } else if (key === "ja") {
    localStorage.setItem("locale", "ja");
    // locale.value = "ja";
  } else {
    localStorage.setItem("locale", "en");
    // locale.value = "en";
  }
  setTimeout(() => {
    location.reload();
  }, 1000);
};

const getSkillTypeList = () => {
  if (skillMap[locale.value]) {
    return Object.values(skillMap[locale.value]).map((item) => item.name);
  } else {
    return [];
  }
};

const toPalConf = () => {
  window.open("/pal-conf");
};

const toGithub = () => {
  window.open("https://github.com/qycnet/palworld-server-tool-main/releases");
};
const serverToolInfo = ref({});
const hasNewVersion = ref(false);
const getServerToolInfo = async () => {
  const { data } = await new ApiService().getServerToolInfo();
  serverToolInfo.value = data.value;
  if (data.value) {
    hasNewVersion.value = isNewVersion(data.value?.version, data.value?.latest);
  }
};
const isNewVersion = (version, latest) => {
  if (version == "Unknown" || version == "Develop" || latest == "") {
    return false;
  }
  const currentVersion = version.split("v")[1];
  const latestVersion = latest?.split("v")[1];
  const currentParts = currentVersion.substring(1).split(".");
  const latestParts = latestVersion.substring(1).split(".");
  for (let i = 0; i < currentParts.length; i++) {
    const currentPart = parseInt(currentParts[i], 10);
    const latestPart = parseInt(latestParts[i], 10);
    if (latestPart > currentPart) {
      return true;
    } else if (latestPart < currentPart) {
      return false;
    }
  }
  return false;
};

// get data
const getServerInfo = async () => {
  const { data } = await new ApiService().getServerInfo();
  serverInfo.value = data.value;
};

const getServerMetrics = async () => {
  const { data } = await new ApiService().getServerMetrics();
  serverMetrics.value = data.value;
};

const getPlayerList = async () => {
  getOnlineList();
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  playerList.value = data.value;
  rconPlayerOptions.value = data.value.map((item) => {
    return {
      label: `${item.nickname}(${item.player_uid})(${item.steam_id})`,
      value: `${item.player_uid}-${item.steam_id}`,
    };
  });
};
const getOnlineList = async () => {
  const { data } = await new ApiService().getOnlinePlayerList();
  onlinePlayerList.value = data.value;
};

// login
const showLoginModal = ref(false);
const password = ref("");
const handleLogin = async () => {
  const { data, statusCode } = await new ApiService().login({
    password: password.value,
  });
  if (statusCode.value === 401) {
    message.error(t("message.autherr"));
    password.value = "";
    return;
  }
  let token = data.value.token;
  localStorage.setItem(PALWORLD_TOKEN, token);
  userStore().setIsLogin(true, token);
  await getWhiteList();
  authToken.value = token;
  message.success(t("message.authsuccess"));
  showLoginModal.value = false;
  isLogin.value = true;
};
const showRconDrawer = ref(false);
const rconCommands = ref([]);
const rconSelectedPlayer = ref(null);
const rconPlayerOptions = ref([]);
const rconSelectedItem = ref(null);
const rconItemOptions = ref([]);
const rconSelectedPal = ref(null);
const rconPalOptions = ref([]);
const rconCommandsExtra = ref({});
const showCustomPalModal = ref(false);

const openCustomPalModal = () => {
  showCustomPalModal.value = true;
};

const closeCustomPalModal = () => {
  showCustomPalModal.value = false;
};

const goToCustomPalWebsite = () => {
  window.open("http://givepal_j.apiqy.cn/", "_blank");
};

// 自定义帕鲁模态框
const customPalModal = () => {
  dialog.info({
    title: t("modal.customPal"),
    content: t("message.customPalDesc"),
    positiveText: t("button.confirm"),
    negativeText: t("button.cancel"),
    onPositiveClick: () => {
      goToCustomPalWebsite();
    },
    onNegativeClick: () => {
      // 用户取消操作
    }
  });
};
const copyText = async (text) => {
  if (text == "" || text == null) {
    message.error(t("message.copyempty"));
    return;
  }
  if (navigator.clipboard) {
    try {
      await navigator.clipboard.writeText(text);
      message.success(t("message.copysuccess"));
    } catch (err) {
      message.error(t("message.copyerr", { err }));
    }
  } else {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    document.body.appendChild(textarea);
    textarea.select();
    try {
      document.execCommand("copy");
      message.success(t("message.copysuccess"));
    } catch (err) {
      message.error(t("message.copyerr", { err }));
    }
    document.body.removeChild(textarea);
  }
};
const handleRconDrawer = () => {
  if (checkAuthToken()) {
    showRconDrawer.value = true;
    getRconCommands();
  } else {
    message.error(t("message.requireauth"));
    showRconDrawer.value = true;
  }
};
const getRconCommands = async () => {
  if (checkAuthToken()) {
    const { data, statusCode } = await new ApiService().getRconCommands();
    if (statusCode.value === 200) {
      rconCommands.value = data.value;
      rconCommands.value.forEach((item) => {
        rconCommandsExtra.value[item.uuid] = "";
      });
    }
  }
};
const sendRconCommand = async (uuid) => {
  const content = rconCommandsExtra.value[uuid];
  if (checkAuthToken()) {
    const { data, statusCode } = await new ApiService().sendRconCommand({
      uuid,
      content,
    });
    if (statusCode.value === 200) {
      message.info(data.value?.message);
    } else {
      message.error(t("message.rconfail", { err: data.value?.error }));
    }
  }
};

const showRconAddModal = ref(false);
const newRconCommand = ref("");
const newRconPlaceholder = ref("");
const newRconRemark = ref("");
const handleAddRconCommand = () => {
  showRconAddModal.value = true;
  newRconCommand.value = "";
  newRconPlaceholder.value = "";
  newRconRemark.value = "";
};
const handleImportRconFinish = (options) => {
  getRconCommands();
  setTimeout(() => {
    message.success(t("message.importRconSuccess"));
    showRconAddModal.value = false;
  }, 500);
};
const handleImportRconError = (options) => {
  let err = options.event?.target?.response
    ? JSON.parse(options.event?.target?.response).error
    : "";
  message.error(t("message.importRconFail", { err }));
};
const addRconCommand = async () => {
  if (checkAuthToken()) {
    const { data, statusCode } = await new ApiService().addRconCommand({
      command: newRconCommand.value,
      placeholder: newRconPlaceholder.value,
      remark: newRconRemark.value,
    });
    if (statusCode.value === 200) {
      message.success(t("message.addrconsuccess"));
      await getRconCommands();
      newRconCommand.value = "";
      newRconPlaceholder.value = "";
      newRconRemark.value = "";
    } else {
      message.error(t("message.addrconfail", { err: data.value?.error }));
    }
  }
};
const removeRconCommand = async (uuid) => {
  if (checkAuthToken()) {
    const { data, statusCode } = await new ApiService().removeRconCommand(uuid);
    if (statusCode.value === 200) {
      message.success(t("message.removerconsuccess"));
      await getRconCommands();
    } else {
      message.error(t("message.removerconfail", { err: data.value?.error }));
    }
  }
};

// 控制中心（下拉菜单）
// 包含：白名单管理、RCON 命令、游戏内广播、关闭服务器、退出管理员登录
const renderIcon = (icon, color = "#666") => {
  return () => {
    return h(
      NIcon,
      {
        color: color,
      },
      {
        default: () => h(icon),
      }
    );
  };
};
const controlCenterOption = [
  // {
  //   label: () => {
  //     return h("div", null, {
  //       default: () => t("button.backup"),
  //     });
  //   },
  //   key: "backup",
  //   icon: renderIcon(ArchiveOutlined),
  // },
  {
    label: () => {
      return h("div", null, {
        default: () => t("button.palconf"),
      });
    },
    key: "palconf",
    icon: renderIcon(Settings),
  },
  {
    label: () => {
      return h("div", null, {
        default: () => t("button.whitelist"),
      });
    },
    key: "whitelist",
    icon: renderIcon(ShieldCheckmarkSharp),
  },
  // {
  //   label: () => {
  //     return h("div", null, {
  //       default: () => t("button.rcon"),
  //     });
  //   },
  //   key: "rcon",
  //   icon: renderIcon(Terminal),
  // },
  {
    label: () => {
      return h("div", null, {
        default: () => t("button.broadcast"),
      });
    },
    key: "broadcast",
    icon: renderIcon(BroadcastTower),
  },
  {
    label: () => {
      return h(
        "div",
        {
          style: { color: "#cc2d48" },
        },
        {
          default: () => t("button.shutdown"),
        }
      );
    },
    key: "shutdown",
    icon: renderIcon(SettingsPowerRound, "#cc2d48"),
  },
  {
    label: () => {
      return h(
        "div",
        {
          style: { color: "#cc2d48" },
        },
        {
          default: () => t("button.quitpst"),
        }
      );
    },
    key: "quitpst",
    icon: renderIcon(SettingsPowerRound, "#cc2d48"),
  },
];
const handleSelectControlCenter = (key) => {
  if (key === "palconf") {
    toPalConf();
  } else if (key === "whitelist") {
    handleWhiteList();
  } else if (key === "rcon") {
    handleRconDrawer();
  } else if (key === "broadcast") {
    handleStartBrodcast();
  } else if (key === "shutdown") {
    handleShutdown();
  } else if (key === "quitpst") {
    handleStartQuitpst();
  } else {
    message.error("错误");
  }
};

// 白名单
const showWhiteListModal = ref(false);
const whiteList = ref([]);
const handleWhiteList = () => {
  if (checkAuthToken()) {
    showWhiteListModal.value = true;
    getWhiteList();
  } else {
    message.error(t("message.requireauth"));
    showWhiteListModal.value = true;
  }
};
const getWhiteList = async () => {
  if (checkAuthToken()) {
    const { data, statusCode } = await new ApiService().getWhitelist();
    if (statusCode.value === 200) {
      if (data.value) {
        whitelistStore().setWhitelist(data.value);
        whiteList.value = [];
        data.value.forEach((item) => {
          whiteList.value.push({
            ...item,
            isNew: false,
          });
        });
      } else {
        whiteList.value = [];
      }
    }
  }
};
// 查看白名单中的该玩家
const showWhitelistPlayer = ref(null);
const showCurrentPlayer = (id) => {
  showWhitelistPlayer.value = id;
  showWhiteListModal.value = false;
};
// 从白名单中移除该玩家
const removeWhiteList = async (player) => {
  if (!player.player_uid && !player.steam_id) {
    message.error(
      t("message.removewhitefail", {
        err: "player_uid or steam_id is required",
      })
    );
    return;
  }
  if (player.isNew) {
    const index = whiteList.value.findIndex(
      (e) => e.player_uid === player.player_uid
    );
    whiteList.value.splice(index, 1);
  } else {
    const { data, statusCode } = await new ApiService().removeWhitelist(player);
    if (statusCode.value === 200) {
      message.success(t("message.removewhitesuccess"));
      await getWhiteList();
    } else {
      message.error(t("message.removewhitefail", { err: data.value?.error }));
    }
  }
};
// 添加一项到白名单中
const virtualListInst = ref();
const handleAddNewWhiteList = () => {
  whiteList.value.unshift({
    name: "",
    player_uid: "",
    steam_id: "",
    isNew: true,
  });
  virtualListInst.value?.scrollTo({ index: 0 });
};
// 保存修改白名单
const putWhiteList = async () => {
  if (whiteList.value.length === 0) {
    return;
  }
  const whiteListData = JSON.stringify(whiteList.value);
  const { data, statusCode } = await new ApiService().putWhitelist(
    whiteListData
  );
  if (statusCode.value === 200) {
    message.success(t("message.addwhitesuccess"));
    getWhiteList();
    showWhiteListModal.value = false;
  } else {
    message.error(t("message.addwhitefail", { err: data.value?.error }));
  }
};
// 接受玩家加入到黑名单信息
const getSonWhitelistStatus = () => {
  getWhiteList();
};

// 广播
const showBroadcastModal = ref(false);
const broadcastText = ref("");
const handleStartBrodcast = () => {
  // 开始广播
  if (checkAuthToken()) {
    showBroadcastModal.value = true;
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};
const handleBroadcast = async () => {
  const { data, statusCode } = await new ApiService().sendBroadcast({
    message: broadcastText.value,
  });
  if (statusCode.value === 200) {
    message.success(t("message.broadcastsuccess"));
    showBroadcastModal.value = false;
    broadcastText.value = "";
  } else {
    message.error(t("message.broadcastfail", { err: data.value?.error }));
  }
};

// 退出登录
const handleStartQuitpst = () => {
  if (checkAuthToken()) {
    dialog.warning({
	   title: t("message.warn"),
	   content: t("message.Signoutprompst"),
	   positiveText: t("button.confirm"),
	   negativeText: t("button.cancel"),
	   onPositiveClick: async () => {
	       const token = localStorage.removeItem(PALWORLD_TOKEN);
           userStore().setIsLogin(false, nall);
           isLogin.value = false;
		   if (token === false && isLogin.value === false) {
               message.success(t("message.quitpstsuccess"));
	           showLoginModal.value = true;
			   return;
           } else {
             message.error(t("message.quitpstfail"));
           }
	   },
	   onNegativeClick: () => {},
    });
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};

const doShutdown = async () => {
  return await new ApiService().shutdownServer({
    seconds: 60,
    message: "Server Will Shutdown After 60 Seconds",
  });
};

// 关机
const handleShutdown = () => {
  if (checkAuthToken()) {
    dialog.warning({
      title: t("message.warn"),
      content: t("message.shutdowntip"),
      positiveText: t("button.confirm"),
      negativeText: t("button.cancel"),
      onPositiveClick: async () => {
        const { data, statusCode } = await doShutdown();
        if (statusCode.value === 200) {
          message.success(t("message.shutdownsuccess"));
          return;
        } else {
          message.error(t("message.shutdownfail", { err: data.value?.error }));
        }
      },
      onNegativeClick: () => {},
    });
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};

const toPlayers = async () => {
  if (currentDisplay.value === "players") {
    return;
  }
  currentDisplay.value = "players";
  playerToGuildStore().setUpdateStatus("players");
};
const toGuilds = async () => {
  if (currentDisplay.value === "guilds") {
    return;
  }
  currentDisplay.value = "guilds";
  playerToGuildStore().setUpdateStatus("guilds");
};

const toMap = async () => {
  if (currentDisplay.value === "map") {
    return;
  }
  currentDisplay.value = "map";
  playerToGuildStore().setUpdateStatus("map");
};

const playerToGuildStatus = computed(() =>
  playerToGuildStore().getUpdateStatus()
);

watch(
  () => playerToGuildStatus.value,
  (newVal) => {
    currentDisplay.value = newVal;
    if (newVal === "players") {
    } else if (newVal === "guilds") {
    }
  }
);

/**
 * 检测 token
 */
const checkAuthToken = () => {
  const token = localStorage.getItem(PALWORLD_TOKEN);
  if (token && token !== "") {
    if (isTokenExpired(token)) {
      localStorage.removeItem(PALWORLD_TOKEN);
      return false;
    }
    isLogin.value = true;
    authToken.value = token;
    return true;
  }
  return false;
};
const isTokenExpired = (token) => {
  const payload = JSON.parse(atob(token.split(".")[1]));
  return payload.exp < Date.now() / 1000;
};

const backupModal = ref(false);
const backupList = ref([]);

const handleBackupList = () => {
  if (checkAuthToken()) {
    backupModal.value = true;
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};
const getBackupList = async () => {
  if (checkAuthToken()) {
    const { data, statusCode } = await new ApiService().getBackupList({
      startTime: range.value[0],
      endTime: range.value[1],
    });
    if (statusCode.value === 200) {
      backupList.value = data.value;
    }
  }
};
const getBackupListWithRange = async (selectRange) => {
  let startTime = selectRange[0] ? selectRange[0] : 0;
  let endTime = selectRange[1] ? selectRange[1] : 0;
  if (checkAuthToken()) {
    const { data, statusCode } = await new ApiService().getBackupList({
      startTime,
      endTime,
    });
    if (statusCode.value === 200) {
      backupList.value = data.value;
    }
  }
};
const backupColumns = [
  {
    title: t("item.time"),
    key: "save_time",
    width: "200px",
    render: (row) => {
      return dayjs(row.save_time).format("YYYY-MM-DD HH:mm:ss");
    },
  },
  // {
  //   title: t("item.backupFile"),
  //   key: "path",
  //   render: (row) => {
  //     return row.path;
  //   },
  // },
  {
    title: "",
    key: "action",
    width: "200px",
    render: (row) => {
      return [
        h(
          NButton,
          {
            type: "primary",
            size: "small",
            renderIcon: () => h(CloudDownloadOutlined),
            onClick: () => downloadBackup(row),
          },
          { default: () => t("button.download") }
        ),
        h(
          NButton,
          {
            type: "error",
            size: "small",
            renderIcon: () => h(DeleteOutlineTwotone),
            style: "margin-left: 20px",
            onClick: () => removeBackup(row),
          },
          { default: () => t("button.remove") }
        ),
      ];
    },
  },
];

const range = ref([Date.now() - 1 * 24 * 60 * 60 * 1000, Date.now()]);
const isDownloading = ref(false);
const removeBackup = async (item) => {
  if (checkAuthToken()) {
    isDownloading.value = true;
    const { data, statusCode } = await new ApiService().removeBackup(
      item.backup_id
    );
    if (statusCode.value === 200) {
      message.success(t("message.removebackupsuccess"));
      await getBackupList();
    } else {
      message.error(t("message.removebackupfail", { err: data.value?.error }));
    }
    isDownloading.value = false;
  }
};

const downloadBackup = async (item) => {
  if (checkAuthToken()) {
    isDownloading.value = true;
    try {
      const { data: blobData, execute: fetchBlob } =
        await new ApiService().downloadBackup(item.backup_id);
      await fetchBlob();
      const url = URL.createObjectURL(blobData.value);
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", item.path);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
      message.success(t("message.downloadsuccess"));
    } catch (error) {
      console.error("Download failed", error);
    }
    isDownloading.value = false;
  }
};

onMounted(async () => {
  locale.value = localStorage.getItem("locale");
  languageOptions.value = [
    {
      label: "简体中文",
      key: "zh",
      disabled: locale.value == "zh",
    },
    {
      label: "English",
      key: "en",
      disabled: locale.value == "en",
    },
    {
      label: "日本語",
      key: "ja",
      disabled: locale.value == "ja",
    },
  ];
  const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
  mediaQuery.addEventListener("change", updateDarkMode);
  isDarkMode.value = mediaQuery.matches;

  rconItemOptions.value = itemMap[locale.value].map((item) => {
    return {
      label: `${item.name}-${item.key}`,
      value: item.key,
    };
  });

  rconPalOptions.value = Object.entries(palMap[locale.value]).map(
    ([key, value]) => {
      return {
        label: `${value}-${key}`,
        value: key,
      };
    }
  );
  skillTypeList.value = getSkillTypeList();
  loading.value = true;
  checkAuthToken();
  getServerInfo();
  getServerMetrics();
  getServerToolInfo();
  getPlayerList();
  await getWhiteList();
  loading.value = false;
  await getBackupList();
  setInterval(async () => {
    await getPlayerList();
    await getServerMetrics();
  }, 60000);
  // 调试用
  // currentDisplay.value = "map";
  // playerToGuildStore().setUpdateStatus("map");
});
</script>

<template>
  <div class="home-page">
    <div
      :class="isDarkMode ? 'bg-#18181c text-#fff' : 'bg-#fff text-#18181c'"
      class="flex justify-between items-center p-3"
    >
      <n-space class="flex items-center">
        <!-- 显示标题，根据语言环境动态显示 -->
		<span
          class="line-clamp-1"
          :class="smallScreen ? 'text-lg' : 'text-2xl'"
          >{{ $t("title") }}</span
        >
		<!-- 显示服务器工具版本信息，如果有新版本则显示'new'标签 -->
        <n-badge
          v-if="serverToolInfo?.version"
          :value="hasNewVersion ? 'new' : ''"
        >
          <n-tag
            type="warning"
            :size="smallScreen ? 'mini' : 'medium'"
            round
            @click="toGithub"
            style="cursor: pointer"
            >{{ serverToolInfo.version }}</n-tag
          >
        </n-badge>
		<!-- 显示服务器信息，包括名称和版本，使用工具提示显示更多信息 -->
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-tag type="default" :size="smallScreen ? 'medium' : 'large'">{{
              serverInfo?.name
                ? `${serverInfo.name + " " + serverInfo.version}`
                : $t("message.loading")
            }}</n-tag>
          </template>
          <div>
            <p>{{ $t("item.serverFps") }}: {{ serverMetrics?.server_fps }}</p>
            <p>{{ $t("item.serverUptime") }}: {{ serverMetrics?.uptime }}(s)</p>
            <p>{{ $t("item.serverDays") }}: {{ serverMetrics?.days }}</p>
            <p>
              {{ $t("item.serverFrameTime") }}:
              {{ serverMetrics?.server_frame_time }}(ms)
            </p>
            <p>
              {{ $t("item.maxPlayerNum") }}: {{ serverMetrics?.max_player_num }}
            </p>
          </div>
        </n-tooltip>
      </n-space>

      <n-space>
        <!-- 下拉菜单，用于选择语言 -->
		<n-dropdown
          trigger="hover"
          :options="languageOptions"
          @select="handleSelectLanguage"
        >
          <n-button type="default" secondary strong circle>
            <template #icon>
              <n-icon><LanguageSharp /></n-icon>
            </template>
          </n-button>
        </n-dropdown>

        <!-- 登录按钮或已认证标签，根据登录状态动态显示 -->
		<n-button
          type="primary"
          secondary
          strong
          @click="showLoginModal = true"
          v-if="!isLogin"
        >
          <template #icon>
            <n-icon>
              <AdminPanelSettingsOutlined />
            </n-icon>
          </template>
          {{ $t("button.auth") }}
        </n-button>
        <n-tag v-else type="success" size="large" round>
          <template #icon>
            <n-icon>
              <AdminPanelSettingsOutlined />
            </n-icon>
          </template>
          {{ $t("status.authenticated") }}
        </n-tag>
      </n-space>
    </div>
    <div class="w-full">
      <div class="rounded-lg" v-if="!loading">
        <n-layout style="height: calc(100vh - 64px)">
          <n-layout-header class="p-3 flex justify-between h-16" bordered>
            <n-button-group :size="smallScreen ? 'medium' : 'large'">
			<!-- 未登录隐藏玩家按钮 -->
			<!-- <n-space v-if="isLogin" class="button.players"> -->
			<!-- 玩家按钮，根据当前显示状态设置按钮类型 -->
              <n-button
                @click="toPlayers"
                :type="currentDisplay === 'players' ? 'primary' : 'tertiary'"
                secondary
                strong
                round
              >
                <template #icon>
                  <n-icon>
                    <GameController />
                  </n-icon>
                </template>
                {{ $t("button.players") }}
              </n-button>
			  <!-- </n-space> -->
			  
			<!-- 未登录隐藏公会按钮 -->
			<!-- <n-space v-if="isLogin" class="button.guilds"> -->
			<!-- 公会按钮，根据当前显示状态设置按钮类型 -->
              <n-button
                @click="toGuilds()"
                :type="currentDisplay === 'guilds' ? 'primary' : 'tertiary'"
                secondary
                strong
                round
              >
                <template #icon>
                  <n-icon>
                    <SupervisedUserCircleRound />
                  </n-icon>
                </template>
                {{ $t("button.guilds") }}
              </n-button>
			<!-- </n-space> -->
			
			<!-- 地图按钮，根据当前显示状态设置按钮类型 -->
              <!-- <n-space class="button.map"> -->
			  <n-button
                @click="toMap()"
                :type="currentDisplay === 'map' ? 'primary' : 'tertiary'"
                secondary
                strong
                round
              >
                <template #icon>
                  <n-icon>
                    <PublicRound />
                  </n-icon>
                </template>
                {{ $t("button.map") }}
              </n-button>
			  <!-- </n-space> -->
            </n-button-group>
			
			<!-- 状态标签，显示玩家数量和在线数量 -->
            <n-space>
			<!-- 玩家数量标签 -->
              <n-tag type="info" round size="large">{{
                $t("status.player_number", { number: playerList?.length })
              }}</n-tag>
			<!-- 在线数量标签 -->
              <n-tag type="success" round size="large">{{
                $t("status.online_number", {
                  number: serverMetrics?.current_player_num,
                })
              }}</n-tag>
            </n-space>
			
			<!-- 登录后的操作按钮和控制中心下拉菜单 -->
            <n-space v-if="isLogin" class="flex items-center">
              <!-- 备份按钮 -->
			  <n-button
                :size="smallScreen ? 'medium' : 'large'"
                type="success"
                secondary
                strong
                round
                @click="handleBackupList"
              >
                <template #icon>
                  <n-icon>
                    <ArchiveOutlined />
                  </n-icon>
                </template>
                {{ $t("button.backup") }}
              </n-button>
			  <!-- RCON按钮 -->
              <n-button
                :size="smallScreen ? 'medium' : 'large'"
                type="primary"
                secondary
                strong
                round
                @click="handleRconDrawer"
              >
                <template #icon>
                  <n-icon>
                    <Terminal />
                  </n-icon>
                </template>
                {{ $t("button.rcon") }}
              </n-button>
			  <!-- 控制中心下拉菜单 -->
              <n-dropdown
                trigger="click"
                size="large"
                :options="controlCenterOption"
                @select="handleSelectControlCenter"
              >
                <n-button
                  :size="smallScreen ? 'medium' : 'large'"
                  type="error"
                  secondary
                  strong
                  round
                >
                  <template #icon>
                    <n-icon>
                      <GuiManagement />
                    </n-icon>
                  </template>
                  {{ $t("button.controlCenter") }}</n-button
                >
              </n-dropdown>
			  <!--注释掉的按钮-->
              <!-- <n-button
                :size="smallScreen ? 'medium' : 'large'"
                type="default"
                secondary
                strong
                round
                @click="toPalConf"
              >
                <template #icon>
                  <n-icon>
                    <Settings />
                  </n-icon>
                </template>
                {{ $t("button.palconf") }}
              </n-button>
              <n-button
                :size="smallScreen ? 'medium' : 'large'"
                type="warning"
                secondary
                strong
                round
                @click="handleWhiteList"
              >
                <template #icon>
                  <n-icon>
                    <ShieldCheckmarkSharp />
                  </n-icon>
                </template>
                {{ $t("button.whitelist") }}
              </n-button>
              <n-button
                :size="smallScreen ? 'medium' : 'large'"
                type="success"
                secondary
                strong
                round
                @click="handleStartBrodcast"
              >
                <template #icon>
                  <n-icon>
                    <BroadcastTower />
                  </n-icon>
                </template>
                {{ $t("button.broadcast") }}
              </n-button>
              <n-button
                :size="smallScreen ? 'medium' : 'large'"
                type="error"
                secondary
                strong
                round
                @click="handleShutdown"
              >
                <template #icon>
                  <n-icon>
                    <SettingsPowerRound />
                  </n-icon>
                </template>
                {{ $t("button.shutdown") }}
              </n-button> -->
            </n-space>
          </n-layout-header>
		  <!-- 内容区域，根据当前显示状态显示不同的组件 -->
          <div class="overflow-hidden" style="height: calc(100% - 64px)">
            <player-list
              v-if="currentDisplay === 'players'"
              :showWhitelistPlayer="showWhitelistPlayer"
              @onWhitelistStatus="getSonWhitelistStatus"
            ></player-list>
            <guild-list v-if="currentDisplay === 'guilds'"></guild-list>
            <map-view v-if="currentDisplay === 'map'"></map-view>
          </div>
        </n-layout>
      </div>
    </div>
  </div>
  <!-- 登录模态框 -->
  <n-modal
    v-model:show="showLoginModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 600px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('modal.auth')"
    size="huge"
    :bordered="false"
    :segmented="segmented"
  >
    <div>
	  <!-- 描述性文本 -->
      <span class="block pb-2">{{ $t("message.authdesc") }}</span>
	  <!-- 密码输入框 -->
      <n-input
        type="password"
        show-password-on="click"
        size="large"
        v-model:value="password"
      ></n-input>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <!-- 取消按钮 -->
		<n-button
          type="tertiary"
          @click="
            () => {
              showLoginModal = false;
              password = '';
            }
          "
          >{{ $t("button.cancel") }}</n-button
        >
		<!-- 确认按钮，带有一定的边距 -->
        <n-button class="ml-3 w-40" type="primary" @click="handleLogin">{{
          $t("button.confirm")
        }}</n-button>
      </div>
    </template>
  </n-modal>

  <!-- broadcast modal -->
  <n-modal
    v-model:show="showBroadcastModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 600px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('modal.broadcast')"
    size="huge"
    :bordered="false"
    :segmented="segmented"
  >
    <div>
      <!-- 输入框用于输入广播文本 -->
	  <n-input
        type="text"
        show-password-on="click"
        v-model:value="broadcastText"
      ></n-input>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <!-- 取消按钮，点击后关闭模态框并清空广播文本 -->
		<n-button
          type="tertiary"
          @click="
            () => {
              showBroadcastModal = false;
              broadcastText = '';
            }
          "
          >{{ $t("button.cancel") }}</n-button
        >
		<!-- 确认按钮，点击后处理广播逻辑 -->
        <n-button class="ml-3 w-40" type="primary" @click="handleBroadcast">{{
          $t("button.confirm")
        }}</n-button>
      </div>
    </template>
  </n-modal>

  <!-- 自定义RCON抽屉组件 -->
  <n-modal
    v-model:show="showRconAddModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 600px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('button.addRcon')"
    size="huge"
    :bordered="false"
    :segmented="segmented"
  >
    <!-- 标签页组件，用于切换导入和添加RCON命令的视图 -->
	<n-tabs default-value="import" size="large" justify-content="space-evenly">
      <!-- 导入标签页 -->
	  <n-tab-pane name="import" :tab="$t('button.import')">
        <!-- 文件上传组件，用于导入RCON文件 -->
		<n-upload
          multiple
          directory-dnd
          action="/api/rcon/import"
          :headers="{ Authorization: `Bearer ${authToken}` }"
          :max="1"
          @finish="handleImportRconFinish"
          @error="handleImportRconError"
        >
          <n-upload-dragger>
            <!-- 拖拽区域内容 -->
			<div style="margin-bottom: 12px">
              <n-icon size="48" :depth="3">
                <ArchiveOutline />
              </n-icon>
            </div>
            <n-text style="font-size: 16px">
              {{ $t("message.importRconTitle") }}
            </n-text>
            <n-p depth="3" style="margin: 8px 0 0 0">
              {{ $t("message.importRconDesc") }}
            </n-p>
          </n-upload-dragger>
        </n-upload>
      </n-tab-pane>
      <n-tab-pane name="add" :tab="$t('button.add')">
        <!-- 输入框组件，用于输入RCON命令、备注和占位符 -->
		<n-input
          v-model:value="newRconCommand"
          size="large"
          round
          :placeholder="$t('input.rcon')"
        ></n-input>
        <n-input
          class="mt-5"
          v-model:value="newRconRemark"
          size="large"
          round
          :placeholder="$t('input.remark')"
        ></n-input>
		<!-- 按钮组件，用于添加RCON命令 -->
        <n-input
          class="mt-5"
          v-model:value="newRconPlaceholder"
          size="large"
          round
          :placeholder="$t('input.placeholder')"
        ></n-input>
        <n-button
          class="mt-5"
          style="width: 100%"
          type="primary"
          @click="addRconCommand"
          strong
          secondary
        >
          {{ $t("button.add") }}
        </n-button>
      </n-tab-pane>
    </n-tabs>
  </n-modal>
  <!-- 抽屉组件，用于显示RCON命令详情 -->
  <n-drawer v-model:show="showRconDrawer" :width="502" placement="right">
    <n-drawer-content :title="t('modal.rcon')">
      <!-- 抽屉底部内容，包含一个添加RCON命令的按钮 -->
      <template #footer>
        <n-space>
          <n-button type="primary" strong secondary @click="customPalModal">
            {{ $t("button.customPal") }}
          </n-button>
          <n-button type="primary" strong secondary @click="handleAddRconCommand">
            {{ $t("button.addRcon") }}
          </n-button>
        </n-space>
      </template>
	  <!-- 玩家选择器 -->
      <div class="flex w-full items-center">
        <n-select
          :placeholder="$t('input.selectPlayer')"
          v-model:value="rconSelectedPlayer"
          filterable
          :options="rconPlayerOptions"
        />
		<!-- 复制Steam64 -->
        <n-button
          class="ml-2"
          type="primary"
          strong
          secondary
          @click="copyText('steam_' + (rconSelectedPlayer?.split('-')[1]))"
        >
          {{ $t("button.copysid") }}
        </n-button>
		<!-- 复制玩家ID -->
        <n-button
          class="ml-2"
          type="primary"
          strong
          secondary
          @click="copyText(rconSelectedPlayer?.split('-')[0])"
        >
          {{ $t("button.copypid") }}
        </n-button>
      </div>
	  <!-- 玩家ID和Steam64显示 -->
      <div class="flex w-full items-center mt-3 justify-between">
        <n-text
          >PlayerID: {{ rconSelectedPlayer?.split("-")[0] || "-" }}</n-text
        >
        <n-text>Steam64: {{ rconSelectedPlayer?.split("-")[1] || "-" }}</n-text>
      </div>

      <!-- 项目选择器 -->
	  <div class="flex w-full items-center mt-3">
        <n-select
          :placeholder="$t('input.selectItem')"
          v-model:value="rconSelectedItem"
          filterable
          :options="rconItemOptions"
        />
		<!-- 复制物品ID -->
        <n-button
          class="ml-2"
          type="primary"
          strong
          secondary
          @click="copyText(rconSelectedItem)"
        >
          {{ $t("button.copyitem") }}
        </n-button>
      </div>
	  <!-- 项目ID显示 -->
      <div class="flex w-full items-center mt-3 justify-between">
        <!-- 显示项目ID，如果rconSelectedItem为空则显示"-" -->
		<n-text>ID: {{ rconSelectedItem || "-" }}</n-text>
      </div>

      <div class="flex w-full items-center mt-3">
        <!-- 下拉选择框，用于选择项目 -->
		<n-select
          :placeholder="$t('input.selectPal')"
          v-model:value="rconSelectedPal"
          filterable
          :options="rconPalOptions"
        />
		
		<!-- 
         注释：原本存在一个按钮用于复制选中的项目ID，
         但此功能目前被禁用或待实现，因此注释掉相关代码。
         如果需要重新启用，可以取消以下代码的注释。
        -->
        <!-- <n-button
          class="ml-2"
          type="primary"
          strong
          secondary
          @click="copyText(rconSelectedPal)"
        >
          {{ $t("button.copypal") }}
        </n-button> -->
		<!-- 复制帕鲁ID -->
		<n-button
          class="ml-2"
          type="primary"
          strong
          secondary
          @click="copyText(rconSelectedPal)"
        >
          {{ $t("button.copypal") }}
        </n-button>
		
      </div>
	  <!-- ID 显示区域，如果 rconSelectedPal 为空则显示 "-" -->
      <div class="flex w-full items-center mt-3 justify-between">
        <n-text>ID: {{ rconSelectedPal || "-" }}</n-text>
      </div>
	  <!-- 如果没有 rconCommands，则显示空内容提示 -->
      <n-empty class="mt-3" v-if="rconCommands.length == 0"> </n-empty>
	  <!-- 可折叠的命令列表区域 -->
      <n-collapse class="mt-3">
        <n-collapse-item
          v-for="rconCommand in rconCommands"
          :key="rconCommand.uuid"
          :title="rconCommand.command"
          :name="rconCommand.uuid"
        >
          <!-- 命令的额外说明信息 -->
		  <template #header-extra> {{ rconCommand.remark }} </template>
          <!-- 输入框和按钮组合区域 -->
		  <n-input-group>
            <n-input
              round
              :placeholder="rconCommand.placeholder"
              v-model:value="rconCommandsExtra[rconCommand.uuid]"
            >
              <!-- 输入框前缀，显示命令和占位符信息 -->
			  <template #prefix>
                <n-text>{{
                  rconCommand.command + (rconCommand.placeholder ? "  +" : "")
                }}</n-text>
              </template>
            </n-input>
			<!-- 执行按钮，点击后发送命令 -->
            <n-button
              type="primary"
              ghost
              round
              @click="sendRconCommand(rconCommand.uuid)"
            >
              {{ $t("button.execute") }}
            </n-button>
          </n-input-group>
		  <!-- 删除按钮，点击后移除命令 -->
          <n-button
            class="mt-3"
            style="width: 100%"
            type="error"
            dashed
            @click="removeRconCommand(rconCommand.uuid)"
          >
            <template #icon>
              <n-icon>
                <DeleteFilled />
              </n-icon>
            </template>
            {{ $t("button.remove") }}
          </n-button>
        </n-collapse-item>
      </n-collapse>
    </n-drawer-content>
  </n-drawer>

  <!-- whitelist modal -->
  <n-modal
    v-model:show="showWhiteListModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 700px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('modal.whitelist')"
    size="large"
    :bordered="false"
    :mask-closable="false"
    :close-on-esc="false"
    :segmented="segmented"
  >
    <div>
      <!-- 当白名单列表为空时显示空状态 -->
	  <n-empty v-if="whiteList.length == 0"> </n-empty>
      <!-- 使用虚拟列表来渲染白名单项，提高性能 -->
	  <n-virtual-list
        v-else
        ref="virtualListInst"
        style="height: 320px"
        :item-size="42"
        :items="whiteList"
      >
        <template #default="{ item }">
          <div
            :key="item.player_uid"
            class="flex flex-col item mlr-3 mb-3"
            style="height: 42px"
          >
            <n-grid>
              <n-gi span="19">
                <n-input-group>
                  <!-- 玩家昵称输入框 -->
				  <n-input
                    v-model:value="item.name"
                    :style="{ width: '33%' }"
                    :placeholder="$t('input.nickname')"
                  />
				  <!-- 玩家UID输入框 -->
                  <n-input
                    v-model:value="item.player_uid"
                    :style="{ width: '33%' }"
                    :placeholder="$t('input.player_uid')"
                  />
				  <!-- Steam ID输入框 -->
                  <n-input
                    v-model:value="item.steam_id"
                    :style="{ width: '33%' }"
                    :placeholder="$t('input.steam_id')"
                  />
                </n-input-group>
              </n-gi>
              <n-gi span="5">
                <div class="flex justify-end mr-3">
                  <n-space v-if="item.player_uid || item.steam_id">
                    <!-- 查看玩家按钮 -->
					<n-button
                      strong
                      secondary
                      type="primary"
                      @click="showCurrentPlayer(item.player_uid)"
                    >
                      <template #icon>
                        <n-icon><RemoveRedEyeTwotone /></n-icon>
                      </template>
                    </n-button>
					<!-- 删除白名单项按钮 -->
                    <n-button
                      @click="removeWhiteList(item)"
                      strong
                      secondary
                      type="error"
                    >
                      <template #icon>
                        <n-icon><DeleteOutlineTwotone /></n-icon>
                      </template>
                    </n-button>
                  </n-space>
                </div>
              </n-gi>
            </n-grid>
          </div>
        </template>
      </n-virtual-list>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <n-space>
          <!-- 添加新白名单项按钮 -->
		  <n-button type="primary" @click="handleAddNewWhiteList">
            {{ $t("button.addNew") }}
          </n-button>

          <!-- 取消按钮 -->
		  <n-button
            type="tertiary"
            @click="
              () => {
                showWhiteListModal = false;
              }
            "
          >
            {{ $t("button.cancel") }}
          </n-button>

          <!-- 保存按钮，当白名单为空时禁用 -->
		  <n-button
            :disabled="whiteList.length === 0"
            @click="putWhiteList"
            strong
            secondary
            type="success"
          >
            {{ $t("button.save") }}
          </n-button>
        </n-space>
      </div>
    </template>
  </n-modal>
  <!-- backup modal -->
  <!-- 备份模态框 -->
  <n-modal
    v-model:show="backupModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 700px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('modal.backup')"
    size="small"
    :bordered="false"
    :mask-closable="false"
    :close-on-esc="false"
    :segmented="segmented"
  >
    <div>
      <!-- 当没有备份列表时显示空状态 -->
	  <n-empty description="empty" v-if="backupList.length == 0"> </n-empty>
	  <!-- 当有备份列表时显示内容 -->
      <div class="flex flex-col item mlr-3 mb-3 p-1" v-else>
        <!-- 日期选择器，用于选择备份的时间范围 -->
		<n-date-picker
          class="mb-4"
          v-model:value="range"
          type="datetimerange"
          @confirm="getBackupListWithRange"
        />
		<!-- 滚动条容器，用于包裹数据表格 -->
        <n-scrollbar style="max-height: 320px">
          <!-- 数据表格，用于显示备份列表 -->
		  <n-data-table
            :columns="backupColumns"
            :data="backupList"
            :bordered="false"
          />
        </n-scrollbar>
      </div>
    </div>
	<!-- 模态框底部操作区域 -->
    <template #footer>
      <div class="flex justify-end">
        <n-space>
          <!-- 关闭按钮，用于关闭模态框 -->
		  <n-button
            type="tertiary"
            @click="
              () => {
                backupModal = false;
              }
            "
          >
            {{ $t("button.close") }}
          </n-button>
        </n-space>
      </div>
    </template>
  </n-modal>
</template>
