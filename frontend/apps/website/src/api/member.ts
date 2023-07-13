import { queryStringify } from '@/commons'
import Ajax, { QuietAjax } from './Ajax'
import { UserInfo } from '@/typings/user'

// 获取成员列表
export const getMembers = async (params: Record<string, any>) => Ajax.get('/members' + queryStringify(params))
// 添加成员
export const addMember = async (user: Partial<UserInfo>) => Ajax.post('/members', user)
// 修改成员信息
export const updateMember = async ({ id, ...params }: Partial<UserInfo>) => QuietAjax.put(`/members/${id}`, params)
// 删除成员
export const deleteMember = async (id: number) => Ajax.delete(`/members/${id}`)
