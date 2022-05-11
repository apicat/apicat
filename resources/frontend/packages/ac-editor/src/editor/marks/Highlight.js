import { markInputRule } from '../utils'
import Mark from '../lib/Mark'
import applyMark from '../commands/applyMark'

const getAttrs = (dom) => {
    let attrs

    try {
        attrs = JSON.parse(dom.getAttribute('data-color'))
    } catch (e) {
        attrs = {}
    }

    return attrs
}

const setAttrs = (node) => {
    const { fontColor, bgColor } = node.attrs

    let style = ''
    let attrs = {}
    if (bgColor) {
        style += `background-color: ${bgColor};`
    }
    if (fontColor) {
        style += `color: ${fontColor};`
    }
    attrs.style = style
    attrs['data-color'] = JSON.stringify(node.attrs)
    return attrs
}

export default class Highlight extends Mark {
    get name() {
        return 'highlight'
    }

    get schema() {
        return {
            attrs: {
                fontColor: {
                    default: '',
                },
                bgColor: {
                    default: '',
                },
            },
            inline: true,
            group: 'inline',
            parseDOM: [{ tag: 'mark[data-color]', getAttrs }],
            toDOM: (node) => ['mark', setAttrs(node)],
        }
    }

    inputRules({ type }) {
        return [markInputRule(/(?:==)([^=]+)(?:==)$/, type)]
    }

    commands({ type }) {
        return (attrs) => (state, dispatch) => {
            let { tr } = state
            tr = applyMark(tr.setSelection(state.selection), type, attrs)
            if (tr.docChanged || tr.storedMarksSet) {
                dispatch && dispatch(tr)
                return true
            }
        }
    }

    get toMarkdown() {
        return {
            open: '==',
            close: '==',
            mixable: true,
            expelEnclosingWhitespace: true,
        }
    }

    parseMarkdown() {
        return { mark: 'highlight' }
    }
}
