/**
 * 目标成员在项目中的权限发生变更，错误
 */

export class TargetMemberPermissionError extends Error {
  constructor(message?: string) {
    super(message)
    this.name = 'TargetMemberPermissionError'
  }
}
