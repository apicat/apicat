<template>
    <el-card shadow="never" :body-style="{ padding: 0 }">
        <template #header>
            <div class="flex items-center justify-between">
                <span>迭代列表</span>
                <div class="absolute flex right-2">
                    <el-select
                        v-show="isActiveIterateList"
                        v-model="value"
                        filterable
                        placeholder="所有项目"
                        clearable
                        :loading="isLoadingForProject"
                        loading-text="加载中..."
                        class="w-52"
                        popper-class="w-52"
                    >
                        <el-option v-for="item in options" :key="item.id" :label="item.name" :value="item.id" />
                        <template #empty>
                            <p class="p-2 text-center">暂无项目</p>
                        </template>
                    </el-select>

                    <el-button class="ml-2" type="primary" @click="onAddIterateBtnClick()">新建迭代</el-button>
                </div>
            </div>
        </template>

        <AcTable :loading="isLoading" :table-data="iterations" :page-total="total" v-model:current-page="currentPage" :columns="columns" class="pb-3" />
    </el-card>
    <AddIterateModal ref="addIterateModal" @ok="refreshIterateList" />
    <PlanIterateModal ref="planIterateModal" @ok="getTableData" />
</template>

<script setup lang="tsx">
    import AddIterateModal from './components/AddIterateModal.vue'
    import PlanIterateModal from './components/PlanIterateModal.vue'

    import useTable from '@/hooks/useTable'
    import { deleteIteration, getIterations, toIterateDocumentPath } from '@/api/iterate'
    import { getAllProjectList } from '@/api/project'
    import { ref, toRaw, watch, onMounted } from 'vue'
    import useApi from '@/hooks/useApi'

    import { useIterateStore } from '@/stores/iterate'
    import { storeToRefs } from 'pinia'
    import { AsyncMsgBox } from '@/components/AsyncMessageBox'
    import { ElMessage as $Message } from 'element-plus'
    import { PROJECT_ROLES_KEYS } from '@/common/constant'

    const searchParam = { project_id: '' }
    const iterateStore = useIterateStore()
    const { isActiveIterateList, activeTab } = storeToRefs(iterateStore)

    const transformData = (item: any) => {
        item.to = toIterateDocumentPath(item.id_public)
        return item
    }
    const {
        isLoading,
        currentPage,
        data: iterations,
        total,
        getTableData,
    } = useTable(getIterations, { searchParam, totalKey: 'total_iterations', dataKey: 'iterations', isLoaded: false, transform: transformData })

    const [isLoadingForProject, executeGetProjectList] = useApi(getAllProjectList)

    const columns = ref([
        {
            title: '迭代名称',
            ellipsis: true,
            tooltip: true,
            key: 'title',
            render(row: any) {
                return (
                    <a href={row.to} class="text-blue-600 cursor-pointer">
                        {row.title}
                    </a>
                )
            },
        },
        {
            title: '所属项目',
            ellipsis: true,
            tooltip: true,
            width: 200,
            key: 'project_title',
            render(row: any) {
                return (
                    <div class="flex items-center collect-item">
                        <div class="overflow-hidden">
                            <p class="truncate ">{row.project_title}</p>
                        </div>

                        {isActiveIterateList.value && (
                            <i
                                class={['cursor-pointer collect-icon iconfont h-[21px]', { 'iconstar-fill': row.star, iconstar: !row.star }]}
                                onClick={() => onCollectIconClick(row)}
                            ></i>
                        )}
                    </div>
                )
            },
        },
        {
            title: 'API数量',
            ellipsis: true,
            tooltip: true,
            key: 'api_num',
        },
        {
            title: '创建时间',
            key: 'created_at',
        },
        {
            title: '操作',
            width: 150,
            ellipsis: false,
            tooltip: false,
            render(row: any) {
                if (row.authority === PROJECT_ROLES_KEYS.READER) {
                    return []
                }
                return (
                    <div>
                        <span class="mr-3 cursor-pointer" onClick={() => onEditIterateBtnClick(row)}>
                            编辑
                        </span>
                        <span class="mr-3 cursor-pointer" onClick={() => onIterateBtnClick(row)}>
                            规划
                        </span>
                        <span class="mr-3 text-red-400 cursor-pointer" onClick={() => onRemoveIterateBtnClick(row)}>
                            删除
                        </span>
                    </div>
                )
            },
        },
    ])

    const value = ref('')
    const options: any = ref([])
    const addIterateModal = ref()
    const planIterateModal = ref()

    const onAddIterateBtnClick = (iterate?: any) => {
        addIterateModal.value?.show(iterate)
    }

    const onCollectIconClick = async (row: any) => {
        const api = row.star ? iterateStore.removeProjectFromCollect : iterateStore.addProjectToCollect
        row.star = !row.star
        try {
            await api(row)
            markIterateIsStar(row.project_id, row.star)
        } catch (e) {
            row.star = !row.star
            markIterateIsStar(row.project_id, row.star)
        }
    }

    const markIterateIsStar = (project_id: any, isStar: boolean) => {
        iterations.value.forEach((item: any) => {
            if (item.project_id === project_id) {
                item.star = isStar
            }
        })
    }

    const onEditIterateBtnClick = (row: any) => onAddIterateBtnClick(row)

    const onIterateBtnClick = (row: any) => {
        planIterateModal.value?.show(toRaw(row))
    }

    const onRemoveIterateBtnClick = (row: any) => {
        AsyncMsgBox({
            title: '删除提示',
            content: (
                <div class="truncate" title={row.title}>
                    确定删除迭代「{row.title}」吗？
                </div>
            ),
            onOk: () => {
                return deleteIteration({ iteration_id: row.id }).then((res: any) => {
                    $Message.success(res.msg || '删除成功！')
                    getTableData()
                })
            },
        })
    }

    const refreshIterateList = async () => {
        // currentPage.value = 1
        await getTableData()
    }

    // 项目筛选
    watch(
        () => value.value,
        async () => {
            // 仅在迭代列表中进行筛选
            if(activeTab.value){
                return
            }
            searchParam.project_id = value.value
            await refreshIterateList()
        }
    )

    // 收藏项目 切换
    watch(activeTab, async () => {
        value.value = ''
        currentPage.value = 1
        searchParam.project_id = activeTab.value
        await refreshIterateList()
    })

    const init = async () => {
        const [{ data }] = await Promise.all([executeGetProjectList()])
        options.value = data || []
    }

    onMounted(init)

</script>

<style lang="scss">
    .collect-item {
        .collect-icon {
            display: none;
        }

        &:hover .collect-icon {
            display: block;
        }
    }
</style>
