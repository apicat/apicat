export function useTextOverflowTooltip(textVNodeClass: string) {
  const tooltipData = {
    tooltipVisible: ref(false),
    tooltipRef: undefined as any,
    tooltipContent: '',
  }
  let currentTask: Promise<any> | undefined

  const getPadding = (el: HTMLElement) => {
    const style = window.getComputedStyle(el, null)
    const paddingLeft = Number.parseInt(style.paddingLeft, 10) || 0
    const paddingRight = Number.parseInt(style.paddingRight, 10) || 0
    const paddingTop = Number.parseInt(style.paddingTop, 10) || 0
    const paddingBottom = Number.parseInt(style.paddingBottom, 10) || 0
    return {
      left: paddingLeft,
      right: paddingRight,
      top: paddingTop,
      bottom: paddingBottom,
    }
  }
  const handleCellMouseEnter = (() => {
    let timer: any
    return async (event: MouseEvent) => {
      while (currentTask) {
        await currentTask
        currentTask = undefined
        await nextTick()
      }
      if (timer) clearTimeout(timer)
      timer = setTimeout(() => {
        currentTask = (async (event: MouseEvent) => {
          const cellChild = (event.target as HTMLElement).querySelector(textVNodeClass) as HTMLElement
          const range = document.createRange()
          range.setStart(cellChild, 0)
          range.setEnd(cellChild, cellChild.childNodes.length)
          let rangeWidth = range.getBoundingClientRect().width
          let rangeHeight = range.getBoundingClientRect().height
          const offsetWidth = rangeWidth - Math.floor(rangeWidth)
          if (offsetWidth < 0.001) {
            rangeWidth = Math.floor(rangeWidth)
          }
          const offsetHeight = rangeHeight - Math.floor(rangeHeight)
          if (offsetHeight < 0.001) {
            rangeHeight = Math.floor(rangeHeight)
          }

          const { top, left, right, bottom } = getPadding(cellChild)
          const horizontalPadding = left + right
          const verticalPadding = top + bottom
          if (
            rangeWidth + horizontalPadding > cellChild.offsetWidth ||
            rangeHeight + verticalPadding > cellChild.offsetHeight ||
            cellChild.scrollWidth > cellChild.offsetWidth
          ) {
            tooltipData.tooltipRef = event.target!
            tooltipData.tooltipContent = cellChild.innerText || cellChild.textContent || ''
            tooltipData.tooltipVisible.value = true
          }
        })(event)
      }, 0)
    }
  })()

  return {
    handleCellMouseEnter,
    handleCellMouseLeave() {
      tooltipData.tooltipVisible.value = false
      tooltipData.tooltipRef = undefined
      tooltipData.tooltipContent = ''
    },
    tooltipData,
  }
}
