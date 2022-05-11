import Mark from '../lib/Mark'
import { toggleMark, markInputRule } from '../commands'

export default class Italic extends Mark {
    get name() {
        return 'italic'
    }

    get schema() {
        return {
            parseDOM: [{ tag: 'i' }, { tag: 'em' }, { style: 'font-style=italic' }],
            toDOM: () => ['em', 0],
        }
    }

    keys({ type }) {
        return {
            'Mod-i': toggleMark(type),
        }
    }

    inputRules({ type }) {
        return [markInputRule(/(?:^|[^_])(_([^_]+)_)$/, type), markInputRule(/(?:^|[^*])(\*([^*]+)\*)$/, type)]
    }

    get toMarkdown() {
        return {
            open: '*',
            close: '*',
            mixable: true,
            expelEnclosingWhitespace: true,
        }
    }

    parseMarkdown() {
        return { mark: 'em' }
    }
}
