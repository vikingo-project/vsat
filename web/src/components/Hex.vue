<template>
  <div class="col-auto">

    <!-- <div class="mt-3 mb-2">00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f</div> -->
    <div class="mt-3 ow d-flex">
      <div class="col">
        <el-input type="textarea" class="" spellcheck="false" :autosize="{ minRows: 2}" v-model="proxy"
          @keydown="checkEvent">
        </el-input>
      </div>
      <div class="col p-3">{{ raw }}</div>
    </div>

  </div>
</template>
<script>
export default {
  props: ["readOnly", "content"],
  data() {
    return {
      value: "",
    };
  },
  mounted() {
    this.value = this.filter(this.content);
  },

  methods: {
    filter(s) {
      return s
        .replace(/[^0-9A-F]/gi, "")
        .replace(/(..)/g, "$1 ")
        .replace(/ $/, "")
        .toUpperCase();
    },
    checkEvent(e) {
      if (this.readOnly) {
        e.preventDefault();
        return;
      }

      var c = this.filter(e.key);
      if (c.length == 0) {
        e.preventDefault();
        return;
      }
    },
  },
  computed: {
    proxy: {
      get() {
        return this.value;
      },
      set(v) {
        this.value = this.filter(v);
      },
    },

    raw() {
      var h = "";
      for (var i = 0; i < this.value.length; i += 3) {
        var c = parseInt(this.value.substr(i, 2), 16);
        h = 63 < c && 127 > c ? h + String.fromCharCode(c) : h + ".";
      }
      return h.replace(/(.{16})/g, "$1 ");
    },
    areaHight() {
      return 1.5 + this.value.length / 47 + "em";
    },
  },
};
</script>
<style scoped>
* {
  margin: 0;
  padding: 0;
  vertical-align: top;
  font: 1em/1em courier;
}
</style>