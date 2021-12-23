<template>
  <div class="login">
    <div class="login-page">
      <div class="title">
        <h1 class="text-center">Vikingo Satellite</h1>

        <el-form :inline="true" :model="form" ref="formRef" class="form-login">
          <el-form-item label="Access token">
            <el-input v-model="form.token" placeholder="Enter token"></el-input>
          </el-form-item>
          <el-form-item class="ml-md-3 mt-3 mt-md-0 w-100 w-md-auto">
            <el-button
              :loading="loading"
              class="btn btn-primary w-100"
              @click="checkToken"
              >Sign in<i class="ml-2 mr-0 el-icon-right"></i
            ></el-button>
          </el-form-item>
        </el-form>

        <!-- <a href="/" class="text-center link-login">Where can I find token?</a>-->
      </div>
    </div>
  </div>
</template>

<script>
import { checkAuth } from "@/auth";
export default {
  data() {
    return {
      loading: false,
      form: {
        token: "",
      },
    };
  },
  methods: {
    checkToken() {
      let that = this;
      this.loading = true;
      this.$refs.formRef.validate((valid) => {
        if (valid) {
          checkAuth(that.form.token).then((res) => {
            console.log("check auth =>", res);
            that.loading = false;
          });
        }
      });
    },
  },
};
</script>

