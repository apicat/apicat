import { Side } from "../GapCursorSelection";

export const isLeftCursor = (side) => side === Side.LEFT;

export function getMediaNearPos(doc, $pos, schema, dir = -1) {
  let $currentPos = $pos;
  let currentNode = null;
  const { mediaSingle, media, mediaGroup } = schema.nodes;

  do {
    $currentPos = doc.resolve(
      dir === -1 ? $currentPos.before() : $currentPos.after()
    );

    if (!$currentPos) {
      return null;
    }

    currentNode =
      (dir === -1 ? $currentPos.nodeBefore : $currentPos.nodeAfter) ||
      $currentPos.parent;

    if (!currentNode || currentNode.type === schema.nodes.doc) {
      return null;
    }

    if (
      currentNode.type === mediaSingle ||
      currentNode.type === media ||
      currentNode.type === mediaGroup
    ) {
      return currentNode;
    }
  } while ($currentPos.depth > 0);

  return null;
}

export const isTextBlockNearPos = (doc, schema, $pos, dir) => {
  let $currentPos = $pos;
  let currentNode = dir === -1 ? $currentPos.nodeBefore : $currentPos.nodeAfter;

  // If next node is a text or a text block bail out early.
  if (currentNode && (currentNode.isTextblock || currentNode.isText)) {
    return true;
  }

  while ($currentPos.depth > 0) {
    $currentPos = doc.resolve(
      dir === -1 ? $currentPos.before() : $currentPos.after()
    );

    if (!$currentPos) {
      return false;
    }

    currentNode =
      (dir === -1 ? $currentPos.nodeBefore : $currentPos.nodeAfter) ||
      $currentPos.parent;

    if (!currentNode || currentNode.type === schema.nodes.doc) {
      return false;
    }

    if (currentNode.isTextblock) {
      return true;
    }
  }

  let childNode = currentNode;

  while (childNode && childNode.firstChild) {
    childNode = childNode.firstChild;
    if (childNode && (childNode.isTextblock || childNode.isText)) {
      return true;
    }
  }

  return false;
};

export function getBreakoutModeFromTargetNode(node) {
  let layout;
  if (node.attrs.layout) {
    layout = node.attrs.layout;
  }

  if (node.marks && node.marks.length) {
    layout = (
      node.marks.find((mark) => mark.type.name === "breakout") || {
        attrs: { mode: "" },
      }
    ).attrs.mode;
  }

  if (["wide", "full-width"].indexOf(layout) === -1) {
    return "";
  }

  return layout;
}

//
// export const isIgnoredClick = (elem) => {
//     if (elem.nodeName === 'BUTTON' || elem.closest('button')) {
//         return true;
//     }
//
//     // check if we're clicking an image caption placeholder
//     if (elem.closest(`[data-id="${CAPTION_PLACEHOLDER_ID}"]`)) {
//         return true;
//     }
//
//     // check if target node has a parent table node
//     let tableWrap;
//     let node = elem;
//     while (node) {
//         if (
//             node.className &&
//             (node.getAttribute('class') || '').indexOf(
//                 TableCssClassName.TABLE_CONTAINER,
//             ) > -1
//         ) {
//             tableWrap = node;
//             break;
//         }
//         node = node.parentNode;
//     }
//
//     if (tableWrap) {
//         const rowControls = tableWrap.querySelector(
//             `.${TableCssClassName.ROW_CONTROLS_WRAPPER}`,
//         );
//         const isColumnControlsDecoration =
//             elem &&
//             elem.classList &&
//             elem.classList.contains(TableCssClassName.COLUMN_CONTROLS_DECORATIONS);
//         return (
//             (rowControls && rowControls.contains(elem)) || isColumnControlsDecoration
//         );
//     }
//
//     // Check if unsupported node selection
//     // (without this, selection requires double clicking in FF due to posAtCoords differences)
//     if (elem.closest(`.${UnsupportedSharedCssClassName.BLOCK_CONTAINER}`)) {
//         return true;
//     }
//
//     return false;
// };
