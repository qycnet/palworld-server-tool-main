import { defineStore } from "pinia";
import { ref } from "vue";

const whitelistStore = defineStore(
  "whitelist",
  () => {
    const whitelist = ref([]);
    // Set
    const setWhitelist = (list) => {
      // 设置白名单
      whitelist.value = list;
    };
    // Get
    const getWhitelist = () => {
      // 返回白名单的值
      return whitelist.value;
    };

    return {
      setWhitelist,
      getWhitelist,
    };
  },
  {
    persist: true,
  }
);

export default whitelistStore;
