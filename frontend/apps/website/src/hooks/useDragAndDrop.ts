import { noop } from '@vueuse/core'

export type DragDropOptions = {
  dragClass?: string
  direction?: 'vertical' | 'horizontal'
  onDragStartHandle?: (e?: DragEvent) => { dragElement: HTMLElement } | void
  onDrop?: (dragIndex: number, dropIndex: number) => void
}

export const useDragAndDrop = (options?: DragDropOptions) => {
  const { dragClass = 'dragging', direction = 'vertical', onDrop = noop, onDragStartHandle = noop } = options || {}

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

  const updateDropIndicatorPosition = (dragEle: HTMLElement) => {
    if (!dropIndicator || !dragEle) {
      return
    }

    const { top, left } = dragEle.getBoundingClientRect()
    dropIndicator.style.display = 'block'
    dropIndicator.style.top = top + 'px'
    dropIndicator.style.left = `${left - 10}px`
    dropIndicator.style[direction === 'vertical' ? 'height' : 'width'] = `${dragEle[direction === 'vertical' ? 'offsetHeight' : 'offsetWidth']}px`
  }

  const onDragStart = (e: DragEvent, index: string | number) => {
    e.dataTransfer!.dropEffect = 'move'
    const nodeEle = e.currentTarget as HTMLElement
    nodeEle && nodeEle.classList.add(dragClass)
    const data = onDragStartHandle(e)
    e.dataTransfer?.setDragImage(data?.dragElement ?? nodeEle, 0, 0)
    e.dataTransfer?.setData(DRAG_KEY, String(index))
    nodeEle.style.opacity = '0.5'
    createDropIndicator(nodeEle)
  }

  const onDragEnd = (e: DragEvent) => {
    const nodeEle = e.currentTarget as HTMLElement
    nodeEle.classList.remove(dragClass)
    nodeEle.style.opacity = ''
    dropIndicator && dropIndicator.remove()
  }

  const onDragLeave = (e: DragEvent, index: string | number) => {
    dropIndicator && (dropIndicator.style.display = 'none')
  }

  const onDragOver = (e: DragEvent, index: string | number) => {
    e.preventDefault()
    if (dropTest(e)) {
      updateDropIndicatorPosition(e.currentTarget as HTMLElement)
    }
  }

  function dropTest(ev: DragEvent) {
    ev.dataTransfer!.dropEffect = 'move'
    const dom = ev.currentTarget as HTMLElement
    if (ev[direction === 'vertical' ? 'offsetX' : 'offsetY'] < dom[direction === 'vertical' ? 'clientWidth' : 'clientHeight'] / 2) {
      return 1
    }

    return 0
  }

  const onDropHandler = (e: DragEvent, index: number) => {
    if (e.dataTransfer?.getData(DRAG_KEY)) {
      onDrop && onDrop(parseInt(e.dataTransfer?.getData(DRAG_KEY), 10), index)
    }
  }

  return {
    onDragStart,
    onDragLeave,
    onDragOver,
    onDragEnd,
    onDropHandler,
  }
}
