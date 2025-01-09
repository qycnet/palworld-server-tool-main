import { defineStore } from "pinia";
import { ref } from "vue";

const pageStore = defineStore(
  "page",
  () => {
    let screenWidth = ref(0);
    // Set
    const setScreenWidth = (info) => {
      // 将传入的屏幕宽度信息赋值给 screenWidth.value
      screenWidth.value = info;
    };
    // Get
    const getScreenWidth = () => {
      // 返回屏幕宽度
      return screenWidth.value;
    };

    return {
      setScreenWidth,
      getScreenWidth,
    };
  },
  {
    persist: true,
  }
);

export default pageStore;
