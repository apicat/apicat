<template>
    <div class="ac-project-info ac-project-info--hover">
        <template v-if="isManager || isDeveloper">
            <template v-if="isIterateRoute">
                <div class="ac-project-info__img">
                    <img :src="project.icon" :alt="project.name" />
                    <a :href="isIterateRoute ? '/iterations' : '/home'">
                        <el-icon class="ac-project-info__back"><ArrowLeftBold /></el-icon>
                    </a>
                </div>
            </template>

            <el-popover v-else placement="bottom" popper-class="ac-popper-menu ac-popper-menu--large" width="250px">
                <template #reference>
                    <div class="ac-project-info__img">
                        <img :src="project.icon" :alt="project.name" />
                        <a :href="isIterateRoute ? '/iterations' : '/home'">
                            <el-icon class="ac-project-info__back"><ArrowLeftBold /></el-icon>
                        </a>
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
                        <i class="mr-1 icon iconfont" :class="menu.icon" />{{ menu.text }}
                    </li>
                </ul>
            </el-popover>
        </template>

        <template v-else>
            <div class="ac-project-info__img">
                <img :src="project.icon" :alt="project.name" />
                <a :href="isIterateRoute ? '/iterations' : '/home'">
                    <el-icon class="ac-project-info__back"><ArrowLeftBold /></el-icon>
                </a>
            </div>
        </template>

        <template v-if="isIterateRoute && iterateInfo">
            <div class="pr-2 overflow-hidden">
                <p class="text-base truncate" :title="project.name">{{ project.name }}</p>
                <p class="text-sm truncate" :title="iterateInfo.title">{{ iterateInfo.title }}</p>
            </div>
            <el-tooltip effect="dark" content="私有项目" placement="bottom" v-if="isPrivate">
                <el-icon :size="16"><Lock /></el-icon>
            </el-tooltip>
        </template>

        <div v-else class="ac-project-info__title" :title="project.name">
            {{ project.name }}
            <el-tooltip effect="dark" content="私有项目" placement="bottom" v-if="isPrivate">
                <el-icon class="ac-project-info__icon"><Lock /></el-icon>
            </el-tooltip>
        </div>
    </div>
    <ProjectShareModal ref="projectShareModal" />
</template>
<script setup lang="ts">
    import { ArrowLeftBold, Lock } from '@element-plus/icons-vue'
    import { ref, inject, watch } from 'vue'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'
    import { API_PROJECT_EXPORT_ACTION_MAPPING } from '@/api/exportFile'
    import { generateProjectMembersUrl, generateProjectParamsUrl, generateProjectSettingUrl, generateProjectTrashUrl } from '@/api/project'
    import { useIterateStore } from '@/stores/iterate'
    import { getIdPublicByRouter } from '@/hooks/useIdPublicParam'

    const projectStore = useProjectStore()
    const iterateStore = useIterateStore()
    const { isIterateRoute, iterateInfo } = storeToRefs(iterateStore)
    const { projectInfo: project, isManager, isPrivate, isDeveloper } = storeToRefs(projectStore)

    const projectShareModal = ref()
    const projectExportModal: any = inject('projectExportModal')
    const id_public = getIdPublicByRouter()

    const onPopperItemClick = (menu: any) => {
        if (menu.hrefFn) {
            location.href = menu.hrefFn(id_public, false)
            return
        }
    }

    const onExportBtnClick = () => {
        projectExportModal.value.title = '导出项目'
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
        { text: '项目设置', icon: 'iconIconPopoverSetting', hrefFn: generateProjectSettingUrl },
        { text: '公共参数', icon: 'iconIconPopoverConfig', hrefFn: generateProjectParamsUrl },
        { text: '项目成员', icon: 'iconIconPopoverUser', hrefFn: generateProjectMembersUrl },
        { text: '分享项目', icon: 'iconshare2', onClick: () => onShareProjectBtnClick() },
        { text: '导出项目', icon: 'iconIconPopoverUpload', onClick: () => onExportBtnClick() },
        { text: '回收站', icon: 'icontrash', hrefFn: generateProjectTrashUrl },
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
