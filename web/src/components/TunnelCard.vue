<template>
  <div>
    <article class="mini-card" v-loading="loading">
      <div class="mini-card-type native">
        <span>{{ listenAddr }}</span>
      </div>
      <div class="mini-card-left-side">
        <div class="mini-article-card-title">
          <h3 class="d-flex align-items-center">
            {{ tunnel.type
            }}<span
              class="status-tag"
              v-bind:style="{ 'background-color': colorByStatus }"
            ></span>
          </h3>

          <span
            class="atom-source"
            style="margin: auto"
            v-if="tunnel.type == `HTTP` && tunnel.publicAddr !== ``"
            ><a
              :key="p"
              v-for="p in ['http://', 'https://']"
              :href="p + tunnel.publicAddr"
              target="_blank"
              >{{ p }}{{ tunnel.publicAddr }}<br /></a
          ></span>
        </div>
      </div>

      <div class="mini-card-right-side">
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
            @click="toggleTunnelState"
          >
          </el-button>
        </div>
      </div>
    </article>

    <el-dialog
      :title="`Edit tunnel `"
      :visible.sync="editTunnelDialogVisible"
      class="modal-new-modules"
    >
      <form-wizard finishButtonText="Update tunnel">
        <tab-content title="Settings" :before-change="validateTunnelSettings">
          <el-form
            :model="localTunnel"
            ref="form"
            class="service-settings"
            :rules="ruleValidate"
            @submit.native.prevent
          >
            <el-form-item label="Destination host" prop="dstHost">
              <el-select
                filterable
                allow-create
                v-model="localTunnel.dstHost"
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
                      v-model="localTunnel.dstPort"
                      :min="1"
                      :max="65535"
                    ></el-input-number>
                  </el-form-item>
                </div>
                <div class="mr-4">
                  <el-form-item label="Destination TLS">
                    <el-checkbox v-model="localTunnel.dstTLS"></el-checkbox>
                  </el-form-item>
                </div>
                <div class="">
                  <el-form-item label="Auto start">
                    <el-checkbox v-model="localTunnel.autoStart"></el-checkbox>
                  </el-form-item>
                </div>
              </div>
            </el-form-item>
          </el-form>
        </tab-content>
        <template v-slot:custom-buttons-left>
          <button
            role="button"
            tabindex="0"
            v-loading="removing"
            class="btn btn-link del"
            @click="removeTunnel"
          >
            <i class="vik vik-delete"></i> Remove tunnel
          </button>
        </template>
      </form-wizard>
    </el-dialog>
  </div>
</template>
<script>
import { FormWizard, TabContent } from "@/components/vue-form-wizard";
import "@/components/vue-form-wizard/assets/wizard.scss";
import { bus } from "@/bus.js";
import * as utils from "@/utils.js";

export default {
  props: ["tunnel"],
  components: {
    FormWizard,
    TabContent,
  },
  data() {
    return {
      finishBtnText: "Save",
      networks: [],
      loading: false,
      removing: false,
      ctlBtnLoading: false,
      editTunnelDialogVisible: false,
      localTunnel: { ...this.tunnel },
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
  mounted() {
    this.tunnelInfo = Object.assign({}, this.tunnel);
  },
  computed: {
    colorByStatus() {
      return this.tunnel.connected ? "#B9E8BC" : "#EA8C8C";
    },
    listenAddr() {
      return `${this.tunnel.dstHost}:${this.tunnel.dstPort}`;
    },

    ctlButtonStyle() {
      let icons = ["vik"];
      if (this.tunnel.connected) icons.push("vik-stop");
      if (!this.tunnel.connected) icons.push("vik-run");
      return icons.join(" ");
    },
  },
  methods: {
    toggleTunnelState() {
      let that = this;
      that.ctlBtnLoading = true;
      utils
        .$post(`/api/tunnels/${this.tunnel.connected ? "stop" : "start"}/`, {
          hash: this.tunnel.hash,
        })
        .then((data) => {
          if (data.status === "error") {
            that.$notify.error({
              title: "Error",
              message: data.error,
            });
          }
          that.ctlBtnLoading = false;
          bus.$emit("refresh-tunnels");
        })
        .catch((err) => {
          that.ctlBtnLoading = false;
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },

    showEditDialog() {
      this.loadNetworks();
      this.editTunnelDialogVisible = true;
    },
    removeTunnel() {
      let that = this;
      this.removing = true;
      utils
        .$post(`/api/tunnels/remove/`, this.localTunnel)
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
              message: "Tunnel removed",
            });
            setTimeout(() => {
              bus.$emit("refresh-tunnels");
            }, 1500);
            that.editTunnelDialogVisible = false;
          }
        })
        .catch((err) => {
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
    updateTunnel() {
      let that = this;
      utils
        .$post(`/api/tunnels/update/`, this.localTunnel)
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
            that.editTunnelDialogVisible = false;
          }
        })
        .catch((err) => {
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
    validateTunnelSettings() {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(async (valid) => {
          if (valid) {
            await this.updateTunnel();
            bus.$emit("refresh-tunnels");
            this.editTunnelDialogVisible = false;
          }
          resolve(valid);
        });
      });
    },

    async validateModuleSettings() {
      let that = this;
      this.finishBtnText = "Saving...";
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(async (valid) => {
          if (valid) {
            this.editTunnelDialogVisible = false;
            await this.updateTunnel();
            setTimeout(() => {
              bus.$emit("refresh-tunnels");
            }, 2000);
          }
          that.finishBtnText = "Save";
          resolve(valid);
        });
      });
    },
    async loadNetworks() {
      utils
        .$get(`/api/networks/`)
        .then((data) => {
          if (data.status == "ok") {
            this.networks = data.networks;
          }
        })
        .catch((err) => {
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },
  },
};
</script>
