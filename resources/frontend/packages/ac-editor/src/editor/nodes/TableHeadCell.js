import Node from '../lib/Node'
import { getCellAttrs, setCellAttrs } from '../utils/tableCellAttr'

// import { DecorationSet, Decoration } from "prosemirror-view";
// import { Plugin } from "prosemirror-state";
// import {
//   isColumnSelected,
//   getCellsInRow,
//   // selectColumn,
// } from "prosemirror-utils";
// import {selectColumn} from "@/editor/plugins/table/select-nodes";
// // import {selectCol} from "@/editor/plugins/table/commands";

export default class TableHeadCell extends Node {
    get name() {
        return 'th'
    }

    get schema() {
        return {
            // draggable: true,
            // selectable: true,

            attrs: {
                colspan: { default: 1 },
                rowspan: { default: 1 },
                alignment: { default: null },
            },

            content: 'paragraph+',
            tableRole: 'header_cell',
            isolating: true,
            parseDOM: [{ tag: 'th', getAttrs: (dom) => getCellAttrs(dom) }],
            toDOM(node) {
                return ['th', setCellAttrs(node), 0]
            },
        }
    }

    toMarkdown(state, node) {
        state.renderContent(node)
    }

    parseMarkdown() {
        return {
            block: 'th',
            getAttrs: (tok) => ({ alignment: tok.info }),
        }
    }

    // get plugins() {
    //   if (this.editor.options.readonly) {
    //     return [];
    //   }
    //
    //   return [
    //     new Plugin({
    //       props: {
    //         decorations: (state) => {
    //           const { doc, selection } = state;
    //           const { view } = this.editor;
    //           const decorations = [];
    //           const cells = getCellsInRow(0)(selection);
    //           console.log(cells)
    //           if (cells) {
    //             cells.forEach(({ pos }, index) => {
    //               decorations.push(
    //                 Decoration.widget(pos + 1, () => {
    //                   const colSelected = isColumnSelected(index)(selection);
    //                   let className = "grip-column";
    //                   if (colSelected) {
    //                     className += " selected";
    //                   }
    //                   if (index === 0) {
    //                     className += " first";
    //                   } else if (index === cells.length - 1) {
    //                     className += " last";
    //                   }
    //                   const grip = document.createElement("a");
    //                   grip.className = className;
    //                   grip.addEventListener("mousedown", (event) => {
    //                     event.preventDefault();
    //                     console.log(index)
    //                     view.dispatch(selectColumn(index)(state.tr));
    //                     // selectCol(state,view.dispatch,index)
    //                   });
    //                   return grip;
    //                 })
    //               );
    //             });
    //           }
    //
    //           return DecorationSet.create(doc, decorations);
    //         },
    //       },
    //     }),
    //   ];
    // }
}
