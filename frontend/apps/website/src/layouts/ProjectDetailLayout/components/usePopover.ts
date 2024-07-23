export function usePopover(options?: { ignore?: string[], onHide?: () => void }) {
  const { ignore = [], onHide } = options || {}

  const popoverRefEl = ref<HTMLElement | null>(null)
  const isShowPopoverMenu = ref(false)

  function show(el: HTMLElement) {
    popoverRefEl.value = el
    isShowPopoverMenu.value = true
  }

  function hide() {
    popoverRefEl.value = null
    isShowPopoverMenu.value = false
  }

  onClickOutside(popoverRefEl, () => {
    hide()
    onHide?.()
  }, {
    ignore,
  })

  return {
    popoverRefEl,
    isShowPopoverMenu,

    show,
    hide,
  }
}
