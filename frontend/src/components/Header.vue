<template>
  <header>
    <nav class="navbar container navbar-expand-sm">
      <button class="menu-btn btn-sm py-0 px-2 btn-primary" @click="headerDrawer = true" ><i class="el-icon-menu mx-0"></i></button>
      <router-link class="navbar-brand" :to="{ name: 'Interactions' }">
        <img src="@/assets/logo.svg" alt />
      </router-link>
      <button class="menu-btn  btn-sm py-0 px-2 " @click="logout"><i class="vik vik-logout mx-0"></i></button>

      <div class="collapse navbar-collapse" id="navbar1">
        <ul class="navbar-nav mr-auto">
          <router-link class="nav-item" tag="li" :to="{ name: 'Interactions' }">
            <a class="nav-link">Interactions</a>
          </router-link>

          <router-link class="nav-item" tag="li" :to="{ name: 'Services' }">
            <a class="nav-link">Services</a>
          </router-link>

          <router-link class="nav-item" tag="li" :to="{ name: 'Tunnels' }">
            <a class="nav-link">Tunnels</a>
          </router-link>

          <router-link class="nav-item" tag="li" :to="{ name: 'Files' }">
            <a class="nav-link">Files</a>
          </router-link>
        </ul>

        <div class="nav-item ml-left">
          <a
            @click="logout"
            href="javascript:;"
            class="btn btn-sm btn-icon"
            title="Logout"
          >
            <i class="vik vik-logout"></i>
          </a>
        </div>
      </div>
    </nav>

    <el-drawer
      direction="ltr"
      :append-to-body="true"
      :visible.sync="headerDrawer"
      custom-class="header-drawer"
      :with-header="false"
    >
      <div class="drawer-body">
        <router-link class="navbar-brand" :to="{ name: 'Interactions' }">
          <img src="@/assets/logo.svg" alt />
        </router-link>

        <ul class="navbar-nav mr-auto">
          <router-link class="nav-item" tag="li" :to="{ name: 'Interactions' }">
            <a class="nav-link">Interactions</a>
          </router-link>

          <router-link class="nav-item" tag="li" :to="{ name: 'Services' }">
            <a class="nav-link">Services</a>
          </router-link>

          <router-link class="nav-item" tag="li" :to="{ name: 'Tunnels' }">
            <a class="nav-link">Tunnels</a>
          </router-link>

          <router-link class="nav-item" tag="li" :to="{ name: 'Files' }">
            <a class="nav-link">Files</a>
          </router-link>
        </ul>
      </div>

      <div class="drawer-footer mt-auto w-100 ">
       <button @click="logout" class="btn btn-primary w-100" title="Logout">
         <i class="vik vik-logout mr-2"></i>Logout
       </button>
      </div>
    </el-drawer>

  </header>
  
</template>
<script>
import { logout } from "@/auth";
import button from './element/packages/button/src/button.vue';
export default {
  components: { button },
  data() {
    return {
      headerDrawer: false,
    };
  },
  methods: {
    logout() {
      logout();
    },
  },
};
</script>