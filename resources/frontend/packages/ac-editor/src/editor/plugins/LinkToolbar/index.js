import { h, render } from 'vue'
import tippy from 'tippy.js'
import Component from './LinkEditor.vue'
import { $on } from '@ac/shared'

export default class LinkToolbar {
    constructor(editor) {
        this.editor = editor
        this.view = editor.view

        this.element = this.renderVm()
        this.popper = this.createPopper()

        this.bindEvent()
    }

    renderVm() {
        const dom = document.createElement('div')
        const vNode = h(Component, {
            mark: this.editor.schema.marks.link.create(),
            isCreate: true,
        })

        render(vNode, dom)
        this.$vm = vNode.component.proxy
        return this.$vm.$el
    }

    bindEvent() {
        $on(this.$vm, 'on-create', (attr) => this.onCreateLink(attr))
        $on(this.$vm, 'on-close', () => this.close())
        $on(this.$vm, 'toggle-blank', (attr) => this.onToggleBlank(attr))
        window.addEventListener('mousedown', this.handleClickOutside)
    }

    handleClickOutside = (ev) => {
        if (ev.target && this.element && this.element.contains(ev.target)) {
            return
        }
        this.hide()
    }

    onToggleBlank(attrs) {
        this.onCreateLink(attrs)
    }

    close() {
        this.popper.hide()
    }

    onCreateLink(attrs) {
        this.createLink(attrs)
        this.view.focus()
        this.close()
    }

    createLink(attrs) {
        const { dispatch, state } = this.view
        const { from, to } = state.selection
        const title = attrs.href

        dispatch(this.view.state.tr.insertText(attrs.href, from, to).addMark(from, to + title.length, state.schema.marks.link.create(attrs)))
    }

    createPopper() {
        return tippy(document.body, {
            appendTo: () => document.body,
            duration: 0,
            getReferenceClientRect: null,
            content: this.element,
            interactive: true,
            theme: 'light',
            trigger: 'manual',
            placement: 'top',
            hideOnClick: 'toggle',
        })
    }

    show() {
        let { view } = this.editor
        let { selection } = this.editor.state
        const { node: element } = view.domAtPos(selection.from)

        this.popper.setProps({
            getReferenceClientRect: () => {
                let rect = element.getBoundingClientRect().toJSON()
                rect.width = 1
                return rect
            },
        })

        this.popper.show()

        // todo 为什么会影响滚动行为？
        setTimeout(() => {
            this.$vm.focus()
        }, 0)
    }

    hide() {
        this.close()
    }

    destroy() {
        this.popper && this.popper.destroy()
        window.removeEventListener('mousedown', this.handleClickOutside)
    }
}
