<template>
    <el-popover v-model:visible="visible" ref="popoverRef" popper-class="ac-popper-menu" :virtual-ref="popoverRefEl" virtual-triggering width="auto">
        <ul>
            <li class="ac-popper-menu__item" v-for="item in actions" :key="item.selector" @click="onPopperItemClick(item)" v-html="item.menuHtml"></li>
        </ul>
    </el-popover>

    <ProjectExportModal ref="projectExportModal" />
    <ProjectShareModal ref="projectShareModal" />
    <ChangeProjectGroupModal ref="changeProjectGroupModal" @on-ok="onChangeProjectGroupSuccess" />
</template>
<script setup lang="ts">
    import ProjectExportModal from '@/components/ProjectExportModal.vue'
    import ProjectShareModal from '@/components/ProjectShareModal.vue'
    import ChangeProjectGroupModal from './ChangeProjectGroupModal.vue'
    import { ElMessage as $Message } from 'element-plus'
    import { AsyncMsgBox } from '@/components/AsyncMessageBox'
    import { onClickOutside } from '@vueuse/core'
    import { ref, watch } from 'vue'
    import { deleteProject, quitProject } from '@/api/project'
    import { API_PROJECT_EXPORT_ACTION_MAPPING } from '@/api/exportFile'
    import { useRouter } from 'vue-router'
    import { PROJECT_ALL_ROLE_LIST, PROJECT_ROLES_KEYS, PROJECT_VISIBLE_TYPES } from '@/common/constant'

    const props = withDefaults(defineProps<{ ignoreEle: any }>(), {
        ignoreEle: [],
    })

    const emit = defineEmits(['changeGroup', 'quitProject', 'deleteProject'])
    const router = useRouter()

    const visible = ref(false)
    const popoverRefEl = ref(null)
    const projectInfo: any = ref({})

    const projectExportModal = ref()
    const projectShareModal = ref()
    const changeProjectGroupModal = ref()

    // 下拉菜单选项
    const actions: any = ref([])

    onClickOutside(popoverRefEl, (event) => {
        const composedPath = event.composedPath()
        const ignore = props.ignoreEle

        if (ignore && ignore.length > 0) {
            if (ignore.some((target: any) => target && (event.target === target || composedPath.includes(target)))) return
        }

        visible.value = false
    })

    const onChangeGroupBtnClick = () => {
        changeProjectGroupModal.value.show(projectInfo.value)
    }

    const onChangeProjectGroupSuccess = () => {
        emit('changeGroup')
    }

    const onDeleteProjectBtnClick = () => {
        AsyncMsgBox({
            title: '删除提示',
            content: '确定删除该项目吗？',
            onOk: (done: any) => {
                return deleteProject(projectInfo.value.id)
                    .then((res: any) => {
                        $Message.success(res.msg || '删除成功！')
                        emit('deleteProject')
                    })
                    .catch(() => {
                        done()
                        emit('quitProject')
                    })
            },
        })
    }

    const onQuitProjectBtnClick = () => {
        AsyncMsgBox({
            title: '退出提示',
            content: '确定退出该项目吗？',
            onOk: (done: any) => {
                return quitProject(projectInfo.value.id)
                    .then((res: any) => {
                        $Message.success(res.msg || '删除成功！')
                        emit('quitProject')
                    })
                    .catch(() => {
                        done()
                        emit('quitProject')
                    })
            },
        })
    }

    const onShareProjectBtnClick = () => {
        projectShareModal.value.show(projectInfo.value)
    }

    const onExportProjectBtnClick = () => {
        projectExportModal.value.show(
            {
                project_id: projectInfo.value.id,
            },
            API_PROJECT_EXPORT_ACTION_MAPPING
        )
    }

    const ALL_ACTIONS = [
        {
            text: '预览项目',
            selector: 'project_preview',
            isNewOpen: true,
            href: '{id}',
            field: 'preview_link',
            icon: 'iconIconPopoverPlay',
            roles: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER],
        },
        {
            text: '项目成员',
            selector: 'project_members',
            href: '/project/{id}/members',
            icon: 'iconIconPopoverUser',
            route: { name: 'project.members' },
            roles: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER, PROJECT_ROLES_KEYS.READER],
        },
        {
            text: '分享项目',
            selector: 'project_share',
            clickFn: onShareProjectBtnClick,
            icon: 'iconshare',
            roles: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER, PROJECT_ROLES_KEYS.READER],
        },
        {
            text: '导出项目',
            selector: 'project_export',
            clickFn: onExportProjectBtnClick,
            icon: 'iconIconPopoverUpload',
            roles: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER],
        },
        {
            text: '项目设置',
            selector: 'project_edit',
            href: '/project/{id}/setting',
            icon: 'iconIconPopoverSetting',
            route: { name: 'project.setting' },
            roles: [PROJECT_ROLES_KEYS.MANAGER],
        },
        {
            text: '项目分组',
            selector: 'project_group',
            clickFn: onChangeGroupBtnClick,
            icon: 'iconfenzu',
            roles: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER, PROJECT_ROLES_KEYS.READER],
        },
        {
            text: '删除项目',
            selector: 'project_delete',
            clickFn: onDeleteProjectBtnClick,
            icon: 'iconIconPopoverAlert',
            roles: [PROJECT_ROLES_KEYS.MANAGER],
        },
        {
            text: '退出项目',
            selector: 'project_quit',
            clickFn: onQuitProjectBtnClick,
            icon: 'iconIconPopoverExit',
            roles: [PROJECT_ROLES_KEYS.DEVELOPER, PROJECT_ROLES_KEYS.READER],
        },
    ]

    const ACTIONS_MAPPING = {} as any

    ;[PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER, PROJECT_ROLES_KEYS.READER].forEach((key) => {
        const menus = (ACTIONS_MAPPING[key] = ACTIONS_MAPPING[key] || [])
        ALL_ACTIONS.forEach((item) => {
            const keyword = item.roles.join(',')
            if (keyword.indexOf(key) !== -1) {
                menus.push(item)
            }
        })
    })

    // 生成下拉菜单项
    const generateActionsByProject = () => {
        const roleInfo = PROJECT_ALL_ROLE_LIST.find((item: any) => item.key === projectInfo.value.authority) as any
        let key = (roleInfo || {}).key

        let actionArray = ACTIONS_MAPPING[key] || []

        // 阅读者&私有项目不能分享
        if (key === PROJECT_ROLES_KEYS.READER && projectInfo.value.visibility === PROJECT_VISIBLE_TYPES.PRIVATE) {
            actionArray = actionArray.filter((item: any) => item.selector !== 'project_share')
        }

        actions.value = actionArray.map((item: any) => {
            let cp = { ...item }
            if (cp.href) {
                cp.href = cp.href.replace('{id}', projectInfo.value[item.field || 'id'])
            }
            cp.menuHtml = `<i class="icon iconfont mr-1 ${cp.icon || ''}"></i>${cp.text}`
            return cp
        })
    }

    const onPopperItemClick = (item: any) => {
        if (item.clickFn) {
            item.clickFn()
            return
        }

        if (item.isNewOpen) {
            window.open(item.href)
            return
        }

        if (item.href && !item.isNewOpen) {
            item.route.params = { project_id: projectInfo.value[item.field || 'id'] }
            router.push(item.route)
        }
    }

    const show = (el: any, project: any) => {
        popoverRefEl.value = el
        projectInfo.value = project
        visible.value = true
    }

    watch(projectInfo, () => {
        generateActionsByProject()
    })

    defineExpose({
        show,
    })
</script>
