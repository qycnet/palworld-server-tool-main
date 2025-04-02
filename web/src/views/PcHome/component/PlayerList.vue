<script setup>
// 导入ApiService用于API请求
import ApiService from "@/service/api";
// 导入pageStore用于页面状态管理
import pageStore from "@/stores/model/page.js";
// 从vue导入ref, onMounted, computed等响应式API
import { ref, onMounted, computed } from "vue";
// 导入dayjs库用于日期处理
import dayjs from "dayjs";
// 导入vue-i18n用于国际化
import { useI18n } from "vue-i18n";
// 导入skill.json技能数据
import skillMap from "@/assets/skill.json";
// 从naive-ui导入NAvatar和NTag组件
import { NAvatar, NTag } from "naive-ui";
// 导入PlayerDetail子组件
import PlayerDetail from "./PlayerDetail.vue";
// 导入playerToGuildStore和whitelistStore用于玩家和白名单状态管理
import playerToGuildStore from "@/stores/model/playerToGuild";
import whitelistStore from "@/stores/model/whitelist";

// 使用国际化功能
const { t, locale } = useI18n();

// 定义接收的props
const props = defineProps(["showWhitelistPlayer"]);
// 计算属性，用于获取showWhitelistPlayer的值
const showWhitelistPlayer = computed(() => props.showWhitelistPlayer);

// 响应式变量，用于判断是否为暗色模式
const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches
);

// 计算属性，获取页面宽度
const pageWidth = computed(() => pageStore().getScreenWidth());
// 计算属性，判断是否为小屏幕
const smallScreen = computed(() => pageWidth.value < 1024);

// 响应式变量，用于加载玩家列表的状态
const loadingPlayer = ref(false);
// 响应式变量，用于加载玩家详情的状态
const loadingPlayerDetail = ref(false);
// 响应式变量，存储玩家列表
const playerList = ref([]);
// 响应式变量，存储玩家详情信息
const playerInfo = ref(null);
// 响应式变量，存储玩家好友列表
const playerPalsList = ref([]);
// 响应式变量，存储技能类型列表
const skillTypeList = ref([]);

// 获取玩家列表
const getPlayerList = async () => {
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  playerList.value = data.value;
};

// 获取玩家详情信息
const getPlayerInfo = async (player_uid) => {
  const { data } = await new ApiService().getPlayer({ playerUid: player_uid });
  playerInfo.value = data.value;
  playerPalsList.value = playerInfo?.value.pals
    ? JSON.parse(JSON.stringify(playerInfo?.value.pals))
    : [];
  nextTick(() => {
    const playerInfoEL = document.getElementById("player-info");
    if (playerInfoEL) {
      playerInfoEL.scrollIntoView({ behavior: "smooth" });
    }
  });
};

// 点击事件处理函数，用于获取玩家详情信息
const clickGetPlayerInfo = async (id) => {
  if (playerInfo.value.player_uid !== id) {
    loadingPlayerDetail.value = true;
    await getPlayerInfo(id);
    loadingPlayerDetail.value = false;
  }
};

// 监听showWhitelistPlayer的变化，用于获取对应的玩家详情信息
watch(
  () => showWhitelistPlayer.value,
  async (newVal) => {
    if (playerInfo.value.player_uid !== newVal) {
      loadingPlayerDetail.value = true;
      await getPlayerInfo(newVal);
      loadingPlayerDetail.value = false;
    }
  }
);

// 白名单
const whiteList = computed(() => whitelistStore().getWhitelist());
// 计算属性，判断玩家是否在白名单中
const isWhite = computed(() => (player) => {
  if (player) {
    return whiteList.value.some((whitelistItem) => {
      return (
        (whitelistItem.player_uid &&
          whitelistItem.player_uid === player.player_uid) ||
        (whitelistItem.steam_id && whitelistItem.steam_id === player.steam_id)
      );
    });
  } else {
    return false;
  }
});

// 生命周期钩子，组件挂载时执行
onMounted(async () => {
  loadingPlayerDetail.value = true;
  loadingPlayer.value = true;
  await getPlayerList();
  skillTypeList.value = getSkillTypeList();
  loadingPlayer.value = false;
  if (playerList.value.length > 0) {
    const currentUid = playerToGuildStore().getCurrentUid();
    await getPlayerInfo(
      currentUid ? currentUid : playerList.value[0].player_uid
    );
    playerToGuildStore().setCurrentUid(null);
  }
  loadingPlayerDetail.value = false;
});

// 其他操作
const getUserAvatar = () => {
  return new URL("@/assets/avatar.webp", import.meta.url).href;
};
// 获取技能类型列表的函数
const getSkillTypeList = () => {
  if (skillMap[locale.value]) {
    return Object.values(skillMap[locale.value]).map((item) => item.name);
  } else {
    return [];
  }
};
// 判断玩家是否在线的函数
const isPlayerOnline = (last_online) => {
  return dayjs() - dayjs(last_online) < 80000;
};
// 格式化显示玩家最后在线时间的函数
const displayLastOnline = (last_online) => {
  if (dayjs(last_online).year() < 1970) {
    return "Unknown";
  }
  return dayjs(last_online).format("YYYY-MM-DD HH:mm:ss");
};
</script>
<template>
  <div class="paler-list h-full">
    <n-layout has-sider class="h-full">
      <!-- 侧边栏布局 -->
	  <n-layout-sider
        :width="smallScreen ? 360 : 400"
        content-style="padding: 24px;"
        :native-scrollbar="false"
        bordered
        class="relative"
      >
        <!-- 玩家列表 -->
		<n-list hoverable clickable>
          <n-list-item
            v-for="player in playerList"
            :key="player.player_uid"
            style="padding: 12px 8px"
            @click="clickGetPlayerInfo(player.player_uid)"
          >
            <!-- 玩家头像 -->
			<template #prefix>
              <n-avatar :src="getUserAvatar()" round></n-avatar>
            </template>
			<!-- 玩家信息 -->
            <div>
              <div class="flex">
                <!-- 玩家在线状态 -->
				<n-tag
                  :bordered="false"
                  size="small"
                  :type="
                    isPlayerOnline(player.last_online) ? 'success' : 'error'
                  "
                  round
                >
                  {{
                    isPlayerOnline(player.last_online)
                      ? $t("status.online")
                      : $t("status.offline")
                  }}
                </n-tag>
				<!-- 玩家等级 -->
                <n-tag class="ml-2" type="primary" size="small" round>
                  Lv.{{ player.level }}
                </n-tag>
				<!-- 白名单标签 -->
                <n-tag
                  v-if="isWhite(player)"
                  class="ml-2"
                  :bordered="false"
                  round
                  size="small"
                  :color="{
                    color: isDarkMode ? '#fff' : '#d9c36c',
                    textColor: isDarkMode ? '#d9c36c' : '#fff',
                  }"
                >
                  {{ $t("status.whitelist") }}
                </n-tag>
				<!-- 玩家昵称 -->
                <span class="flex-1 pl-2 font-bold line-clamp-1">{{
                  player.nickname
                }}</span>
              </div>
			  <!-- 最后在线时间 -->
              <n-tag :bordered="false" round size="small" class="mt-2">
                {{ $t("status.last_online") }}:
                {{ displayLastOnline(player.last_online) }}
              </n-tag>
            </div>
          </n-list-item>
        </n-list>
		<!-- 加载中提示 -->
        <n-spin
          size="small"
          v-if="loadingPlayer"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description> 加载中... </template>
        </n-spin>
      </n-layout-sider>
	  <!-- 玩家详情区域 -->
      <n-layout :native-scrollbar="false" class="relative">
        <player-detail
          :playerInfo="playerInfo"
          :playerPalsList="playerPalsList"
          :whiteList="whiteList"
        ></player-detail>
		<!-- 加载中提示 -->
        <n-spin
          size="small"
          v-if="loadingPlayerDetail"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description> 加载中... </template>
        </n-spin>
      </n-layout>
    </n-layout>
  </div>
</template>