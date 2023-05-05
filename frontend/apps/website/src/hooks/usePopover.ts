import { noop } from '@vueuse/core'

export type PopoverOptions = {
  ignore?: []
  onHide?: () => void
}

export const usePopover = (options?: PopoverOptions) => {
  let { onHide = noop, ignore = [] } = options || {}
  const popoverRefEl = ref<Nullable<HTMLElement>>(null)
  const isShow = ref(false)

  const showPopover = (el: HTMLElement) => {
    popoverRefEl.value = el
    isShow.value = true
  }

  const hidePopover = () => {
    isShow.value = false
    setTimeout(() => {
      popoverRefEl.value = null
    }, 100)

    onHide()
  }

  onClickOutside(popoverRefEl, () => hidePopover(), {
    ignore: ['.ac-popper-menu', '.ignore-popper'].concat(ignore),
  })

  return {
    popoverRefEl,
    isShow,
    showPopover,
    hidePopover,
  }
}
