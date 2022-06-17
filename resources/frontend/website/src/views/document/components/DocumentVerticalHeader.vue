<template>
    <header class="ac-doc-header is-vertical">
        <div class="ac-doc-header-inner">
            <div class="ac-doc-header--left">left</div>
            <div class="ac-doc-header--right">right</div>
        </div>
    </header>

    <ProjectShareModal ref="projectShareModal" />
</template>
<script setup lang="ts">
    import { ArrowLeftBold, ArrowDownBold, Search, View, Share } from '@element-plus/icons-vue'
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
