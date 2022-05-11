export { default as getColumnIndex } from "./getColumnIndex";
export { default as getDataTransferFiles } from "./getDataTransferFiles";
export { default as getMarkAttrs } from "./getMarkAttrs";
export { default as getNodeAttrs } from "./getNodeAttrs";
export { default as getMarkRange } from "./getMarkRange";
export { default as getParentListItem } from "./getParentListItem";
export { default as getRowIndex } from "./getRowIndex";
export { default as markInputRule } from "./markInputRule";
export { default as markIsActive } from "./markIsActive";
export { default as nodeEqualsType } from "./nodeEqualsType";
export { default as nodeIsActive } from "./nodeIsActive";
export { default as renderToHtml } from "./renderToHtml";
export { default as findImagePlaceholder } from "./findImagePlaceholder";
export { default as clamp } from "./clamp";
export { default as loadImage } from "./loadImage";

export { default as isInCode } from "./isInCode";
export { default as isInList } from "./isInList";
export { default as isList } from "./isList";
export { default as isMarkActive } from "./isMarkActive";
export { default as isModKey, isMac } from "./isModKey";
export { default as isNodeActive } from "./isNodeActive";
export { default as isUrl } from "./isUrl";
export { default as backticksFor } from "./backticksFor";
export { default as isTextSelection } from "./isTextSelection";


export function cellWrapping($pos) {
  for (let d = $pos.depth; d > 0; d--) { // Sometimes the cell can be in the same depth.
    const role = $pos.node(d).type.spec.tableRole;
    if (role === "cell" || role === 'header_cell') return $pos.node(d)
  }
  return null
}