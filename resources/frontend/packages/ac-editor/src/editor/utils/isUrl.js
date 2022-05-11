export default function isUrl(text) {
  if (text.match(/\n/)) {
    return false;
  }

  try {
    new URL(text);
    return true;
  } catch (err) {
    return false;
  }
}
