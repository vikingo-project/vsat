<template>
  <div>
    <pre v-for="record in parsedData">{{ record.key }}: {{ record.value }}</pre>
  </div>
</template>
<script>
import { bus } from "@/bus.js";
export default {
  data() {
    return {};
  },
  mounted() {},
  props: ["value"],
  methods: {
    viewFullText() {
      if (this.value.length > 100) this.showModal = true;
    },
  },
  computed: {
    parsedData() {
      let data = [];
      try {
        let decoded = JSON.parse(this.value);
        for (let k in decoded) {
          data.push({ key: k, value: decoded[k] });
        }
      } catch (e) {
        console.log(e);
      }
      return data;
    },
  },
};
</script>
