<template>
  <div>
    <article class="mini-card" v-loading="loading">
      <div class="mini-card-type native">
        <span>{{ baseProto }}</span>
      </div>
      <div class="mini-card-left-side">
        <div class="mini-article-card-title">
          <h3 class="d-flex align-items-center">
            {{ service.serviceName
            }}<span
              class="status-tag"
              v-bind:style="{ 'background-color': colorByStatus }"
            ></span>
          </h3>
          <span class="atom-source" style="margin: auto">{{ listenAddr }}</span>
        </div>
      </div>

      <div class="mini-card-right-side">
        <div class="tags">
          <!-- 
          <div v-if="parsedTags(service).length">
            <span v-for="tag in parsedTags(service)" :key="tag">{{ tag }}</span>
          </div>
          -->
        </div>
        <div class="btn-group">
          <el-button
            class="btn-sm btn-icon btn"
            :loading="loading"
            icon="vik vik-settings"
            @click="showEditDialog"
          >
          </el-button>
          <el-button
            class="btn-sm btn-icon btn"
            :loading="ctlBtnLoading"
            :icon="ctlButtonStyle"
            @click="toggleServiceState"
          >
          </el-button>
        </div>
      </div>
    </article>

    <el-dialog
      :title="`Edit service ` + localService.serviceName"
      :visible.sync="editServiceDialogVisible"
      class="modal-new-modules"
    >
      <form-wizard :finishButtonText="finishBtnText">
        <tab-content title="Settings" :before-change="validateModuleSettings">
          <!-- global settings -->
          <el-form
            :model="localService"
            ref="localServiceForm"
            class="service-settings"
            :rules="ruleValidate"
            @submit.native.prevent
          >
            <el-form-item label="Service name" prop="serviceName">
              <el-input
                v-model="localService.serviceName"
                placeholder="Service name"
              ></el-input>
            </el-form-item>
            <el-form-item label="Listen IP" prop="listenIP">
              <el-select v-model="localService.listenIP" placeholder="IP">
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
                      v-model="localService.listenPort"
                      suffix-icon="el-icon-date"
                      placeholder="Port number"
                      :min="1"
                      :key="65535"
                    ></el-input-number>
                  </el-form-item>
                </div>
                <div class="">
                  <el-form-item label="Auto start">
                    <el-checkbox v-model="localService.autoStart"></el-checkbox>
                  </el-form-item>
                </div>
              </div>
            </el-form-item>
          </el-form>
          <!-- module settings -->
          <component
            v-bind:is="settingsView"
            :preset="moduleSettings"
            ref="moduleSettings"
          ></component>
        </tab-content>

        <template v-slot:custom-buttons-left>
          <button
            role="button"
            tabindex="0"
            v-loading="removing"
            class="btn btn-link del"
            @click="removeService"
          >
            <i class="vik vik-delete"></i> Remove service
          </button>
        </template>
      </form-wizard>
    </el-dialog>
  </div>
</template>
<script>
import { FormWizard, TabContent } from "@/components/vue-form-wizard";
import "@/components/vue-form-wizard/assets/wizard.scss";
import Empty from "@/components/settings/Empty.vue";
import DNS from "@/components/settings/DNS.vue";
import RogueMysql from "@/components/settings/RogueMysql.vue";
import HTTP from "@/components/settings/HTTP.vue";
import TCP from "@/components/settings/TCP.vue";
import FTP from "@/components/settings/FTP.vue";
import { bus } from "@/bus.js";
import * as utils from "@/utils.js";
import { getNetworks } from "@/api.js";

export default {
  props: ["service"],
  components: {
    FormWizard,
    TabContent,
  },
  data() {
    return {
      finishBtnText: "Save",
      moduleSettingsComp: {
        Empty: Empty,
        DNS: DNS,
        HTTP: HTTP,
        TCP: TCP,
        FTP: FTP,
        Rogue_MySQL_Server: RogueMysql,
      },
      networks: [],
      loading: false,
      removing: false,
      ctlBtnLoading: false,
      editServiceDialogVisible: false,
      localService: { ...this.service },
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
            message: "Port should be set",
            trigger: "blur",
          },
        ],
      },
    };
  },
  mounted() {
    this.serviceInfo = Object.assign({}, this.service);
  },
  computed: {
    colorByStatus() {
      return this.service.active ? "#B9E8BC" : "#EA8C8C";
    },
    listenAddr() {
      return `${this.service.listenIP}:${this.service.listenPort}`;
    },
    baseProto() {
      return this.service.baseProto.toUpperCase();
    },
    getType() {
      return this.service["type"];
    },
    ctlButtonStyle() {
      let icons = ["vik"];
      if (this.service.active) icons.push("vik-stop");
      if (!this.service.active) icons.push("vik-run");
      return icons.join(" ");
    },

    settingsView() {
      let view = this.moduleSettingsComp[this.localService.moduleName];
      return view ? view : Empty; // if module has no settings view render default Empty view
    },
    moduleSettings() {
      if (this.localService.settings && this.localService.settings.length > 0) {
        try {
          let s = JSON.parse(this.localService.settings);
          return s;
        } catch (e) {
          console.log("Failed to decode module preset", e);
        }
      }
      return {};
    },
    localServiceDefaultPort() {
      return this.$refs.moduleSettings.defaultPort
        ? this.$refs.moduleSettings.defaultPort
        : 8080;
    },
  },
  methods: {
    toggleServiceState() {
      let that = this;
      that.ctlBtnLoading = true;
      utils
        .$post(`/api/services/${this.service.active ? "stop" : "start"}/`, {
          hash: this.service.hash,
        })
        .then((data) => {
          if (data.status === "error") {
            that.$notify.error({
              title: "Error",
              message: data.error,
            });
          }
          that.ctlBtnLoading = false;
          bus.$emit("refresh-services");
        })
        .catch((err) => {
          that.ctlBtnLoading = false;
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
    fixModuleName(name) {
      return name.replace(/_/g, " ");
    },
    showEditDialog() {
      this.loadNetworks();
      this.editServiceDialogVisible = true;
    },
    parsedTags(raw) {
      return ["tcp", "DNS"];
      /*
      try {
        let p = JSON.parse(raw);
        return p;
      } catch (e) {
        return [];
      }
      */
    },
    removeService() {
      let that = this;
      this.removing = true;
      utils
        .$post(`/api/services/remove/`, this.localService)
        .then((data) => {
          this.removing = false;
          if (data.status == "error") {
            that.$notify.error({
              title: "Error",
              message: data.error,
            });
            return;
          }
          if (data.status == "ok") {
            that.$notify.success({
              title: "Success",
              message: "Service removed",
            });
            setTimeout(() => {
              bus.$emit("refresh-services");
            }, 1500);
            that.editServiceDialogVisible = false;
          }
        })
        .catch((err) => {
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
    updateService() {
      let that = this;
      utils
        .$post(`/api/services/update/`, this.localService)
        .then((data) => {
          if (data.status == "error") {
            that.$notify.error({
              title: "Error",
              message: data.error,
            });
            return;
          }
          if (data.status == "ok") {
            that.$notify.success({
              title: "Success",
              message: "Service updated",
            });
            that.editServiceDialogVisible = false;
          }
        })
        .catch((err) => {
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
    validateModuleSettings() {
      let that = this;
      this.finishBtnText = "Saving...";
      return new Promise((resolve, reject) => {
        this.$refs.localServiceForm.validate((valid) => {
          if (valid) {
            if (this.$refs.moduleSettings.validate) {
              this.$refs.moduleSettings.validate().then(
                (settings) => {
                  this.localService.moduleSettings = settings;

                  this.updateService();
                  // close dialog
                  this.localServiceDialogVisible = false;
                  // update services list
                  setTimeout(() => {
                    bus.$emit("refresh-services");
                  }, 2500);
                  resolve(true);
                },
                () => {
                  resolve(false);
                }
              );
            } else {
              // if module has no custom settings
              this.localService.moduleSettings = {};
              this.updateService(); // updateService
              this.localServiceDialogVisible = false;
              // update services list
              bus.$emit("refresh-services");
              resolve(true);
            }
          }
          that.finishBtnText = "Save";
          resolve(valid);
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
