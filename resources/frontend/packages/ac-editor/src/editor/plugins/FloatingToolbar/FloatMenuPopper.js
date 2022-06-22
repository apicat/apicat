import tippy from 'tippy.js'
import { $on } from '@natosoft/shared'
import { h, render } from 'vue'

export default class FloatMenuPopper {
    constructor(elements, menus, editor) {
        this.refs = elements
        this.menus = menus
        this.editor = editor
        this.$vms = {}

        this.createVms()

        this.createPopover()
    }

    createContent(menuKey) {
        const $vm = this.$vms[menuKey]
        if ($vm) {
            return $vm.$el
        }
        return ''
    }

    createVms() {
        ;(this.menus || []).forEach((item) => {
            if (item && item.popper) {
                const dom = document.createElement('div')
                const vNode = h(item.popper)
                render(vNode, dom)

                let $vm = vNode.component.proxy
                $vm.type = item.markType || item.nodeType

                this.bindEvent($vm)
                this.$vms[item.name] = $vm
            }
        })
    }

    bindEvent($vm) {
        $on($vm, 'on-update', (type, attrs) => this.updateNodeOrMarkAttrs(type, attrs))
    }

    updateNodeOrMarkAttrs(type, attrs) {
        let command = this.editor.commands[type.name]
        command && command(attrs)
    }

    createPopover() {
        if (!this.refs) {
            return
        }

        this.popper = tippy(this.refs, {
            arrow: false,
            content: (reference) => this.createContent(reference.getAttribute('data-key')),
            moveTransition: 'transform 0.4s cubic-bezier(0.22, 1, 0.36, 1)',
            theme: 'light',
            interactive: true,
        })
    }

    updateMenus(menus) {
        this.menus = menus
        ;(this.menus || []).forEach((item) => {
            const $vm = this.$vms[item.name] || null
            if ($vm && $vm.attrs) {
                $vm.attrs = item.markType ? (item.mark ? item.mark.attrs : {}) : {}
            }
        })
    }

    destroy() {
        if (this.popper && this.popper.destroy) {
            this.popper.destroy()
        }
    }
}
