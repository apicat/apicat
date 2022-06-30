<template>
    <el-card shadow="never" :body-style="{ padding: 0 }" v-loading="isLoading">
        <template #header>
            <span>{{ title }}</span>
            <div class="absolute right-2" style="top: 7px" v-if="isAdmin">
                <el-button type="primary" @click="onModifyMemberInfoTextClick()">添加成员</el-button>
            </div>
        </template>

        <el-table :data="members" class="pb-3" header-row-class-name="table__gray_header" empty-text="暂无成员">
            <el-table-column prop="name" label="姓名" show-overflow-tooltip />
            <el-table-column prop="email" label="邮箱" show-overflow-tooltip />
            <el-table-column prop="authority_name" label="角色" />
            <el-table-column label="操作" v-if="!isNormal">
                <template #default="{ row }">
                    <template v-if="!row.isSelf && isAdmin">
                        <span class="text-blue-600 hover:text-blue-500 cursor-pointer mr-3" @click="onModifyPasswordTextClick(row)">修改密码</span>
                        <span class="text-blue-600 hover:text-blue-500 cursor-pointer mr-3" @click="onModifyMemberInfoTextClick(row)">编辑成员</span>
                        <span class="text-red-500 hover:text-red-400 cursor-pointer" @click="onRemoveMemberTextClick(row)">移除成员</span>
                    </template>

                    <template v-if="!row.isSelf && isManager && userInfo.authority < row.authority">
                        <span class="text-blue-600 hover:text-blue-500 cursor-pointer mr-3" @click="onModifyPasswordTextClick(row)">修改密码</span>
                        <span class="text-blue-600 hover:text-blue-500 cursor-pointer mr-3" @click="onModifyMemberInfoTextClick(row)">编辑成员</span>
                        <span class="text-red-500 hover:text-red-400 cursor-pointer" @click="onRemoveMemberTextClick(row)">移除成员</span>
                    </template>
                </template>
            </el-table-column>
        </el-table>

        <el-pagination class="justify-end pb-3" :page-size="PAGE_SIZE" layout="prev, pager, next" v-model:current-page="page" :page-count="total" />
    </el-card>

    <AddMember ref="addMemberRef" @on-ok="loadMembers" />
    <ModifyMemberPassword ref="modifyMemberPwdRef" />
</template>
<script setup lang="ts">
    import { ref, computed, onMounted, shallowRef } from 'vue'
    import { storeToRefs } from 'pinia'
    import NProgress from 'nprogress'
    import { ElMessage as $Message } from 'element-plus'
    import { AsyncMsgBox } from '@/components/AsyncMessageBox'
    import AddMember from './components/AddMember.vue'
    import ModifyMemberPassword from './components/ModifyMemberPassword.vue'

    import { getMembers, removeMember } from '@/api/team'
    import { PAGE_SIZE } from '@/common/constant'
    import { useApi } from '@/hooks/useApi'
    import { usePage } from '@/hooks/usePage'
    import { useUserStore } from '@/stores/user'

    NProgress.configure({ showSpinner: false })

    const { page } = usePage(() => loadMembers())
    const [isLoading, api] = useApi(getMembers, { isShowMessage: false })
    const userStore = useUserStore()

    const { isAdmin, isManager, isNormal, userInfo } = storeToRefs(userStore)

    const addMemberRef = ref()
    const modifyMemberPwdRef = ref()
    const members = shallowRef([])
    const totalMembers = ref(0)
    const total = ref(0)

    const title = computed(() => '成员列表' + (totalMembers.value ? `(${totalMembers.value})` : ''))

    const onRemoveMemberTextClick = (member: any) => {
        AsyncMsgBox({
            title: '删除提示',
            content: '确定移除该成员吗？',
            onOk: () => {
                return removeMember(member.user_id).then((res: any) => {
                    $Message.success(res.msg || '移除成功！')
                    loadMembers()
                })
            },
        })
    }

    const onModifyMemberInfoTextClick = (member?: any) => {
        addMemberRef.value?.show(member)
    }

    const onModifyPasswordTextClick = (user: any) => {
        modifyMemberPwdRef.value?.show(user.user_id)
    }

    const loadMembers = async () => {
        const { data } = await api({ page: page.value })
        const { members: _members, total_members, total_page } = data || {}
        members.value = (_members || []).map((member: any) => {
            member.to = { name: 'member.detail', params: { uid: member.user_id } }
            member.isSelf = member.user_id === userInfo.value.user_id
            member.name = member.isSelf ? `${member.name}(我)` : member.name
            return member
        })
        totalMembers.value = total_members || 0
        total.value = total_page || 0
    }

    onMounted(async () => {
        await loadMembers()
    })
</script>
