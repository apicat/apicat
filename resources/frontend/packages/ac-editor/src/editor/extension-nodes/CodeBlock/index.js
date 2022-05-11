import Node from '../../lib/Node'
import CodeBlockNodeView from '../../extension-nodes/CodeBlock/CodeBlockNodeView'
import languages from './languages'
import { Selection } from 'prosemirror-state'

function arrowHandler(dir) {
    return (state, dispatch, view) => {
        if (state.selection.empty && view.endOfTextblock(dir)) {
            let side = dir === 'left' || dir === 'up' ? -1 : 1,
                $head = state.selection.$head
            let nextPos = Selection.near(state.doc.resolve(side > 0 ? $head.after() : $head.before()), side)
            if (nextPos.$head && nextPos.$head.parent.type.name === 'code_block') {
                dispatch(state.tr.setSelection(nextPos))
                return true
            }
        }
        return false
    }
}

export default class CodeBlock extends Node {
    get name() {
        return 'code_block'
    }

    get defaultOptions() {
        return {
            languages,
        }
    }

    get schema() {
        return {
            attrs: {
                language: {
                    default: 'json',
                },
            },
            content: 'text*',
            group: 'block',
            code: true,
            defining: true,
            isolating: true,
            marks: '',
            selectable: true,
            draggable: true,
            parseDOM: [
                {
                    tag: 'pre',
                    preserveWhitespace: 'full',
                    getAttrs: function (node) {
                        return { language: node.getAttribute('data-language') || '' }
                    },
                },
            ],
            toDOM: () => ['pre', ['code', 0]],
        }
    }

    // inputRules({ type }) {
    //   return [codeBlockRule(type)];
    // }

    keys() {
        return {
            ArrowLeft: arrowHandler('left'),
            ArrowRight: arrowHandler('right'),
            ArrowUp: arrowHandler('up'),
            ArrowDown: arrowHandler('down'),
        }
    }

    get nodeView() {
        return CodeBlockNodeView
    }
}
