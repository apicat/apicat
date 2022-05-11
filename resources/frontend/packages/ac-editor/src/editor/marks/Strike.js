import Mark from '../lib/Mark'
import { toggleMark, markInputRule } from '../commands'

export default class Strike extends Mark {
    get name() {
        return 'strike'
    }

    get schema() {
        return {
            parseDOM: [
                {
                    tag: 's',
                },
                {
                    tag: 'del',
                },
                {
                    tag: 'strike',
                },
                {
                    style: 'text-decoration',
                    getAttrs: (value) => value === 'line-through',
                },
            ],
            toDOM: () => ['s', 0],
        }
    }

    keys({ type }) {
        return {
            'Mod-d': toggleMark(type),
        }
    }

    inputRules({ type }) {
        return [markInputRule(/~([^~]+)~$/, type)]
    }

    get markdownToken() {
        return 's'
    }

    parseMarkdown() {
        return { mark: 'strikethrough' }
    }
}
