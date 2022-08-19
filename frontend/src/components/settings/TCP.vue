<template>
  <div>
    <el-form :model="settings" ref="settingsForm" class="service-settings">
      <el-form-item label="Logging">
        <el-checkbox v-model="settings.log_request">Log request</el-checkbox>
        <el-checkbox v-model="settings.log_response">Log response</el-checkbox>
      </el-form-item>

      <h3 class="mt-4">Mode</h3>

      <el-form-item
        class="rule-group"
        prop="mode"
        :rules="{
          required: true,
          message: 'Mode is required',
          trigger: 'blur',
        }"
      >
        <el-radio-group
          v-model="settings.mode"
          size="small"
          class="radio-button-group mb-3"
        >
          <el-radio-button label="response">Custom response</el-radio-button>
          <el-radio-button label="proxy">Proxy</el-radio-button>
        </el-radio-group>
        <el-form-item v-if="settings.mode == 'proxy'" label="Destination addr"
          ><el-input
            v-model="settings.proxy_settings.destination"
            placeholder="for ex: 1.3.3.7:22"
          ></el-input>
        </el-form-item>
        <el-form-item v-if="settings.mode == 'response'" label="Template">
          <el-input
            type="textarea"
            :rows="2"
            placeholder="Response body"
            v-model="settings.response_settings.response"
          >
          </el-input>
        </el-form-item>
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

      settings: {
        log_request: true,
        log_response: false,
        mode: "response",
        proxy_settings: {},
        response_settings: { response: "" },
      },
    };
  },
  async mounted() {
    if (this.preset) {
      this.settings = this.preset;
    }
  },
  watch: {
    preset() {
      if (this.preset) {
        this.settings = this.preset;
      }
    },
  },

  methods: {
    defaultPort() {
      return 31337;
    },
    validate() {
      return new Promise((resolve, reject) => {
        this.$refs.settingsForm.validate((valid) => {
          if (valid) resolve(this.settings);
          else reject();
        });
      });
    },
  },
};
</script>