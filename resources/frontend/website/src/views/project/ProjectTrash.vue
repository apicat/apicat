<template>
    <el-card shadow="never" :body-style="{ padding: 0 }" v-loading="isLoading">
        <template #header>
            <span>回收站</span>
        </template>

        <el-table :data="list" empty-text="暂无数据">
            <el-table-column prop="title" label="文档名称" show-overflow-tooltip />
            <el-table-column prop="deleted_at" label="删除时间" />
            <el-table-column prop="remaining" label="剩余" />
            <el-table-column label="操作">
                <template #default="{ row }">
                    <a class="cursor-pointer mr-3 text-blue-600" target="_blank" :href="row.previewUrl">预览</a>
                    <span class="cursor-pointer mr-3 text-blue-600" href="javascript:void(0)" @click="onRestoreBtnClick(row)">恢复</span>
                </template>
            </el-table-column>
        </el-table>
    </el-card>

    <RestoreDocumentModal ref="restoreDocumentModal" @on-ok="loadProjectTrashList" />
</template>

<script setup lang="ts">
    import { ref, watch, h } from 'vue'
    import { storeToRefs } from 'pinia'
    import { ElMessage as $Message } from 'element-plus'
    import NProgress from 'nprogress'
    import { getProjectTrashList, restoreApiDocument } from '@/api/project'
    import { HTTP_STATUS, timestampFormat } from '@ac/shared'
    import { useProjectStore } from '@/stores/project'
    import { toPreviewTrashDocumentPath } from '@/router/preview.router'

    import RestoreDocumentModal from './components/RestoreDocumentModal.vue'

    const projectStore = useProjectStore()
    const { projectInfo: project } = storeToRefs(projectStore)

    const restoreDocumentModal = ref()
    const list = ref([])
    const isLoading = ref(false)

    const onRestoreBtnClick = (doc: any) => {
        NProgress.start()

        const data = {
            project_id: projectStore.projectInfo.id,
            doc_id: doc.id,
        }

        restoreApiDocument(data)
            .then(({ status }) => {
                if (status === HTTP_STATUS.NO_PARENT_DIR) {
                    restoreDocumentModal.value && restoreDocumentModal.value.show(data)
                    return
                }

                if (status === HTTP_STATUS.OK) {
                    $Message({
                        type: 'success',
                        showClose: true,
                        message: () => {
                            return h('span', null, [
                                '文档恢复成功，',
                                h(
                                    'a',
                                    {
                                        class: 'text-blue-600',
                                        href: `/editor/${projectStore.projectInfo.id}/doc/${doc.id}`,
                                    },
                                    '查看详情'
                                ),
                            ])
                        },
                    })
                    loadProjectTrashList()
                }
            })
            .catch((e) => {
                //
            })
            .finally(() => {
                NProgress.done()
            })
    }

    const loadProjectTrashList = async () => {
        isLoading.value = true
        try {
            const { data } = await getProjectTrashList(projectStore.projectInfo.id)
            list.value = (data || []).map((item: any) => {
                item.remaining = timestampFormat(item.deleted_at) + '天'
                item.previewUrl = toPreviewTrashDocumentPath({ project_id: projectStore.projectInfo.id, doc_id: item.id })
                return item
            })
        } catch (e) {
            //
        } finally {
            isLoading.value = false
        }
    }

    watch(
        () => project.value,
        async () => {
            if (project.value && project.value.id) {
                await loadProjectTrashList()
            }
        },
        { immediate: true }
    )
</script>
