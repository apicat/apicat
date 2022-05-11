import { Plugin } from "prosemirror-state";
import { Decoration, DecorationSet } from "prosemirror-view";
import Extension from "../lib/Extension";
import { findParentNode } from "prosemirror-utils";

export default class Placeholder extends Extension {
  get name() {
    return "empty-placeholder";
  }

  get defaultOptions() {
    return {
      emptyNodeClass: "placeholder",
      placeholder: "写点什么吧",
    };
  }

  get plugins() {
    return [
      new Plugin({
        props: {
          decorations: (state) => {
            const parent = findParentNode((node) => node.type.name === "paragraph")(state.selection);

            const isTopNode = parent.depth === 1;
            const isEmpty = parent.node.content.size === 0;

            if (isTopNode && isEmpty) {
              const decoration = Decoration.node(
                parent.pos,
                parent.pos + parent.node.nodeSize,
                {
                  class: this.options.emptyNodeClass,
                  "data-empty-text": this.options.placeholder,
                }
              );
              return DecorationSet.create(state.doc, [decoration]);
            }
            return null;
          },
        },
      }),
    ];
  }
}
