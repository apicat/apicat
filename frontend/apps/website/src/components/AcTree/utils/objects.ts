import { get, set } from 'lodash-unified'

export { hasOwn } from '@vue/shared'

export const keysOf = <T>(arr: any) => Object.keys(arr) as Array<keyof T>

export const getProp = <T = any>(obj: Record<string, any>, path: Arrayable<string>, defaultValue?: any): { value: T } => {
  return {
    get value() {
      return get(obj, path, defaultValue)
    },
    set value(val: any) {
      set(obj, path, val)
    },
  }
}
