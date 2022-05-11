<template>
    <AcTable :loading="isLoading" :table-data="members" :page-total="total" v-model:current-page="currentPage" :columns="columns" class="pb-3" />
    <el-popover ref="popoverRef" popper-class="ac-popper-menu" width="auto" v-model:visible="isShowPopover" :virtual-ref="rolePopperRef" virtual-triggering>
        <ul>
            <li v-for="role in roles" :key="role.value" class="ac-popper-menu__item" @click="onRoleItemClick(role)">{{ role.text }}</li>
        </ul>
    </el-popover>
</template>

<script lang="tsx">
    import { changeMemberAuthority, getProjectMembers, removeMember, transferProject } from '@/api/project'
    import { PROJECT_ROLE_LIST } from '@ac/shared'
    import { useProjectStore } from '@/stores/project'
    import { useUserStore } from '@/stores/user'
    import { storeToRefs } from 'pinia'
    import { useTable } from '@/hooks/useTable'
    import { watch, ref, defineComponent, Ref } from 'vue'
    import { ArrowDown } from '@element-plus/icons-vue'
    import { onClickOutside } from '@vueuse/core'
    import NProgress from 'nprogress'
    import { ElMessage as $Message } from 'element-plus'
    import { AsyncMsgBox } from '@/components/AsyncMessageBox'
    import { useRouter } from 'vue-router'

    NProgress.configure({ showSpinner: false })

    export default defineComponent({
        emits: ['on-remove', 'on-success'],

        setup(props, { emit }) {
            const searchParam = { project_id: '' }
            const projectStore = useProjectStore()
            const { projectInfo: project } = storeToRefs(projectStore)
            const userStore = useUserStore()
            const rolePopperRef = ref()
            const isShowPopover = ref(false)
            const popoverRef = ref()
            const roles = PROJECT_ROLE_LIST
            const $router = useRouter()

            let currentMember: any = null

            const {
                isLoading,
                currentPage,
                data: members,
                total,
                getTableData,
            } = useTable(getProjectMembers, { searchParam, totalKey: 'total_members', dataKey: 'project_members', isLoaded: false })

            const columns = [
                { title: '姓名', key: 'name' },
                { title: '邮箱', key: 'email' },
                {
                    title: '角色',
                    render: (row: any) => {
                        if (userStore.userInfo.user_id === row.user_id) {
                            return <span>{row.authority_name}</span>
                        }

                        return (
                            <a
                                class="inline-flex items-center el-icon__more"
                                onClick={(e: any) => onChangeRoleBtnClick(e, row)}
                                data-role={row.authority}
                                href="javascript:void(0)"
                            >
                                <span>{row.authority_name}</span>
                                <el-icon>
                                    <ArrowDown />
                                </el-icon>
                            </a>
                        )
                    },
                },
                {
                    title: '操作',
                    width: 160,
                    render: (row: any) => {
                        if (row.authority === 0) {
                            return []
                        }

                        let vDom = []
                        if (row.authority !== 2) {
                            vDom.push(
                                <span class="cursor-pointer text-blue-600 mr-2" onClick={() => onTransferProjectBtnClick(row)}>
                                    移交项目
                                </span>
                            )
                        }

                        if (userStore.userInfo.user_id !== row.user_id) {
                            vDom.push(
                                <span class="cursor-pointer text-red-400" onClick={() => onRemoveMemberBtnClick(row)}>
                                    移除成员
                                </span>
                            )
                        }
                        return vDom
                    },
                },
            ]

            const onChangeRoleBtnClick = (e: any, row: any) => {
                currentMember = row
                rolePopperRef.value = e.currentTarget
                isShowPopover.value = true
            }

            const onTransferProjectBtnClick = (row: any) => {
                AsyncMsgBox({
                    title: '移交提示',
                    content: '确定将项目移交给该成员吗？',
                    onOk: (done: any) => {
                        return transferProject(project.value.id, row.user_id)
                            .then((res: any) => {
                                $Message.success(res.msg || '移交成功！')
                                $router.replace({ name: 'projects' })
                            })
                            .catch(() => done())
                    },
                })
            }

            const onRemoveMemberBtnClick = (row: any) => {
                AsyncMsgBox({
                    title: '删除提示',
                    content: '确定移除该成员吗？',
                    onOk: () => {
                        return removeMember(project.value.id, row.user_id).then((res: any) => {
                            $Message.success(res.msg || '移除成功！')
                            emit('on-remove')
                            getTableData()
                        })
                    },
                })
            }

            const onRoleItemClick = (role: any) => {
                NProgress.start()

                changeMemberAuthority({
                    authority: role.value,
                    user_id: currentMember.user_id,
                    project_id: project.value.id,
                })
                    .then((res: any) => {
                        $Message.success(res.msg || '权限修改成功！')
                        currentMember.authority = role.value
                        currentMember.authority_name = role.text
                    })
                    .catch(() => {
                        //
                    })
                    .finally(() => {
                        NProgress.done()
                    })
            }

            onClickOutside(rolePopperRef, (e) => {
                const target = e.target as HTMLElement
                const parent = target.parentNode as HTMLElement
                const gParent = parent.parentNode as HTMLElement

                if (parent?.classList?.contains('el-icon__more') || gParent?.classList?.contains('el-icon__more')) {
                    return
                }

                isShowPopover.value = false
            })

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

            return {
                isLoading: isLoading as Ref<boolean>,
                roles,
                project,
                members,
                currentPage,
                total,
                columns,
                popoverRef,
                isShowPopover,
                rolePopperRef,

                onRoleItemClick,
                getTableData,
            }
        },
    })
</script>
