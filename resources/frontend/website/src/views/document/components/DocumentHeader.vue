<template>
    <header class="ac-doc-header">
        <div class="ac-doc-header-inner">
            <div class="w-full relative text-center">
                <router-link class="absolute left-0 top-0 h-full flex items-center text-zinc-500" :to="{ name: 'projects' }">
                    <el-icon class="mr-1"><arrow-left-bold /></el-icon>
                    返回
                </router-link>

                <div class="h-14 text-base font-medium">
                    <el-popover placement="bottom" trigger="click" popper-class="ac-popper-menu ac-popper-menu--small" width="auto">
                        <template #reference>
                            <div class="flex items-center justify-center h-full max-w-sm m-auto cursor-pointer hover:text-zinc-600">
                                <p class="truncate">{{ project?.name }}</p>
                                <el-icon v-show="project" :size="12" class="ml-1"><arrow-down-bold /></el-icon>
                            </div>
                        </template>

                        <ul>
                            <li
                                v-for="menu in popperMenus"
                                :key="menu.icon"
                                class="ac-popper-menu__item"
                                :class="{ 'border-t': menu.divided }"
                                @click="menu.onClick ? menu.onClick(menu) : onPopperItemClick(menu)"
                            >
                                <i class="icon iconfont mr-1" :class="menu.icon" />{{ menu.text }}
                            </li>
                        </ul>
                    </el-popover>
                </div>

                <div class="absolute right-0 top-0 h-full flex items-center text-zinc-500">
                    <a class="flex items-center h-full" href="javascript:void(0)" @click="onSearchTextClick">
                        <el-icon class="mr-2"><Search /></el-icon>搜索
                    </a>

                    <el-divider class="mx-5" direction="vertical" />

                    <router-link class="flex items-center" :to="{ path: previewUrl }" target="_blank">
                        <el-icon class="mr-2"><View /></el-icon>预览项目
                    </router-link>

                    <el-divider class="mx-5" direction="vertical" />

                    <a class="flex items-center" href="javascript:void(0)" @click="onShareProjectBtnClick">
                        <el-icon class="mr-2"><share /></el-icon>分享
                    </a>
                </div>
            </div>
        </div>
    </header>

    <ProjectShareModal ref="projectShareModal" />
    <SearchDocumentPopover ref="searchDocumentPopoverRef" />
</template>
<script setup lang="ts">
    import { ArrowLeftBold, ArrowDownBold, Search, View, Share } from '@element-plus/icons-vue'
    import { ref, inject, watch } from 'vue'
    import { useRouter } from 'vue-router'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'
    import { toPreviewProjectPath } from '@/router/preview.router'
    import { API_PROJECT_EXPORT_ACTION_MAPPING } from '@/api/exportFile'
    import SearchDocumentPopover from './SearchDocumentPopover.vue'

    const { push, currentRoute } = useRouter()
    const projectStore = useProjectStore()
    const { projectInfo: project, isManager } = storeToRefs(projectStore)

    const projectShareModal = ref()
    const searchDocumentPopoverRef = ref()
    const projectExportModal: any = inject('projectExportModal')
    const previewUrl = toPreviewProjectPath({ project_id: currentRoute.value.params.project_id })

    const onPopperItemClick = (menu: any) => {
        if (menu.href) {
            push(menu.href.replace('{project_id}', currentRoute.value.params.project_id))
            return
        }
    }

    const onSearchTextClick = (e: any) => {
        searchDocumentPopoverRef.value?.show(e.currentTarget)
    }

    const onExportBtnClick = () => {
        projectExportModal.value?.show(
            {
                project_id: project.value.id,
            },
            API_PROJECT_EXPORT_ACTION_MAPPING
        )
    }

    const onShareProjectBtnClick = () => {
        projectShareModal.value?.show(project.value)
    }

    const allMenus = [
        { text: '公共参数', icon: 'iconIconPopoverConfig', href: '/project/{project_id}/params' },
        { text: '项目设置', icon: 'iconIconPopoverSetting', href: '/project/{project_id}/setting' },
        { text: '项目成员', icon: 'iconIconPopoverUser', href: '/project/{project_id}/members' },
        { text: '导出项目', icon: 'iconIconPopoverUpload', onClick: (menu?: any) => onExportBtnClick(), divided: true },
        { text: '回收站', icon: 'icontrash', href: '/project/{project_id}/trash' },
    ]

    const popperMenus: any = ref(allMenus)

    watch(
        project,
        () => {
            if (project.value && !isManager.value) {
                popperMenus.value = allMenus.filter((item: any) => item.icon !== 'iconIconPopoverSetting')
            }
        },
        {
            immediate: true,
        }
    )
</script>
