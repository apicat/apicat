import { InputRule } from 'prosemirror-inputrules'
import { Plugin } from 'prosemirror-state'
import { findParentNode } from 'prosemirror-utils'
import { Decoration, DecorationSet } from 'prosemirror-view'
import Extension from '../../lib/Extension'

const MAX_MATCH = 500
const OPEN_REGEX = /^\/(\w+|[@\w]*)?$/
const CLOSE_REGEX = /(^(?!\/(\w+)?)(.*)$|^\/((\w+)\s.*|\s)$)/

/**
 * 控制 回车键 是否可以执行
 * based on the input rules code in Prosemirror, here:
 * https://github.com/ProseMirror/prosemirror-inputrules/blob/master/src/inputrules.js
 */

export const BLOCK_MENU_TRIGGER_EVENT = {
    OPEN: 'block.menu.open',
    CLOSE: 'block.menu.close',
}

function run(view, from, to, regex, handler) {
    if (view.composing) {
        return false
    }
    const state = view.state
    const $from = state.doc.resolve(from)
    const node = $from.node(1)
    if ($from.parent.type.spec.code || (node && node.type.name !== 'paragraph')) {
        return false
    }

    const textBefore = $from.parent.textBetween(Math.max(0, $from.parentOffset - MAX_MATCH), $from.parentOffset, null, '\ufffc')

    const match = regex.exec(textBefore)
    const tr = handler(state, match)
    if (!tr) return false
    return true
}

export default class BlockMenuTrigger extends Extension {
    get allowCreate() {
        return !this.editor.options.readonly
    }

    get name() {
        return 'blockmenu'
    }

    get plugins() {
        return [
            new Plugin({
                props: {
                    handleClick: () => {
                        this.editor.emit(BLOCK_MENU_TRIGGER_EVENT.CLOSE)
                        return false
                    },
                    handleKeyDown: (view, event) => {
                        // 模拟过滤功能
                        // Prosemirror input rules are not triggered on backspace, however
                        // we need them to be evaluted for the filter trigger to work
                        // correctly. This additional handler adds inputrules-like handling.
                        if (event.key === 'Backspace') {
                            // timeout ensures that the delete has been handled by prosemirror
                            // and any characters removed, before we evaluate the rule.
                            setTimeout(() => {
                                const { pos } = view.state.selection.$from
                                return run(view, pos, pos, OPEN_REGEX, (state, match) => {
                                    this.editor.emit(match ? BLOCK_MENU_TRIGGER_EVENT.OPEN : BLOCK_MENU_TRIGGER_EVENT.CLOSE, match ? match[1] : null)
                                    return null
                                })
                            })
                        }

                        // If the query is active and we're navigating the block menu then
                        // just ignore the key events in the editor itself until we're done
                        if (event.key === 'Enter' || event.key === 'ArrowUp' || event.key === 'ArrowDown' || event.key === 'Tab') {
                            const { pos } = view.state.selection.$from
                            // 当提示菜单显示时，拦截键盘事件，不做任何事情
                            return run(view, pos, pos, OPEN_REGEX, (state, match) => {
                                // just tell Prosemirror we handled it and not to do anything
                                return match ? true : null
                            })
                        }

                        return false
                    },
                    decorations: (state) => {
                        const parent = findParentNode((node) => node.type.name === 'paragraph')(state.selection)

                        if (!parent) {
                            return
                        }

                        const decorations = []
                        const isEmpty = parent && parent.node.content.size === 0
                        const isSlash = parent && parent.node.textContent === '/'
                        const isTopLevel = state.selection.$from.depth === 1

                        if (isTopLevel) {
                            if (isEmpty) {
                                // decorations.push(Decoration.widget(parent.pos, () => button));

                                decorations.push(
                                    Decoration.node(parent.pos, parent.pos + parent.node.nodeSize, {
                                        class: 'placeholder',
                                        'data-empty-text': this.options.dictionary.newLineEmpty,
                                    })
                                )
                            }

                            if (isSlash) {
                                decorations.push(
                                    Decoration.node(parent.pos, parent.pos + parent.node.nodeSize, {
                                        class: 'placeholder',
                                        'data-empty-text': `  ${this.options.dictionary.newLineWithSlash}`,
                                    })
                                )
                            }

                            return DecorationSet.create(state.doc, decorations)
                        }
                    },
                },
            }),
        ]
    }

    inputRules() {
        return [
            new InputRule(OPEN_REGEX, (state, match) => {
                if (match && state.selection.$from.node(1).type.name === 'paragraph') {
                    this.editor.emit(BLOCK_MENU_TRIGGER_EVENT.OPEN, match[1])
                }
                return null
            }),
            new InputRule(CLOSE_REGEX, (state, match) => {
                if (match) {
                    this.editor.emit(BLOCK_MENU_TRIGGER_EVENT.CLOSE)
                }
                return null
            }),
        ]
    }
}
