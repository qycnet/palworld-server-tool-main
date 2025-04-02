import { useFetch } from "@vueuse/core";
import router from "@/router";

class Service {
  /**
   * Fetches data from a specified URL.
   *
   * @param {string} url - The URL to fetch data from.
   * @return {Promise<Response>} A Promise that resolves to the response from the server.
   */
  fetch(url) {
    // 调用useFetch函数
    return useFetch(`${url}`, {
      // 在数据更新时允许出现错误
      updateDataOnError: true,
      // 在fetch请求前执行的操作
      beforeFetch({ options }) {
        // 从localStorage中获取token
        const token = localStorage.getItem("palworld_token");
        // 设置请求头
        options.headers = {
          // 设置Authorization头，格式为Bearer token
          Authorization: `Bearer ${token}`,
          // 设置Content-Type头为application/json
          "Content-Type": "application/json",
          // 扩展options.headers，以便添加其他自定义头
          ...options.headers,
          // 设置Remote-Ip-Address头，从localStorage中获取ip，如果未获取到则默认为"127.0.0.1"
          "Remote-Ip-Address": localStorage.getItem("ip") || "127.0.0.1",
        };
        // 返回修改后的options
        return {
          options,
        };
      },
      // 在fetch请求出错时执行的操作
      onFetchError(context) {
        // 如果响应状态码为401，则移除localStorage中的token
        if (context.response.status === 401) {
          localStorage.removeItem("palworld_token");
          // 返回上下文
          return context;
        }
        // 返回上下文
        return context;
      },
    });
  }

  /**
   * Generates a query string from a given credential object.
   *
   * @param {Object} credential - The credential object.
   * @return {string} - The generated query string.
   */
  generateQuery(credential) {
    // 获取 credential 对象的键值对数组
    const entries = Object.entries(credential);
    return entries
      // 对键值对数组进行 reduce 操作
      .reduce((accumulation, [key, value]) => {
        // 如果 value 存在，则将其添加到 accumulation 数组中
        if (value) {
          accumulation.push(`${key}=${value}`);
        }
        // 返回 accumulation 数组
        return accumulation;
      }, [])
      // 将 accumulation 数组中的元素用 "&" 连接成一个字符串
      .join("&");
  }
}

export default Service;
