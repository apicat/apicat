<template>
    <AcTable :loading="isLoading" :table-data="members" :page-total="total" v-model:current-page="currentPage" :columns="columns" class="pb-3" />
</template>

<script setup lang="tsx">
    import { getProjectMembers } from '@/api/project'
    import { useProjectStore } from '@/stores/project'
    import { storeToRefs } from 'pinia'
    import { useTable } from '@/hooks/useTable'
    import { watch } from 'vue'

    const emit = defineEmits(['on-success'])

    const columns = [
        { title: '姓名', key: 'name' },
        { title: '邮箱', key: 'email' },
        {
            title: '角色',
            render: (row: any) => <span class="text dis_hover">{row.authority_name}</span>,
        },
    ]

    const searchParam = { project_id: '' }

    const projectStore = useProjectStore()

    const { projectInfo: project } = storeToRefs(projectStore)

    const {
        isLoading,
        total,
        currentPage,
        data: members,
        getTableData,
    } = useTable(getProjectMembers, { searchParam, totalKey: 'total_members', dataKey: 'project_members', isLoaded: false })

    watch(
        () => project.value,
        async () => {
            if (project.value && project.value.id) {
                searchParam.project_id = project.value.id
                await getTableData()
                emit('on-success', total.value || 0)
            }
        },
        { immediate: true }
    )
</script>
