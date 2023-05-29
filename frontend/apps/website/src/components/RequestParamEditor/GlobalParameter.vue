<template>
  <div :class="[nsRow.b(), { [nsRow.is('disabled')]: !readonly && !item.isUse, [nsRow.m('global')]: !readonly }]" v-for="item in data" :key="item.id">
    <div :class="nsRow.e('content')">
      <div :class="[nsRow.e('item'), nsRow.e('name')]">
        <span :class="{ copy_text: readonly, 'px-12px': !readonly }">{{ item.name }}</span>
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('type')]" :style="{ justifyContent: !readonly ? 'flex-start' : 'center' }">
        <span class="px-8px">{{ item.schema.type }}</span>
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('required')]">
        {{ item.required }}
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('example')]">
        {{ item.schema.example }}
      </div>
      <div :class="[nsRow.e('item'), nsRow.e('description')]">
        {{ item.schema.description }}
      </div>
      <div v-if="allowMock" :class="[nsRow.e('item'), nsRow.e('mock')]">{{ $t('editor.table.paramMock') }}</div>
      <div :class="[nsRow.e('item'), nsRow.e('operation')]" v-if="!readonly">
        <el-switch :model-value="item.isUse" size="small" @change="(v:any) => onSwitch && onSwitch(item.id, v)" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'

defineProps<{ data: Array<any>; onSwitch?: any; readonly?: boolean; allowMock?: boolean }>()

const nsRow = useNamespace('simple-row')
</script>
<style scoped lang="scss">
.disabled td {
  opacity: 0.5;
  &:last-child {
    opacity: 1;
  }
}
</style>
