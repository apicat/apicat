const defaultKey = Symbol('asyncInitKey')

interface Ctx {
  tasks: Promise<void>[]
  addTask(a: (() => Promise<void>) | Promise<void>): void
  done(): Promise<void>
}

export function provideAsyncInitTask(key?: string) {
  const ctx: Ctx = {
    tasks: [],
    addTask(t) {
      if (typeof t === 'function')
        this.tasks.push(t())
      else
        this.tasks.push(t)
    },
    async done() {
      let len = 0
      while (len !== this.tasks.length) {
        len = this.tasks.length
        await Promise.all(this.tasks)
      }
    },
  }
  provide(key ?? defaultKey, ctx)
  return ctx
}

export function injectAsyncInitTask(key?: string) {
  return inject<Ctx>(key ?? defaultKey)
}
