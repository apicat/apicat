import { RouteRecordRaw } from 'vue-router'
import { getRouteNormalInfo } from '@natosoft/shared'
import { PROJECT_ROLES_KEYS } from '@/common/constant'
import { PROJECT_SETTING_PATH, PROJECT_MEMBERS_PATH, PROJECT_PARAMS_PATH, PROJECT_TRASH_PATH } from './constant'

const ProjectNav = () => import('@/views/project/ProjectNav.vue')
const ProjectSetting = () => import('@/views/project/ProjectSetting.vue')
const ProjectMembers = () => import('@/views/project/ProjectMembers.vue')
const ProjectParam = () => import('@/views/project/ProjectParam.vue')
const ProjectTrash = () => import('@/views/project/ProjectTrash.vue')

const ProjectRouteList: RouteRecordRaw[] = [
    {
        path: PROJECT_SETTING_PATH,
        name: 'project.setting',
        component: ProjectSetting,
        meta: { title: '项目设置', role: [PROJECT_ROLES_KEYS.MANAGER] },
    },
    {
        path: PROJECT_MEMBERS_PATH,
        name: 'project.members',
        component: ProjectMembers,
        meta: { title: '项目成员', role: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER, PROJECT_ROLES_KEYS.READER] },
    },
    {
        path: PROJECT_PARAMS_PATH,
        name: 'project.params',
        component: ProjectParam,
        meta: { title: '公共参数', role: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER] },
    },
    {
        path: PROJECT_TRASH_PATH,
        name: 'project.trash',
        component: ProjectTrash,
        meta: { title: '回收站', role: [PROJECT_ROLES_KEYS.MANAGER, PROJECT_ROLES_KEYS.DEVELOPER] },
    },
]

export const ProjectRoutes = getRouteNormalInfo(ProjectRouteList)

export default {
    path: '/project/:project_id',
    name: 'project',
    component: ProjectNav,
    children: ProjectRouteList,
}
