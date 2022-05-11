import refractor from "refractor/core";
import flattenDeep from "lodash/flattenDeep";
import { Plugin, PluginKey } from "prosemirror-state";
import { Decoration, DecorationSet } from "prosemirror-view";
import { findBlockNodes } from "prosemirror-utils";

export const LANGUAGES = {
  none: "None", // additional entry to disable highlighting
  bash: "Bash",
  css: "CSS",
  clike: "C",
  csharp: "C#",
  go: "Go",
  markup: "HTML",
  java: "Java",
  javascript: "JavaScript",
  json: "JSON",
  php: "PHP",
  powershell: "Powershell",
  python: "Python",
  ruby: "Ruby",
  sql: "SQL",
  typescript: "TypeScript",
};

function parseNodes(nodes = [], classNames = []) {
  return nodes.map((node) => {
    if (node.type === "element") {
      const classes = [...classNames, ...(node.properties.className || [])];
      return parseNodes(node.children, classes);
    }
    return {
      text: node.value,
      classes: classNames,
    };
  });
}

function getDecorations({ doc, name }) {
  const decorations = [];
  const blocks = findBlockNodes(doc).filter(
    (item) => item.node.type.name === name
  );

  blocks.forEach((block) => {
    let startPos = block.pos + 1;
    const language = block.node.attrs.language;
    if (!language || language === "none" || !refractor.registered(language)) {
      return;
    }

    const nodes = refractor.highlight(block.node.textContent, language);

    flattenDeep(parseNodes(nodes))
      .map((node) => {
        const from = startPos;
        const to = from + node.text.length;

        startPos = to;

        return {
          ...node,
          from,
          to,
        };
      })
      .forEach((node) => {
        const decoration = Decoration.inline(node.from, node.to, {
          class: (node.classes || []).join(" "),
        });
        decorations.push(decoration);
      });
  });

  return DecorationSet.create(doc, decorations);
}

export default function Prism({ name }) {
  return new Plugin({
    key: new PluginKey("prism"),
    state: {
      init: (_, { doc }) => {
        return getDecorations({ doc, name });
      },
      apply: (transaction, decorationSet, oldState, newState) => {
        const oldNodeName = oldState.selection.$head.parent.type.name;
        const newNodeName = newState.selection.$head.parent.type.name;
        const oldNodes = findBlockNodes(
          oldState.doc,
          (node) => node.type.name === name
        );
        const newNodes = findBlockNodes(
          newState.doc,
          (node) => node.type.name === name
        );

        if (
          transaction.docChanged &&
          // Apply decorations if:
          // selection includes named node,
          ([oldNodeName, newNodeName].includes(name) ||
            // OR transaction adds/removes named node,
            newNodes.length !== oldNodes.length ||
            // OR transaction has changes that completely encapsulte a node
            // (for example, a transaction that affects the entire document).
            // Such transactions can happen during collab syncing via y-prosemirror, for example.
            transaction.steps.some((step) => {
              return (
                step.from !== undefined &&
                step.to !== undefined &&
                oldNodes.some((node) => {
                  return (
                    node.pos >= step.from &&
                    node.pos + node.node.nodeSize <= step.to
                  );
                })
              );
            }))
        ) {
          return getDecorations({ doc: transaction.doc, name });
        }

        return decorationSet.map(transaction.mapping, transaction.doc);
      },
    },
    props: {
      decorations(state) {
        return this.getState(state);
      },
    },
  });
}
