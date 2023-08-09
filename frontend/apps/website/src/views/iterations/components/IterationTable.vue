<template>
  <div>table</div>
</template>
<script lang="ts" setup>
import { Iteration } from '@/typings'

interface Props {
  data: Iteration[]
  page: number
  total: number
  pageSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
  page: 1,
  total: 0,
  pageSize: 15,
})
</script>
