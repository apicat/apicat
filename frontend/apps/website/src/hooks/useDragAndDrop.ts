export type DragDropOptions = {
  dragClass?: string
  direction: 'vertical' | 'horizontal'
}

export const useDragAndDrop = (options?: DragDropOptions) => {
  const { dragClass = 'dragging', direction = 'vertical' } = options || {}

  const DRAG_KEY = `DRAG_AND_DROP_${Date.now()}`

  let dropIndicator: HTMLElement | null = null

  const createDropIndicator = (dragEle: HTMLElement) => {
    const el = document.getElementById(DRAG_KEY)
    if ((el && el.dataset.direction === direction) || !dragEle) {
      return
    }

    const { top, left } = dragEle.getBoundingClientRect()

    dropIndicator = dropIndicator || document.createElement('div')
    dropIndicator.id = DRAG_KEY
    dropIndicator.dataset.direction = direction
    dropIndicator.style.position = 'absolute'
    dropIndicator.style.top = top + 'px'
    dropIndicator.style.left = `${left - 10}px`
    dropIndicator.style[direction === 'vertical' ? 'width' : 'height'] = '1px'
    dropIndicator.style[direction === 'vertical' ? 'height' : 'width'] = `${dragEle[direction === 'vertical' ? 'offsetHeight' : 'offsetWidth']}px`
    dropIndicator.style.border = '1px dashed var(--el-color-primary)'

    document.body.append(dropIndicator)
  }

  const onDragStart = (e: DragEvent, index: string | number) => {
    e.dataTransfer!.dropEffect = 'move'
    const nodeEle = e.currentTarget as HTMLElement
    nodeEle && nodeEle.classList.add(dragClass)
    e.dataTransfer?.setDragImage(nodeEle, 0, 0)
    e.dataTransfer?.setData(DRAG_KEY, String(index))

    console.log(DRAG_KEY, e.dataTransfer?.getData(DRAG_KEY))

    nodeEle.style.opacity = '0.5'
    createDropIndicator(nodeEle)
  }

  const onDragEnd = (e: DragEvent) => {
    const nodeEle = e.currentTarget as HTMLElement
    nodeEle.classList.remove(dragClass)
    nodeEle.style.opacity = ''
    dropIndicator && dropIndicator.remove()
  }

  const onDragLeave = (e: DragEvent, index: string | number) => {}

  const onDragOver = (e: DragEvent, index: string | number) => {
    e.preventDefault()
    const flag = dropTest(e)
  }

  function dropTest(ev: DragEvent) {
    ev.dataTransfer!.dropEffect = 'move'
    const dom = ev.currentTarget as HTMLElement
    if (ev[direction === 'vertical' ? 'offsetX' : 'offsetY'] < dom[direction === 'vertical' ? 'clientWidth' : 'clientHeight'] / 2) {
      return -1
    }
    return 1
  }

  const onDropHandler = (e: DragEvent, index: string | number) => {}

  return {
    onDragStart,
    onDragLeave,
    onDragOver,
    onDragEnd,
    onDropHandler,
  }
}
