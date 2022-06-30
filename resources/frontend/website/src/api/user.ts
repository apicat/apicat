import { useApi } from '@/hooks/useApi'
import { getFirstChar } from '@natosoft/shared'
import Ajax from './Ajax'

// 用户登录
export const userEmailLogin = (data = {}) => Ajax.post('/email_login', { ...data })
export const userLogin = userEmailLogin
// 用户注册
export const userRegister = (data = {}) => Ajax.post('/email_register', { ...data })
// 退出
export const logout = () => Ajax.post('/logout')
// 上传头像
export const settingAvatar = (data = {}) => Ajax.post('/user/change_avatar', data)
// 更新个人信息
export const updateUserProfile = (data = {}) => Ajax.post('/user/change_profile', { ...data })
// 获取个人信息
export const getUserInfo = () =>
    Ajax.get('/user/profile').then((res) => {
        if (res.data) {
            const { name, avatar } = res.data
            if (!avatar) {
                res.data.lastName = getFirstChar(name || '无')
            }
        }
        return res
    })
// 修密码
export const modifyPassword = () => useApi((data = {}) => Ajax.post('/user/change_password', { ...data }))
