<script setup lang="ts">
import { useLocale } from '@/hooks/useLocale'
const { elementPlusLocale } = useLocale()
</script>
<template>
  <el-config-provider :locale="elementPlusLocale">
    <RouterView />
  </el-config-provider>
</template>
