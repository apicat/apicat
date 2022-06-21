import { RouteRecordRaw } from 'vue-router'
import { getRouteNormalInfo } from '@ac/shared'
import { PROJECT_ROLES_MAP } from '@/common/constant'

const ProjectNav = () => import('@/views/project/ProjectNav.vue')
const ProjectSetting = () => import('@/views/project/ProjectSetting.vue')
const ProjectMembers = () => import('@/views/project/ProjectMembers.vue')
const ProjectParam = () => import('@/views/project/ProjectParam.vue')
const ProjectTrash = () => import('@/views/project/ProjectTrash.vue')

const ProjectRouteList: RouteRecordRaw[] = [
    {
        path: 'setting',
        name: 'project.setting',
        component: ProjectSetting,
        meta: { title: '项目设置', role: [PROJECT_ROLES_MAP.MANAGER] },
    },
    {
        path: 'members',
        name: 'project.members',
        component: ProjectMembers,
        meta: { title: '项目成员', role: [PROJECT_ROLES_MAP.MANAGER, PROJECT_ROLES_MAP.DEVELOPER, PROJECT_ROLES_MAP.READER] },
    },
    {
        path: 'params',
        name: 'project.params',
        component: ProjectParam,
        meta: { title: '公共参数', role: [PROJECT_ROLES_MAP.MANAGER, PROJECT_ROLES_MAP.DEVELOPER] },
    },
    {
        path: 'trash',
        name: 'project.trash',
        component: ProjectTrash,
        meta: { title: '回收站', role: [PROJECT_ROLES_MAP.MANAGER, PROJECT_ROLES_MAP.DEVELOPER] },
    },
]

export const ProjectRoutes = getRouteNormalInfo(ProjectRouteList)

export default {
    path: '/project/:project_id',
    name: 'project',
    component: ProjectNav,
    redirect: { name: 'project.setting' },
    children: ProjectRouteList,
}
