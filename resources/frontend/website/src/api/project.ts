import Ajax from './Ajax'
import { compile } from 'path-to-regexp'
import { wrapperOrigin } from '@/common/utils'
import { PROJECT_PREVIEW_PATH, PROJECT_SETTING_PATH, PROJECT_MEMBERS_PATH, PROJECT_PARAMS_PATH, PROJECT_TRASH_PATH } from '@/router/constant'

export const createProject = (project = {}) => Ajax.post('/project/create', project)

export const uploadProjectIcon = (data = {}) => Ajax.post('/project/icon', data)

export const uploadNavLogo = (data = {}) => Ajax.post('/project/change_navigation_logo', data)

export const deleteProject = (project_id: unknown) => Ajax.post('/project/remove', { project_id })

export const quitProject = (project_id: unknown) => Ajax.post('/project/quit', { project_id })

export const addMember = (member = {}) => Ajax.post('/project/add_member', member)

export const getProjectMembers = (data: unknown) => Ajax.get('/project/members', { params: data })

export const removeMember = (project_id: unknown, user_id: unknown) => Ajax.post('/project/remove_member', { project_id, user_id })

export const transferProject = (project_id: unknown, user_id: unknown) => Ajax.post('/project/transfer', { project_id, user_id })

export const changeMemberAuthority = (data = {}) => Ajax.post('/project/change_member_authority', data)

export const settingProject = (project = {}) => Ajax.post('/project/setting', project)

export const navigation = (navigation = {}) => Ajax.post('/project/navigation_setting', navigation)

export const share = (project_id: unknown, share: unknown) => Ajax.post('/project/share', { project_id, share })

export const resetSecretkey = (project_id: unknown) => Ajax.post('/project/reset_share_secretkey', { project_id })

export const changeManager = (data = {}) => Ajax.post('/project/change_manager', data)

export const restoreDocument = (data = {}) => Ajax.post('/project/trash_restore', data)
export const restoreApiDocument = (data = {}) => Ajax.post('/doc/restore_api_doc', data)

// 根据角色获取项目列表
export const getProjectListByRole = (authority: unknown) => Ajax.get('/projects/json', { params: { authority } })

// 获取项目列表
export const getProjectList = (group_id: number) => Ajax.get('/projects', { params: { group_id: group_id == 0 ? '' : group_id } })

// 创建项目分类
export const createProjectGroup = (data: unknown) => Ajax.post('/project_group/create', data)

// 获取项目分类列表
export const getProjectGroupList = () => Ajax.get('/project_groups')

// 删除项目分类
export const removeProjectGroup = (id: unknown) => Ajax.post('/project_group/remove', { id })

// 重命名分类名称
export const renameProjectGroup = (data: unknown) => Ajax.post('/project_group/rename', data)

// 排序
export const sortProjectGroup = (ids: unknown) => Ajax.post('/project_group/change_order', { ids })

// 获取项目详情
export const getProjectDetail = (project_id: unknown, token?: unknown) => Ajax.get('/project', { params: { project_id, token } })

// 获取项目状态
export const getProjectStatus = (project_id: unknown) => Ajax.get('/project/status', { params: { project_id } })

// 更换项目分组
export const changeProjectGroup = (data: unknown) => Ajax.post('/project/change_group', data)

//获取项目导航详情
export const getProjectNavigationDetail = (project_id: unknown) => Ajax.get('/project/navigation', { params: { project_id } })

// 获取项目模板
export const getProjectTemplateList = (project_id: unknown) => Ajax.get('/project/templates', { params: { project_id } })

// 获取项目回收站文档
export const getProjectTrashList = (project_id: unknown) => Ajax.get('/doc/trash', { params: { project_id } })

// 不在此项目中的成员
export const getWithoutProjectMemberList = (project_id: unknown) => Ajax.get('/project/without_members', { params: { project_id } })

// 生成预览链接地址
export const generateProjectPreviewUrl = (project_id: string, hasOrigin?: true) => wrapperOrigin(hasOrigin) + compile(PROJECT_PREVIEW_PATH)({ project_id })
// 生成项目详情链接地址
export const generateProjectDetailUrl = generateProjectPreviewUrl

// 生成项目设置链接地址
export const generateProjectSettingUrl = (project_id: string, hasOrigin?: true) => wrapperOrigin(hasOrigin) + compile(PROJECT_SETTING_PATH)({ project_id })
// 生成项目成员链接地址
export const generateProjectMembersUrl = (project_id: string, hasOrigin?: true) => wrapperOrigin(hasOrigin) + compile(PROJECT_MEMBERS_PATH)({ project_id })
// 生成项目参数链接地址
export const generateProjectParamsUrl = (project_id: string, hasOrigin?: true) => wrapperOrigin(hasOrigin) + compile(PROJECT_PARAMS_PATH)({ project_id })
// 生成项目回收站链接地址
export const generateProjectTrashUrl = (project_id: string, hasOrigin?: true) => wrapperOrigin(hasOrigin) + compile(PROJECT_TRASH_PATH)({ project_id })
