import { noop } from '@vueuse/core'

export type DragDropOptions = {
  dragClass?: string
  direction?: 'vertical' | 'horizontal'
  onDragStartHandle?: (e?: DragEvent) => { dragElement: HTMLElement } | void
  onDrop?: (dragIndex: number, dropIndex: number, flag: number) => void
}

export const useDragAndDrop = (options?: DragDropOptions) => {
  const { dragClass = 'dragging', direction = 'vertical', onDrop = noop, onDragStartHandle = noop } = options || {}

  const DRAG_KEY = `DRAG_AND_DROP_${Date.now()}`

  const onDragStart = (e: DragEvent, index: string | number) => {
    e.dataTransfer!.dropEffect = 'move'
    const nodeEle = e.currentTarget as HTMLElement
    nodeEle && nodeEle.classList.add(dragClass)
    const data = onDragStartHandle(e)
    e.dataTransfer?.setDragImage(data?.dragElement ?? nodeEle, 0, 0)
    e.dataTransfer?.setData(DRAG_KEY, String(index))
    nodeEle.style.opacity = '0.5'
  }

  const onDragEnd = (e: DragEvent) => {
    const nodeEle = e.currentTarget as HTMLElement
    nodeEle.classList.remove(dragClass)
    nodeEle.style.opacity = ''
  }

  const onDragLeave = (e: DragEvent, index: string | number) => {
    const dom = (e.currentTarget as HTMLElement).parentNode as HTMLElement
    dom.style.borderLeft = ''
    dom.style.borderRight = ''
  }

  const onDragOver = (e: DragEvent, index: string | number) => {
    e.preventDefault()
    let direction = dropTest(e, index)
    if (direction === 0) {
      return
    }
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'move'
      const dom = (e.currentTarget as HTMLElement).parentNode as HTMLElement
      switch (direction) {
        case -1:
          dom.style.borderLeft = ''
          dom.style.borderRight = '1px var(--primary-color) dashed'
          break
        case 1:
          dom.style.borderLeft = '1px var(--primary-color) dashed'
          dom.style.borderRight = ''
          break
      }
    }
  }

  function dropTest(e: DragEvent, index: string | number) {
    const dragData = parseInt(e.dataTransfer?.getData(DRAG_KEY) as any, 10)
    if (dragData === index) {
      return 0
    }

    e.dataTransfer!.dropEffect = 'move'
    const dom = e.currentTarget as HTMLElement
    if (e[direction === 'vertical' ? 'offsetX' : 'offsetY'] < dom[direction === 'vertical' ? 'clientWidth' : 'clientHeight'] / 2) {
      return 1
    }
    return -1
  }

  const onDropHandler = (e: DragEvent, index: number) => {
    onDragLeave(e, index)
    const flag = dropTest(e, index)
    if (flag === 0) {
      return
    }

    if (e.dataTransfer?.getData(DRAG_KEY)) {
      onDrop && onDrop(parseInt(e.dataTransfer?.getData(DRAG_KEY), 10), index, flag)
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
