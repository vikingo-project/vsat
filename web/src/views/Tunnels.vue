<template>
  <div class="container">
    <div class="row">
      <div class="col-md-3 col-lg-2 col-xl-2 navbar-collapse sidebar">
        <div class="sidebar-header">
          <a class="vik-brand" href="#">
            <img src="@/assets/logo.svg" alt />
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

        <div class="category-list">
          <h2>Filter</h2>
          <el-form>
            <el-form-item label="Tunnel type">
              <el-select v-model="filter.tunnel_type" clearable>
                <el-option key="http" label="HTTP" value="http"> </el-option>
                <el-option key="tcp" label="TCP" value="TCP"> </el-option>
              </el-select>
            </el-form-item>
          </el-form>
          <el-button
            :loading="false"
            class="btn btn-primary btn-lg"
            @click="applyFilter"
            icon="vik vik-search"
          >
            Apply
          </el-button>
        </div>
      </div>

      <div class="col">
        <div class="content" v-loading="tunnelsLoading">
          <div class="market-page">
            <h1>Tunnels</h1>
          </div>

          <div class="row item-container atom-list">
            <div class="col-lg-6">
              <article
                @click="newTunnelDialogVisible = true"
                class="btn btn-primary btn-new"
              >
                <i class="vik vik-plus"></i>Add new
              </article>
            </div>

            <div
              class="col-lg-6"
              v-for="tunnel in filteredTunnels"
              :key="tunnel.hash"
            >
              <TunnelCard :tunnel="tunnel"></TunnelCard>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog
      title="New tunnel"
      :visible.sync="newTunnelDialogVisible"
      class="modal-new-modules"
    >
      <form-wizard :key="formKey" finishButtonText="Create new tunnel">
        <tab-content title="Tunnel type" :before-change="validateTunnelType">
          <div class="content-mod">
            <div class="list-modules">
              <div class="select-modules">
                <el-radio v-model="newTunnel.type" label="HTTP">
                  <h3>HTTP</h3>
                  <p>
                    Forward HTTP requests from your personal subdomain like
                    example.vkng.cc to custom address or service
                  </p>
                </el-radio>
                <el-radio v-model="newTunnel.type" label="TCP" disabled>
                  <h3>TCP</h3>
                  <p>Forward TCP stream from vkng.cc to your local address</p>
                </el-radio>
              </div>
            </div>
          </div>
        </tab-content>
        <tab-content title="Settings" :before-change="validateTunnelSettings">
          <!-- global settings -->
          <el-form
            :model="newTunnel"
            ref="form"
            class="service-settings"
            :rules="ruleValidate"
            @submit.native.prevent
          >
            <el-form-item label="Destination host" prop="dstHost">
              <el-select
                filterable
                allow-create
                v-model="newTunnel.dstHost"
                placeholder="Destination host"
              >
                <el-option
                  v-for="net in networks"
                  :value="net.ip"
                  :key="net.ip"
                  >{{ net.ip }}</el-option
                >
              </el-select>
            </el-form-item>

            <el-form-item>
              <div class="d-flex">
                <div class="mr-4">
                  <el-form-item label="Destination port" prop="dstPort">
                    <el-input-number
                      v-model="newTunnel.dstPort"
                      :min="1"
                      :max="65535"
                    ></el-input-number>
                  </el-form-item>
                </div>
                <div class="">
                  <el-form-item label="Auto start">
                    <el-checkbox v-model="newTunnel.autoStart"></el-checkbox>
                  </el-form-item>
                </div>
              </div>
            </el-form-item>

            <el-form-item>
              <el-form-item label="Backend uses TLS">
                <el-checkbox v-model="newTunnel.dstTLS"></el-checkbox>
              </el-form-item>
            </el-form-item>
          </el-form>
        </tab-content>
      </form-wizard>
    </el-dialog>
  </div>
</template>
<script>
import { bus } from "@/bus.js";
import * as utils from "@/utils.js";
import TunnelCard from "@/components/TunnelCard.vue";
import { FormWizard, TabContent } from "@/components/vue-form-wizard";
import "@/components/vue-form-wizard/assets/wizard.scss";

export default {
  name: "tunnels",
  components: {
    FormWizard,
    TabContent,
    TunnelCard,
  },

  data() {
    this.tunnelBootstrap = {
      // tunnelName: "",
      type: ``,
      dstHost: ``,
      dstPort: 80,
      dstTLS: false,
    };

    return {
      filteredTunnels: [],
      filter: { tunnel_type: "" },
      modules: [],
      networks: [],
      newTunnel: Object.assign({}, this.tunnelBootstrap),
      formKey: "default",

      tags: [],
      newTunnelDialogVisible: false,
      tunnelsLoading: false,
      ruleValidate: {
        dstHost: [
          {
            required: true,
            message: "Destination host should be set",
            trigger: "blur",
          },
        ],
        dstPort: [
          {
            required: true,
            message: "Destination port should be set",
            trigger: "blur",
          },
        ],
      },
    };
  },

  watch: {
    $route(to, from) {
      this.loadTunnels();
    },
  },
  created() {
    bus.$on("refresh-tunnels", (data) => {
      this.loadTunnels();
    });
  },
  mounted() {
    this.loadTunnels();
    this.loadModules();
    this.loadNetworks();
  },
  methods: {
    validateTunnelType() {
      return new Promise((resolve) => {
        resolve(this.newTunnel.type !== "");
      });
    },
    validateTunnelSettings() {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(async (valid) => {
          if (valid) {
            await this.createTunnel();
            bus.$emit("refresh-tunnels");
            this.newTunnelDialogVisible = false;
            this.resetform();
          }

          resolve(valid);
        });
      });
    },
    resetform() {
      this.newTunnel = Object.assign({}, this.serviceBootstrap);
      this.formKey = utils.randomString(5);
    },

    async createTunnel() {
      let that = this;
      return new Promise((resolve, reject) => {
        utils
          .$post(`/api/tunnels/create/`, this.newTunnel)
          .then((data) => {
            if (data.status == "error") {
              that.$notify.error({
                title: "Error",
                message: data.error,
              });
              reject();
            }
            if (data.status == "ok") {
              that.$notify.success({
                title: "Success",
                message: "Service created",
              });
              resolve();
            }
          })
          .catch((err) => {
            that.$notify.error({
              title: "Error",
              message: err,
            });
            reject();
          });
      });
    },
    async applyFilter() {
      this.applyLoading = true;
      await this.loadTunnels();
      this.applyLoading = false;
    },
    menuFilter(items, name) {
      return items.filter((i) => i.name == name);
    },
    isActiveMenuItem(k, v) {
      return this.$route.query[k] && this.$route.query[k] === v;
    },

    addParam(k, v) {
      utils.addQueryParam(this, k, v);
    },
    removeParam(k) {
      utils.removeQueryParam(this, k);
    },

    async loadTunnels() {
      this.tunnelsLoading = true;
      let merged = Object.assign(this.$route.query, { page: this.pageNum });
      merged = Object.assign(merged, this.filter);
      utils
        .$get(`/api/tunnels/?` + utils.objectToString(merged))
        .then((data) => {
          this.tunnelsLoading = false;
          if (data.status == "ok") {
            this.filteredTunnels = data.tunnels;
          } else {
            this.tunnelsLoading = false;
            this.$notify.error({
              title: "Error",
              message: "Failed to load tunnels",
            });
          }
        });
    },
    async loadModules() {
      let that = this;
      utils
        .$get(`/api/modules/`)
        .then((data) => {
          if (data.status == "ok") {
            this.modules = data.modules;
          }
        })
        .catch((err) => {
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
    async loadNetworks() {
      utils
        .$get(`/api/networks/`)
        .then((data) => {
          if (data && data.status == "ok") {
            this.networks = data.networks;
          }
        })
        .catch((err) => {
          this.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
  },
};
</script>
