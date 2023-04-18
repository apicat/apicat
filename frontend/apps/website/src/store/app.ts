import { defineStore } from 'pinia'

interface AppState {}

export const uesAppStore = defineStore('app', {
  state: (): AppState => ({}),
  actions: {},
})

export default uesAppStore
