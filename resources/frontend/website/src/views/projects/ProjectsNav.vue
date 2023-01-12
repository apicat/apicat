<template>
    <SidebarLayout>
        <div class="flex flex-col bg-white border divide-y">
            <a
                class="relative flex items-center h-12 pl-6 text-neutral-600 hover:text-neutral-900"
                href="javascript:void(0);"
                :class="projectGroupActiveClass(0)"
                @click="onSwitchProjectGroup()"
            >
                <el-icon><icon-menu /></el-icon><span class="ml-1">所有项目</span>
            </a>

            <div class="divide-y ac-group-nav" ref="sortableList" v-show="projectGroupList.length">
                <a
                    v-for="group in projectGroupList"
                    :key="group.name"
                    class="relative flex items-center justify-between pr-4 text-neutral-600 hover:text-neutral-900 ac-group-nav__item"
                    :class="projectGroupActiveClass(group.id)"
                    href="javascript:void(0);"
                >
                    <section class="flex items-center flex-1 h-12 pl-6 overflow-hidden" @click="onSwitchProjectGroup(group)">
                        <el-icon><folder /></el-icon><span class="ml-1 truncate" :title="group.name">{{ group.name }}</span>
                    </section>
                    <el-icon @click="onShowMorePopover($event, group)" class="ml-1 el-icon__more"><more-filled /></el-icon>
                </a>
            </div>

            <a
                class="relative flex items-center h-12 pl-6 text-neutral-600 hover:text-neutral-900"
                href="javascript:void(0);"
                @click="onAddProjectCategoryBtnClick"
            >
                <el-icon><folder-add /></el-icon><span class="ml-1">添加分组</span>
            </a>
        </div>
    </SidebarLayout>

    <AddProjectCategoryModal ref="addProjectCategoryModal" />

    <el-popover popper-class="ac-popper-menu" width="auto" v-model:visible="isShowPopover" :virtual-ref="moreIconRef" virtual-triggering>
        <ul>
            <li class="ac-popper-menu__item" @click="onRenameGroupBtnClick">
                <el-icon class="mr-1"><edit /></el-icon>重命名
            </li>
            <li class="ac-popper-menu__item" @click="onDeleteGroupBtnClick">
                <el-icon class="mr-1"><delete /></el-icon>删除
            </li>
        </ul>
    </el-popover>
</template>

<script setup lang="ts">
    import SidebarLayout from '@/layout/SidebarLayout.vue'
    import AddProjectCategoryModal from './components/AddProjectCategoryModal.vue'
    import { Folder, FolderAdd, Menu as IconMenu, MoreFilled, Edit, Delete } from '@element-plus/icons-vue'
    import { onClickOutside } from '@vueuse/core'
    import { onMounted, onUnmounted, ref } from 'vue'
    import { storeToRefs } from 'pinia'
    import { useProjectsStore } from '@/stores/projects'
    import { AsyncMsgBox } from '@/components/AsyncMessageBox'
    import { ElMessage as $Message } from 'element-plus'
    import { useSortable } from '@/hooks/useSortable'

    const projectStore = useProjectsStore()
    const { activeGroup, projectGroupList } = storeToRefs(projectStore)

    const addProjectCategoryModal = ref<any>(null)
    const isShowPopover = ref(false)
    const moreIconRef = ref()
    const sortableList = ref()

    let currentGroup: any = null

    const projectGroupActiveClass = (groupId: number) => {
        return [
            {
                active: activeGroup.value.id === groupId,
            },
        ]
    }

    const onAddProjectCategoryBtnClick = () => {
        addProjectCategoryModal.value?.show()
    }

    const onSwitchProjectGroup = (group?: any) => {
        projectStore.switchProjectGroup(group)
    }

    const onShowMorePopover = (e: MouseEvent, group: any) => {
        currentGroup = group
        moreIconRef.value = e.currentTarget
        isShowPopover.value = true
    }

    const onRenameGroupBtnClick = () => {
        if (currentGroup) {
            addProjectCategoryModal.value?.show(currentGroup)
        }
    }

    const onDeleteGroupBtnClick = () => {
        if (currentGroup) {
            AsyncMsgBox({
                title: '删除提示',
                content: `确定删除${currentGroup.name}分组？`,
                onOk: () => {
                    return projectStore.deleteProjectGroup(currentGroup).then((res: any) => {
                        $Message.success(res.msg || '删除成功！')
                    })
                },
            })
        }
    }

    onClickOutside(moreIconRef, (e) => {
        const target = e.target as HTMLElement
        const parent = target.parentNode as HTMLElement
        const gParent = parent.parentNode as HTMLElement

        if (parent?.classList?.contains('el-icon__more') || gParent?.classList?.contains('el-icon__more')) {
            return
        }

        isShowPopover.value = false
    })

    const { initSortable } = useSortable(sortableList, {
        draggable: '.ac-group-nav__item',
        onEnd(e: any) {
            projectStore.sortProjectGroup(e.oldIndex, e.newIndex)
        },
    })

    onMounted(async () => {
        await projectStore.getProjectGroupList()
        initSortable()
    })

    onUnmounted(() => {
        currentGroup = null
    })
</script>

<style lang="scss" scoped>
    @use '@/assets/stylesheet/mixins/mixins.scss' as *;
    @include b(group-nav) {
        @include e(item) {
            &.sortable-drag {
                background-color: rgb(250, 250, 250);
            }

            .el-icon__more {
                opacity: 0;
            }
            &:hover .el-icon__more {
                opacity: 1;
            }
        }
    }
</style>
