export type PopoverOptions = {
  clickOutSide: () => void
}

export const usePopover = (options?: PopoverOptions) => {
  const popoverRefEl = ref<Nullable<HTMLElement>>(null)
  const isShow = ref(false)

  const showPopover = (el: HTMLElement) => {
    popoverRefEl.value = el
    isShow.value = true
  }

  const hidePopover = () => {
    popoverRefEl.value = null
    isShow.value = false
    options?.clickOutSide()
  }

  onClickOutside(
    popoverRefEl,
    () => {
      popoverRefEl.value = null
      isShow.value = false
      options?.clickOutSide()
    },
    {
      ignore: ['.ac-popper-menu'],
    }
  )

  return {
    popoverRefEl,
    isShow,
    showPopover,
    hidePopover,
  }
}
