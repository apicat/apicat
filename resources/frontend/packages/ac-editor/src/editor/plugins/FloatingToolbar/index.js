import tippy from 'tippy.js'
import { PluginKey, Plugin, AllSelection } from 'prosemirror-state'
import { h, render } from 'vue'
import Component from './SelectionMenus.vue'
import { isNodeSelection } from 'prosemirror-utils'
import { isNodeActive, isTextSelection } from '../../utils'
import posToDOMRect from '../../utils/posToDOMRect'
import { $once } from '@ac/shared'

export const FloatingToolbarPluginKey = new PluginKey('FloatingToolbarPluginKey')

export default class FloatingToolbar {
    constructor(editor, options) {
        this.editor = editor
        this.options = options

        this.element = this.renderVm()

        this.element.addEventListener('mousedown', this.mousedownHandler, {
            capture: true,
        })

        this.element.style.visibility = 'visible'

        this.createPopper()

        this.editor.registerPlugin(this.createPlugin())
    }

    renderVm() {
        const dom = document.createElement('div')
        const vNode = h(Component, {
            ...this.options,
            editor: this.editor,
        })

        render(vNode, dom)

        this.$vm = vNode.component.proxy

        $once(this.$vm, 'on-close', () => this.hide())

        return vNode.component.proxy.$el
    }

    createPopper() {
        this.popper = tippy(document.body, {
            appendTo: () => document.body,
            duration: 0,
            maxWidth: 'none',
            getReferenceClientRect: null,
            content: this.element,
            interactive: true,
            theme: 'light',
            trigger: 'manual',
            placement: 'top',
            hideOnClick: 'toggle',
        })
    }

    createPlugin() {
        return new Plugin({
            key: FloatingToolbarPluginKey,
            view: () => ({
                update: (view, oldState) => this.update(view, oldState),
                destroy: () => this.destroy(),
            }),
            props: {
                handleDOMEvents: {
                    focus: (view, event) => {
                        this.focusHandler(event)
                        return false
                    },
                    blur: (view, event) => {
                        this.blurHandler(event)
                        return false
                    },
                },
            },
        })
    }

    mousedownHandler = () => {
        this.preventHide = true
    }

    focusHandler = () => {
        setTimeout(() => this.update(this.editor.view))
    }

    blurHandler = (event) => {
        if (this.preventHide) {
            this.preventHide = false

            return
        }

        let _a
        if (
            (event === null || event === void 0 ? void 0 : event.relatedTarget) &&
            ((_a = this.element.parentNode) === null || _a === void 0 ? void 0 : _a.contains(event.relatedTarget))
        ) {
            return
        }

        this.hide()
    }

    update(view, oldState) {
        if (!this.editor.isEditable || oldState === view.state) {
            return
        }

        const { state, composing } = view
        const { doc, selection } = state
        const isSame = oldState && oldState.doc.eq(doc) && oldState.selection.eq(selection)

        if (composing || isSame) {
            return
        }

        const { empty, ranges } = selection

        // support for CellSelections
        const from = Math.min(...ranges.map((range) => range.$from.pos))
        const to = Math.max(...ranges.map((range) => range.$to.pos))

        const isCodeSelection = isNodeActive(state.schema.nodes.code_block)(state)
        const isImageSelection = isNodeActive(state.schema.nodes.image)(state)
        const isEmptyTextBlock = !doc.textBetween(from, to).length && isTextSelection(view.state.selection)

        // console.log(empty,isEmptyTextBlock,isCodeSelection,(selection instanceof AllSelection), (isNodeSelection(selection) && !isImageSelection))
        if (empty || isEmptyTextBlock || isCodeSelection || selection instanceof AllSelection || (isNodeSelection(selection) && !isImageSelection)) {
            this.hide()
            return
        }

        this.$vm.updateFloatMenus()

        this.popper.setProps({
            getReferenceClientRect: () => {
                return posToDOMRect(view, selection, from, to)
            },
        })

        this.show()
    }

    show() {
        this.popper.show()
    }

    hide() {
        this.popper.hide()
    }

    destroy() {
        this.popper.destroy()
        this.element.removeEventListener('mousedown', this.mousedownHandler)
    }
}
