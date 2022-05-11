import { CellSelection, isInTable } from 'prosemirror-tables'
import { cellWrapping } from '../utils'

export default function tableRowMenuItems(state, index, dictionary) {
    const isTable = isInTable(state)
    return [
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignLeft,
            icon: 'editor-text-align-left',
            attrs: { alignment: 'left' },
            active: () => false,
        },
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignCenter,
            icon: 'editor-text-align-center',
            attrs: { alignment: 'center' },
            active: () => false,
        },
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignRight,
            icon: 'editor-text-align-right',
            attrs: { alignment: 'right' },
            active: () => false,
        },
        {
            name: 'separator',
        },
        {
            name: 'addRowAfter',
            tooltip: dictionary.addRowBefore,
            icon: 'editor-tablerowplusbefore',
            attrs: { index: index - 1 },
            active: () => false,
            visible: index !== 0,
        },
        {
            name: 'addRowAfter',
            tooltip: dictionary.addRowAfter,
            icon: 'editor-table-row-plus-after',
            attrs: { index },
            active: () => false,
        },
        {
            name: 'deleteRow',
            tooltip: dictionary.deleteRow,
            icon: 'editor-tablerowremove',
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
