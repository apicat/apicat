<template>
  <tr v-for="item in data" :key="item.id" :class="{ disabled: !item.isUse }">
    <td></td>
    <td class="px-12px">{{ item.name }}</td>
    <td class="px-12px">{{ item.schema.type }}</td>
    <td class="text-center">{{ item.required }}</td>
    <td class="px-12px">{{ item.schema.example }}</td>
    <td class="px-12px">{{ item.schema.description }}</td>
    <td class="text-center">
      <el-switch :model-value="item.isUse" size="small" @change="(v) => onSwitch(item.id, v)" />
    </td>
  </tr>
</template>
<script setup lang="ts">
defineProps<{ data: Array<any>; onSwitch: any }>()
</script>
<style scoped lang="scss">
.disabled td {
  opacity: 0.5;
  &:last-child {
    opacity: 1;
  }
}
</style>
