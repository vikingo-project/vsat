<template>
  <div class="container">
    <div class="row">
      <div class="col-md-3 col-lg-2 col-xl-2 navbar-collapse sidebar">
        <div class="category-list">
          <h2>Filter</h2>
          <el-form>
            <el-form-item label="Service name">
              <el-input v-model="filter.service_name"></el-input>
            </el-form-item>
            <el-form-item label="Base protocol">
              <el-select v-model="filter.base_proto" clearable>
                <el-option key="udp" label="UDP" value="UDP"> </el-option>
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
        <div class="content" v-loading="servicesLoading">
          <div class="market-page">
            <h1>Services</h1>
            <button
              class="filter-btn btn-sm btn-primary"
              @click="filterDrawer = true"
            >
              <i class="vik vik-adjust mr-2"></i>Filter
            </button>
          </div>
          <div class="filter-result" v-if="paramsToTags().length > 0">
            <div class="filter-result-text">
              <span>filtered by</span>
            </div>
            <div class="tags">
              <span
                class="action"
                v-for="(tag, idx) in paramsToTags()"
                :key="idx"
              >
                {{ tag.value }}
                <a
                  type="button"
                  class="tag-delete"
                  @click="removeParam(tag.name)"
                  >×</a
                >
              </span>
            </div>
          </div>

          <div class="row item-container atom-list">
            <div class="col-lg-6">
              <article
                @click="newServiceDialogVisible = true"
                class="btn btn-primary btn-new"
              >
                <i class="vik vik-plus"></i>Add new
              </article>
            </div>

            <div
              class="col-lg-6"
              v-for="service in filteredServices"
              :key="service.hash"
            >
              <ServiceCard :service="service"></ServiceCard>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog
      title="New service"
      :visible.sync="newServiceDialogVisible"
      class="modal-new-modules"
    >
      <form-wizard :key="formKey" finishButtonText="Create service">
        <tab-content title="Choose module" :before-change="validateModule">
          <div class="content-mod">
            <div class="list-modules">
              <div class="select-modules">
                <el-radio
                  v-model="newService.moduleName"
                  :label="module.name"
                  v-for="(module, idx) in modules"
                  :key="idx"
                >
                  <h3>{{ fixModuleName(module.name) }}</h3>
                  <p>{{ module.description }}</p>
                </el-radio>
              </div>
            </div>
          </div>
        </tab-content>
        <tab-content title="Settings" :before-change="validateModuleSettings">
          <!-- global settings -->
          <el-form
            :model="newService"
            ref="newServiceForm"
            class="service-settings"
            :rules="ruleValidate"
            @submit.native.prevent
          >
            <el-form-item label="Service name" prop="serviceName">
              <el-input
                v-model="newService.serviceName"
                placeholder="Service name"
              ></el-input>
            </el-form-item>
            <el-form-item label="Listen IP" prop="listenIP">
              <el-select
                v-model="newService.listenIP"
                placeholder="IP"
                suffix-icon="el-icon-date"
              >
                <el-option value="0.0.0.0">All interfaces (0.0.0.0)</el-option>
                <el-option
                  v-for="net in networks"
                  :value="net.ip"
                  :key="net.ip"
                  >{{ net.name + ` (` + net.ip + `)` }}</el-option
                >
              </el-select>
            </el-form-item>

            <el-form-item>
              <div class="d-flex">
                <div class="mr-4">
                  <el-form-item label="Port" prop="listenPort">
                    <el-input-number
                      v-model="newService.listenPort"
                      suffix-icon="el-icon-date"
                      placeholder="Set port"
                      :min="1"
                      :max="65535"
                    ></el-input-number>
                  </el-form-item>
                </div>
                <div class="">
                  <el-form-item label="Auto start">
                    <el-checkbox v-model="newService.autoStart"></el-checkbox>
                  </el-form-item>
                </div>
              </div>
            </el-form-item>
          </el-form>
          <!-- module settings -->
          <component v-bind:is="settingsView" ref="moduleSettings"></component>
        </tab-content>
      </form-wizard>
    </el-dialog>

    <el-drawer
      title="Filter"
      direction="rtl"
      :append-to-body="true"
      :visible.sync="filterDrawer"
      custom-class="filter-drawer"
    >
      <div class="drawer-body">
        <div class="category-list">
          <el-form>
            <el-form-item label="Service name">
              <el-input v-model="filter.service_name"></el-input>
            </el-form-item>
            <el-form-item label="Base protocol">
              <el-select v-model="filter.base_proto" clearable>
                <el-option key="udp" label="UDP" value="UDP"> </el-option>
                <el-option key="tcp" label="TCP" value="TCP"> </el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </div>
      </div>

      <div class="drawer-footer mt-auto w-100">
        <el-button
          :loading="false"
          class="btn btn-primary btn-lg"
          @click="applyFilter"
          icon="vik vik-search"
        >
          Apply
        </el-button>
      </div>
    </el-drawer>
  </div>
</template>
<script>
import { bus } from "@/bus.js";
import * as utils from "@/utils.js";
import ServiceCard from "@/components/ServiceCard.vue";
import { FormWizard, TabContent } from "@/components/vue-form-wizard";
import "@/components/vue-form-wizard/assets/wizard.scss";
import { getServices, getModules, getNetworks, createService } from "@/api.js";

import Empty from "@/components/settings/Empty.vue";
import DNS from "@/components/settings/DNS.vue";
import RogueMysql from "@/components/settings/RogueMysql.vue";
import HTTP from "@/components/settings/HTTP.vue";
import TCP from "@/components/settings/TCP.vue";
import FTP from "@/components/settings/FTP.vue";

export default {
  components: {
    FormWizard,
    TabContent,
    ServiceCard,
  },

  data() {
    this.serviceBootstrap = {
      serviceName: "",
      moduleName: "",
      listenIP: ``,
      listenPort: 8080,
      autoStart: false,
      moduleSettings: {},
    };
    return {
      filteredServices: [],
      filter: { service_name: "", base_proto: "" },
      modules: [],
      networks: [],
      newService: { ...this.serviceBootstrap },
      filterDrawer: false,

      formKey: "default",
      moduleSettingsComp: {
        Empty: Empty,
        DNS: DNS,
        HTTP: HTTP,
        TCP: TCP,
        FTP: FTP,
        Rogue_MySQL_Server: RogueMysql,
      },

      menuItems: [
        { name: "type", value: "TCP" },
        { name: "type", value: "UDP" },
      ],
      tags: [],
      newServiceDialogVisible: false,
      servicesLoading: false,
      ruleValidate: {
        serviceName: [
          {
            required: true,
            message: "Service name cannot be empty",
            trigger: "blur",
          },
        ],
        listenIP: [
          {
            required: true,
            message: "Listen IP should be set",
            trigger: "blur",
          },
        ],
        listenPort: [
          {
            required: true,
            message: "Listen port should be set",
            trigger: "blur",
          },
        ],
      },
    };
  },
  computed: {
    settingsView() {
      let view = this.moduleSettingsComp[this.newService.moduleName];
      return view ? view : Empty; // if module has no settings view render default Empty view
    },
  },
  watch: {
    $route(to, from) {
      this.loadServices();
    },
  },
  created() {
    bus.$on("refresh-services", (data) => {
      this.loadServices();
    });
  },
  mounted() {
    this.loadServices();
    this.loadModules();
    this.loadNetworks();
  },
  methods: {
    validateModule() {
      return new Promise((resolve) => {
        let status =
          this.newService.moduleName && this.newService.moduleName !== "";
        if (status) {
          this.newService.serviceName =
            this.fixModuleName(this.newService.moduleName) + " service";
          if (this.$refs.moduleSettings.defaultPort) {
            this.newService.listenPort =
              this.$refs.moduleSettings.defaultPort();
          }
        }
        resolve(status);
      });
    },
    validateModuleSettings() {
      return new Promise((resolve, reject) => {
        this.$refs.newServiceForm.validate(async (valid) => {
          if (valid) {
            if (this.$refs.moduleSettings.validate) {
              this.$refs.moduleSettings.validate().then(
                async (settings) => {
                  this.newService.moduleSettings = settings;
                  await this.createService();
                  // close dialog
                  this.newServiceDialogVisible = false;
                  // reset form fields
                  this.resetNewServiceForm();
                  // update services list
                  bus.$emit("refresh-services");
                  resolve(true);
                },
                () => {
                  resolve(false);
                }
              );
            } else {
              // if module has no custom settings
              this.newService.moduleSettings = {};
              await this.createService();
              this.newServiceDialogVisible = false;
              this.resetNewServiceForm();

              // update services list
              bus.$emit("refresh-services");
              resolve(true);
            }
          }
          resolve(valid);
        });
      });
    },
    resetNewServiceForm() {
      this.newService = Object.assign({}, this.serviceBootstrap);
      this.formKey = utils.randomString(5);
    },

    async createService() {
      return new Promise((resolve, reject) => {
        createService(this.newService)
          .then((data) => {
            console.log("data=", data);
            this.$notify.success({
              title: "Success",
              message: "Service created",
            });
            resolve();
          })
          .catch((err) => {
            this.$notify.error({
              title: "Error",
              message: err,
            });
            reject();
          });

        /*
        utils
          .$post(`/api/services/create/`, this.newService)
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
            
          });
          */
      });
    },
    async applyFilter() {
      this.applyLoading = true;
      await this.loadServices();
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
    paramsToTags() {
      let tags = [];
      let obj = Object.assign({}, this.$route.query);
      for (let name in obj) {
        if (name == "cat" || name == "type" || name == "tag") {
          tags.push({ name: name, value: obj[name] });
        }
      }
      return tags;
    },

    fixModuleName(name) {
      return name.replace(/_/g, " ");
    },

    async loadServices() {
      this.servicesLoading = true;
      let params = Object.assign(this.$route.query, { page: this.pageNum });
      params = Object.assign(params, this.filter);

      getServices(utils.objectToString(params))
        .then((data) => {
          this.servicesLoading = false;
          this.filteredServices = data.Records;
        })
        .catch((e) => {
          this.servicesLoading = false;
          this.$notify.error({
            title: "Error",
            message: "Failed to load services",
          });
        });
    },
    async loadModules() {
      getModules()
        .then((data) => {
          this.modules = data.Records;
        })
        .catch((e) => {
          this.$notify.error({
            title: "Error",
            message: e,
          });
        });
    },
    async loadNetworks() {
      getNetworks()
        .then((data) => {
          this.networks = data.Records;
        })
        .catch((e) => {
          this.$notify.error({
            title: "Error",
            message: e,
          });
        });
    },
  },
};
</script>

<style>
.time {
  font-size: 13px;
  color: #999;
}

.bottom {
  margin-top: 13px;
  line-height: 12px;
}

.button {
  padding: 0;
  float: right;
}

.image {
  width: 100%;
  display: block;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both;
}
</style>
