import { defineStore } from "pinia";
import { ref } from "vue";

const playerToGuildStore = defineStore(
  "playerToGuild",
  () => {
    const currentUid = ref(null);
    const updateStatus = ref("players");
    // Set
    const setCurrentUid = (uid) => {
      // 设置当前用户ID
      currentUid.value = uid;
    };
    const setUpdateStatus = (status) => {
      // 更新状态值
      updateStatus.value = status;
    };
    // Get
    const getCurrentUid = () => {
      // 返回当前用户的UID
      return currentUid.value;
    };
    const getUpdateStatus = () => {
      // 返回updateStatus.value的值
      return updateStatus.value;
    };

    return {
      setCurrentUid,
      setUpdateStatus,
      getCurrentUid,
      getUpdateStatus,
    };
  },
  {
    persist: true,
  }
);

export default playerToGuildStore;
