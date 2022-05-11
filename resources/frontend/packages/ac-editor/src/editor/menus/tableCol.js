import { CellSelection, isInTable } from 'prosemirror-tables'
import { cellWrapping, isNodeActive } from '../utils'

export default function tableColMenuItems(state, index, dictionary) {
    const { schema } = state
    const isTable = isInTable(state)

    return [
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignLeft,
            icon: 'editor-text-align-left',
            attrs: { index, alignment: 'left' },
            active: isNodeActive(schema.nodes.th, {
                alignment: 'left',
            }),
        },
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignCenter,
            icon: 'editor-text-align-center',
            attrs: { index, alignment: 'center' },
            active: isNodeActive(schema.nodes.th, {
                alignment: 'center',
            }),
        },
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignRight,
            icon: 'editor-text-align-right',
            attrs: { index, alignment: 'right' },
            active: isNodeActive(schema.nodes.th, {
                alignment: 'right',
            }),
        },
        {
            name: 'separator',
        },
        {
            name: 'addColumnBefore',
            tooltip: dictionary.addColumnBefore,
            icon: 'editor-tablecolumnplusbefore',
            active: () => false,
        },
        {
            name: 'addColumnAfter',
            tooltip: dictionary.addColumnAfter,
            icon: 'editor-table-column-plus-after',
            active: () => false,
        },
        {
            name: 'deleteColumn',
            tooltip: dictionary.deleteColumn,
            icon: 'editor-tablecolumnremove',
            active: () => false,
        },
        {
            name: 'separator',
        },
        {
            name: 'splitCell',
            tooltip: '分割单元格',
            icon: 'editor-split-cells-horizontal',
            visible: isTable,
            disable: (state) => {
                const sel = state.selection
                let cellNode
                if (sel instanceof CellSelection) {
                    if (sel.$anchorCell.pos !== sel.$headCell.pos) return true
                    cellNode = sel.$anchorCell.nodeAfter
                } else {
                    cellNode = cellWrapping(sel.$from)
                    if (!cellNode) return true
                }
                return cellNode.attrs.colspan === 1 && cellNode.attrs.rowspan === 1
            },
        },

        {
            name: 'mergeCells',
            tooltip: '合并单元格',
            icon: 'editor-merge-cells-horizontal',
            visible: isTable,
            disable: (state) => !(state.selection instanceof CellSelection) || state.selection.$anchorCell.pos === state.selection.$headCell.pos,
        },
    ]
}
