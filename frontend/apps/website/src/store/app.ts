import { defineStore } from 'pinia'

export type CreateModeType = 'document' | 'schema'

export enum CreateModeEnum {
  document = 'document',
  schema = 'schema',
}

interface AppState {
  createMode: CreateModeType | null
}

export const uesAppStore = defineStore('app', {
  state: (): AppState => ({
    // 标识创建文档|模型
    createMode: null,
  }),
  actions: {
    setCreateMode(mode: CreateModeType | null) {
      this.createMode = mode
    },
  },
})

export default uesAppStore
