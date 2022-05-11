import refractor from 'refractor/core'
import bash from 'refractor/lang/bash'
import css from 'refractor/lang/css'
import clike from 'refractor/lang/clike'
import csharp from 'refractor/lang/csharp'
import go from 'refractor/lang/go'
import java from 'refractor/lang/java'
import javascript from 'refractor/lang/javascript'
import json from 'refractor/lang/json'
import markup from 'refractor/lang/markup'
import php from 'refractor/lang/php'
import python from 'refractor/lang/python'
import powershell from 'refractor/lang/powershell'
import ruby from 'refractor/lang/ruby'
import sql from 'refractor/lang/sql'
import typescript from 'refractor/lang/typescript'

import Prism, { LANGUAGES } from '../plugins/Prism'
import Node from '../lib/Node'
import { isInCode } from '../utils'
;[bash, css, clike, csharp, go, java, javascript, json, markup, php, python, powershell, ruby, sql, typescript].forEach(refractor.register)

/**
 * 代码块
 * @type {string}
 */
export default class CodeBlockHighlight extends Node {
    constructor(options = {}) {
        super(options)
    }

    get languageOptions() {
        return Object.entries(LANGUAGES)
    }

    get defaultOptions() {
        return {
            languages: this.languageOptions,
        }
    }

    get name() {
        return 'code_block'
    }

    get schema() {
        return {
            attrs: {
                language: {
                    default: 'json',
                },
            },
            content: 'text*',
            marks: '',
            group: 'block',
            code: true,
            defining: true,
            draggable: false,
            parseDOM: [
                { tag: 'pre', preserveWhitespace: 'full' },
                {
                    tag: '.code-block',
                    preserveWhitespace: 'full',
                    contentElement: 'code',
                    getAttrs: (dom) => {
                        return {
                            language: dom.dataset.language,
                        }
                    },
                },
            ],
            toDOM: (node) => {
                const button = document.createElement('button')
                button.innerText = '复制'
                button.className = 'copy_text'
                button.type = 'button'
                button.addEventListener('click', this.handleCopyToClipboard)

                const select = document.createElement('select')
                select.addEventListener('change', (e) => this.handleLanguageChange(e))

                this.languageOptions.forEach(([key, label]) => {
                    const option = document.createElement('option')
                    const value = key === 'none' ? '' : key
                    option.value = value
                    option.innerText = label
                    option.selected = node.attrs.language === value
                    select.appendChild(option)
                })

                return [
                    'div',
                    { class: 'code-block', 'data-language': node.attrs.language },
                    ['div', { contentEditable: false }, select, button],
                    ['pre', ['code', { spellCheck: false }, 0]],
                ]
            },
        }
    }

    keys() {
        return {
            // "Shift-Ctrl-\\": setBlockType(type),
            'Shift-Enter': (state, dispatch) => {
                if (!isInCode(state)) return false

                const { tr, selection } = state
                dispatch(tr.insertText('\n', selection.from, selection.to))
                return true
            },
            Tab: (state, dispatch) => {
                if (!isInCode(state)) return false

                const { tr, selection } = state
                dispatch(tr.insertText('    ', selection.from, selection.to))
                return true
            },
        }
    }

    handleLanguageChange(event) {
        const { view } = this.editor
        const { tr } = view.state
        const element = event.target
        const { top, left } = element.getBoundingClientRect()
        const result = view.posAtCoords({ top, left })

        if (result) {
            const transaction = tr.setNodeMarkup(result.inside, undefined, {
                language: element.value,
            })
            view.dispatch(transaction)
        }
    }

    handleCopyToClipboard = (event) => {
        const { view } = this.editor
        const element = event.target
        const { top, left } = element.getBoundingClientRect()
        const result = view.posAtCoords({ top, left })

        if (result) {
            const node = view.state.doc.nodeAt(result.pos)
            node && element.setAttribute('data-text', node.textContent)
        }
    }

    get plugins() {
        return [Prism({ name: this.name })]
    }
}
