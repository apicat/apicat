import { type Ref, ref } from 'vue'

interface Options {
  onlyOne?: boolean
  atleastOne?: boolean
  defaults?: string[]
}

interface Emits {
  onChange?(key: string, openState: boolean): void
}

interface CollapseItemCtx {
  open(): void
  close(emit?: boolean): void
}

type OpendItems = Set<string>
export interface UseCollapse {
  id: number
  opendItems: Ref<OpendItems>
  ctx: {
    current: Ref<OpendItems>
    open(key: string): void
    close(key: string): void
    register(name: string, itemCtx: CollapseItemCtx): CollapseItemCtx
  }
}

let ctxID = 0
export function useCollapse(props: Options, cbs?: Emits): UseCollapse {
  const opendItems = ref<Set<string>>(new Set())
  const cards: Record<string, CollapseItemCtx> = {}
  const defaults = props.defaults ?? []

  return {
    id: ctxID++,
    opendItems,
    ctx: {
      current: opendItems,
      open(key: string) {
        const c = cards[key]
        if (!c) defaults.push(key)
        else c.open()
      },
      close(key: string) {
        const c = cards[key]
        if (c) c.close()
      },
      register(name: string, itemCtx: CollapseItemCtx): CollapseItemCtx {
        const open = itemCtx.open
        const close = itemCtx.close
        itemCtx.open = () => {
          open()
          opendItems.value.add(name)
          if (props.onlyOne) {
            Object.keys(cards).forEach((n) => {
              const itemCtx = cards[n]
              if (n !== name) itemCtx.close(false)
            })
          }
          cbs?.onChange?.(name, true)
        }
        itemCtx.close = (emit = true) => {
          if (opendItems.value.size === 1 && props.atleastOne) return
          else {
            close()
            opendItems.value.delete(name)
            if (emit) cbs?.onChange?.(name, false)
          }
        }
        cards[name] = itemCtx

        if (props.defaults && props.defaults.indexOf(name) !== -1) itemCtx.open()
        return itemCtx
      },
    },
  }
}
