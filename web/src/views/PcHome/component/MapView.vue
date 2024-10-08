<script setup>
import {useI18n} from "vue-i18n";
import "leaflet/dist/leaflet.css";
import {LCircle, LIcon, LMap, LMarker, LPopup, LTileLayer, LTooltip} from "@vue-leaflet/vue-leaflet";
import ApiService from "@/service/api.js";
import IconBase from "@/assets/map/base.webp";
import IconPlayer from "@/assets/map/player.webp";
import IconBossTower from "@/assets/map/boss_tower.webp";
import IconFastTravel from "@/assets/map/fast_travel.webp";
import playerToGuildStore from "@/stores/model/playerToGuild.js";
import points from "@/assets/map/points.json";

const {t} = useI18n();

const LAND_SCAPE = [335112, 617000, -582888, -301000];

const api = new ApiService();

const zoom = ref(2);
const tiles = ref("");
const playerList = ref([]);
const guildList = ref([]);
const showPlayer = ref(true);
const showBaseCamp = ref(true);
const showBossTower = ref(true);
const showFastTravel = ref(true);

let timer = null;

const toMapPosition = (position) => {
  const x = -256 + 256 * (position[0] - LAND_SCAPE[2]) / (LAND_SCAPE[0] - LAND_SCAPE[2]);
  const y = 256 * (position[1] - LAND_SCAPE[3]) / (LAND_SCAPE[1] - LAND_SCAPE[3]);
  return [x, y];
};

const toMapDistance = (distance) => {
  return 256 * (distance / (LAND_SCAPE[0] - LAND_SCAPE[2]))
};

const ToPlayers = async (uid) => {
  playerToGuildStore().setCurrentUid(uid);
  playerToGuildStore().setUpdateStatus("players");
};

const refreshPlayer = async () => {
  const { data } = await api.getOnlinePlayerList();
  for (const i of data.value) {
    for (const j of playerList.value) {
      if(i.player_uid === j.player_uid) {
        j.location_x = i.location_x;
        j.location_y = i.location_y;
        break;
      }
    }
  }
  timer = setTimeout(refreshPlayer, 5000);
}

onMounted(async () => {
  let res = await api.getPlayerList({});
  playerList.value = res.data.value;
  res = await api.getGuildList();
  guildList.value = res.data.value;
  res = await api.getServerToolConfig();
  tiles.value = res.data.value.map_tiles_url;

  refreshPlayer();
});

onUnmounted(async () => {
  clearTimeout(timer);
});

</script>

<template>
  <div class="map-view h-full">
    <l-map
      ref="map"
      style="width: 100%; height: 100%;"
      crs="Simple"
      v-model:zoom="zoom"
      :use-global-leaflet="false"
      :center="[-128, 128]"
      :min-zoom="0"
      :max-zoom="6"
      :options="{zoomControl: false, attributionControl: false}"
    >
      <l-tile-layer
        :url="tiles"
        no-wrap
        layer-type="base"
        :options="{'bounds': [[0, 0], [-256, 256]]}"
      ></l-tile-layer>
      <l-marker v-if="showFastTravel" v-for="i in points.fast_travel" :lat-lng="toMapPosition([i[0], i[1]])">
        <l-icon :icon-url="IconFastTravel" :icon-size="[48, 48]"/>
      </l-marker>
      <l-marker v-if="showBossTower" v-for="i in points.boss_tower" :lat-lng="toMapPosition([i[0], i[1]])">
        <l-icon :icon-url="IconBossTower" :icon-size="[48, 48]"/>
      </l-marker>
      <l-marker v-if="showPlayer" v-for="i in playerList" :lat-lng="toMapPosition([i.location_x, i.location_y])">
        <l-icon :icon-url="IconPlayer" :icon-size="[45, 45]"/>
        <l-tooltip :options="{direction: 'top', permanent: true, offset: [0, -15]}">{{ i.nickname }}</l-tooltip>
      </l-marker>
      <template v-if="showBaseCamp" v-for="i in guildList">
        <template v-for="j in i.base_camp">
          <l-marker :lat-lng="toMapPosition([j.location_x, j.location_y])">
            <l-icon :icon-url="IconBase" :icon-size="[55, 55]"/>
            <l-popup :options="{interactive: true}">
              <div style="padding-bottom: 3px; font-size: 16px">
                {{ $t("map.baseCampTitle", {name: i.name}) }}
              </div>
              <div style="line-height: 25px">
                {{ $t('map.guildMember') }}
                <span v-for="k in i.players" class="player_name" @click="ToPlayers(k.player_uid)">
                  {{ k.nickname }}
                </span>
              </div>
            </l-popup>
          </l-marker>
          <l-circle :lat-lng="toMapPosition([j.location_x, j.location_y])" :radius="toMapDistance(j.area)"/>
        </template>
      </template>
    </l-map>
    <div class="control">
      <div>
        <span>{{ $t('map.showFastTravel')}}</span>
        <n-switch v-model:value="showFastTravel"/>
      </div>
      <div>
        <span>{{ $t('map.showBossTower')}}</span>
        <n-switch v-model:value="showBossTower"/>
      </div>
      <div>
        <span>{{ $t('map.showPlayer')}}</span>
        <n-switch v-model:value="showPlayer"/>
      </div>
      <div>
        <span>{{ $t('map.showBaseCamp')}}</span>
        <n-switch v-model:value="showBaseCamp"/>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.leaflet-container {
  background: #102536;
  outline: 0;
}

.player_name {
  margin: 0 3px;
  padding: 3px;
  color: #fff;
  background-color: #009f5d;
  border-radius: 3px;
}

.control {
  width: 200px;
  height: 180px;
  position: absolute;
  bottom: 20px;
  right: 20px;
  background-color: rgb(24, 24, 28);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  z-index: 999;
}

.control > div {
  display: flex;
  justify-content: space-between;
  margin: 10px;
}
</style>