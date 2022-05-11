import Node from '../lib/Node'
import { getCellAttrs, setCellAttrs } from '../utils/tableCellAttr'
// import { DecorationSet, Decoration } from "prosemirror-view";
// import { Plugin } from "prosemirror-state";
// import {
//   isRowSelected,
//   getCellsInColumn,
//   // selectRow,
//   // isTableSelected,
//   // selectTable,
// } from "prosemirror-utils";
// import {selectRow} from "@/editor/plugins/table/commands";

export default class TableCell extends Node {
    get name() {
        return 'td'
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
            tableRole: 'cell',
            isolating: true,

            parseDOM: [{ tag: 'td', getAttrs: (dom) => getCellAttrs(dom) }],
            toDOM(node) {
                return ['td', setCellAttrs(node), 0]
            },
        }
    }

    toMarkdown(state, node) {
        state.renderContent(node)
    }

    parseMarkdown() {
        return {
            block: 'td',
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
    //           const cells = getCellsInColumn(0)(selection);
    //
    //           if (cells) {
    //             cells.forEach(({ pos }, index) => {
    //               decorations.push(
    //                 Decoration.widget(pos + 1, () => {
    //                   let className = "grip-row";
    //                   if (isRowSelected(index)(selection)) {
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
    //                     // view.dispatch(selectRow(index)(state.tr));
    //                     selectRow(state,view.dispatch,index)
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
