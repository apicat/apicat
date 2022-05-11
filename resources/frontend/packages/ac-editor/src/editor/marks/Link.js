import { Plugin, TextSelection } from 'prosemirror-state'
import { getMarkRange } from '../utils'
import { InputRule } from 'prosemirror-inputrules'
import Mark from '../lib/Mark'

import { updateMark, removeMark } from '../commands'

const LINK_INPUT_REGEX = /\[(.+)]\((\S+)\)/

function isPlainURL(link, parent, index, side) {
    if (link.attrs.title || !/^\w+:/.test(link.attrs.href)) {
        return false
    }

    const content = parent.child(index + (side < 0 ? -1 : 0))
    if (!content.isText || content.text !== link.attrs.href || content.marks[content.marks.length - 1] !== link) {
        return false
    }

    if (index === (side < 0 ? 1 : parent.childCount - 1)) {
        return true
    }

    const next = parent.child(index + (side < 0 ? -2 : 1))
    return !link.isInSet(next.marks)
}

function getAttrs(dom) {
    return {
        href: dom.getAttribute('href'),
        openInNewTab: dom.getAttribute('target') === '_blank',
    }
}

function toDOM(mark) {
    const { href, openInNewTab } = mark.attrs

    const attrs = {}
    attrs.href = href

    let ref = 'nofollow'

    if (openInNewTab) {
        attrs.target = '_blank'
        ref += ' noopener noreferrer'
    }

    attrs.ref = ref.trim()

    return ['a', attrs, 0]
}

export default class Link extends Mark {
    get name() {
        return 'link'
    }

    get schema() {
        return {
            attrs: {
                href: {
                    default: '',
                },
                openInNewTab: {
                    default: true,
                },
            },
            inclusive: false,
            parseDOM: [
                {
                    tag: 'a[href]',
                    getAttrs,
                },
            ],
            toDOM,
        }
    }

    inputRules({ type }) {
        return [
            new InputRule(LINK_INPUT_REGEX, (state, match, start, end) => {
                const [okay, alt, href] = match
                const { tr } = state

                if (okay) {
                    tr.replaceWith(start, end, this.editor.schema.text(alt)).addMark(start, start + alt.length, type.create({ href }))
                }

                return tr
            }),
        ]
    }

    commands({ type }) {
        return (attrs) => {
            if (attrs) {
                return updateMark(type, attrs)
            }
            return removeMark(type)
        }
    }

    get plugins() {
        return [
            new Plugin({
                props: {
                    handleClick(view, pos) {
                        const { schema, doc, tr } = view.state

                        const range = getMarkRange(doc.resolve(pos), schema.marks.link)

                        if (!range) return false

                        const $start = doc.resolve(range.from)
                        const $end = doc.resolve(range.to)

                        const transaction = tr.setSelection(new TextSelection($start, $end))

                        view.dispatch(transaction)
                        return true
                    },
                },
            }),
        ]
    }

    get toMarkdown() {
        return {
            open(_state, mark, parent, index) {
                return isPlainURL(mark, parent, index, 1) ? '<' : '['
            },
            close(state, mark, parent, index) {
                return isPlainURL(mark, parent, index, -1)
                    ? '>'
                    : '](' + state.esc(mark.attrs.href) + (mark.attrs.title ? ' ' + state.quote(mark.attrs.title) : '') + ')'
            },
        }
    }

    parseMarkdown() {
        return {
            mark: 'link',
            getAttrs: (tok) => ({
                href: tok.attrGet('href'),
                title: tok.attrGet('title') || null,
            }),
        }
    }
}
