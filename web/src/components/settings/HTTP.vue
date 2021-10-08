<template>
  <div>
    <el-form :model="settings" ref="settingsForm" class="service-settings">
      <el-form-item label="Logging">
        <el-checkbox v-model="settings.log_request">Log request</el-checkbox>
        <el-checkbox v-model="settings.log_response">Log response</el-checkbox>
      </el-form-item>
      <el-form-item label="Encryption">
        <el-checkbox v-model="settings.tls.enabled">Use TLS</el-checkbox>
        <el-checkbox v-if="settings.tls.enabled" v-model="settings.tls.autocert"
          >Use autocert</el-checkbox
        >
      </el-form-item>
      <el-form-item label="Files">
        <el-checkbox v-model="settings.allow_file_upload"
          >Allow file upload</el-checkbox
        >
      </el-form-item>

      <h3 class="mt-5">Hosts</h3>
      <div v-for="(host, idx) in settings.hosts" :key="idx">
        <div class="rule-group">
          <div class="control-rule-item w-100 d-flex justify-content-end">
            <button
              class="btn btn-icon del"
              @click.prevent="removeRules(rules)"
            >
              <i class="vik vik-delete"></i>
            </button>
          </div>
          <el-form-item
            label="Hostname"
            :prop="'hosts.' + idx + '.hostname'"
            :rules="{
              required: true,
              message: 'Hostname can not be empty',
              trigger: 'blur',
            }"
          >
            <el-input v-model="host.hostname" placeholder="Hostname">
            </el-input>
          </el-form-item>
          <div
            class="location-group"
            v-for="(location, idx2) in host.locations"
            :key="idx2"
          >
            <div
              class="control-rule-item w-100 d-flex justify-content-end"
              v-if="host.locations.length > 1"
            >
              <button
                class="btn btn-icon del"
                @click.prevent="removeLocation(idx, idx2)"
              >
                <i class="vik vik-delete"></i>
              </button>
            </div>
            <el-form-item
              label="Location"
              :key="idx2"
              :prop="'hosts.' + idx + '.locations.' + idx2 + '.path'"
              :rules="{
                required: true,
                message: 'Path can not be empty',
                trigger: 'blur',
              }"
            >
              <el-input
                v-model="location.path"
                suffix-icon="el-icon-edit"
                placeholder="path"
              ></el-input>
            </el-form-item>

            <el-form-item
              label="Action"
              :prop="'hosts.' + idx + '.locations.' + idx2 + '.action_name'"
              :rules="{
                validator: checkAction,
                trigger: 'blur',
              }"
            >
              <el-radio-group
                @change="changeAction(idx, idx2)"
                v-model="location.action_name"
                size="small"
                class="radio-button-group"
              >
                <el-radio-button label="template"
                  >Render template</el-radio-button
                >
                <el-radio-button label="file">Serve file</el-radio-button>
                <el-radio-button label="folder">Serve folder</el-radio-button>
                <el-radio-button label="proxy">Proxy</el-radio-button>
              </el-radio-group>
            </el-form-item>
            <span v-if="location.action_name == 'template'">
              <el-form-item label="Status code"
                ><el-input-number
                  v-model="location.action_data.status"
                  suffix-icon="el-icon-date"
                  placeholder=""
                  :min="100"
                  :key="599"
                ></el-input-number>
              </el-form-item>
              <el-form-item label="Headers"
                ><div
                  v-for="(header, hidx) in location.action_data.headers"
                  :key="hidx"
                  class="headers-group"
                >
                  <div class="row">
                    <div class="col">
                      <el-input placeholder="Name" v-model="header.name">
                      </el-input>
                    </div>
                    <div class="col">
                      <el-input placeholder="Value" v-model="header.value">
                      </el-input>
                    </div>
                    <div class="col-auto d-flex align-content-center">
                      <el-button
                        class="btn btn-mute del"
                        @click="removeHeader(idx, idx2, hidx)"
                      >
                        <i class="vik vik-delete"></i>
                      </el-button>
                    </div>
                  </div>
                </div>

                <el-button
                  class="btn btn-add-new-group"
                  @click="addHeader(idx, idx2)"
                >
                  <i class="vik vik-plus"></i>Add header
                </el-button>
              </el-form-item>

              <el-form-item label="Template">
                <el-input
                  type="textarea"
                  :rows="2"
                  placeholder="Response body"
                  v-model="location.action_data.template"
                >
                </el-input>
                <el-alert class="mt-2" :closable="false" type="info" show-icon>
                  You can use variables like
                  <div v-pre>
                    {{ req.headers.User_Agent }}, {{ req.uri }},
                    {{ req.remote_addr }}
                  </div>
                </el-alert>
              </el-form-item>
            </span>

            <el-form-item v-if="location.action_name == 'file'" label="File">
              <el-select
                v-model="location.action_data.hash"
                placeholder="Choose a file"
                suffix-icon="el-icon-file"
                filterable
              >
                <el-option
                  :key="file.hash"
                  v-for="file in files"
                  :label="file.file_name"
                  :value="file.hash"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item
              v-if="location.action_name == 'folder'"
              label="Folder"
            >
              <el-input
                v-model="location.action_data.folder"
                suffix-icon="el-icon-edit"
                placeholder="Path"
              ></el-input>
            </el-form-item>

            <span v-if="location.action_name == 'proxy'">
              <el-form-item label="Destination">
                <el-input
                  v-model="location.action_data.destination"
                  suffix-icon="el-icon-edit"
                  placeholder="https://desti.nation"
                ></el-input>
              </el-form-item>
              <el-form-item label="Custom host header (if needed)">
                <el-input v-model="location.action_data.custom_host"></el-input>
              </el-form-item>
            </span>
            <!-- end of proxy tab -->
          </div>
          <el-form-item>
            <el-button class="btn btn-add-new-group" @click="addLocation(idx)">
              <i class="vik vik-plus"></i>Add location
            </el-button>
          </el-form-item>
        </div>
      </div>

      <el-form-item>
        <el-button class="btn btn-outline-secondary" @click="addHostname">
          <i class="vik vik-plus"></i>Add new host
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>
<script>
import * as utils from "@/utils.js";

export default {
  name: "http-settings",
  props: ["preset"],
  data() {
    return {
      files: [],
      custom_settings: {
        template: { template: "", status: 200, headers: [] },
        file: { hash: "" },
        folder: { folder: "" },
        proxy: { destination: "", custom_host: "" },
      },

      hostBootstrap: {
        hostname: "host.name",
        locations: [
          {
            path: "/",
            action_name: "template",
            action_data: { template: "hi", status: 200, headers: [] }, // default
          },
        ],
      },
      settings: {
        allow_file_upload: true,
        log_request: true,
        log_response: false,
        tls: { enabled: false, autocert: false },
        hosts: [
          {
            hostname: "*",
            locations: [
              {
                path: "/",
                action_name: "template",
                action_data: {
                  template: "Hello {{ client.ip }}",
                  status: 200,
                  headers: [],
                },
              },
            ],
          },
        ],
      },
    };
  },
  async mounted() {
    if (this.preset) {
      console.log("preset changed..", this.preset);
      // todo: load file by hash
      this.settings = this.preset;
    }
    utils
      .$get(`/api/files/?` + utils.objectToString({ page: this.currentPage }))
      .then((data) => {
        if (data.status == "ok") {
          this.files = data.files;
          this.total = data.total;
        }
      })
      .catch((err) => {
        this.$notify.error({
          title: "Error",
          message: err,
        });
      });
  },
  watch: {
    preset() {
      if (this.preset) {
        console.log("preset changed..", this.preset);
        this.settings = this.preset;
      }
    },
  },

  methods: {
    defaultPort() {
      return 80;
    },
    validate() {
      return new Promise((resolve, reject) => {
        this.$refs.settingsForm.validate((valid) => {
          if (valid) resolve(this.settings);
          else reject();
        });
      });
    },
    removeLocation(hostID, locationID) {
      this.settings.hosts[hostID].locations.splice(locationID, 1);
    },

    addHeader(hostID, locationID) {
      this.settings.hosts[hostID].locations[
        locationID
      ].action_data.headers.push({
        name: "",
        value: "",
      });
    },
    removeHeader(hostID, locationID, headerID) {
      this.settings.hosts[hostID].locations[
        locationID
      ].action_data.headers.splice(headerID, 1);
    },
    removeRules(item) {
      var index = this.newServiceSettingsForm.rules.indexOf(item);
      if (index !== -1) {
        this.form.rules.splice(index, 1);
      }
    },
    addRules() {
      this.form.rules.push({
        key: Date.now(),
        value: "",
      });
    },
    addLocation(index) {
      this.settings.hosts[index].locations.push({
        path: "/" + utils.randomString(6),
        action_name: "",
        action_data: {},
      });
    },
    addHostname() {
      this.settings.hosts.push(this.hostBootstrap);
    },
    getActionNode(obj, path) {
      for (
        var i = 0, path = path.split("."), len = path.length;
        i < len - 1;
        i++
      ) {
        obj = obj[path[i]];
      }
      return obj;
    },

    changeAction(host, location) {
      this.settings.hosts[host].locations[location].action_data = Object.assign(
        {},
        this.custom_settings[
          this.settings.hosts[host].locations[location].action_name
        ]
      );
    },
    checkAction(rule, value, callback) {
      console.log("call checkAction", rule.field, value);
      switch (value) {
        case "template":
          callback();
          break;
        case "file":
          var action = this.getActionNode(this.settings, rule.field);
          if (!action.action_data.hash || action.action_data.hash === "") {
            callback(new Error("Choose a file"));
          } else {
            callback();
          }
          break;
        case "folder":
          var action = this.getActionNode(this.settings, rule.field);
          if (!action.action_data.folder || action.action_data.folder === "") {
            callback(new Error("Choose a folder"));
            return;
          }
          break;
        case "proxy":
          var action = this.getActionNode(this.settings, rule.field);
          if (
            !action.action_data.destination ||
            action.action_data.destination === ""
          ) {
            callback(new Error("Destination URL shold be set"));
            return;
          }
          break;
        default:
          callback(new Error("Action is required"));
          break;
      }
      return true;
    },
  },
};
</script>