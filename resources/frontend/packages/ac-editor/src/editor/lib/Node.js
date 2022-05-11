import Extension from './Extension'

export default class Node extends Extension {
    get type() {
        return 'node'
    }

    get schema() {
        return {}
    }

    get markdownToken() {
        return ''
    }

    commands({ type }) {
        return (attrs) => (state, dispatch) => {
            const tr = state.tr
            tr.replaceSelectionWith(type.create(attrs))
            dispatch(tr)
            return true
        }
    }

    toMarkdown(state, node) {
        console.error('toMarkdown not implemented', state, node)
    }

    parseMarkdown() {
        return
    }
}
