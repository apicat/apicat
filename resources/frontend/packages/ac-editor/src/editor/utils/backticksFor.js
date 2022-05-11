export default function backticksFor(node, side) {
  const ticks = /`+/g;
  let match = null;
  let len = 0;

  if (node.isText) {
    while ((match = ticks.exec(node.text))) {
      len = Math.max(len, match[0].length);
    }
  }

  let result = len > 0 && side > 0 ? " `" : "`";
  for (let i = 0; i < len; i++) {
    result += "`";
  }
  if (len > 0 && side < 0) {
    result += " ";
  }
  return result;
}
