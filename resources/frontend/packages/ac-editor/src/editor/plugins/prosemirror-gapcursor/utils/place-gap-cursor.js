import { getBreakoutModeFromTargetNode, isLeftCursor } from "./common";
import { Side } from "../GapCursorSelection";

/**
 * We have a couple of nodes that require us to compute style
 * on different elements, ideally all nodes should be able to
 * compute the appropriate styles based on their wrapper.
 */
const nestedCases = {
  "tableView-content-wrap": "table",
  "mediaSingleView-content-wrap": ".rich-media-item",
};
const computeNestedStyle = (dom) => {
  const foundKey = Object.keys(nestedCases).find((className) =>
    dom.classList.contains(className)
  );
  const nestedSelector = foundKey && nestedCases[foundKey];

  if (nestedSelector) {
    const nestedElement = dom.querySelector(nestedSelector);
    if (nestedElement) {
      return window.getComputedStyle(nestedElement);
    }
  }
};

const measureHeight = (style) => {
  return measureValue(style, [
    "height",
    "padding-top",
    "padding-bottom",
    "border-top-width",
    "border-bottom-width",
  ]);
};

const measureWidth = (style) => {
  return measureValue(style, [
    "width",
    "padding-left",
    "padding-right",
    "border-left-width",
    "border-right-width",
  ]);
};

const measureValue = (style, measureValues) => {
  const [base, ...contentBoxValues] = measureValues;
  const measures = [style.getPropertyValue(base)];

  const boxSizing = style.getPropertyValue("box-sizing");
  if (boxSizing === "content-box") {
    contentBoxValues.forEach((value) => {
      measures.push(style.getPropertyValue(value));
    });
  }

  let result = 0;
  for (let i = 0; i < measures.length; i++) {
    result += parseFloat(measures[i]);
  }
  return result;
};

const mutateElementStyle = (element, style, side) => {
  if (isLeftCursor(side)) {
    element.style.marginLeft = style.getPropertyValue("margin-left");
  } else {
    const marginRight = parseFloat(style.getPropertyValue("margin-right"));
    if (marginRight > 0) {
      element.style.marginLeft = `-${Math.abs(marginRight)}px`;
    } else {
      element.style.paddingRight = `${Math.abs(marginRight)}px`;
    }
  }
};

export const toDOM = (view, getPos) => {
  const selection = view.state.selection;
  const { $from, side } = selection;
  const isRightCursor = side === Side.RIGHT;
  const node = isRightCursor ? $from.nodeBefore : $from.nodeAfter;
  const nodeStart = getPos();
  const dom = view.nodeDOM(nodeStart);

  const element = document.createElement("span");
  element.className = `ProseMirror-gapcursor ${
    isRightCursor ? "-right" : "-left"
  }`;
  element.appendChild(document.createElement("span"));

  if (dom instanceof HTMLElement && element.firstChild) {
    const style = computeNestedStyle(dom) || window.getComputedStyle(dom);

    const gapCursor = element.firstChild;
    gapCursor.style.height = `${measureHeight(style)}px`;

    // TODO remove this table specific piece. need to figure out margin collapsing logic
    // if (nodeStart !== 0 || (node && node.type.name === "table")) {
    //   gapCursor.style.marginTop = style.getPropertyValue("margin-top");
    // }

    const breakoutMode = node && getBreakoutModeFromTargetNode(node);
    if (breakoutMode) {
      gapCursor.setAttribute("layout", breakoutMode);
      gapCursor.style.width = `${measureWidth(style)}px`;
    } else {
      mutateElementStyle(gapCursor, style, selection.side);
    }
  }

  return element;
};
