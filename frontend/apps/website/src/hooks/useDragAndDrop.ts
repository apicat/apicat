export const useDragAndDrop = ({ dragClass = 'draging' }) => {
  const DRAG_KEY = `DRAG_AND_DROP_${Date.now()}`

  const onDragStart = (e: DragEvent, data: string) => {
    e.dataTransfer!.dropEffect = 'move'
    console.log(e)
    const nodeEle = e.target as Element
    nodeEle && nodeEle.classList.add(dragClass)
    e.dataTransfer?.setDragImage(nodeEle as HTMLElement, 0, 0)
    e.dataTransfer?.setData(DRAG_KEY, data)
  }

  const onDragEnd = (e: DragEvent) => {
    const nodeEle = e.target as Element
    nodeEle && nodeEle.classList.remove(dragClass)
  }

  const onDropHandler = (e: DragEvent) => {}

  return {
    onDragStart,
    onDragEnd,
    onDropHandler,
  }
}
