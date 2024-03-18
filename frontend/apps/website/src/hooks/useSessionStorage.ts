/* eslint-disable no-restricted-properties */
export const useSessionStorage = (function () {
  const listeners: Record<string, Ref<string | null>> = {}
  const keys = {
    setItem: (key: any, val: any) => {
      sessionStorage[key] = val
      const r = listeners[key]
      if (r)
        r.value = val
    },
    removeItem: (key: any) => {
      delete sessionStorage[key]
      const r = listeners[key]
      if (r)
        r.value = null
    },
    clear: () => {
      const items = Object.entries(sessionStorage)
      for (let i = 0; i < items.length; i++) {
        const key = items[i][0]
        sessionStorage.removeItem(key)
      }
    },
  } as any
  for (const key in keys) {
    // eslint-disable-next-line no-proto
    Object.defineProperty(window.sessionStorage.__proto__, key, {
      get() {
        return keys[key]
      },
    })
  }

  function listen(key: string): Ref<string | null> {
    if (listeners[key]) {
      return listeners[key]
    }
    else {
      const r = ref<string | null>(sessionStorage.getItem(key))
      listeners[key] = r
      return r
    }
  }
  function get(key: string): string | null {
    const val = listeners[key]
    if (val)
      return val.value
    else return null
  }

  return () => ({ listen, get })
})()
