import { isIgnored } from "./is-ignored";

export const isValidTargetNode = (node) => {
  return !!node && !isIgnored(node);
};
