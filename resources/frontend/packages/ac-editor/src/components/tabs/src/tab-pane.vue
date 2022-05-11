<template>
  <div v-show="show">
    <slot></slot>
  </div>
</template>
<script>
export default {
  name: "TabPane",
  inject: ["tabsInstance"],
  props: {
    name: {
      type: String,
    },
    label: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      show: true,
    };
  },
  mounted() {
    this.tabsInstance.initTabs();
  },
  watch: {
    label() {
      this.tabsInstance.initTabs();
    },
    name() {
      this.tabsInstance.initTabs();
    },
  },
};
</script>
