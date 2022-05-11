import { h, render } from 'vue'
import tippy from 'tippy.js'
import { $on } from '@ac/shared'

import BlockMenuCtor from './BlockMenu.vue'

import { BLOCK_MENU_TRIGGER_EVENT } from './BlockMenuTrigger'

export default class BlockMenu {
    constructor(editor, options) {
        this.editor = editor
        this.options = options

        this.renderVm()
        this.createPopper()

        this.bindEvent()

        this.editor.on(BLOCK_MENU_TRIGGER_EVENT.OPEN, (searchKeyword) => this.open(searchKeyword))
        this.editor.on(BLOCK_MENU_TRIGGER_EVENT.CLOSE, () => this.close())

        this.editor.on('destroy', () => this.destroy())
    }

    renderVm() {
        const dom = document.createElement('div')

        const BlockMenuComponent = h(BlockMenuCtor, {
            ...this.options,
            isActive: false,
            searchKeyword: null,
            editor: this.editor,
        })

        render(BlockMenuComponent, dom)

        this.$vm = BlockMenuComponent.component.proxy
        this.$props = BlockMenuComponent.component.props
    }

    bindEvent() {
        $on(this.$vm, 'close', () => this.close())
        $on(this.$vm, 'onCreateLinkTrigger', () => this.onCreateLinkTrigger())
    }

    createPopper() {
        this.popper = tippy(document.body, {
            content: this.$vm.$el,
            appendTo: () => document.body,
            theme: 'light',
            placement: 'bottom-start',
            maxWidth: 'none',
            trigger: 'manual',
            arrow: false,
            interactive: true,
        })

        this.popper.hide()
    }

    open(searchKeyword) {
        document.body.style.overflow = 'hidden'

        const { view, state } = this.editor
        const { selection } = state
        const paragraph = view.domAtPos(selection.$from.pos)

        if (paragraph.node && paragraph.node.getBoundingClientRect) {
            this.popper.setProps({
                getReferenceClientRect: () => {
                    let { width, height, x, y } = paragraph.node.getBoundingClientRect()
                    if (!width && !height && !x && !y) {
                        return {
                            width,
                            height,
                            left: -1000,
                            top: -1000,
                        }
                    }
                    return paragraph.node.getBoundingClientRect()
                },
            })
            this.popper.show()

            this.$props.searchKeyword = searchKeyword
            this.$props.isActive = true
        }
    }

    close() {
        this.$props && (this.$props.isActive = false)
        this.popper.hide()
        document.body.style.overflow = null
    }

    hide() {
        this.close()
    }

    onCreateLinkTrigger() {
        this.editor.linkToolbar && this.editor.linkToolbar.show()
    }

    destroy() {
        this.popper && this.popper.destroy()
        this.$vm = null
    }
}
