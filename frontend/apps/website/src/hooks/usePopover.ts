import { noop } from '@vueuse/core'

export type PopoverOptions = {
  ignore?: Array<any>
  onHide?: () => void
}

export const usePopover = (options?: PopoverOptions) => {
  let { onHide = noop, ignore = [] } = options || {}
  const popoverRefEl = ref<Nullable<HTMLElement>>(null)
  const isShow = ref(false)
  ignore = ['.ac-popper-menu', '.ignore-popper'].concat(ignore)

  const showPopover = (el: HTMLElement) => {
    popoverRefEl.value = el
    isShow.value = true
  }

  const hidePopover = () => {
    isShow.value = false
    popoverRefEl.value = null
    onHide()
  }

  const shouldIgnore = (event: PointerEvent) =>
    ignore.some((target2) => Array.from(window.document.querySelectorAll(target2)).some((el) => el === event.target || event.composedPath().includes(el)))

  const stop: any = onClickOutside(popoverRefEl, (e) => {
    if (shouldIgnore(e)) {
      return
    }
    hidePopover()
  })

  onUnmounted(() => stop())

  return {
    popoverRefEl,
    isShow,
    showPopover,
    hidePopover,
  }
}
