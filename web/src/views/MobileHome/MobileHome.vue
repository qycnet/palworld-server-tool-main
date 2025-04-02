<script setup>
import {
  AdminPanelSettingsOutlined,
  SupervisedUserCircleRound,
  SettingsPowerRound,
} from "@vicons/material";
import { ChevronsLeft } from "@vicons/tabler";
import { GameController, LanguageSharp } from "@vicons/ionicons5";
import { BroadcastTower } from "@vicons/fa";
import { onMounted, ref } from "vue";
import { NTag, NButton, useMessage, useDialog } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import palMap from "@/assets/pal.json";
import skillMap from "@/assets/skill.json";
import PlayerList from "./component/PlayerList.vue";
import GuildList from "./component/GuildList.vue";
import PlayerDetail from "./component/PlayerDetail.vue";
import GuildDetail from "./component/GuildDetail.vue";
import userStore from "@/stores/model/user";

const { t, locale } = useI18n();

const message = useMessage();
const dialog = useDialog();

const PALWORLD_TOKEN = "palworld_token";

const loading = ref(false);
const serverInfo = ref({});
const localeLowerPalMap = ref({});
const currentDisplay = ref("players"); // 当前显示的内容（players或guilds）
const isShowDetail = ref(false);
const playerList = ref([]);
const onlinePlayerList = ref([]);
const guildList = ref([]);
const playerInfo = ref({});
const playerPalsList = ref([]);
const currentPlayerPalsList = ref([]);
const guildInfo = ref({});
const skillTypeList = ref([]);
const languageOptions = ref([]);

const contentRef = ref(null);

const isLogin = ref(false);
const authToken = ref("");

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches
);

const updateDarkMode = (e) => {
  isDarkMode.value = e.matches;
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

/**
 * 获取技能类型列表
 * @returns {string[]} 技能类型的名称数组
 */
const getSkillTypeList = () => {
  // 检查是否存在当前语言环境的技能映射
  if (skillMap[locale.value]) {
    // 返回技能映射中所有技能项的名称
	return Object.values(skillMap[locale.value]).map((item) => item.name);
  } else {
    // 如果没有找到，则返回一个空数组
	return [];
  }
};

// get data
const getServerInfo = async () => {
  // 通过ApiService获取服务器信息
  const { data } = await new ApiService().getServerInfo();
  // 将获取到的数据赋值给serverInfo.value
  serverInfo.value = data.value;
};
/**
 * 获取玩家列表
 * @param {boolean} [is_update_info=true] 是否更新玩家信息
 */
const getPlayerList = async (is_update_info = true) => {
  // 获取在线玩家列表
  getOnlineList();
  // 通过ApiService获取玩家列表，按最后在线时间降序排列
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  // 将获取到的数据赋值给playerList.value
  playerList.value = data.value;
};
/**
 * 获取公会列表
 */
const getGuildList = async () => {
  // 通过ApiService获取公会列表
  const { data } = await new ApiService().getGuildList();
  // 将获取到的数据赋值给guildList.value
  guildList.value = data.value;
};

/**
 * 获取玩家信息
 * @param {string} player_uid 玩家UID
 */
const getPlayerInfo = async (player_uid) => {
  // 通过ApiService获取指定玩家的信息
  const { data } = await new ApiService().getPlayer({ playerUid: player_uid });
  // 将获取到的玩家信息赋值给playerInfo.value
  playerInfo.value = data.value;
  // 复制玩家的好友列表到playerPalsList.value
  playerPalsList.value = JSON.parse(JSON.stringify(playerInfo.value.pals));
  // 将当前显示的好友列表限制为pageSize.value的长度
  currentPlayerPalsList.value = playerPalsList.value.slice(0, pageSize.value);
  // 设置显示详情标志为true
  isShowDetail.value = true;
  // 将内容区域滚动到顶部
  contentRef.value.scrollTo(0, 0);
};

/**
 * 获取公会信息
 * @param {string} admin_player_uid 公会管理员UID
 */
const getGuildInfo = async (admin_player_uid) => {
  // 通过ApiService获取指定公会的信息
  const { data } = await new ApiService().getGuild({
	adminPlayerUid: admin_player_uid,
  });
  // 将获取到的公会信息赋值给guildInfo.value
  guildInfo.value = data.value;
  // 设置显示详情标志为true
  isShowDetail.value = true;
  // 将内容区域滚动到顶部
  contentRef.value.scrollTo(0, 0);
};

// 接受子组件
const getChoosePlayer = (uid) => {
  getPlayerInfo(uid);
};
const getChooseGuild = (uid) => {
  getGuildInfo(uid);
};

const getPalName = (name) => {
  const lowerName = name.toLowerCase();
  return localeLowerPalMap.value[lowerName]
    ? localeLowerPalMap.value[lowerName]
    : name;
};

// 游戏用户的帕鲁列表分页，搜索
const clickSearch = (searchValue) => {
  const pattern = /^\s*$|(\s)\1/;
  if (searchValue && !pattern.test(searchValue)) {
    playerPalsList.value = playerInfo.value.pals.filter((item) => {
      return (
        item.skills.some((skill) => {
          return (
            skillMap[locale.value][skill]
              ? skillMap[locale.value][skill].name
              : skill
          ).includes(searchValue);
        }) || getPalName(item.type).includes(searchValue)
      );
    });
  } else {
    playerPalsList.value = JSON.parse(JSON.stringify(playerInfo.value.pals));
  }
  currentPage.value = 1;
  if (playerPalsList.value.length <= 10) {
    finished.value = true;
    currentPlayerPalsList.value = playerPalsList.value ?? [];
  } else {
    finished.value = false;
    currentPlayerPalsList.value = playerPalsList.value.slice(0, pageSize.value);
  }
};
// 滚动加载更多
const palsLoading = ref(false); // 是否正在加载更多好友列表
const currentPage = ref(1); // 当前页码
const pageSize = ref(10); // 每页显示的好友数量
const finished = ref(false); // 是否已经加载完所有好友列表
// 加载更多好友列表的函数
const onLoadPals = () => {
  // 判断是否已经加载完所有好友列表
  if (playerPalsList.value.length <= currentPage.value * pageSize.value) {
    finished.value = true;
  } else {
    // 更新当前页码，并重新计算当前显示的好友列表
	currentPage.value += 1;
    currentPlayerPalsList.value = playerPalsList.value.slice(
      0,
      pageSize.value * currentPage.value
    );
  }
};
// 滚动事件处理函数
const onContentScroll = () => {
  // 判断当前显示的内容是否为玩家列表且详情面板是否显示
  if (currentDisplay.value === "players" && isShowDetail.value) {
    // 获取滚动容器DOM元素
	const dom = document.getElementsByClassName("n-layout-scroll-container");
    // 判断是否滚动到底部
	if (dom[1].scrollTop + dom[1].clientHeight > dom[1].scrollHeight - 6) {
      onLoadPals(); // 加载更多好友列表
    }
  }
};

// 获取在线玩家列表的函数
const getOnlineList = async () => {
  // 调用API获取在线玩家列表数据
  const { data } = await new ApiService().getOnlinePlayerList();
  // 更新在线玩家列表数据
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
  authToken.value = token;
  message.success(t("message.authsuccess"));
  showLoginModal.value = false;
  isLogin.value = true;
};

// broadcast
const showBroadcastModal = ref(false);
const broadcastText = ref("");
const handleStartBrodcast = () => {
  // broadcast start
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

//shutdown
const doShutdown = async () => {
  return await new ApiService().shutdownServer({
    seconds: 60,
    message: "Server Will Shutdown After 60 Seconds",
  });
};

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
/**
 * 切换到玩家列表页面
 */
const toPlayers = async () => {
  // 如果当前已经显示玩家列表，则直接返回
  if (currentDisplay.value === "players") {
    return;
  }
  // 获取玩家列表
  await getPlayerList();
  // 更新当前显示状态为玩家列表
  currentDisplay.value = "players";
  // 隐藏详情
  isShowDetail.value = false;

  // 设置加载状态为完成
  palsLoading.value = false;
  // 设置完成状态为未完成
  finished.value = false;
  // 重置当前页码为第一页
  currentPage.value = 1;

  // 滚动到内容区域的顶部
  contentRef.value.scrollTo(0, 0);
};
/**
 * 切换到公会列表页面
 */
const toGuilds = async () => {
  // 如果当前已经显示公会列表，则直接返回
  if (currentDisplay.value === "guilds") {
    return;
  }
  // 获取公会列表
  await getGuildList();
  // 更新当前显示状态为公会列表
  currentDisplay.value = "guilds";
  // 隐藏详情
  isShowDetail.value = false;

  // 设置加载状态为完成
  palsLoading.value = false;
  // 设置完成状态为未完成
  finished.value = false;
  // 重置当前页码为第一页
  currentPage.value = 1;

  // 滚动到内容区域的顶部
  contentRef.value.scrollTo(0, 0);
};
/**
 * 返回列表页面（无论是玩家列表还是公会列表）
 */
const returnList = () => {
  // 隐藏详情
  isShowDetail.value = false;

  // 设置加载状态为完成
  palsLoading.value = false;
  // 设置完成状态为未完成
  finished.value = false;
  // 重置当前页码为第一页
  currentPage.value = 1;

  // 滚动到内容区域的顶部
  contentRef.value.scrollTo(0, 0);
};

/**
 * 检查授权令牌是否有效
 */
const checkAuthToken = () => {
  // 从本地存储中获取令牌
  const token = localStorage.getItem(PALWORLD_TOKEN);
  // 如果令牌存在且不为空
  if (token && token !== "") {
    // 检查令牌是否过期
	if (isTokenExpired(token)) {
      // 如果过期，则从本地存储中移除令牌
	  localStorage.removeItem(PALWORLD_TOKEN);
      return false;
    }
	// 设置登录状态为已登录
    isLogin.value = true;
	// 设置授权令牌
    authToken.value = token;
    return true;
  }
  // 如果没有令牌或令牌为空，则返回false
  return false;
};
/**
 * 检查令牌是否过期
 * @param {string} token - 要检查的令牌
 * @returns {boolean} - 如果令牌过期，则返回true；否则返回false
 */
const isTokenExpired = (token) => {
  // 解析令牌中的负载
  const payload = JSON.parse(atob(token.split(".")[1]));
  // 检查令牌的过期时间是否小于当前时间
  return payload.exp < Date.now() / 1000;
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
  localeLowerPalMap.value = Object.keys(palMap[locale.value]).reduce(
    (acc, key) => {
      acc[key.toLowerCase()] = palMap[locale.value][key];
      return acc;
    },
    {}
  );
  const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
  mediaQuery.addEventListener("change", updateDarkMode);
  isDarkMode.value = mediaQuery.matches;

  skillTypeList.value = getSkillTypeList();
  loading.value = true;
  checkAuthToken();
  getServerInfo();
  await getPlayerList();
  loading.value = false;
  setInterval(() => {
    getPlayerList(false);
  }, 60000);
});
</script>

<template>
  <div class="home-page overflow-hidden">
    <div
      :class="isDarkMode ? 'bg-#18181c text-#fff' : 'bg-#fff text-#18181c'"
      class="flex justify-between items-center p-3"
    >
      <div>
        <span class="line-clamp-1 text-base">{{ $t("title") }}</span>
        <n-tag type="default" size="small">{{
          serverInfo?.name
            ? `${serverInfo.name + " " + serverInfo.version}`
            : $t("message.loading")
        }}</n-tag>
      </div>
      <n-space vertical>
        <n-space justify="end">
          <n-tag type="info" round size="small">{{
            $t("status.player_number", { number: playerList?.length })
          }}</n-tag>
          <n-tag type="success" round size="small">{{
            $t("status.online_number", { number: onlinePlayerList?.length })
          }}</n-tag>
        </n-space>
        <n-space justify="end" class="flex items-center">
          <n-dropdown
            trigger="hover"
            :options="languageOptions"
            @select="handleSelectLanguage"
          >
            <n-button type="default" secondary strong circle size="small">
              <template #icon>
                <n-icon><LanguageSharp /></n-icon>
              </template>
            </n-button>
          </n-dropdown>

          <n-button
            type="primary"
            size="small"
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
          <n-tag v-else type="success" size="small" round>
            <template #icon>
              <n-icon>
                <AdminPanelSettingsOutlined />
              </n-icon>
            </template>
            {{ $t("status.authenticated") }}
          </n-tag>
        </n-space>
      </n-space>
    </div>
    <div class="w-full">
      <div class="rounded-lg" v-if="!loading && playerList.length > 0">
        <n-layout style="height: calc(100vh - 86px)" has-sider>
          <n-layout-header
            class="flex flex-col justify-between"
            :class="isLogin ? 'h-16' : 'h-10'"
            bordered
          >
            <div v-if="isLogin" class="flex justify-center items-center px-3">
              <n-button
                size="small"
                type="success"
                class="mr-2"
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
                size="small"
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
              </n-button>
            </div>
            <div v-else></div>
            <div class="flex justify-end">
              <n-button-group size="small" class="w-full">
                <n-button
                  v-if="isShowDetail"
                  class="w-20%"
                  @click="returnList"
                  type="tertiary"
                  strong
                  secondary
                >
                  <n-icon size="24">
                    <ChevronsLeft />
                  </n-icon>
                </n-button>
                <n-button
                  :class="isShowDetail ? 'w-40%' : 'w-50%'"
                  @click="toPlayers"
                  :type="currentDisplay === 'players' ? 'primary' : 'tertiary'"
                  secondary
                  strong
                >
                  <template #icon>
                    <n-icon>
                      <GameController />
                    </n-icon>
                  </template>
                  {{ $t("button.players") }}
                </n-button>
                <n-button
                  :class="isShowDetail ? 'w-40%' : 'w-50%'"
                  @click="toGuilds"
                  :type="currentDisplay === 'guilds' ? 'primary' : 'tertiary'"
                  secondary
                  strong
                >
                  <template #icon>
                    <n-icon>
                      <SupervisedUserCircleRound />
                    </n-icon>
                  </template>
                  {{ $t("button.guilds") }}
                </n-button>
              </n-button-group>
            </div>
          </n-layout-header>
          <n-layout
            position="absolute"
            style="top: 64px"
            ref="contentRef"
            @scroll="onContentScroll"
          >
            <!-- 当 isShowDetail 为 false 时显示列表 -->
			<div v-if="!isShowDetail">
              <!-- list -->
              <player-list
                v-if="currentDisplay === 'players'"
                :playerList="playerList"
                @onGetInfo="getChoosePlayer"
              ></player-list>
              <guild-list
                v-if="currentDisplay === 'guilds'"
                :guildList="guildList"
                @onGetInfo="getChooseGuild"
              >
              </guild-list>
            </div>
            <!-- detail -->
			<!-- 当 isShowDetail 为 true 时显示详情 -->
            <div v-else class="relative">
              <player-detail
                v-if="currentDisplay === 'players'"
                :playerInfo="playerInfo"
                :currentPlayerPalsList="currentPlayerPalsList"
                :finished="finished"
                @onSearch="clickSearch"
              ></player-detail>
              <guild-detail
                v-if="currentDisplay === 'guilds'"
                :guildInfo="guildInfo"
              ></guild-detail>
            </div>
          </n-layout>
        </n-layout>
      </div>
    </div>
  </div>
  <!-- 登录 modal -->
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
      <span class="block pb-2">{{ $t("message.authdesc") }}</span>
      <n-input
        type="password"
        show-password-on="click"
        size="large"
        v-model:value="password"
      ></n-input>
    </div>
    <template #footer>
      <div class="flex justify-end">
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
        <n-button class="ml-3 w-40" type="primary" @click="handleLogin">{{
          $t("button.confirm")
        }}</n-button>
      </div>
    </template>
  </n-modal>
  <!--  广播 modal -->
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
      <n-input
        type="text"
        show-password-on="click"
        v-model:value="broadcastText"
      ></n-input>
    </div>
    <template #footer>
      <div class="flex justify-end">
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
        <n-button class="ml-3 w-40" type="primary" @click="handleBroadcast">{{
          $t("button.confirm")
        }}</n-button>
      </div>
    </template>
  </n-modal>
</template>
<style scoped lang="less">
:deep .n-layout-scroll-container {
  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
