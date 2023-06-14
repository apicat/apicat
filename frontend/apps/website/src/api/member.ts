import Ajax, { QuietAjax } from './Ajax'
import { UserInfo } from '@/typings/user'

// 获取成员列表
export const getMembers = async (data: Record<string, any>) => Ajax.get('/members?' + new URLSearchParams(data).toString())
// 添加成员
export const addMember = async (user: Partial<UserInfo>) => Ajax.post('/members', user)
// 修改成员信息
export const updateMember = async ({ id, ...other }: Partial<UserInfo>) => QuietAjax.put(`/members/${id}`, other)
// 删除成员
export const deleteMember = async (id: number) => Ajax.delete(`/members/${id}`)
