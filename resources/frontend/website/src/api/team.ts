import Ajax from './Ajax'

// 获取团队成员
export const getMembers = (data = {}) => Ajax.get('/team/members', { params: data })

// 获取成员信息
export const getMemberInfo = (user_id: number) =>
    Ajax.get('/team/member?user_id=' + user_id).then((res) => {
        if (res.data) {
            const { name, avatar } = res.data
            if (!avatar) {
                res.data.avatar = (name || '无').substr(0, 1).toUpperCase()
            }
        }
        return res
    })

// 获取成员所参与的项目
export const getMemberJoinProjectList = (user_id: number) => Ajax.get('/team/member_projects?user_id=' + user_id)

// 添加成员
export const addMember = (data = {}) => Ajax.post('/team/add_member', data)

// 修改成员信息
export const modifyMember = (data = {}) => Ajax.post('/team/edit_member_info', data)

// 修改成员密码
export const modifyMemberPassword = (data = {}) => Ajax.post('/team/edit_member_password', data)

// 移出成员
export const removeMember = (user_id: number) => Ajax.post('/team/remove_member', { user_id })
