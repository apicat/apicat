import { isInTable, CellSelection } from 'prosemirror-tables'
import ColorPanel from '../plugins/FloatingToolbar/ColorPanel.vue'
import { HIGHLIGHT_COLOR } from '../../common/constants'
import { isInList, isMarkActive, isNodeActive, cellWrapping } from '../utils'
import { markRaw } from 'vue'

export default function formattingMenuItems(state, dictionary) {
    const { schema } = state
    const isTable = isInTable(state)
    const isList = isInList(state)
    const allowBlocks = !isTable && !isList

    return [
        {
            name: 'heading',
            tooltip: dictionary.heading,
            icon: 'editor-h1',
            active: isNodeActive(schema.nodes.heading, { level: 1 }),
            attrs: { level: 1 },
            visible: allowBlocks,
        },
        {
            name: 'heading',
            tooltip: dictionary.subheading,
            icon: 'editor-h2',
            active: isNodeActive(schema.nodes.heading, { level: 2 }),
            attrs: { level: 2 },
            visible: allowBlocks,
        },
        {
            name: 'heading',
            tooltip: dictionary.subheading,
            icon: 'editor-h3',
            active: isNodeActive(schema.nodes.heading, { level: 3 }),
            attrs: { level: 3 },
            visible: allowBlocks,
        },
        {
            name: 'bold',
            tooltip: dictionary.strong,
            icon: 'editor-bold',
            active: isMarkActive(schema.marks.bold),
        },
        {
            name: 'italic',
            tooltip: dictionary.em,
            icon: 'editor-italic',
            active: isMarkActive(schema.marks.italic),
        },
        {
            name: 'strike',
            tooltip: dictionary.strikethrough,
            icon: 'editor-strike-through',
            active: isMarkActive(schema.marks.strike),
        },
        {
            name: 'underline',
            tooltip: dictionary.underline,
            icon: 'editor-text-underline',
            active: isMarkActive(schema.marks.underline),
        },
        {
            name: 'highlight',
            tooltip: dictionary.mark,
            icon: 'editor-text-bg-color',
            attrs: { fontColor: '', bgColor: HIGHLIGHT_COLOR[9] },
            popper: markRaw(ColorPanel),
            markType: schema.marks.highlight,
            active: isMarkActive(schema.marks.highlight),
            getStyle(attrs) {
                let style = {}
                if (attrs.bgColor) {
                    style['background-color'] = attrs.bgColor
                }

                if (attrs.fontColor) {
                    style['color'] = attrs.fontColor
                }

                return style
            },
        },
        {
            name: 'code',
            tooltip: dictionary.codeInline,
            icon: 'editor-codeblock',
            active: isMarkActive(schema.marks.code),
        },
        {
            name: 'separator',
            visible: allowBlocks,
        },
        {
            name: 'blockquote',
            tooltip: dictionary.quote,
            icon: 'editor-quote-left',
            active: isNodeActive(schema.nodes.blockquote),
            attrs: { level: 2 },
            visible: allowBlocks,
        },
        {
            name: 'separator',
        },
        {
            name: 'link',
            tooltip: dictionary.createLink,
            icon: 'editor-link',
            active: isMarkActive(schema.marks.link),
            attrs: { href: '' },
        },
        {
            name: 'separator',
        },
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignLeft,
            icon: 'editor-text-align-left',
            attrs: { alignment: 'left' },
            visible: isTable,
            active: () => false,
        },
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignCenter,
            icon: 'editor-text-align-center',
            attrs: { alignment: 'center' },
            visible: isTable,
            active: () => false,
        },
        {
            name: 'setColumnAttr',
            tooltip: dictionary.alignRight,
            icon: 'editor-text-align-right',
            attrs: { alignment: 'right' },
            visible: isTable,
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
