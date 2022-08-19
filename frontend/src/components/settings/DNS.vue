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

          <div class="rule-group">

            <div class="row mb-0 mb-md-3">
              <div class="col-12 col-md-6 mb-3 mb-md-0">
                <el-form-item label="Resource records">
                  <el-select v-model="record.type" placeholder="Choose records">
                    <el-option v-for="type in recordTypes" :key="type" :value="type">
                    </el-option>
                  </el-select>
                </el-form-item>
              </div>
              <div class="col-12 col-md-6 mb-3 mb-md-0">
                <div v-if="record.type == 'CAA'">
                  <el-form-item label="Tag">
                    <el-input v-model="record.arg1" placeholder="Value tag"> </el-input>
                  </el-form-item>
                </div>
              </div>
            </div>

            <div class="row">
              <div class="col-12 col-md-6 mb-3 mb-md-0">
                <el-form-item label="Name">
                  <el-input v-model="record.name" placeholder="Write name"></el-input>
                </el-form-item>
              </div>
              <div class="col-12 col-md-6">
                <el-form-item label="Content">
                  <el-input v-model="record.content" placeholder="Write content"> </el-input>
                </el-form-item>
              </div>
            </div>

            <div class="control-rule-item mt-3 w-100 d-flex justify-content-end">
              <button class="btn btn-icon del" @click.prevent="removeRecord(idx)">
                Delete
                <i class="vik vik-delete"></i>
              </button>
            </div>

          </div>
        </div>

        <el-button class="btn btn-add-new-group" @click="addRecord()">
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
      recordBootstrap: { type: "", name: "", content: "", arg1: "" },
      settings: {
        recursive: false,
        resolvers: ["8.8.8.8", "1.1.1.1"],
        records: [],
      },
      recordTypes: ["A", "AAAA", "CNAME", "MX", "TXT", "CAA"],
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