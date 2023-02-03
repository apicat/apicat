<template>
    <div class="ac-header-operate">
        <div class="ac-header-operate__main">
            <p v-if="!isEdit" :class="titleClass">{{ title }}</p>
            <p v-else><i class="iconfont" :class="isSaving ? 'iconcloud-upload' : 'iconcloud'"></i> {{ isSaving ? '保存中...' : '已保存在云端' }}</p>
        </div>

        <div class="ac-header-operate__btns" v-if="!isGuest">
            <el-button type="primary" @click="onSaveOrEditBtnClick" :loading="isLoading"> {{ isEdit ? '预览' : '编辑' }}</el-button>
            <el-tooltip effect="dark" content="分享该文档" placement="bottom">
                <i class="iconfont iconshare2" @click="onShareBtnClick"></i>
            </el-tooltip>
            <el-tooltip effect="dark" content="导出该文档" placement="bottom">
                <i class="iconfont iconIconPopoverUpload" @click="onExportBtnClick"></i>
            </el-tooltip>
            <el-tooltip effect="dark" content="历史记录" placement="bottom">
                <i class="iconfont iconhistory" @click="goDocumentHistoryRecord"></i>
            </el-tooltip>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, computed, inject } from 'vue'
    import emitter, * as EVENT from '@/common/emitter'
    import { useRoute, useRouter } from 'vue-router'
    import {
        DOCUMENT_EDIT_NAME,
        DOCUMENT_DETAIL_NAME,
        ITERATE_DOCUMENT_EDIT_NAME,
        ITERATE_DOCUMENT_DETAIL_NAME,
        DOCUMENT_HISTORY_DETAIL_NAME,
    } from '@/router/constant'
    import { API_SINGLE_EXPORT_ACTION_MAPPING } from '@/api/exportFile'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'
    import { useIterateStore } from '@/stores/iterate'

    defineProps({
        title: {
            type: String,
            default: '',
        },
    })

    const documentShareModal: any = inject('documentShareModal')
    const projectExportModal: any = inject('projectExportModal')

    const { params } = useRoute()

    const { currentRoute, push } = useRouter()
    const projectStore = useProjectStore()
    const { isIterateRoute } = useIterateStore()
    const { isGuest, projectInfo } = storeToRefs(projectStore)

    const isSaving = ref(false)
    const isLoading = ref(false)

    const isShowTitle = ref(false)

    const titleClass = computed(() => {
        return [
            'ac-header-operate__title',
            {
                hidden: !isShowTitle.value,
                animate__slideInDown: isShowTitle.value,
                animate__slideOutUp: !isShowTitle.value,
            },
        ]
    })

    const isEdit = computed(() => currentRoute.value.name === DOCUMENT_EDIT_NAME || currentRoute.value.name === ITERATE_DOCUMENT_EDIT_NAME)

    const getCommonParams = () => {
        const { params } = currentRoute.value
        const node_id = parseInt(params.node_id as string, 10)
        return {
            iterate_id: params.iterate_id,
            project_id: params.project_id,
            node_id,
        }
    }
    const onSaveOrEditBtnClick = () => {
        const { node_id } = getCommonParams()

        // 编辑 -> 预览 点击
        if (isEdit.value) {
            isLoading.value = true
            emitter.emit(EVENT.DOCUMENT_SAVE_BTN_ING)
            return
        }

        // 预览 -> 编辑 点击
        push({ name: isIterateRoute ? ITERATE_DOCUMENT_EDIT_NAME : DOCUMENT_EDIT_NAME, params: { ...params, node_id } })
    }

    const onShareBtnClick = () => {
        const { node_id } = getCommonParams()
        documentShareModal.value?.show({
            docId: node_id,
            nodeId: node_id,
        })
    }

    const onExportBtnClick = () => {
        const { node_id } = getCommonParams()
        projectExportModal.value.title = '导出文档'
        projectExportModal.value.show({ project_id: projectInfo.value.id, doc_id: node_id }, API_SINGLE_EXPORT_ACTION_MAPPING)
    }

    const goDocumentHistoryRecord = () => {
        const { node_id, iterate_id } = getCommonParams()
        const routeParams: any = {
            name: DOCUMENT_HISTORY_DETAIL_NAME,
            params: { project_id: projectInfo.value.id, doc_id: node_id },
        }
        // 迭代路由
        if (isIterateRoute) {
            routeParams.query = { from: iterate_id }
        }
        push(routeParams)
    }

    const saveDocumentDone = () => {
        isSaving.value = false
    }

    const onSaveDocumentDone = () => {
        isLoading.value = false
        const { node_id } = getCommonParams()
        push({ name: isIterateRoute ? ITERATE_DOCUMENT_DETAIL_NAME : DOCUMENT_DETAIL_NAME, params: { ...params, node_id } })
    }

    emitter.on(EVENT.IS_SHOW_DOCUMENT_TITLE, (isShow: any) => {
        isShowTitle.value = isShow
    })
    emitter.on(EVENT.DOCUMENT_SAVE_ING, () => {
        isSaving.value = true
    })

    emitter.on(EVENT.DOCUMENT_SAVE_DONE, saveDocumentDone)
    emitter.on(EVENT.DOCUMENT_SAVE_ERROR, saveDocumentDone)

    emitter.on(EVENT.DOCUMENT_SAVE_BTN_DONE, onSaveDocumentDone)
    emitter.on(EVENT.DOCUMENT_SAVE_BTN_ERROR, () => {
        isLoading.value = false
    })
</script>
