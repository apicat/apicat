import { toggleMark } from "prosemirror-commands";
import Extension from "./Extension";

export default class Node extends Extension {
  get type() {
    return "mark";
  }

  get schema() {
    return {};
  }

  commands({ type }) {
    return () => toggleMark(type);
  }
}
