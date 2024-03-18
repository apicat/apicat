declare module UserAPI {
  interface ResponseUserInfo {
    id: number
    email: string
    name: string
    github: boolean
    language: string
    avatar?: string
    role: 'user' | 'admin'
  }

  interface RequestGeneral {
    language: string
    name: string
  }

  interface RequestChangeAvatar {
    avatar?: File
    croppedX: number
    croppedY: number
    croppedWidth: number
    croppedHeight: number
  }

  interface ResponseChangeAvatar {
    avatar: string
  }

  interface RequestResetPassword {
    password: string
    newPassword: string
    reNewPassword: string
  }
}
