import some from "lodash/some";
import { getColumnIndex, getRowIndex } from "@/editor/utils/index";

export default function isVisible(vm) {
  const { view } = vm.editor;
  const { selection } = view.state;
  if (!selection) return false;
  if (selection.empty) return false;
  if (selection.node && selection.node.type.name === "image") {
    return true;
  }

  const colIndex = getColumnIndex(view.state.selection);
  const rowIndex = getRowIndex(view.state.selection);

  const isTableSelection = colIndex !== undefined && rowIndex !== undefined;
  if (isTableSelection) {
    return false;
  }

  if (selection.node) return false;

  const slice = selection.content();
  const fragment = slice.content;
  const nodes = fragment.content;

  return some(nodes, (n) => n.content.size);
}
