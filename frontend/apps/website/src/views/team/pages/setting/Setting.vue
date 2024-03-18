<script setup lang="ts">
import { useTeamStore } from '@/store/team'
import Page30pLayout from '@/layouts/Page30pLayout.vue'

const teamStore = useTeamStore()
const compo = computed(() =>
  defineAsyncComponent(
    teamStore.isOwner ? () => import('./Owner.vue') : () => import('./Normal.vue'),
  ),
)
</script>

<template>
  <Page30pLayout>
    <div class="w-full bg-white">
      <h1>{{ $t('app.team.setting.title') }}</h1>
      <component :is="compo" />
    </div>
  </Page30pLayout>
</template>

<style scoped>
:deep(.el-form-item) {
  margin-bottom: 0;
}
</style>
