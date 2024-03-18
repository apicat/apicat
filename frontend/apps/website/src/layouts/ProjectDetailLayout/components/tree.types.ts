export interface TreeData<T> {
  id: string | number
  name: string
  isLeaf: boolean
  status: 'none' | 'loading' | 'success' | 'error'
  children?: TreeData<T>[]
  data: T
}

export type TreeNodeWrapper<T> = T & {
  _isTemp: boolean
  _status: 'none' | 'loading' | 'success' | 'error'
}
