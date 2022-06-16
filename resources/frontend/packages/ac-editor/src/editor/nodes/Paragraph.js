import { setBlockType } from '../commands'
import Node from '../lib/Node'
import { NodeSelection } from 'prosemirror-state'

export default class Paragraph extends Node {
    get name() {
        return 'paragraph'
    }

    get schema() {
        return {
            content: 'inline*',
            group: 'block',
            draggable: false,
            parseDOM: [
                {
                    tag: 'p',
                },
            ],
            toDOM: () => ['p', 0],
        }
    }

    commands({ type }) {
        return () => setBlockType(type)
    }

    keys({ type }) {
        return {
            // 删除第一个段落
            Backspace: (state, dispatch) => {
                const { empty } = state.selection
                const node = state.selection.$anchor.node()
                if (empty && node.type === type && !node.textContent && state.selection.anchor === 1) {
                    const tr = state.tr
                    tr.setSelection(NodeSelection.create(state.doc, 0))
                    dispatch(tr)
                }
                return false
            },
        }
    }

    toMarkdown(state, node) {
        if (node.textContent.trim() === '' && node.childCount === 0 && !state.inTable) {
            state.write('\\\n')
        } else {
            state.renderInline(node)
            state.closeBlock(node)
        }
    }

    parseMarkdown() {
        return { block: 'paragraph' }
    }
}
