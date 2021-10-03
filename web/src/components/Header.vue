<template>
  <header>
    <nav class="navbar container navbar-expand-sm">
      <router-link class="navbar-brand" :to="{ name: 'Interactions' }">
        <img src="@/assets/logo.svg" alt />
      </router-link>

      <div class="collapse navbar-collapse" id="navbar1">
        <ul class="navbar-nav mr-auto">
          <router-link class="nav-item" tag="li" :to="{ name: 'Interactions' }">
            <a class="nav-link">Interactions</a>
          </router-link>

          <router-link class="nav-item" tag="li" :to="{ name: 'Services' }">
            <a class="nav-link">Services</a>
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
  </header>
</template>
<script>
import { logout } from "@/auth";
export default {
  methods: {
    logout() {
      logout();
    },
  },
};
</script>