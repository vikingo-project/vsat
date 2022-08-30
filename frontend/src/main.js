// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import "babel-polyfill";
import Vue from "vue";
import axios from "axios";
import VueAxios from "vue-axios";

import "@/components/element/packages/theme-chalk/src/common/transition.scss";
import "@/components/element/packages/theme-chalk/src/select.scss";
import "@/components/element/packages/theme-chalk/src/input-number.scss";
import "@/components/element/packages/theme-chalk/src/input.scss";
import "@/components/element/packages/theme-chalk/src/alert.scss";
import "@/components/element/packages/theme-chalk/src/checkbox.scss";
import "@/components/element/packages/theme-chalk/src/switch.scss";
import "@/components/element/packages/theme-chalk/src/button.scss";
import "@/components/element/packages/theme-chalk/src/icon.scss";
import "@/components/element/packages/theme-chalk/src/radio.scss";
import "@/components/element/packages/theme-chalk/src/radio-button.scss";
import "@/components/element/packages/theme-chalk/src/dialog.scss";
import "@/components/element/packages/theme-chalk/src/table.scss";
import "@/components/element/packages/theme-chalk/src/tabs.scss";
import "@/components/element/packages/theme-chalk/src/pagination.scss";
import "@/components/element/packages/theme-chalk/src/loading.scss";
import "@/components/element/packages/theme-chalk/src/notification.scss";
import "@/components/element/packages/theme-chalk/src/upload.scss";
import "@/components/element/packages/theme-chalk/src/tag.scss";
import "@/components/element/packages/theme-chalk/src/popover.scss";
import "@/components/element/packages/theme-chalk/src/popconfirm.scss";
import "@/components/element/packages/theme-chalk/src/dropdown.scss";
import "@/components/element/packages/theme-chalk/src/drawer.scss";
import "@/components/element/packages/theme-chalk/src/date-picker.scss";

import "@/assets/styles/input.css";
import "@/assets/styles/button.css";
import "@/assets/styles/icon.css";
import "@/assets/styles/table.css";
import "@/assets/styles/color.css";
import "@/assets/styles/trash.css";
import "@/assets/fonts/vikingo.css";
import "@/assets/addons.css";
import "@/assets/styles/style.css";
import "@/assets/styles/sat.css";

import App from "./App";
import router from "./router";
import VueSocketIO from "vue-socket.io";
import Elem from "@/components/element/src/index";
import lang from "@/components/element/src/locale/lang/en";
import locale from "@/components/element/src/locale";
import { checkAuth } from "@/auth";
import config from "@/config";

Vue.use(Elem);
Vue.use(VueAxios, axios);
if (!config.desktop_mode) {
  Vue.use(
    new VueSocketIO({
      debug: false,
      connection:
        location.protocol +
        "//" +
        location.hostname +
        (location.port ? ":" + location.port : "") // http://127.0.0.1:1025
    })
  );
}

locale.use(lang);
Vue.config.lang = "en-US";
Vue.config.productionTip = false;

/* eslint-disable no-new */
checkAuth().then(result => {
  new Vue({
    el: "#app",
    router,
    components: { App },
    template: "<App/>"
  });
});
