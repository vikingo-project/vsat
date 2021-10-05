<template>
  <div>
    <el-form
      :model="settings"
      ref="settingsForm"
      :rules="rules"
      @submit.native.prevent
    >
      <el-form-item label="Path to file" prop="filepath">
        <el-input
          v-model="settings.filepath"
          placeholder="for example: /etc/passwd"
        ></el-input>
      </el-form-item>
    </el-form>
  </div>
</template>
<script>
export default {
  props: ["preset"],
  data() {
    return {
      settings: {
        filepath: "",
      },

      rules: {},
    };
  },
  mounted() {
    if (this.preset) {
      this.settings = this.preset;
    }
  },
  methods: {
    defaultPort() {
      return 3306;
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