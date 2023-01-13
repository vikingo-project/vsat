<template>
  <div>
    <pre :key="idx" v-for="(record, idx) in parsedData"
      >{{ record.key }}: {{ record.value }}</pre
    >
  </div>
</template>
<script>
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
        console.error(e);
      }
      return data;
    },
  },
};
</script>
