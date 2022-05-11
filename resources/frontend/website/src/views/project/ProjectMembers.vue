<template>
    <el-card shadow="never" :body-style="{ padding: 0 }">
        <template #header>
            <span>成员列表{{ totalMembers ? `(${totalMembers})` : '' }}</span>
            <div class="absolute right-2" style="top: 7px" v-if="project && project.authority === 0">
                <el-button @click="onAddMemberBtnClick" type="primary">添加成员</el-button>
            </div>
        </template>

        <ProjectMembersManage
            ref="projectMembersManage"
            v-if="project && project.authority === 0"
            @on-remove="onRemoveMemberSuccess"
            @on-success="onGetMemberListSuccess"
        />
        <ProjectMembersDeveloper v-if="project && project.authority !== 0" @on-success="onGetMemberListSuccess" />
    </el-card>

    <AddProjectMemberModal :members="withoutProjectMemberList" ref="addProjectMemberModal" @on-ok="onAddMemberSuccess" />
</template>

<script setup lang="tsx">
    import ProjectMembersDeveloper from './components/ProjectMembersDeveloper.vue'
    import ProjectMembersManage from './components/ProjectMembersManage.vue'
    import AddProjectMemberModal from './components/AddProjectMemberModal.vue'

    import { getWithoutProjectMemberList } from '@/api/project'
    import { useProjectStore } from '@/stores/project'
    import { storeToRefs } from 'pinia'
    import { watch, ref } from 'vue'

    const projectStore = useProjectStore()
    const { projectInfo: project } = storeToRefs(projectStore)
    const withoutProjectMemberList = ref([])
    const addProjectMemberModal = ref()
    const projectMembersManage = ref()
    const totalMembers = ref(0)

    const onAddMemberBtnClick = () => {
        addProjectMemberModal.value?.show(project.value)
    }

    const onAddMemberSuccess = async () => {
        projectMembersManage.value?.getTableData()
        await loadWithoutProjectMemberList()
    }

    const onRemoveMemberSuccess = async () => {
        await loadWithoutProjectMemberList()
    }

    const onGetMemberListSuccess = (total: number) => {
        totalMembers.value = total
    }

    const loadWithoutProjectMemberList = async () => {
        const { authority, id } = project.value || {}

        if (authority !== 0) {
            return
        }

        // 非 管理者不获取成员
        try {
            const { data } = await getWithoutProjectMemberList(id)
            withoutProjectMemberList.value = data || []
        } catch (e) {
            withoutProjectMemberList.value = []
        }
    }

    watch(
        () => project.value,
        async () => {
            if (project.value && project.value.id) {
                await loadWithoutProjectMemberList()
            }
        },
        { immediate: true }
    )
</script>
