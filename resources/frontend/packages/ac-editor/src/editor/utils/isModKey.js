export const isMac = window.navigator.platform === "MacIntel";

export default function isModKey(event) {
  return isMac ? event.metaKey : event.ctrlKey;
}
