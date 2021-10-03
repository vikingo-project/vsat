import Vue from "vue";
import Router from "vue-router";
Vue.use(Router);

let router = new Router({
  linkActiveClass: "activeItem",
  linkExactActiveClass: "activeItem",
  mode: "history",
  scrollBehavior(to, from, savedPosition) {
    return { x: 0, y: 0 };
  },
  routes: [
    {
      path: "/auth",
      name: "Auth",
      component: () => import("@/views/Auth.vue"),
      meta: {
        clearLayout: true
      }
    },
    {
      path: "/",
      name: "Interactions",
      component: () => import("@/views/Interactions.vue")
    },

    {
      name: "Files",
      path: "/files",
      component: () => import("@/views/Files.vue")
    },

    {
      name: "Services",
      path: "/services",
      component: () => import("@/views/Services.vue")
    }
  ]
});

export default router;
