// we don't show gap cursor for those nodes
const IGNORED_NODES = [
  "paragraph",
  "bulletList",
  "orderedList",
  "listItem",
  "taskItem",
  "decisionItem",
  "heading",
  "blockquote",
  "layoutColumn",
  "caption",
  "media",
];

export const isIgnored = (node) => {
  return !!node && IGNORED_NODES.indexOf(node.type.name) !== -1;
};
