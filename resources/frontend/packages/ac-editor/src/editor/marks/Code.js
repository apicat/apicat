import Mark from '../lib/Mark'
import { toggleMark, markInputRule } from '../commands'
import { backticksFor } from '../utils'
import { InputRule } from 'prosemirror-inputrules'
import moveLeft from '../commands/moveLeft'
import moveRight from '../commands/moveRight'

export default class Code extends Mark {
    get name() {
        return 'code'
    }

    get schema() {
        return {
            excludes: '_',
            parseDOM: [{ tag: 'code' }],
            toDOM: () => ['code', { class: 'code' }, 0],
        }
    }

    keys({ type }) {
        return {
            'Mod-`': toggleMark(type),
            ArrowLeft: moveLeft(),
            ArrowRight: moveRight(),
        }
    }

    inputRules({ type }) {
        return [
            markInputRule(/(?:`)([^`]+)(?:`)$/, type),
            new InputRule(/\s\s$/, (state) => {
                const { $cursor } = state.selection
                if (type.isInSet(state.storedMarks || $cursor.marks())) {
                    let tr = state.tr.insertText('', $cursor.pos - 1, $cursor.pos)
                    return tr.removeStoredMark(type)
                }
            }),
        ]
    }

    get toMarkdown() {
        return {
            open(_state, _mark, parent, index) {
                return backticksFor(parent.child(index), -1)
            },
            close(_state, _mark, parent, index) {
                return backticksFor(parent.child(index - 1), 1)
            },
            escape: false,
        }
    }

    parseMarkdown() {
        return { mark: 'code_inline' }
    }
}
