import { textblockTypeInputRule } from 'prosemirror-inputrules'
import { backspaceToParagraph } from '../commands'
import toggleBlockType from '../commands/toggleBlockType'
import Node from '../lib/Node'

export default class Heading extends Node {
    className = 'heading-name'

    get name() {
        return 'heading'
    }

    get defaultOptions() {
        return {
            levels: [1, 2, 3, 4, 5, 6],
        }
    }

    get schema() {
        return {
            attrs: {
                level: {
                    default: 1,
                },
            },
            content: 'inline*',
            group: 'block',
            defining: true,
            draggable: false,
            parseDOM: this.options.levels.map((level) => ({
                tag: `h${level}`,
                attrs: { level },
                // contentElement: "span",
            })),
            toDOM: (node) => {
                // const button = document.createElement("button");
                // button.innerText = "#";
                // button.type = "button";
                // button.className = "heading-anchor";
                return [`h${node.attrs.level}`, { class: 'handle-middle' }, 0]
            },
        }
    }

    toMarkdown(state, node) {
        state.write(state.repeat('#', node.attrs.level) + ' ')
        state.renderInline(node)
        state.closeBlock(node)
    }

    parseMarkdown() {
        return {
            block: 'heading',
            getAttrs: (token) => ({
                level: +token.tag.slice(1),
            }),
        }
    }

    commands({ type, schema }) {
        return (attrs) => {
            return toggleBlockType(type, schema.nodes.paragraph, attrs)
        }
    }

    keys({ type }) {
        // const options = this.options.levels.reduce(
        //   (items, level) => ({
        //     ...items,
        //     ...{
        //       [`Shift-Ctrl-${level}`]: setBlockType(type, { level }),
        //     },
        //   }),
        //   {}
        // );

        return {
            // ...options,
            Backspace: backspaceToParagraph(type),
        }
    }

    // get plugins() {
    //   const getAnchors = (doc) => {
    //     const decorations = [];
    //     const previouslySeen = {};
    //
    //     doc.descendants((node, pos) => {
    //       if (node.type.name !== this.name) return;
    //
    //       // calculate the optimal id
    //       const slug = headingToSlug(node);
    //       let id = slug;
    //
    //       // check if we've already used it, and if so how many times?
    //       // Make the new id based on that number ensuring that we have
    //       // unique ID's even when headings are identical
    //       if (previouslySeen[slug] > 0) {
    //         id = headingToSlug(node, previouslySeen[slug]);
    //       }
    //
    //       // record that we've seen this slug for the next loop
    //       previouslySeen[slug] =
    //         previouslySeen[slug] !== undefined ? previouslySeen[slug] + 1 : 1;
    //
    //       decorations.push(
    //         Decoration.widget(
    //           pos,
    //           () => {
    //             const anchor = document.createElement("a");
    //             anchor.id = id;
    //             anchor.className = this.className;
    //             return anchor;
    //           },
    //           {
    //             side: -1,
    //             key: id,
    //           }
    //         )
    //       );
    //     });
    //
    //     return DecorationSet.create(doc, decorations);
    //   };
    //
    //   const plugin = new Plugin({
    //     state: {
    //       init: (config, state) => {
    //         return getAnchors(state.doc);
    //       },
    //       apply: (tr, oldState) => {
    //         return tr.docChanged ? getAnchors(tr.doc) : oldState;
    //       },
    //     },
    //     props: {
    //       decorations: (state) => plugin.getState(state),
    //     },
    //   });
    //
    //   return [plugin];
    // }

    inputRules({ type }) {
        return this.options.levels.map((level) =>
            textblockTypeInputRule(new RegExp(`^(#{1,${level}})\\s$`), type, () => ({
                level,
            }))
        )
    }
}
