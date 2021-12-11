<template>
  <div class="container">
    <div class="row">
      <div class="col-md-3 col-lg-2 col-xl-2 navbar-collapse sidebar">
        <div class="sidebar-header">
          <a class="vik-brand" href="#">
            <img src="@/assets/logo.svg" />
          </a>
          <div id="dismissSidebar">
            <svg
              viewBox="0 0 24 24"
              width="24"
              height="24"
              xmlns="http://www.w3.org/2000/svg"
              class="SideNav__closeIcon__56eIa"
              stroke="currentColor"
              tabindex="0"
            >
              <path d="M1 1L23 23M1 23L23 1" />
            </svg>
          </div>
        </div>

        <div class="head-page d-none">
          <div class="title">
            <h1 class="small-title">Your link</h1>
            <el-input placeholder="Your link will be here" v-model="yourLink">
              <el-button
                class="btn btn-link"
                slot="append"
                icon="el-icon-copy-document"
                >Copy</el-button
              >
            </el-input>
            <el-button class="btn btn-primary w-100 mt-3"
              >Generate<i class="el-icon-right"></i
            ></el-button>
          </div>
        </div>

        <div class="category-list">
          <h2>Filter</h2>
          <el-form>
            <el-form-item label="Service">
              <el-select
                v-model="filter.service"
                placeholder="Service"
                multiple
              >
                <el-option-group
                  v-for="group in groupsAndServices"
                  :key="group.label"
                  :label="group.label"
                >
                  <el-option
                    v-for="item in group.services"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  >
                  </el-option>
                </el-option-group>
              </el-select>
            </el-form-item>
            <el-form-item label="Client IP">
              <el-input v-model="filter.client_ip" placeholder=""></el-input>
            </el-form-item>
            <el-form-item label="Local address">
              <el-input
                v-model="filter.local_addr"
                placeholder="for example: :80"
              ></el-input>
            </el-form-item>
            <el-form-item label="Description">
              <el-input
                v-model="filter.description"
                placeholder="type a word..."
              ></el-input>
            </el-form-item>
          </el-form>

          <el-button
            :loading="applyLoading"
            class="btn btn-primary btn-lg"
            @click="applyFilter"
            icon="vik vik-search"
          >
            Apply
          </el-button>
        </div>
      </div>
      <div class="col">
        <div class="content h-100">
          <div class="market-page">
            <h1>Interactions</h1>
          </div>

          <div class="aux-block h-100" v-loading="loading">
            <el-alert
              class="mb-3"
              show-icon
              type="warning"
              v-if="new_interactions > 0"
              >There are {{ new_interactions }} new interactions
              <div class="tags" style="float: right; margin-left: 12px">
                <span @click="reloadSessions"
                  ><i class="el-icon-refresh" style="margin-right: 4px"></i
                  >Refresh</span
                >
              </div>
            </el-alert>
            <div class="event-list" v-if="sessions.length > 0">
              <div class="event-head">
                <ul class="event-head-list">
                  <li style="width: 1%"><el-checkbox></el-checkbox></li>
                  <li style="width: 0%">Service</li>
                  <li style="width: 5%">When</li>
                  <li class="d-none d-md-block" style="width: 2%">
                    Local address
                  </li>
                  <li class="d-none d-md-block" style="width: 2%">Client IP</li>
                  <li>Description</li>
                </ul>
              </div>

              <el-collapse @change="toggle()" v-model="openedPannels">
                <div
                  class="event-card status-bar status-event-info"
                  v-bind:class="{ new: !session.visited }"
                  :key="session.hash"
                  v-for="session in sessions"
                >
                  <el-checkbox
                    style="position: absolute; top: 2px"
                  ></el-checkbox>
                  <el-collapse-item
                    :name="session.hash"
                    class="event-body"
                    style="margin-left: 20px"
                  >
                    <template slot="title" style="display: flex; width: 100%">
                      <ul class="event-body-list w-100">
                        <li style="width: 0%">{{ session.service_name }}</li>
                        <li style="width: 5%">
                          <DateRenderer :date="session.date"></DateRenderer>
                        </li>
                        <li class="d-none d-md-block" style="width: 2%">
                          {{ session.local_addr }}
                        </li>
                        <li class="d-none d-md-block" style="width: 2%">
                          {{ session.client_ip }}
                        </li>
                        <li>{{ session.description }}</li>
                      </ul>
                    </template>
                    <el-dropdown
                      class="event-drop"
                      trigger="click"
                      @command="handleCommand"
                    >
                      <span
                        class="
                          vik-button
                          btn-sm btn-icon btn
                          vik-button--default
                        "
                        ><i class="el-icon-more"></i
                      ></span>
                      <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item
                          class="del"
                          icon="vik vik-delete"
                          :command="{ action: 'remove', hash: session.hash }"
                          >Remove</el-dropdown-item
                        >
                      </el-dropdown-menu>
                    </el-dropdown>
                    <div class="d-md-none mt-2 mb-3">
                      <div class="hidden-event-prop">
                        <ul class="event-body-list event-head-list w-100">
                          <li class="d-block d-md-none" style="width: 2%">
                            Local address
                          </li>
                          <li class="d-block d-md-none" style="width: 2%">
                            Client IP
                          </li>
                        </ul>
                        <ul class="event-body-list w-100">
                          <li class="d-block d-md-none" style="width: 2%">
                            {{ session.local_addr }}
                          </li>
                          <li class="d-block d-md-none" style="width: 2%">
                            {{ session.client_ip }}
                          </li>
                        </ul>
                      </div>
                    </div>
                    <Extra
                      :key="session.hash"
                      :events="inCache(session.hash)"
                    ></Extra>
                  </el-collapse-item>
                </div>
              </el-collapse>

              <el-pagination
                style="margin-top: 10px"
                background
                layout="prev, pager, next, sizes"
                :total="total"
                :page-sizes="[15, 30, 100, 1000]"
                :page-size="pageSize"
                @current-change="changePage"
                @size-change="changePageSize"
                :current-page.sync="currentPage"
              ></el-pagination>
            </div>
            <div class="nothing" v-else>
              <i>No interactions</i>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import Vue from "vue";
import Extra from "@/components/Extra";
import DateRenderer from "@/components/DateRenderer";
import * as utils from "@/utils";

export default {
  components: { Extra, DateRenderer },
  data() {
    return {
      yourLink: "",
      sessions: [],
      currentPage: 1,
      pageSize: 15,
      cache: {},
      new_interactions: 0,

      loading: false,
      applyLoading: false,
      services: [],
      filter: { service: [], client_ip: "", local_addr: "", description: "" },
      total: 0,
      openedPannels: [],
    };
  },

  created() {
    this.$socket.emit("auth", localStorage.getItem("token"));
    this.sockets.subscribe("notifications", (msg) => {
      if (msg.name && msg.name == "new_interaction") {
        this.new_interactions++;
      }
    });
  },
  mounted() {
    this.loadSessions();
    this.loadServices();
  },
  destroyed() {
    this.sockets.unsubscribe("notifications");
  },

  computed: {
    groupsAndServices() {
      let groups = [];
      for (let s of this.services) {
        if (groups.some((elem) => elem.label === s.moduleName)) {
          for (let gindex in groups) {
            if (groups[gindex].label == s.moduleName) {
              groups[gindex].services.push({
                label: s.serviceName,
                value: s.hash,
              });
            }
          }
        } else {
          let service = { label: s.serviceName, value: s.hash };
          groups.push({ label: s.moduleName, services: [service] });
        }
      }
      return groups;
    },
  },
  watch: {
    $route(to, from) {
      this.loadSessions();
    },
  },
  methods: {
    handleCommand(args) {
      let { action, hash } = args;
      if (action === "remove") {
        this.removeSession(hash);
      }
    },
    removeParam(k) {
      utils.removeQueryParam(this, k);
    },
    checkVisited(session) {
      return session.visited ? true : false;
    },
    changePageSize(newSize) {
      this.pageSize = newSize;
      this.changePage(1);
    },
    changePage(newPage) {
      this.currentPage = newPage;
      if (newPage == 1) {
        this.new_interactions = 0; // reset updates counter
      }
      this.loadSessions();
    },

    reloadSessions() {
      this.new_interactions = 0;
      this.changePage(1);
    },
    async applyFilter() {
      this.applyLoading = true;
      await this.loadSessions();
      this.applyLoading = false;
    },

    inCache(hash) {
      if (hash in this.cache) {
        return this.cache[hash];
      }
      return null;
    },
    async toggle() {
      for (let h of this.openedPannels) {
        if (!(h in this.cache)) {
          await this.loadSessionEvents(h);
          for (let s of this.sessions) {
            if (s.hash === h) {
              s.visited = true;
            }
          }
        }
      }
    },
    loadServices() {
      utils.$get(`/api/services/`).then((data) => {
        if (data.status == "ok") {
          this.services = data.services;
        }
      });
    },

    async loadSessions() {
      return new Promise((resolve) => {
        this.openedPannels = [];
        this.cache = {};
        this.loading = true;
        let merged = Object.assign(this.$route.query, {
          page: this.currentPage,
          size: this.pageSize,
        });
        merged = Object.assign(merged, this.filter);
        utils
          .$get(`/api/sessions/?` + utils.objectToString(merged))
          .then((data) => {
            resolve();
            this.loading = false;
            if (data.status == "ok") {
              this.sessions = data.sessions;
              this.total = data.total;
            } else {
              this.$notify.error({
                title: "Error",
                message: "Failed to load files",
              });
            }
          });
      });
    },
    async loadSessionEvents(hash) {
      utils.$get(`/api/sessions/view/${hash}/`).then((data) => {
        if (data.status == "ok") {
          Vue.set(this.cache, hash, data.events);
        }
      });
    },
    async removeSession(hash) {
      utils.$post(`/api/sessions/remove/`, { hash: hash }).then((data) => {
        if (data.status === "ok") {
          this.$notify.success({
            title: "Success",
            message: "Interaction removed",
          });
          this.loadSessions();
        }
      });
    },
  },
};
</script>
<style >
.scroll-area {
  position: relative;
  margin: auto;
  height: 300px;
}
@media only screen and (min-width: 0) and (max-width: 1023px) {
  .scroll-area {
    height: 400px;
  }
}
@media only screen and (min-width: 1024px) and (max-width: 4096px) {
  .scroll-area {
    height: 500px;
  }
}
.nintenpercent {
  width: 95%;
}
.badger {
  margin-top: -20px;
}
.el-badge__content {
  border-radius: 10px;
  font-size: 10px;
  height: 13px;
  line-height: 13px;
  padding: 0 5px;
  background: #9bc7f5;
}
.tb {
  display: flex;
}
</style>
