<template>
    <header class="ac-doc-header is-vertical">
        <div class="ac-doc-header-inner">
            <div class="ac-doc-header--left">
                <div class="ac-project-info">
                    <el-popover placement="bottom" popper-class="ac-popper-menu ac-popper-menu--large" width="250px">
                        <template #reference>
                            <div class="ac-project-info__img">
                                <img :src="project.icon" :alt="project.name" />
                                <router-link :to="{ name: 'projects' }">
                                    <el-icon class="ac-project-info__back"><ArrowLeftBold /></el-icon>
                                </router-link>
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

                    <div class="ac-project-info__title">{{ project.name }}</div>
                    <el-tooltip effect="dark" content="私有项目" placement="bottom">
                        <el-icon class="ac-project-info__icon"><Lock /></el-icon>
                    </el-tooltip>
                </div>
            </div>
            <div class="ac-doc-header--right">
                <div class="ac-document__operate">
                    <div class="ac-document__fixed_title">{{ project.name }}</div>

                    <div class="ac-document__btns">
                        <el-button type="primary"> 保存 </el-button>
                    </div>
                </div>
            </div>
        </div>
    </header>

    <ProjectShareModal ref="projectShareModal" />
</template>
<script setup lang="ts">
    import { ArrowLeftBold, Lock } from '@element-plus/icons-vue'
    import { ref, inject, watch } from 'vue'
    import { useRouter } from 'vue-router'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'
    import { toPreviewProjectPath } from '@/router/preview.router'
    import { API_PROJECT_EXPORT_ACTION_MAPPING } from '@/api/exportFile'

    const { push, currentRoute } = useRouter()
    const projectStore = useProjectStore()
    const { projectInfo: project, isManager } = storeToRefs(projectStore)

    const projectShareModal = ref()
    const projectExportModal: any = inject('projectExportModal')
    const previewUrl = toPreviewProjectPath({ project_id: currentRoute.value.params.project_id })

    const onPopperItemClick = (menu: any) => {
        if (menu.href) {
            push(menu.href.replace('{project_id}', currentRoute.value.params.project_id))
            return
        }
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
        { text: '项目设置', icon: 'iconIconPopoverSetting', href: '/project/{project_id}/setting' },
        { text: '公共参数', icon: 'iconIconPopoverConfig', href: '/project/{project_id}/params' },
        { text: '成员管理', icon: 'iconIconPopoverUser', href: '/project/{project_id}/members' },
        { text: '分享项目', icon: 'iconshare2', onClick: (menu?: any) => onShareProjectBtnClick() },
        { text: '导出项目', icon: 'iconIconPopoverUpload', onClick: (menu?: any) => onExportBtnClick() },
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
