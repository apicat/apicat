import Mark from '../lib/Mark'
import { toggleMark } from '../commands'

export default class Underline extends Mark {
    get name() {
        return 'underline'
    }

    get schema() {
        return {
            parseDOM: [
                {
                    tag: 'u',
                },
                {
                    style: 'text-decoration',
                    getAttrs: (value) => value === 'underline',
                },
            ],
            toDOM: () => ['u', 0],
        }
    }

    keys({ type }) {
        return {
            'Mod-u': toggleMark(type),
        }
    }

    get toMarkdown() {
        return {
            open: '__',
            close: '__',
            mixable: true,
            expelEnclosingWhitespace: true,
        }
    }

    parseMarkdown() {
        return { mark: 'underline' }
    }
}
