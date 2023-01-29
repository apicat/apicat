<template>
    <main :class="layoutClass">
        <div class="ac-project-info text-[16px]">
            <router-link :to="backRoute" class="flex items-center">
                <el-icon class="mr-[10px] w-[32px] h-[32px] border border-gray-400 bg-white rounded-[4px] hover:bg-gray-6"><ArrowLeftBold /></el-icon>
            </router-link>
            历史记录
        </div>

        <DocumentHistoryOperateHeader :title="title" />

        <div class="ac-doc-layout__left">
            <DocumentHistoryRecordTree />
        </div>
        <div class="ac-doc-layout__right scroll-content">
            <router-view />
        </div>
    </main>

    <div v-html="zoomTemplate" />
</template>
<script setup lang="ts">
    import { ArrowLeftBold } from '@element-plus/icons-vue'
    import DocumentHistoryOperateHeader from '../views/document/components/DocumentHistoryOperateHeader.vue'
    import DocumentHistoryRecordTree from '../views/document/components/DocumentHistoryRecordTree.vue'
    import { ref, provide, computed } from 'vue'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'
    import { useIterateStore } from '@/stores/iterate'
    import { DOCUMENT_DETAIL_NAME, ITERATE_DOCUMENT_DETAIL_NAME } from '@/router/constant'
    import { useRoute, useRouter } from 'vue-router'

    const { params, query } = useRoute()
    const { push } = useRouter()
    const projectStore = useProjectStore()
    const { isGuest, isReader, projectInfo } = storeToRefs(projectStore)
    const iterateStore = useIterateStore()
    const { isIterateRoute, iterateInfo } = storeToRefs(iterateStore)

    const title = ref('')
    const layoutClass = computed(() => ['ac-doc-layout', { readonly: isGuest.value || isReader.value }])

    const backRoute: any =
        query.from && query.from.length === 32
            ? { name: ITERATE_DOCUMENT_DETAIL_NAME, params: { iterate_id: query.from, node_id: params.doc_id } }
            : { name: DOCUMENT_DETAIL_NAME, params: { project_id: params.project_id, node_id: params.doc_id } }

    const zoomTemplate = `<template id="template-zoom-image">
                            <div class="zoom-image-wrapper">
                                <div class="zoom-image-container" data-zoom-container></div>
                            </div>
                          </template>`

    provide('setDocumentTitle', (t: string) => {
        document.title = projectInfo.value.name + (isIterateRoute.value ? `(${iterateInfo.value.title})` : '') + '-' + t
        title.value = t
    })

    provide('goBack', () => push(backRoute))
</script>
