<template>
  <footer>
    <div class="container text-center text-md-left">
      <div class="row">
        <div class="col mb-3 mb-md-0 footer-block">
          <div class="footline">
            <a href="/" style="vertical-align: text-bottom">
              <img src="@/assets/logo.svg" alt />
            </a>
          </div>
        </div>
        <div class="col-md-3 mb-2 mb-md-0 text-center footer-block">
          <div class="footline">
            <a target="_blank" href="https://vikingo.org/satellite/about"> Docs </a>
            |
            <a href="javascript:;" @click="showAbout = true"> About </a>
          </div>
        </div>

        <div class="col footer-block">
          <div class="copyright">
            <span>© 2020-2025 Vikingo Satellite</span>
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
import { getAbout } from "@/api.js";
export default {
  async mounted() {
    getAbout()
      .then((data) => {
        this.version = data.Version;
        this.build = data.Build;
      })
      .catch((e) => {
        console.log(e);
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
