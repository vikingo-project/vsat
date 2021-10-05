<template>
  <div>
    <el-form
      :model="settings"
      ref="settingsForm"
      :rules="rules"
      @submit.native.prevent
    >
      <el-form-item label="DNS records" prop="records">
        <div v-for="(record, idx) in settings.records" :key="idx">
          <div class="control-rule-item w-100 d-flex justify-content-end">
            <button class="btn btn-icon del" @click.prevent="removeRecord(idx)">
              <i class="vik vik-delete"></i>
            </button>
          </div>
          <el-select v-model="record.type"
            ><el-option
              v-for="type in recordTypes"
              :key="type"
              :value="type"
            ></el-option
          ></el-select>
          <el-input v-model="record.name" placeholder="Name"></el-input>
          <el-input v-model="record.content" placeholder="Content"> </el-input>
        </div>
        <el-button class="btn w-100 btn-link" @click="addRecord()">
          <i class="vik vik-plus"></i>Add record
        </el-button>
        <!-- 
        <el-input
          v-model="settings.records"
          type="textarea"
          :rows="3"
          :placeholder="`# for example:\n127.0.0.1 vikingo.org\n127.0.0.1 www.vikingo.org`"
        >
        </el-input>
        -->
      </el-form-item>
      <el-form-item> </el-form-item>
      <el-form-item label="Recursive mode" prop="recursive">
        <el-checkbox v-model="settings.recursive">Enabled</el-checkbox>
      </el-form-item>
      <el-form-item label="Resolvers" prop="resolvers">
        <el-select
          v-model="settings.resolvers"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="Resolver address"
        >
          <el-option
            v-for="(resolver, idx) in settings.resolvers"
            :key="idx"
            :value="resolver"
          >
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
  </div>
</template>
<script>
export default {
  props: ["preset"],
  data() {
    return {
      recordBootstrap: { type: "", name: "", content: "" },
      settings: {
        recursive: false,
        resolvers: ["8.8.8.8", "1.1.1.1"],
        records: [],
      },
      recordTypes: ["A", "AAAA", "CNAME", "MX", "TXT"],
      rules: {},
    };
  },
  mounted() {
    if (this.preset) {
      this.settings = this.preset;
    }
  },
  methods: {
    addRecord() {
      this.settings.records.push(Object.assign({}, this.recordBootstrap));
    },

    removeRecord(idx) {
      this.settings.records.splice(idx, 1);
    },

    defaultPort() {
      return 53;
    },
    validate() {
      console.log("subform validate");
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