import Mark from '../lib/Mark'
import { toggleMark, markInputRule } from '../commands'

export default class Bold extends Mark {
    get name() {
        return 'bold'
    }

    get schema() {
        return {
            parseDOM: [
                {
                    tag: 'strong',
                },
                {
                    tag: 'b',
                    getAttrs: (node) => node.style.fontWeight !== 'normal' && null,
                },
                {
                    style: 'font-weight',
                    getAttrs: (value) => /^(bold(er)?|[5-9]\d{2,})$/.test(value) && null,
                },
            ],
            toDOM: () => ['strong', 0],
        }
    }

    keys({ type }) {
        return {
            'Mod-b': toggleMark(type),
        }
    }

    inputRules({ type }) {
        return [markInputRule(/(?:\*\*|__)([^*_]+)(?:\*\*|__)$/, type)]
    }

    get toMarkdown() {
        return {
            open: '**',
            close: '**',
            mixable: true,
            expelEnclosingWhitespace: true,
        }
    }

    parseMarkdown() {
        return { mark: 'strong' }
    }
}
