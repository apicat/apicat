<template>
  <main :class="ns.b()">
    <div :class="historyInfo.b()">
      <div :class="historyInfo.e('img')">
        <a href="javascript:void(0)" @click="handleGoBack">
          <el-icon :class="historyInfo.e('back')"><ac-icon-ep-arrow-left-bold /></el-icon>
        </a>
      </div>
      <div :class="historyInfo.e('title')">历史记录</div>
    </div>

    <div :class="ns.e('left')">
      <div class="flex flex-col h-full overflow-y-scroll scroll-content">
        <slot name="left"></slot>
      </div>
    </div>
    <div :class="ns.e('right')" class="scroll-content">
      <router-view />
    </div>
  </main>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks/useNamespace'
import { useParams } from '@/hooks/useParams'
import { useDefinitionSchemaStore } from '@/store/definitionSchema'
import useDefinitionResponseStore from '@/store/definitionResponse'
import uesGlobalParametersStore from '@/store/globalParameters'
import useProjectStore from '@/store/project'
import { MAIN_PATH } from '@/router'

const ns = useNamespace('doc-layout')
const historyInfo = useNamespace('history-info')
const projectStore = useProjectStore()
const globalParametersStore = uesGlobalParametersStore()
const definitionSchemaStore = useDefinitionSchemaStore()
const definitionResponseStore = useDefinitionResponseStore()
const { project_id } = useParams()
const router = useRouter()
const backPath = router.options.history.state.back
const handleGoBack = () => {
  router.push(backPath ? (backPath as string) : MAIN_PATH)
}

provide('goBack', handleGoBack)

onMounted(async () => {
  await projectStore.getUrlServers(project_id as string)
  await globalParametersStore.getGlobalParameters(project_id as string)
  await definitionSchemaStore.getDefinitions(project_id as string)
  await definitionResponseStore.getDefinitions(project_id as string)
})
</script>

<style lang="scss">
@use '@/styles/mixins/mixins' as *;
@use '@/styles/variable' as *;

// 项目信息
@include b(history-info) {
  height: $doc-header-height;
  width: $doc-layout-left-width;
  padding: 0 $doc-layout-padding;
  @apply flex items-center fixed left-0 top-0 z-50 bg-gray-100;

  @include e(img) {
    @apply flex-none w-32px h-32px mr-10px cursor-pointer;
  }

  @include e(back) {
    @apply w-32px h-32px rounded-4px  text-12px border-1px border-gray border-solid bg-white hover:bg-gray-100;
  }

  @include e(title) {
    @apply truncate text-16px relative pr-20px;
  }
}
</style>
