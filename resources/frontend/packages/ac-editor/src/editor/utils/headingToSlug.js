import escape from "lodash/escape";
import slugify from "slugify";

// Slugify, escape, and remove periods from headings so that they are
// compatible with both url hashes AND dom ID's (querySelector does not like
// ID's that begin with a number or a period, for example).
function safeSlugify(text) {
  return `h-${escape(
    slugify(text, {
      // eslint-disable-next-line no-useless-escape
      remove: /[!"#$%&'\.()*+,\/:;<=>?@\[\]\\^_`{|}~]/g,
      lower: true,
    })
  )}`;
}

// calculates a unique slug for this heading based on it's text and position
// in the document that is as stable as possible
export default function headingToSlug(node, index = 0) {
  const slugified = safeSlugify(node.textContent);
  if (index === 0) return slugified;
  return `${slugified}-${index}`;
}
