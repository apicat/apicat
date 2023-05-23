<template>
  <tr v-for="item in data" :key="item.id" :class="{ disabled: !readonly && !item.isUse }">
    <td v-if="!readonly"></td>
    <td class="px-12px" :class="{ copy_text: readonly }">{{ item.name }}</td>
    <td class="px-12px" :class="{ 'text-center': readonly }">{{ item.schema.type }}</td>
    <td class="text-center">{{ item.required }}</td>
    <td class="px-12px">{{ item.schema.example }}</td>
    <td class="px-12px">{{ item.schema.description }}</td>
    <td v-if="!readonly" class="text-center">
      <el-switch :model-value="item.isUse" size="small" @change="(v:any) => onSwitch && onSwitch(item.id, v)" />
    </td>
  </tr>
</template>
<script setup lang="ts">
defineProps<{ data: Array<any>; onSwitch?: any; readonly?: boolean }>()
</script>
<style scoped lang="scss">
.disabled td {
  opacity: 0.5;
  &:last-child {
    opacity: 1;
  }
}
</style>
