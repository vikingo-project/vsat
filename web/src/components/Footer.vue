<template>
  <footer>
    <div class="container text-center text-md-left">
      <div class="row">
        <div class="col-md-1 footer-block">
          <div class="footline">
            <a href="/" style="vertical-align: text-bottom">
              <img src="@/assets/logo.svg" alt />
            </a>
          </div>
        </div>
        <div class="col-md-3 footer-block">
          <div class="footline">
            <a href="https://docs.vikingo.org/v/en/vsat"> Docs </a>
            |
            <a href="javascript:;" @click="showAbout = true"> About </a>
          </div>
        </div>

        <div class="col-md-4 footer-block ml-auto">
          <div class="copyright">
            <span>© 2020-2021 Vikingo Satellite</span>
            <br />
          </div>
        </div>
      </div>
    </div>
    <el-dialog title="About" :visible.sync="showAbout" class="modalSave">
      Vikingo Satellite {{ version }} / build {{ build }}</el-dialog
    >
  </footer>
</template>
<script>
import * as utils from "@/utils.js";
export default {
  async mounted() {
    utils.$get(`/api/about/`).then((data) => {
      this.version = data.version;
      this.build = data.build;
    });
  },
  data() {
    return {
      showAbout: false,
      version: "",
      build: "",
    };
  },
};
</script>
