<template>
    <div class="ac-project-info">
        <el-popover placement="bottom" popper-class="ac-popper-menu ac-popper-menu--large" width="250px">
            <template #reference>
                <div class="ac-project-info__img">
                    <img :src="project.icon" :alt="project.name" />
                </div>
            </template>

            <ul>
                <li v-for="menu in popperMenus" :key="menu.icon" class="ac-popper-menu__item">
                    <a :href="menu.href" target="_blank">{{ menu.text }}</a>
                </li>
            </ul>
        </el-popover>

        <div class="ac-project-info__title" :title="project.name">{{ project.name }}</div>
        <el-tooltip effect="dark" content="私有项目" placement="bottom" v-if="isPrivate">
            <el-icon class="ac-project-info__icon"><Lock /></el-icon>
        </el-tooltip>
    </div>
</template>
<script setup lang="ts">
    import { Lock } from '@element-plus/icons-vue'
    import { storeToRefs } from 'pinia'
    import { useProjectStore } from '@/stores/project'

    const projectStore = useProjectStore()
    const { projectInfo: project, isPrivate } = storeToRefs(projectStore)

    // 获取项目所设置的导航列表
    const popperMenus = [
        { text: '项目设置', icon: 'iconIconPopoverSetting', href: '/project/{project_id}/setting' },
        { text: '公共参数', icon: 'iconIconPopoverConfig', href: '/project/{project_id}/params' },
        { text: '成员管理', icon: 'iconIconPopoverUser', href: '/project/{project_id}/members' },
    ]
</script>
