/**
 * 目标成员在项目中的权限发生变更，错误
 */
export class TargetMemberPermissionError extends Error {
  constructor(message?: string) {
    super(message)
    this.name = 'TargetMemberPermissionError'
  }
}

// 未找到被分享的文档
export class NotFountDocumentShare extends Error {
  constructor(message?: string) {
    super(message)
    this.name = 'NotFountDocumentShare'
  }
}

// 分享密钥错误
export class ShareSecretKeyError extends Error {
  constructor(message?: string) {
    super(message)
    this.name = 'ShareSecretKeyError'
  }
}
