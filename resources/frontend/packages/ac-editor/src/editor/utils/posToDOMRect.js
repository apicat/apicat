export default function posToDOMRect(view, selection, from, to) {
  const start = view.coordsAtPos(from);
  const end = view.coordsAtPos(to, -1);

  // ensure that start < end for the menu to be positioned correctly
  const selectionBounds = {
    top: Math.min(start.top, end.top),
    bottom: Math.max(start.bottom, end.bottom),
    left: Math.min(start.left, end.left),
    right: Math.max(start.right, end.right),
  };

  // tables are an oddity, and need their own positioning logic
  const isColSelection = selection.isColSelection && selection.isColSelection();
  const isRowSelection = selection.isRowSelection && selection.isRowSelection();
  const isImageSelection = selection.node && selection.node.type.name === "image";

  if (isColSelection) {
    const { node: element } = view.domAtPos(selection.from);
    const { width } = element.getBoundingClientRect();
    selectionBounds.top -= 20;
    selectionBounds.right = selectionBounds.left + width;
  }

  if (isRowSelection) {
    selectionBounds.right = selectionBounds.left = selectionBounds.left - 18;
  }

  if(isImageSelection){
    const element = view.nodeDOM(selection.from);
    const imageElement = element && element.querySelector('img')
    if(imageElement){
      const { left, top, bottom,right } = imageElement.getBoundingClientRect();
      selectionBounds.right = right
      selectionBounds.left = left
      selectionBounds.top = top
      selectionBounds.bottom = bottom
    }
  }
  const { right, left, bottom, top } = selectionBounds;

  selectionBounds.width = right - left;
  selectionBounds.height = bottom - top;
  selectionBounds.x = left;
  selectionBounds.y = top;

  // selectionBounds.left = Math.round(left + window.scrollX);
  // selectionBounds.top = Math.round(top + window.scrollY);

  return {
    ...selectionBounds,
    toJSON: () => selectionBounds,
  };
}
