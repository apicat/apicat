import { GapCursorSelection } from "./GapCursorSelection";
import { NodeSelection } from "prosemirror-state";

export function atTheEndOfDoc(state) {
  const { selection, doc } = state;
  return doc.nodeSize - selection.$to.pos - 2 === selection.$to.depth;
}

export function atTheBeginningOfDoc(state) {
  const { selection } = state;
  return selection.$from.pos === selection.$from.depth;
}

export function atTheEndOfBlock(state) {
  const { selection } = state;
  const { $to } = selection;
  if (selection instanceof GapCursorSelection) {
    return false;
  }
  if (selection instanceof NodeSelection && selection.node.isBlock) {
    return true;
  }
  return endPositionOfParent($to) === $to.pos + 1;
}

export function atTheBeginningOfBlock(state) {
  const { selection } = state;
  const { $from } = selection;
  if (selection instanceof GapCursorSelection) {
    return false;
  }
  if (selection instanceof NodeSelection && selection.node.isBlock) {
    return true;
  }
  return startPositionOfParent($from) === $from.pos;
}

export function startPositionOfParent(resolvedPos) {
  return resolvedPos.start(resolvedPos.depth);
}

export function endPositionOfParent(resolvedPos) {
  return resolvedPos.end(resolvedPos.depth) + 1;
}
