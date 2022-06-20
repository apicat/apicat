<template>
    <div class="ac-header-operate">
        <div class="ac-header-operate__main">
            <p v-if="!isEdit" :class="titleClass">13123</p>
            <p v-else><i class="iconfont" :class="isSaving ? 'iconcloud-upload' : 'iconcloud'"></i> {{ isSaving ? '保存中...' : '已保存在云端' }}</p>
        </div>

        <div class="ac-header-operate__btns">
            <el-button type="primary" @click="onSaveOrEditBtnClick"> {{ isEdit ? '预览' : '编辑' }}</el-button>
            <i class="iconfont iconshare2" @click="onShareBtnClick"></i>
            <i class="iconfont iconIconPopoverUpload" @click="onExportBtnClick"></i>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, computed, inject } from 'vue'
    import emitter, * as EVENT from '@/common/emitter'
    import { useRouter } from 'vue-router'
    import { DOCUMENT_EDIT_NAME } from '@/router/constant'
    import { API_SINGLE_EXPORT_ACTION_MAPPING } from '@/api/exportFile'

    const documentShareModal: any = inject('documentShareModal')
    const projectExportModal: any = inject('projectExportModal')

    const { currentRoute, push } = useRouter()

    const isSaving = ref(false)

    const isShowTitle = ref(false)
    const titleClass = computed(() => {
        return [
            'ac-header-operate__title animate__animated1 animate__faster1',
            {
                hidden: !isShowTitle.value,
                // animate__slideInDown: isShowTitle.value,
                // animate__slideOutUp: !isShowTitle.value,
            },
        ]
    })

    const isEdit = computed(() => currentRoute.value.name === DOCUMENT_EDIT_NAME)

    const getCommonParams = () => {
        const { params } = currentRoute.value
        const project_id = params.project_id
        const node_id = parseInt(params.node_id as string, 10)
        return {
            project_id,
            node_id,
        }
    }
    const onSaveOrEditBtnClick = () => {
        const { project_id, node_id } = getCommonParams()

        if (isEdit.value) {
            push({ name: 'document.api.detail', params: { project_id, node_id } })
            return
        }
        push({ name: 'document.api.edit', params: { project_id, node_id } })
    }

    const onShareBtnClick = () => {
        const { node_id } = getCommonParams()
        documentShareModal.value?.show({
            docId: node_id,
            nodeId: node_id,
        })
    }

    const onExportBtnClick = () => {
        const { project_id, node_id } = getCommonParams()
        projectExportModal.value.title = '导出文档'
        projectExportModal.value.show({ project_id, doc_id: node_id }, API_SINGLE_EXPORT_ACTION_MAPPING)
    }

    const saveDocumentDone = () => {
        isSaving.value = false
    }

    emitter.on(EVENT.IS_SHOW_DOCUMENT_TITLE, (isShow: any) => {
        isShowTitle.value = isShow
    })
    emitter.on(EVENT.DOCUMENT_SAVE_ING, () => {
        isSaving.value = true
    })
    emitter.on(EVENT.DOCUMENT_SAVE_DONE, saveDocumentDone)
    emitter.on(EVENT.DOCUMENT_SAVE_ERROR, saveDocumentDone)
</script>
