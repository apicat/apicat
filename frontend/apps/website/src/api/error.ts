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

// 没有团队错误
export class NoTeamError extends Error {
  constructor(message?: string) {
    super(message)
    this.name = 'NoTeamError'
  }
}

class ResponseBaseError<T> extends Error {
  name: string = 'BaseError'
  response: T | undefined
  constructor(response?: T, message?: string) {
    super(message)
    this.response = response
  }
}

// http status 400
export class BadRequestError<T> extends ResponseBaseError<T> {
  name = 'BadRequestError'
}
// 无权限访问错误 403
export class NoPermissionError<T> extends ResponseBaseError<T> {
  name = 'NoPermissionError'
}
// 资源未找到 404
export class NotFoundError<T> extends ResponseBaseError<T> {
  name = 'NotFoundError'
}
// 缺失或错误认证 401
export class UnauthorizedError<T> extends ResponseBaseError<T> {
  name = 'UnauthorizedError'
}
