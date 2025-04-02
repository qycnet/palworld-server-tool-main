import { defineStore } from "pinia";
import { ref } from "vue";

const userStore = defineStore(
  "user",
  () => {
    const isLogin = ref(false);
    const token = ref(null);

    const setIsLogin = (status, key) => {
      // 如果状态为true
      if (status) {
        // 将登录状态设置为true
        isLogin.value = true;
        // 将token值设置为传入的key
        token.value = key;
      } else {
        // 如果状态为false
        // 将登录状态设置为false
        isLogin.value = false;
        // 将token值设置为null
        token.value = null;
      }
    };
    const getLoginInfo = () => {
      // 返回一个对象，包含用户的登录状态和token
      return {
        // 用户的登录状态
        isLogin: isLogin.value,
        // 用户的token
        token: token,
      };
    };

    return {
      setIsLogin,
      getLoginInfo,
    };
  },
  {
    persist: true,
  }
);

export default userStore;
