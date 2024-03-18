<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useLocaleProvider as acComponentLocaleProvider } from '@apicat/components'
import { useLocaleProvider as acEditorLocaleProvider } from '@apicat/editor'
import { useAppStore } from './store/app'
import { useLocaleStore } from './store/locale'
import Loading from '@/components/Loading.vue'

const appStore = useAppStore()
const { isShowGlobalLoading } = storeToRefs(appStore)
const { elementPlusLocaleMessage, acCompLocaleMessage, acEditorLocaleMessage, locale } = storeToRefs(useLocaleStore())

onBeforeMount(async () => {
  await appStore.initAppConfig()
})

acComponentLocaleProvider(acCompLocaleMessage, locale)
acEditorLocaleProvider(acEditorLocaleMessage, locale)
</script>

<template>
  <el-config-provider :locale="elementPlusLocaleMessage">
    <Loading v-if="isShowGlobalLoading" />
    <router-view />
  </el-config-provider>
</template>
