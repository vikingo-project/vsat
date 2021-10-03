<template>
  <div v-if="full">{{ fullTime }}</div>
  <div v-else :title="fullTime">{{ timeAgo }}</div>
</template>
<script>
import moment from "moment";
import * as utils from "@/utils";

export default {
  props: ["date", "full"],
  computed: {
    timeAgo() {
      return utils.timeAgo(this.date);
    },
    fullTime() {
      return moment(this.date, "YYYY-MM-DD HH:mm:ss ZZ z").format(
        "YYYY-MM-DD HH:mm:ss"
      );
    },
  },
};
</script>