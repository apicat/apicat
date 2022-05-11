<template>
    <el-card shadow="never" :body-style="{ padding: 0 }" v-loading="isLoading">
        <template #header>
            <span v-html="'&nbsp;'"></span>
            <div class="absolute right-2" style="top: 7px">
                <el-button type="primary" @click="onAddBtnClick">添加参数</el-button>
            </div>
        </template>

        <el-table :data="list" empty-text="暂无数据">
            <el-table-column prop="name" label="参数名称" show-overflow-tooltip />
            <el-table-column prop="type_name" label="参数类型" width="100" />
            <el-table-column prop="is_must" label="必传" width="70">
                <template #default="scope">
                    <span>{{ scope.is_must ? '是' : '否' }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="default_value" label="默认值" width="100" show-overflow-tooltip />
            <el-table-column prop="description" label="参数说明" show-overflow-tooltip />
            <el-table-column label="操作" width="120">
                <template #default="{ row }">
                    <span class="cursor-pointer mr-3 text-blue-600" @click="onModifyBtnClick(row)">编辑</span>
                    <span class="cursor-pointer mr-3 text-red-400" @click="onRemoveBtnClick(row)">删除</span>
                </template>
            </el-table-column>
        </el-table>

        <AddParamsModal ref="modifyModal" @on-ok="loadProjectCommonParamList" />
    </el-card>
</template>

<script setup lang="tsx">
    import AddParamsModal from './components/AddParamsModal.vue'
    import { deleteApiParam, getProjectCommonParamList } from '@/api/params'
    import { useProjectStore } from '@/stores/project'
    import { ref, watch } from 'vue'
    import { storeToRefs } from 'pinia'
    import { AsyncMsgBox } from '@/components/AsyncMessageBox'
    import { ElMessage as $Message } from 'element-plus'

    const modifyModal = ref()
    const list = ref([])
    const isLoading = ref(false)
    const projectStore = useProjectStore()
    const { projectInfo: project } = storeToRefs(projectStore)

    const onAddBtnClick = () => {
        modifyModal.value?.show()
    }

    const onModifyBtnClick = (param: any) => {
        modifyModal.value?.show({ ...param })
    }

    const onRemoveBtnClick = (param: any) => {
        AsyncMsgBox({
            title: '删除提示',
            content: <div class="truncate">确定删除参数「{param.name}」吗？</div>,
            onOk: () => {
                return deleteApiParam(projectStore.projectInfo.id, param.id).then((res: any) => {
                    $Message.success(res.msg || '删除成功！')
                    loadProjectCommonParamList()
                })
            },
        })
    }

    const loadProjectCommonParamList = async () => {
        isLoading.value = true
        try {
            const { data } = await getProjectCommonParamList(projectStore.projectInfo.id)
            let map = data.map || {}

            const arr: any = []

            Object.keys(map).forEach((key) => {
                map[key] && arr.push(map[key])
            })

            list.value = arr
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
                await loadProjectCommonParamList()
            }
        },
        { immediate: true }
    )
</script>
