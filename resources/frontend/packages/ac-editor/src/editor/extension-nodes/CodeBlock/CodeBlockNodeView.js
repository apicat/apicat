import CodeMirror from 'codemirror'
import 'codemirror/addon/edit/matchbrackets'
import { exitCode } from 'prosemirror-commands'
import { undo, redo } from 'prosemirror-history'
import { TextSelection, Selection } from 'prosemirror-state'
import { languages } from './languages'

export default class CodeBlockView {
    constructor({ extension, node, view, getPos, innerDecorations, editor, options = {} }) {
        // Store for later
        this.node = node
        this.view = view
        this.getPos = getPos
        this.editor = editor
        this.schema = editor.schema
        this.extension = extension
        this.options = this.mergeOptions(options)

        this.incomingChanges = false

        // Create a CodeMirror instance
        this.cm = new CodeMirror(null, {
            value: this.node.textContent,
            viewportMargin: Infinity,
            lineNumbers: false,
            styleActiveLine: true,
            matchBrackets: true,
            indentUnit: 4,
            // theme: "base16-light",
            theme: 'idle',
            scrollbarStyle: null,
            extraKeys: this.codeMirrorKeymap(),
            readOnly: this.options.readonly,
        })

        this.updateLanguage()

        const container = document.createElement('div')
        container.setAttribute('contenteditable', 'false')
        container.classList.add('codemirror-container')

        this.langSelect = this.renderLanguageSelect()
        this.copyBtn = this.renderCopyBtn()

        // CodeMirror needs to be in the DOM to properly initialize, so
        // schedule it to update itself
        setTimeout(() => this.cm.refresh(), 20)

        // This flag is used to avoid an update loop between the outer and
        // inner editor
        this.updating = false

        this.bindCmEvent()

        innerDecorations.find().map((d) => {
            const elem = typeof d.type.toDOM === 'function' ? d.type.toDOM() : d.type.toDOM
            container.appendChild(elem)
        })

        // 语言选择
        container.appendChild(this.langSelect)
        container.appendChild(this.copyBtn)

        // The editor's outer node is our DOM representation
        container.appendChild(this.cm.getWrapperElement())

        this.dom = container
    }

    bindCmEvent() {
        // Track whether changes are have been made but not yet propagated
        this.cm.on('beforeChange', () => (this.incomingChanges = true))
        // Propagate updates from the code editor to ProseMirror
        this.cm.on('cursorActivity', () => {
            if (!this.updating && !this.incomingChanges) this.forwardSelection()
        })
        this.cm.on('changes', () => {
            if (!this.updating) {
                this.valueChanged()
                this.forwardSelection()
            }
            this.incomingChanges = false
        })
        this.cm.on('focus', () => this.forwardSelection())
        this.cm.on('blur', () => this.formatContent())
    }

    mergeOptions(opt = {}) {
        if (this.extension && this.extension.options) {
            return { ...this.extension.options, ...opt }
        }

        return { ...opt }
    }

    renderLanguageSelect() {
        const select = document.createElement('select')
        select.addEventListener('change', (e) => this.handleLanguageChange(e))
        ;(this.options.languages || []).forEach((key) => {
            const option = document.createElement('option')
            const value = key === 'none' ? '' : key
            option.value = value
            option.innerText = key
            option.selected = this.node.attrs.language === value
            select.appendChild(option)
        })

        return select
    }

    renderCopyBtn() {
        const button = document.createElement('button')
        button.innerText = '复制'
        button.className = 'copy_text'
        button.type = 'button'
        button.addEventListener('click', this.handleCopyToClipboard)
        return button
    }

    handleLanguageChange(e) {
        e.preventDefault()
        const lang = e.target.value
        const tr = this.view.state.tr
        tr.setNodeMarkup(this.getPos(), undefined, {
            ...this.node.attrs,
            language: lang,
        })
        this.view.dispatch(tr)
        this.cm.focus()
    }

    handleCopyToClipboard = (event) => {
        const { view } = this.editor
        const element = event.target
        const node = view.state.doc.nodeAt(this.getPos())
        node && element.setAttribute('data-text', node.textContent)
    }

    forwardSelection() {
        if (!this.cm.hasFocus()) return
        let state = this.view.state
        let selection = this.asProseMirrorSelection(state.doc)
        if (!selection.eq(state.selection)) this.view.dispatch(state.tr.setSelection(selection))
    }

    asProseMirrorSelection(doc) {
        let offset = this.getPos() + 1
        let anchor = this.cm.indexFromPos(this.cm.getCursor('anchor')) + offset
        let head = this.cm.indexFromPos(this.cm.getCursor('head')) + offset
        return TextSelection.create(doc, anchor, head)
    }

    setSelection(anchor, head) {
        this.cm.focus()
        this.updating = true
        this.cm.setSelection(this.cm.posFromIndex(anchor), this.cm.posFromIndex(head))
        this.updating = false
    }

    formatContent() {
        if (this.node.attrs.language === 'json') {
            let json = this.cm.getValue()
            try {
                json = JSON.stringify(JSON.parse(json), null, 4)
            } catch (e) {
                json = this.cm.getValue()
            }

            this.cm.setValue(json)
        }
    }

    valueChanged() {
        let change = computeChange(this.node.textContent, this.cm.getValue())
        if (change) {
            let start = this.getPos() + 1
            let tr = this.view.state.tr.replaceWith(start + change.from, start + change.to, change.text ? this.schema.text(change.text) : null)
            this.view.dispatch(tr)
        }
    }

    codeMirrorKeymap() {
        let view = this.view
        let mod = /Mac/.test(navigator.platform) ? 'Cmd' : 'Ctrl'
        return CodeMirror.normalizeKeyMap({
            Up: () => this.maybeEscape('line', -1),
            Left: () => this.maybeEscape('char', -1),
            Down: () => this.maybeEscape('line', 1),
            Right: () => this.maybeEscape('char', 1),
            'Ctrl-Enter': () => {
                if (exitCode(view.state, view.dispatch)) view.focus()
            },
            [`${mod}-Z`]: () => undo(view.state, view.dispatch),
            [`Shift-${mod}-Z`]: () => redo(view.state, view.dispatch),
            [`${mod}-Y`]: () => redo(view.state, view.dispatch),
        })
    }

    maybeEscape(unit, dir) {
        let pos = this.cm.getCursor()
        if (
            this.cm.somethingSelected() ||
            pos.line !== (dir < 0 ? this.cm.firstLine() : this.cm.lastLine()) ||
            (unit === 'char' && pos.ch !== (dir < 0 ? 0 : this.cm.getLine(pos.line).length))
        )
            return CodeMirror.Pass
        this.view.focus()
        let targetPos = this.getPos() + (dir < 0 ? 0 : this.node.nodeSize)
        let selection = Selection.near(this.view.state.doc.resolve(targetPos), dir)
        this.view.dispatch(this.view.state.tr.setSelection(selection).scrollIntoView())
        this.view.focus()
    }

    update(node) {
        // console.log("node::", node, node === this.node, node.eq(this.node));

        if (node.type !== this.node.type || node === this.node) return false

        node.attrs.language !== this.node.attrs.language && this.updateLanguage(node.attrs.language)
        this.node = node
        let change = computeChange(this.cm.getValue(), node.textContent)
        if (change) {
            this.updating = true
            this.cm.replaceRange(change.text, this.cm.posFromIndex(change.from), this.cm.posFromIndex(change.to))
            this.updating = false
        }
        return true
    }

    updateLanguage(language) {
        this.cm.setOption('mode', languages[language || this.node.attrs.language] || null)
    }

    selectNode() {
        this.cm.focus()
    }
}

function computeChange(oldVal, newVal) {
    if (oldVal === newVal) return null
    let start = 0,
        oldEnd = oldVal.length,
        newEnd = newVal.length
    while (start < oldEnd && oldVal.charCodeAt(start) === newVal.charCodeAt(start)) ++start
    while (oldEnd > start && newEnd > start && oldVal.charCodeAt(oldEnd - 1) === newVal.charCodeAt(newEnd - 1)) {
        oldEnd--
        newEnd--
    }
    return { from: start, to: oldEnd, text: newVal.slice(start, newEnd) }
}
