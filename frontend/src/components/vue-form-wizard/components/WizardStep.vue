<template>
  <li :class="{ active: tab.active }">
    <transition :name="transition" mode="out-in">
      <span>{{ index + 1 }}</span>
    </transition>

    <slot name="title">
      {{ tab.title }}
    </slot>
  </li>
</template>
<script>
export default {
  name: "wizard-step",
  props: {
    tab: {
      type: Object,
      default: () => {},
    },
    transition: {
      type: String,
      default: "",
    },
    index: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    iconActiveStyle() {
      return {
        backgroundColor: this.tab.color,
      };
    },
    stepCheckedStyle() {
      return {
        borderColor: this.tab.color,
      };
    },
    errorStyle() {
      return {
        borderColor: this.tab.errorColor,
        backgroundColor: this.tab.errorColor,
      };
    },
    stepTitleStyle() {
      let isError = this.tab.validationError;
      return {
        color: isError ? this.tab.errorColor : this.tab.color,
      };
    },
    isStepSquare() {
      return this.tab.shape === "square";
    },
    isTabShape() {
      return this.tab.shape === "tab";
    },
  },
};
</script>
<style>
</style>
