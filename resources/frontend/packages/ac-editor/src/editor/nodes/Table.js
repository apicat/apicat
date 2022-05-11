import Node from '../lib/Node'
import {
    addColumnAfter,
    addColumnBefore,
    deleteColumn,
    deleteRow,
    deleteTable,
    goToNextCell,
    mergeCells,
    setCellAttr,
    splitCell,
    tableEditing,
    toggleHeaderCell,
    toggleHeaderColumn,
    toggleHeaderRow,
} from 'prosemirror-tables'
import { addRowAt, createTable, moveRow } from 'prosemirror-utils'
import { TextSelection } from 'prosemirror-state'
import { tableSidebar } from '../plugins/table/tableSidebar'

export default class Table extends Node {
    get name() {
        return 'table'
    }

    get schema() {
        return {
            draggable: true,
            selectable: true,

            content: 'tr+',
            tableRole: 'table',
            isolating: true,
            group: 'block',
            parseDOM: [{ tag: 'table' }],
            toDOM: () => {
                return [
                    'div',
                    { class: 'scrollable-wrapper original-table' },
                    ['div', { class: 'scrollable' }, ['table', { class: 'rme-table' }, ['tbody', 0]]],
                ]
            },
        }
    }

    commands({ schema }) {
        return {
            createTable:
                ({ rowsCount, colsCount }) =>
                (state, dispatch) => {
                    const offset = state.tr.selection.anchor + 1
                    const nodes = createTable(schema, rowsCount, colsCount)
                    const tr = state.tr.replaceSelectionWith(nodes).scrollIntoView()
                    const resolvedPos = tr.doc.resolve(offset)

                    tr.setSelection(TextSelection.near(resolvedPos))
                    dispatch(tr)
                },
            setColumnAttr:
                ({ alignment }) =>
                (state, dispatch) => {
                    setCellAttr('alignment', alignment)(state, dispatch)
                },
            addColumnBefore: () => addColumnBefore,
            addColumnAfter: () => addColumnAfter,
            deleteColumn: () => deleteColumn,
            addRowAfter:
                ({ index }) =>
                (state, dispatch) =>
                    dispatch(addRowAt(index + 1, true)(state.tr)),
            deleteRow: () => deleteRow,
            deleteTable: () => deleteTable,
            toggleHeaderColumn: () => toggleHeaderColumn,
            toggleHeaderRow: () => toggleHeaderRow,
            toggleHeaderCell: () => toggleHeaderCell,
            setCellAttr: () => setCellAttr,
            mergeCells: () => mergeCells,
            splitCell: () => splitCell,
        }
    }

    keys() {
        return {
            Tab: goToNextCell(1),
            'Shift-Tab': goToNextCell(-1),
        }
    }

    toMarkdown(state, node) {
        state.renderTable(node)
        state.closeBlock(node)
    }

    parseMarkdown() {
        return { block: 'table' }
    }

    get plugins() {
        if (this.editor.options.readonly) {
            return []
        }
        return [tableEditing({ allowTableNodeSelection: true }), tableSidebar()]
    }
}
